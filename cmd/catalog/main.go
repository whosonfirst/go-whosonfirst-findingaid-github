package main

import (
	_ "gocloud.dev/blob/fileblob"
)

import (
	_ "github.com/whosonfirst/go-cache-blob"
	_ "github.com/whosonfirst/go-cache-sqlite"
	_ "github.com/whosonfirst/go-whosonfirst-index-git"
)

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-findingaid-github"
	"github.com/whosonfirst/go-whosonfirst-findingaid/repo"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
	"log"
	"net/url"
)

func main() {

	org := flag.String("org", "whosonfirst-data", "The name of the organization to clone repositories from")
	prefix := flag.String("prefix", "whosonfirst-data", "Limit repositories to only those with this prefix")
	exclude := flag.String("exclude", "", "Exclude repositories with this prefix")
	// updated_since := flag.String("updated-since", "", "A valid Unix timestamp or an ISO8601 duration string (months are currently not supported)")
	forked := flag.Bool("forked", false, "Only include repositories that have been forked")
	not_forked := flag.Bool("not-forked", false, "Only include repositories that have not been forked")
	token := flag.String("token", "", "A valid GitHub API access token")

	cache_uri := flag.String("cache-uri", "gocache://", "...")
	git_uri := flag.String("git-uri", "git://", "...")

	flag.Parse()

	list_opts := organizations.NewDefaultListOptions()

	list_opts.Prefix = *prefix
	list_opts.Exclude = *exclude
	list_opts.Forked = *forked
	list_opts.NotForked = *not_forked
	list_opts.AccessToken = *token

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fa_query := url.Values{}
	fa_query.Set("cache", *cache_uri)
	fa_query.Set("indexer", *git_uri)

	fa_uri := fmt.Sprintf("repo://?%s", fa_query.Encode())

	fa, err := repo.NewRepoFindingAid(ctx, fa_uri)

	if err != nil {
		log.Fatal(err)
	}

	err = github.PopulateFindingAidForOrganization(ctx, fa, *org, list_opts)

	if err != nil {
		log.Fatal(err)
	}
}
