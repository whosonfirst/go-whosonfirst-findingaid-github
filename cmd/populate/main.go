// catalog will populate a go-whosonfirst-findingaid index derived from one or more GitHub repositories
package main

import (
	_ "gocloud.dev/blob/fileblob"
)

import (
	_ "github.com/whosonfirst/go-cache-sqlite"
	_ "github.com/whosonfirst/go-whosonfirst-iterate-git"
)

import (
	"context"
	"flag"
	_ "fmt"
	"github.com/aaronland/gocloud-blob-s3"
	"github.com/sfomuseum/go-flags/multi"
	cache_blob "github.com/whosonfirst/go-cache-blob"
	"github.com/whosonfirst/go-whosonfirst-findingaid-github"
	"github.com/whosonfirst/go-whosonfirst-findingaid/repo"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
	"log"
	"net/url"
)

func main() {

	var prefix multi.MultiString
	flag.Var(&prefix, "prefix", "Limit repositories to only those with this prefix")

	var exclude multi.MultiString
	flag.Var(&exclude, "exclude", "Exclude repositories with this prefix")

	org := flag.String("org", "whosonfirst-data", "The name of the organization to clone repositories from")

	forked := flag.Bool("forked", false, "Only include repositories that have been forked")
	not_forked := flag.Bool("not-forked", false, "Only include repositories that have not been forked")
	token := flag.String("token", "", "A valid GitHub API access token")

	cache_uri := flag.String("cache-uri", "gocache://", "A valid whosonfirst/go-cache URI string.")
	iterator_uri := flag.String("iterator-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate URI string.")

	findingaid_uri := flag.String("findingaid-uri", "repo://?cache={cache_uri}&iterator={iterator_uri}", "A valid whosonfirst/go-whosonfirst-findingaid URI string.")

	
	flag.Parse()

	ctx := context.Background()

	list_opts := organizations.NewDefaultListOptions()

	list_opts.Prefix = prefix
	list_opts.Exclude = exclude
	list_opts.Forked = *forked
	list_opts.NotForked = *not_forked
	list_opts.AccessToken = *token

	//

	ctx_key := cache_blob.BlobCacheOptionsKey("options")

	ctx_opts := map[string]interface{}{
		"ACL":         "public-read",
		"ContentType": "application/json",
	}

	ctx, err := s3blob.SetWriterOptionsWithContextAndMap(ctx, ctx_key, ctx_opts)

	if err != nil {
		log.Fatalf("Failed to set writer options, %v", err)
	}

	//

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()


	fa_uri, err := url.Parse(*findingaid_uri)

	if err != nil {
		log.Fatalf("Failed to parse findingaid URI, %v", err)
	}

	fa_q := fa_uri.Query()

	if fa_q.Get("cache") == "{cache_uri}" {
		fa_q["cache"] = []string{*cache_uri}
	}

	if fa_q.Get("iterator") == "{iterator_uri}" {
		fa_q["iterator"] = []string{*iterator_uri}
	}

	if fa_q.Get("iterator") == "" {
		log.Fatalf("Missing '-iterator-uri' flag.")
	}

	fa_uri.RawQuery = fa_q.Encode()

	log.Println(fa_uri.String())
	fa, err := repo.NewIndexer(ctx, fa_uri.String())

	if err != nil {
		log.Fatal("Failed to create finding aid indexer, %v", err)
	}

	err = github.PopulateFindingAidForOrganization(ctx, fa, *org, list_opts)

	if err != nil {
		log.Fatal("Failed to populate finding aid, %v", err)
	}
}
