# go-cache-sqlite

SQLite driver for the whosonfirst/go-cache interface

## Important

Work in progress. Documentation to follow.

## Usage

_Error handling omitted for the sake of brevity._

```
import (
	"context"
	"github.com/whosonfirst/go-cache"
	_ "github.com/whosonfirst/go-cache-sqlite"
	"io"
	"os"
	"strings"
)

func main() {

	ctx := context.Background()
	c, _ := cache.NewCache(ctx, "sqlite://?dsn=test.db")

	cache.SetString(ctx, c, "hello", "world")

	r, _ := c.Get(ctx, "hello")
	io.Copy(os.Stdout, r)
}
```

## See also

* https://github.com/whosonfirst/go-cache
* https://github.com/aaronland/go-sqlite