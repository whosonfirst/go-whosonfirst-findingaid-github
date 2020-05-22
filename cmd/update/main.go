package main

/*

For example:

cat /usr/local/data/sfomuseum-data-flights-2020-05/filelist.txt | sed -e 's/^data\///' | go run -mod vendor cmd/update/main.go -reader-uri github://sfomuseum-data/sfomuseum-data-flights-2020-05 -stdin

*/

import (
	_ "github.com/whosonfirst/go-reader-github"
	_ "github.com/whosonfirst/go-reader-http"
)

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-findingaid/repo"
	"log"
	"net/url"
	"os"
)

func main() {

	cache_uri := flag.String("cache-uri", "gocache://", "A valid whosonfirst/go-cache.Cache URI.")
	reader_uri := flag.String("reader-uri", "", "A valid whosonfirst/go-reader.Reader URI.")
	stdin := flag.Bool("stdin", false, "Read input from STDIN.")

	flag.Parse()

	ctx := context.Background()

	fa_query := url.Values{}
	fa_query.Set("cache", *cache_uri)
	fa_query.Set("indexer", "null://")

	fa_uri := fmt.Sprintf("repo://?%s", fa_query.Encode())

	fa, err := repo.NewRepoFindingAid(ctx, fa_uri)

	if err != nil {
		log.Fatal(err)
	}

	r, err := reader.NewReader(ctx, *reader_uri)

	if err != nil {
		log.Fatalf("Failed to create reader, %v", err)
	}

	process := func(ctx context.Context, path string) error {

		fh, err := r.Read(ctx, path)

		if err != nil {
			return err
		}

		err = fa.IndexReader(ctx, fh)

		if err != nil {
			return err
		}

		return nil
	}

	if *stdin {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {

			path := scanner.Text()
			err := process(ctx, path)

			if err != nil {
				log.Fatalf("Failed to process %s, %v", path, err)
			}
		}

		err := scanner.Err()

		if err != nil {
			log.Fatalf("Failed to read from STDIN, %v", err)
		}

	} else {

		paths := flag.Args()

		for _, path := range paths {

			err := process(ctx, path)

			if err != nil {
				log.Fatalf("Failed to process %s, %v", path, err)
			}
		}
	}

}
