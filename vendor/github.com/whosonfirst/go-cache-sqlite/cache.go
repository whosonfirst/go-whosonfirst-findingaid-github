package sqlite

import (
	"context"
	"fmt"
	aa_sqlite "github.com/aaronland/go-sqlite"
	"io"
	"io/ioutil"
	_ "log"
)

type CacheRecord struct {
	Key  string
	Body io.Reader
}

type CacheTable struct {
	aa_sqlite.Table
	name string
}

func NewCacheTableWithDatabase(ctx context.Context, db aa_sqlite.Database) (aa_sqlite.Table, error) {

	t, err := NewCacheTable(ctx)

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(ctx, db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewCacheTable(ctx context.Context) (aa_sqlite.Table, error) {

	t := CacheTable{
		name: "cache",
	}

	return &t, nil
}

func (t *CacheTable) Name() string {
	return t.name
}

func (t *CacheTable) Schema() string {

	// https://github.com/straup/go-aaronland-iiif/issues/12

	sql := `CREATE TABLE %s (
		key TEXT NOT NULL PRIMARY KEY,
		body TEXT NOT NULL
	);
	`

	return fmt.Sprintf(sql, t.Name())
}

func (t *CacheTable) InitializeTable(ctx context.Context, db aa_sqlite.Database) error {
	return aa_sqlite.CreateTableIfNecessary(ctx, db, t)
}

func (t *CacheTable) IndexRecord(ctx context.Context, db aa_sqlite.Database, i interface{}) error {

	rec := i.(CacheRecord)

	key := rec.Key

	body, err := ioutil.ReadAll(rec.Body)

	if err != nil {
		return err
	}

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO %s (
		key, body
	) VALUES (
		?, ?
	)`, t.Name())

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(key, body)

	if err != nil {
		return err
	}

	return tx.Commit()
}
