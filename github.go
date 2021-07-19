package github

import (
	"context"
	github_api "github.com/google/go-github/v27/github"
	"github.com/whosonfirst/go-whosonfirst-findingaid"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
	"log"
	"time"
)

// PopulateFindingAidForOrganization will update a go-whosonfirst-findingaid index (fa) derived from one or more GitHub repositories defined according to
// criteria in 'org' and 'list_opts'
func PopulateFindingAidForOrganization(ctx context.Context, fa findingaid.Indexer, org string, list_opts *organizations.ListOptions) error {

	t1 := time.Now()

	defer func() {
		log.Printf("Time to index all: %v\n", time.Since(t1))
	}()

	list_cb := func(repo *github_api.Repository) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		repo_url := repo.GetCloneURL()

		if repo_url == "" {
			return nil
		}

		t1 := time.Now()

		defer func() {
			log.Printf("Time to index %s: %v\n", repo_url, time.Since(t1))
		}()

		err := fa.IndexURIs(ctx, repo_url)

		if err != nil {
			log.Printf("Failed to index %s: %v\n", repo_url, err)
			return err
		}

		return nil
	}

	return organizations.ListReposWithCallback(org, list_opts, list_cb)
}
