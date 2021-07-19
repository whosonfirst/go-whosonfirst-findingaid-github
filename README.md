# go-whosonfirst-findingaid-github

Go package providing tools and methods for working with go-whosonfirst-findingaid indices derived from one or more GitHub repositories.

## Important

Work in progress. Documentation to follow.

## Example

```
package main

import (
	_ "github.com/whosonfirst/go-cache"
	_ "github.com/whosonfirst/go-whosonfirst-index-git"
)

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-findingaid-github"
	"github.com/whosonfirst/go-whosonfirst-findingaid/repo"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fa_uri := "repo://?cache=gocache://&indexer=git://"
	fa, _ := repo.NewRepoFindingAid(ctx, fa_uri)

	org := "whosonfirst-data"
	list_opts := organizations.NewDefaultListOptions()

	github.PopulateFindingAidForOrganization(ctx, fa, org, list_opts)
}
```

_Error handling omitted for the sake of brevity._

## Tools

### populate

```
> go run -mod vendor cmd/catalog/main.go \
	-cache-uri 'file:///usr/local/whosonfirst/go-whosonfirst-findingaid-github/test' \
	-org sfomuseum-data \
	-prefix sfomuseum-data-collection
	
2020/05/20 16:52:13 Time to index https://github.com/sfomuseum-data/sfomuseum-data-collection-classifications.git: 2.756128442s
2020/05/20 16:53:07 Time to index https://github.com/sfomuseum-data/sfomuseum-data-collection.git: 54.316060586s
2020/05/20 16:53:07 Time to index all: 1m0.292242361s

> cat test/1511214253 | jq
{
  "id": 1511214253,
  "repo": "sfomuseum-data-collection-classifications",
  "path": "151/121/425/3/1511214253.geojson"
}
```

### update

Update an indexing catalog from a list of files. For example:

```
$> cd /usr/local/data/sfomuseum-data-flights-2020-05
$> git log --name-only --pretty=format:'' HEAD^..HEAD | grep -v alt > filelist.txt

$> cat filelist.txt | \
	sed -e 's/^data\///' | \
	go run -mod vendor cmd/catalog/main.go -stdin \
	-reader-uri github://sfomuseum-data/sfomuseum-data-flights-2020-05 
```

_This tool will likely change still, specifically to make sure it works with the output of the [go-webhookd GitHubCommits](https://github.com/whosonfirst/go-webhookd#githubcommits) transformation._

## Available caching layers

Anything registered by:

* https://github.com/whosonfirst/go-cache
* https://github.com/whosonfirst/go-cache-blob
* https://github.com/whosonfirst/go-cache-sqlite

## Available readers

Anything registered by:

* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-reader-github
* https://github.com/whosonfirst/go-reader-http

## See also

* https://github.com/whosonfirst/go-whosonfirst-findingaid
* https://github.com/whosonfirst/go-whosonfirst-index-git
* https://github.com/whosonfirst/go-whosonfirst-github
* https://en.wikipedia.org/wiki/Finding_aid