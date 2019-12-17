# go-whosonfirst-findingaid-github

## Important

Work in progress.

## Example

```
package main

import (
	"context"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-whosonfirst-findingaid-github"
	"github.com/whosonfirst/go-whosonfirst-findingaid/git"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	org := "whosonfirst-data"

	list_opts := organizations.NewDefaultListOptions()

	c, _ := cache.NewCache(ctx, "gocache://")

	fa, _ := git.NewRepoFindingAidWithCache(ctx, c)

	github.PopulateFindingAidForOrganization(ctx, fa, org, list_opts)
}
```

_Error handling omitted for the sake of brevity._

## Tools

### findingaid

_This is a dumb name. It will be changed (once the dust has settled)._

```
go run -mod vendor cmd/findingaid/main.go -prefix whosonfirst-data-admin-
2019/12/17 13:50:33 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-ad.git: 2.389023226s
2019/12/17 13:50:33 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-ae.git: 735.076745ms
2019/12/17 13:50:43 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-af.git: 9.758202233s
2019/12/17 13:50:44 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-ag.git: 334.155527ms
2019/12/17 13:50:44 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-ai.git: 359.469678ms
2019/12/17 13:50:46 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-al.git: 2.025303835s
2019/12/17 13:50:48 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-am.git: 2.12029687s
2019/12/17 13:50:48 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-an.git: 413.589471ms
...

```

## See also

* https://github.com/whosonfirst/go-whosonfirst-findingaid
* https://github.com/whosonfirst/go-whosonfirst-index-git
* https://github.com/whosonfirst/go-whosonfirst-github
* https://en.wikipedia.org/wiki/Finding_aid