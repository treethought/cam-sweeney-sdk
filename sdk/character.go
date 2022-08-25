package sdk

import (
	"errors"
	"fmt"
)

type characterResponse struct {
	paginatedResponse
	Docs []Character
}

// Character represents a single character
type Character struct {
	ID      string `json:"_id"`
	Birth   string `json:"birth"`
	Death   string `json:"death"`
	Gender  string `json:"gender"`
	Height  string `json:"height"`
	Realm   string `json:"realm"`
	Spouse  string `json:"spouse"`
	Name    string `json:"name"`
	Race    string `json:"race"`
	WikiUrl string `json:"wikiUrl"`
}

// CharactersClient provides methods for interacting with character resources
type CharactersClient struct {
	c OneAPIClient
}

// List returns a list of all characters
func (ch CharactersClient) List(opts ...RequestOption) ([]Character, error) {
	resp := characterResponse{}

	opts = ch.c.appendOptsToAuth(opts...)
	err := ch.c.doRequestInto("/character", &resp, opts...)
	if err != nil {
		return nil, err
	}
	return resp.Docs, err
}

// Get returns a single Character by ID
func (ch CharactersClient) Get(id string, opts ...RequestOption) (Character, error) {
	path := fmt.Sprintf("/character/%s", id)
	resp := characterResponse{}

	opts = ch.c.appendOptsToAuth(opts...)
	err := ch.c.doRequestInto(path, &resp, opts...)
	if err != nil {
		return Character{}, err
	}
	return resp.Docs[0], err
}

// GetQuotes returns a all quotes of a single Character by ID
func (ch CharactersClient) GetQuotes(id string, opts ...RequestOption) ([]Quote, error) {
	path := fmt.Sprintf("/character/%s/quote", id)
	resp := quoteResponse{}

	opts = ch.c.appendOptsToAuth(opts...)
	err := ch.c.doRequestInto(path, &resp, opts...)
	if err != nil {
		return nil, err
	}
	if len(resp.Docs) == 0 {
		return nil, errors.New("no quotes available")
	}
	return resp.Docs, err
}
