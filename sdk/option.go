package sdk

import (
	"fmt"
	"net/http"
)

type RequestOption func(*http.Request)

type PaginationOptions struct {
	Limit  int
	Page   int
	Offset int
}

func withQuery(key string, val string) RequestOption {
	return func(req *http.Request) {
		req.URL.Query().Add(key, val)
	}
}

func WithLimit(limit int) RequestOption {
	return withQuery("limit", fmt.Sprint(limit))
}

func WithPage(p int) RequestOption {
	return withQuery("page", fmt.Sprint(p))
}

func WithOffset(offset int) RequestOption {
	return withQuery("page", fmt.Sprint(offset))
}

func WithPagination(opts PaginationOptions) RequestOption {
	return func(req *http.Request) {
		if opts.Offset > 0 {
			req.URL.Query().Add("offset", fmt.Sprint(opts.Offset))
		}
		if opts.Limit > 0 {
			req.URL.Query().Add("limit", fmt.Sprint(opts.Limit))
		}
		if opts.Page > 0 {
			req.URL.Query().Add("page", fmt.Sprint(opts.Page))
		}
	}
}

func WithSort(field string, dir string) RequestOption {
	return func(req *http.Request) {
		val := fmt.Sprintf("%s:%s", field, dir)
		req.URL.Query().Add("sort", fmt.Sprint(val))
	}
}

func WithAPIKey(apiKey string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	}
}
