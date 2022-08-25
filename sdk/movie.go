package sdk

import (
	"errors"
	"fmt"
)

type moviesResponse struct {
	Docs []Movie
}

// Movie repesetns a single movie resource
type Movie struct {
	ID                         string `json:"_id,omitempty"`
	Name                       string `json:"name,omitempty"`
	RuntimeInMinutes           int
	BudgetInMillions           int
	BoxOfficeRevenueInMillions float32
	AcademyAwardNominations    int
	AcademyAwardWins           int
	RottenTomatoesScore        float32
}

// MoviesClient provides methods for interacting with movie resources
type MoviesClient struct {
	c OneAPIClient
}

// List returns a list of all movies
func (m MoviesClient) List(opts ...RequestOption) ([]Movie, error) {
	resp := moviesResponse{}

	opts = m.c.appendOptsToAuth(opts...)
	err := m.c.doRequestInto("/movie", &resp, opts...)
	if err != nil {
		return nil, err
	}
	return resp.Docs, nil
}

// Get returns a single movie by ID
func (m MoviesClient) Get(id string, opts ...RequestOption) (Movie, error) {
	path := fmt.Sprintf("/movie/%s", id)
	resp := moviesResponse{}

	opts = m.c.appendOptsToAuth(opts...)
	err := m.c.doRequestInto(path, &resp, opts...)
	if err != nil {
		return Movie{}, err
	}
	return resp.Docs[0], err
}

// GetQuotes returns all quotes of a single movie
func (m MoviesClient) GetQuotes(id string, opts ...RequestOption) ([]Quote, error) {
	path := fmt.Sprintf("/movie/%s/quote", id)
	resp := quoteResponse{}

	opts = m.c.appendOptsToAuth(opts...)
	err := m.c.doRequestInto(path, &resp, opts...)
	if err != nil {
		return nil, err
	}

	if len(resp.Docs) == 0 {
		return nil, errors.New("no quotes available")
	}
	return resp.Docs, err
}
