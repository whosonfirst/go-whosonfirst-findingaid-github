package main

/*

For example:

cat /usr/local/data/sfomuseum-data-flights-2020-05/filelist.txt | \
sed -e 's/^data\///' | \
go run -mod vendor cmd/update/main.go -reader-uri github://sfomuseum-data/sfomuseum-data-flights-2020-05 -mode stdin

*/

import (
	_ "github.com/aaronland/go-cloud-s3blob"
	_ "github.com/whosonfirst/go-reader-github"
	_ "github.com/whosonfirst/go-reader-http"
)

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sfomuseum/go-flags/flagset"
	cache_blob "github.com/whosonfirst/go-cache-blob"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-findingaid/repo"
	"gocloud.dev/blob"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {

	fs := flagset.NewFlagSet("findingaid")

	cache_uri := fs.String("cache-uri", "gocache://", "A valid whosonfirst/go-cache.Cache URI.")
	reader_uri := fs.String("reader-uri", "", "A valid whosonfirst/go-reader.Reader URI.")

	mode := fs.String("mode", "cli", "...")

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVarsWithFeedback(fs, "FINDINGAID", true)

	if err != nil {
		log.Fatalf("Failed to set flags from env vars, %v", err)
	}

	ctx := context.Background()

	fa_query := url.Values{}
	fa_query.Set("cache", *cache_uri)
	fa_query.Set("indexer", "null://")

	fa_uri := fmt.Sprintf("repo://?%s", fa_query.Encode())

	fa, err := repo.NewRepoFindingAid(ctx, fa_uri)

	if err != nil {
		log.Fatal(err)
	}

	// START OF please put me in a common function to share with
	// bulk cataloging (20200526/thisisaaronland)

	before := func(asFunc func(interface{}) bool) error {

		req := &s3manager.UploadInput{}
		ok := asFunc(&req)

		if !ok {
			return errors.New("invalid s3 type")
		}

		req.ACL = aws.String("public-read")
		return nil
	}

	wr_opts := &blob.WriterOptions{
		BeforeWrite: before,
		ContentType: "application/json",
	}

	process := func(ctx context.Context, r reader.Reader, paths ...string) error {

		ctx = context.WithValue(ctx, cache_blob.BlobCacheOptionsKey("options"), wr_opts)

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		remaining := len(paths)

		done_ch := make(chan bool)
		err_ch := make(chan error)

		for _, path := range paths {

			go func(path string) {

				defer func() {
					done_ch <- true
				}()

				fh, err := r.Read(ctx, path)

				if err != nil {
					log.Println("READ ERROR", path, err)
					err_ch <- err
					return
				}

				err = fa.IndexReader(ctx, fh)

				if err != nil {
					log.Println("INDEX ERROR", path, err)
					err_ch <- err
					return
				}
			}(path)

		}

		for remaining > 0 {
			select {
			case <-ctx.Done():
				return nil
			case <-done_ch:
				remaining -= 1
			case err := <-err_ch:
				return err
			}
		}

		return nil
	}

	// END OF please put me in a common function to share with

	switch *mode {

	case "cli":

		r, err := reader.NewReader(ctx, *reader_uri)

		if err != nil {
			log.Fatalf("Failed to create reader, %v", err)
		}

		paths := flag.Args()

		err = process(ctx, r, paths...)

		if err != nil {
			log.Fatalf("Failed to process paths, %v", err)
		}

	case "lambda":

		// this expects to be passed the output of the
		// whosonfirst/go-webhookd/transformations/github.commits
		// thingy (20200526/thisisaaronland)

		handler := func(ctx context.Context, payload []byte) error {

			raw := string(payload)

			// TBD - does this CSV parsing code need to be shared
			// with anything else, specifically reading from STDIN
			// or is it enough to expect people to grep/cut/sed all
			// the things? (20200526/thisisaaronland)

			fh := strings.NewReader(raw)
			r := csv.NewReader(fh)

			to_update := make(map[string][]string)

			for {
				row, err := r.Read()

				if err == io.EOF {
					break
				}

				if err != nil {
					return err
				}

				if len(row) != 3 {
					msg := fmt.Sprintf("Invalid row, %v", row)
					return errors.New(msg)
				}

				repo := row[1]
				path := row[2]

				path = strings.Replace(path, "data/", "", 1)

				paths, ok := to_update[repo]

				if !ok {
					paths = make([]string, 0)
				}

				paths = append(paths, path)
				to_update[repo] = paths
			}

			for repo, paths := range to_update {

				// please update to use URI templates
				r_uri := strings.Replace(*reader_uri, "{repo}", repo, 1)

				r, err := reader.NewReader(ctx, r_uri)

				if err != nil {
					return err
				}

				err = process(ctx, r, paths...)

				if err != nil {
					return err
				}
			}

			return nil
		}

		lambda.Start(handler)

	case "stdin":

		// TBD - what if more than one repo? See notes above about
		// CSV code in lambda handler (20200526/thisisaaronland)

		r, err := reader.NewReader(ctx, *reader_uri)

		if err != nil {
			log.Fatalf("Failed to create reader, %v", err)
		}

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {

			path := scanner.Text()
			err := process(ctx, r, path)

			if err != nil {
				log.Fatalf("Failed to process %s, %v", path, err)
			}
		}

		err = scanner.Err()

		if err != nil {
			log.Fatalf("Failed to read from STDIN, %v", err)
		}

	default:
		log.Fatalf("Invalid mode")
	}
}
