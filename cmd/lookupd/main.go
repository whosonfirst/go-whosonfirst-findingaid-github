package main

import (
	_ "gocloud.dev/blob/fileblob"
)

import (
	_ "github.com/whosonfirst/go-cache-blob"
	_ "github.com/whosonfirst/go-cache-sqlite"
)

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-server"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-whosonfirst-findingaid"
	"github.com/whosonfirst/go-whosonfirst-findingaid/repo"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func NewLookupHandler(fa findingaid.FindingAid) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		q := req.URL.Query()
		str_id := q.Get("id")

		if str_id == "" {
			http.Error(rsp, "Missing id parameter", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			http.Error(rsp, "Invalid id parameter", http.StatusBadRequest)
			return
		}

		var fa_rsp repo.FindingAidResponse // but what if some other type of response...TBD

		err = fa.LookupID(ctx, id, &fa_rsp)

		if err != nil {

			if cache.IsCacheMiss(err) {
				http.Error(rsp, "Not found", http.StatusNotFound)
			} else {
				http.Error(rsp, "Failed to lookup ID", http.StatusInternalServerError)
			}

			return
		}

		rsp.Header().Set("Content-Type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(fa_rsp)

		if err != nil {
			http.Error(rsp, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		return
	}

	h := http.HandlerFunc(fn)
	return h
}

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080", "...")
	cache_uri := flag.String("cache-uri", "gocache://", "...")

	flag.Parse()

	ctx := context.Background()

	fa_query := url.Values{}
	fa_query.Set("cache", *cache_uri)
	fa_query.Set("indexer", "null://")

	fa_uri := fmt.Sprintf("repo://?%s", fa_query.Encode())

	fa, err := repo.NewRepoFindingAid(ctx, fa_uri)

	if err != nil {
		log.Fatalf("Failed to create new finding aid, %v", err)
	}

	lookup_handler := NewLookupHandler(fa)

	mux := http.NewServeMux()
	mux.Handle("/", lookup_handler)

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Failed to create new server, %v", err)
	}

	log.Printf("Listening on %s", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
