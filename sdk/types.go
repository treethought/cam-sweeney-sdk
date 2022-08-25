package sdk

import "fmt"

type paginatedResponse struct {
	Total  int
	Limit  int
	Offset int
	Page   int
	Pages  int
}

// APIError represents an error message provided by the API
type APIError struct {
	Success bool
	Message string
}

// SDKError represents an error when interacting with the API or SDK and provided details of any underlying APIError
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
