package sdk

import "fmt"

type quoteResponse struct {
	Docs []Quote
}

type Quote struct {
	ID        string `json:"_id,omitempty"`
	Character string
	Dialog    string
}

type QuotesClient struct {
	c OneAPIClient
}

// ListBooks returns a list of all "Lord of the Rings" books
func (q QuotesClient) List() ([]Quote, error) {
	resp := quoteResponse{}
	err := q.c.doRequestInto("/quote", &resp)
	if err != nil {
		return nil, err
	}
	return resp.Docs, err
}

// Get a book by it's ID
func (q QuotesClient) Get(id string) (Quote, error) {
	path := fmt.Sprintf("/quote/%s", id)
	resp := quoteResponse{}
	err := q.c.doRequestInto(path, &resp)
	if err != nil {
		return Quote{}, err
	}
	return resp.Docs[0], err
}
