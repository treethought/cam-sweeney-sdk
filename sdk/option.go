package sdk

import (
	"fmt"
	"net/http"
)

// RequestOption can be provided to API Calls to modify the request
type RequestOption func(*http.Request)

// PaginationOptions provide paginiation options
type PaginationOptions struct {
	Limit  int
	Page   int
	Offset int
}

func setQueryParam(req *http.Request, key string, val string) {
	values := req.URL.Query()
	values.Set(key, val)
	req.URL.RawQuery = values.Encode()
}

func withQuery(key string, val string) RequestOption {
	return func(req *http.Request) {
		setQueryParam(req, key, val)
	}
}

// WithLimit applies a limit to the number returned resources
func WithLimit(limit int) RequestOption {
	return withQuery("limit", fmt.Sprint(limit))
}

// WithLimit specifies a specific page to be returned
func WithPage(p int) RequestOption {
	return withQuery("page", fmt.Sprint(p))
}

// WithOffset specifies an offset to be applied to pagination
func WithOffset(offset int) RequestOption {
	return withQuery("offset", fmt.Sprint(offset))
}

// WithPagination applies pagination options to the request
func WithPagination(opts PaginationOptions) RequestOption {
	return func(req *http.Request) {
		if opts.Offset > 0 {
			setQueryParam(req, "offset", fmt.Sprint(opts.Offset))
		}
		if opts.Limit > 0 {
			setQueryParam(req, "limit", fmt.Sprint(opts.Limit))
		}
		if opts.Page > 0 {
			setQueryParam(req, "page", fmt.Sprint(opts.Page))
		}
	}
}

// WithSort indicates a field and direction to sort returned resources.
// field represents the field of an API resource, dir must be either "asc" or "dsc"
// for example, WithSort("realm", "asc")
func WithSort(field string, dir string) RequestOption {
	return func(req *http.Request) {
		val := fmt.Sprintf("%s:%s", field, dir)
		setQueryParam(req, "sort", val)
	}
}

// WithAPIKey causes an Authorization header to be applied to access authenticated resources
// Note that if an apiKey is provided to client creation, this is not needed and will be applied automatically.
// This option is useful for setting a new apiKey or turning a read-only client into an authenticated client
func WithAPIKey(apiKey string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	}
}
