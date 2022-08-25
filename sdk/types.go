package sdk

import "fmt"

type Book struct {
	ID   string `json:"_id,omitempty"`
	Name string `json:"name,omitempty"`
}

type listBooksResponse struct {
	APIResponse
	Docs []Book
}

type paginatedResponse struct {
	Total  int
	Limit  int
	Offset int
	Page   int
	Pages  int
}

type APIResponse struct {
	APIError
	paginatedResponse
}

func (r APIResponse) IsError() error {
	return r.APIError
}

type APIError struct {
	Success bool
	Message string
}

type SDKError struct {
	message  string
	endpoint string
	apiError APIError
}

func (e APIError) Error() string {
	return e.Message
}

func (e SDKError) Error() string {
	return fmt.Sprintf("%s %s: %v", e.message, e.endpoint, e.apiError)
}
