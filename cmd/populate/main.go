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
	"fmt"
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

	// updated_since := flag.String("updated-since", "", "A valid Unix timestamp or an ISO8601 duration string (months are currently not supported)")

	forked := flag.Bool("forked", false, "Only include repositories that have been forked")
	not_forked := flag.Bool("not-forked", false, "Only include repositories that have not been forked")
	token := flag.String("token", "", "A valid GitHub API access token")
	
	cache_uri := flag.String("cache-uri", "gocache://", "A valid whosonfirst/go-cache URI.")
	git_uri := flag.String("git-uri", "git://", "A valid whosonfirst/go-whosonfirst-iterate/emitter URI.")

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

	fa_query := url.Values{}
	fa_query.Set("cache", *cache_uri)
	fa_query.Set("iterator", *git_uri)

	fa_uri := fmt.Sprintf("repo://?%s", fa_query.Encode())

	fa, err := repo.NewIndexer(ctx, fa_uri)

	if err != nil {
		log.Fatal("Failed to create finding aid indexer, %v", err)
	}

	err = github.PopulateFindingAidForOrganization(ctx, fa, *org, list_opts)

	if err != nil {
		log.Fatal("Failed to populate finding aid, %v", err)
	}
}
