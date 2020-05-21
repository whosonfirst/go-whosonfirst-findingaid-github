# go-whosonfirst-findingaid-github

## Important

Work in progress.

## Example

```
package main

import (
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

	org := "whosonfirst-data"

	list_opts := organizations.NewDefaultListOptions()

	fa_uri := "repo://?cache=gocache://&indexer=git://"

	fa, _ := repo.NewRepoFindingAid(ctx, fa_uri)

	github.PopulateFindingAidForOrganization(ctx, fa, org, list_opts)
}
```

_Error handling omitted for the sake of brevity._

## Tools

### catalog

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

### lookupd

```
> go run -mod vendor cmd/lookupd/main.go \
	-cache-uri 'file:///usr/local/whosonfirst/go-whosonfirst-findingaid-github/test'
	
2020/05/20 17:21:26 Listening on http://localhost:8080

> curl -s 'http://localhost:8080?id=1511214253' | jq

{
  "id": 1511214253,
  "repo": "sfomuseum-data-collection-classifications",
  "path": "151/121/425/3/1511214253.geojson"
}
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-findingaid
* https://github.com/whosonfirst/go-whosonfirst-index-git
* https://github.com/whosonfirst/go-whosonfirst-github
* https://en.wikipedia.org/wiki/Finding_aid