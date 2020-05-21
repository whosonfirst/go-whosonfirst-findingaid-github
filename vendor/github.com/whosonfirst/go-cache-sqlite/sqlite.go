package sqlite

import (
	"context"
	"fmt"
	aa_sqlite "github.com/aaronland/go-sqlite"
	aa_database "github.com/aaronland/go-sqlite/database"
	"github.com/whosonfirst/go-cache"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"sync/atomic"
)

type SQLiteCache struct {
	cache.Cache
	db        aa_sqlite.Database
	cache     aa_sqlite.Table
	hits      int64
	misses    int64
	evictions int64
}

func init() {

	ctx := context.Background()
	err := cache.RegisterCache(ctx, "sqlite", NewSQLiteCache)

	if err != nil {
		panic(err)
	}
}

func NewSQLiteCache(ctx context.Context, uri string) (cache.Cache, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()
	dsn := q.Get("dsn")

	db, err := aa_database.NewDB(ctx, dsn)

	if err != nil {
		return nil, err
	}

	err = db.LiveHardDieFast()

	if err != nil {
		return nil, err
	}

	cache_tbl, err := NewCacheTableWithDatabase(ctx, db)

	if err != nil {
		return nil, err
	}

	c := &SQLiteCache{
		db:    db,
		cache: cache_tbl,
	}

	return c, nil
}

func (c *SQLiteCache) Name() string {
	return "sqlite"
}

func (c *SQLiteCache) Close(ctx context.Context) error {
	return c.db.Close()
}

func (c *SQLiteCache) Get(ctx context.Context, key string) (io.ReadCloser, error) {

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

	conn, err := c.db.Conn()

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf("SELECT body FROM %s WHERE key = ?", c.cache.Name())

	row := conn.QueryRowContext(ctx, q, key)

	var body string

	err = row.Scan(&body)

	if err != nil {
		atomic.AddInt64(&c.misses, 1)
		return nil, err
	}

	atomic.AddInt64(&c.hits, 1)

	fh := strings.NewReader(body)
	cl := ioutil.NopCloser(fh)

	return cl, nil
}

func (c *SQLiteCache) Set(ctx context.Context, key string, fh io.ReadCloser) (io.ReadCloser, error) {

	/*
		body, err := ioutil.ReadAll(fh)

		if err != nil {
			return nil, err
		}

		br := bytes.NewReader(body)
		cl := ioutil.NopCloser(br)
	*/

	rec := CacheRecord{
		Key:  key,
		Body: fh,
	}

	err := c.cache.IndexRecord(ctx, c.db, rec)

	if err != nil {
		return nil, err
	}

	return c.Get(ctx, key)
}

func (c *SQLiteCache) Unset(ctx context.Context, key string) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	conn, err := c.db.Conn()

	if err != nil {
		return err
	}

	q := fmt.Sprintf("DELETE FROM %s WHERE key = ?", c.cache.Name())

	_, err = conn.ExecContext(ctx, q, key)

	if err != nil {
		return err
	}

	atomic.AddInt64(&c.evictions, -1)
	return nil
}

func (c *SQLiteCache) Hits() int64 {
	return atomic.LoadInt64(&c.hits)
}

func (c *SQLiteCache) Misses() int64 {
	return atomic.LoadInt64(&c.misses)
}

func (c *SQLiteCache) Evictions() int64 {
	return atomic.LoadInt64(&c.evictions)
}

func (c *SQLiteCache) Size() int64 {
	ctx := context.Background()
	return c.SizeWithContext(ctx)
}

func (c *SQLiteCache) SizeWithContext(ctx context.Context) int64 {

	conn, err := c.db.Conn()

	if err != nil {
		return -1
	}

	q := fmt.Sprintf("SELECT COUNT(key) FROM %s", c.cache.Name())

	row := conn.QueryRowContext(ctx, q)

	var size int64

	err = row.Scan(&size)

	if err != nil {
		return -1
	}

	return size
}
