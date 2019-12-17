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
2019/12/17 14:46:44 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-us.git: 8m29.355327565s
2019/12/17 14:46:48 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-uy.git: 3.216698727s
2019/12/17 14:46:54 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-uz.git: 6.639917211s
2019/12/17 14:46:55 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-va.git: 350.202375ms
2019/12/17 14:46:55 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-vc.git: 385.954555ms
2019/12/17 14:47:15 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-ve.git: 19.566929432s
2019/12/17 14:47:15 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-vg.git: 326.524022ms
2019/12/17 14:47:16 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-vi.git: 723.616302ms
2019/12/17 14:47:46 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-vn.git: 30.542778246s
2019/12/17 14:47:48 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-vu.git: 1.706213031s
2019/12/17 14:47:48 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-wf.git: 462.201482ms
2019/12/17 14:47:49 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-ws.git: 879.747161ms
2019/12/17 14:47:54 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-xk.git: 4.601361461s
2019/12/17 14:47:54 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-xn.git: 437.351844ms
2019/12/17 14:47:55 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-xs.git: 485.617434ms
2019/12/17 14:48:12 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-xx.git: 17.098112165s
2019/12/17 14:48:45 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-xy.git: 32.986311102s
2019/12/17 14:48:46 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-xz.git: 568.386627ms
2019/12/17 14:49:21 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-ye.git: 35.655171132s
2019/12/17 14:49:22 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-yt.git: 798.471424ms
2019/12/17 14:49:55 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-za.git: 33.054950953s
2019/12/17 14:50:12 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-zm.git: 16.964481284s
2019/12/17 14:50:17 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-zw.git: 4.409956503s
2019/12/17 14:50:17 Time to index https://github.com/whosonfirst-data/whosonfirst-data-admin-alt.git: 318.480904ms
2019/12/17 14:50:17 Time to index all: 59m55.499986814s
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-findingaid
* https://github.com/whosonfirst/go-whosonfirst-index-git
* https://github.com/whosonfirst/go-whosonfirst-github
* https://en.wikipedia.org/wiki/Finding_aid