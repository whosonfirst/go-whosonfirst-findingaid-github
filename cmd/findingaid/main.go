package main

import (
	"context"
	"flag"

	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-whosonfirst-findingaid-github"
	"github.com/whosonfirst/go-whosonfirst-findingaid/git"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
	"log"
)

func main() {

	org := flag.String("org", "whosonfirst-data", "The name of the organization to clone repositories from")
	prefix := flag.String("prefix", "whosonfirst-data", "Limit repositories to only those with this prefix")
	exclude := flag.String("exclude", "", "Exclude repositories with this prefix")
	// updated_since := flag.String("updated-since", "", "A valid Unix timestamp or an ISO8601 duration string (months are currently not supported)")
	forked := flag.Bool("forked", false, "Only include repositories that have been forked")
	not_forked := flag.Bool("not-forked", false, "Only include repositories that have not been forked")
	token := flag.String("token", "", "A valid GitHub API access token")

	flag.Parse()

	list_opts := organizations.NewDefaultListOptions()

	list_opts.Prefix = *prefix
	list_opts.Exclude = *exclude
	list_opts.Forked = *forked
	list_opts.NotForked = *not_forked
	list_opts.AccessToken = *token

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := cache.NewCache(ctx, "gocache://")

	if err != nil {
		log.Fatal(err)
	}

	fa, err := git.NewRepoFindingAidWithCache(ctx, c)

	if err != nil {
		log.Fatal(err)
	}

	err = github.PopulateFindingAidForOrganization(ctx, fa, *org, list_opts)

	if err != nil {
		log.Fatal(err)
	}
}
