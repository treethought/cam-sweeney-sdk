package sdk

import "fmt"

type quoteResponse struct {
	paginatedResponse
	Docs []Quote
}

// Quote represents a quote spoken by a character
type Quote struct {
	ID        string `json:"_id,omitempty"`
	Character string
	Dialog    string
}

// QuotesClientt provides methods for interacting with quote resources
type QuotesClient struct {
	c OneAPIClient
}

// List returns a list of all quotes
func (q QuotesClient) List() ([]Quote, error) {
	resp := quoteResponse{}
	err := q.c.doRequestInto("/quote", &resp)
	if err != nil {
		return nil, err
	}
	return resp.Docs, err
}

// Get returns a quote by ID
func (q QuotesClient) Get(id string) (Quote, error) {
	path := fmt.Sprintf("/quote/%s", id)
	resp := quoteResponse{}
	err := q.c.doRequestInto(path, &resp)
	if err != nil {
		return Quote{}, err
	}
	return resp.Docs[0], err
}
