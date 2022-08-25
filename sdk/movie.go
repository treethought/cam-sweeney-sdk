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
func (m MoviesClient) List() ([]Movie, error) {
	resp := moviesResponse{}
	err := m.c.doRequestInto("/movie", &resp, WithAPIKey(m.c.apiKey))
	if err != nil {
		return nil, err
	}
	return resp.Docs, nil
}

// Get returns a single movie by ID
func (m MoviesClient) Get(id string) (Movie, error) {
	path := fmt.Sprintf("/movie/%s", id)
	resp := moviesResponse{}
	err := m.c.doRequestInto(path, &resp, WithAPIKey(m.c.apiKey))
	if err != nil {
		return Movie{}, err
	}
	return resp.Docs[0], err
}

// GetQuotes returns all quotes of a single movie
func (m MoviesClient) GetQuotes(id string) ([]Quote, error) {
	path := fmt.Sprintf("/movie/%s/quote", id)
	resp := quoteResponse{}
	err := m.c.doRequestInto(path, &resp, WithAPIKey(m.c.apiKey))
	if err != nil {
		return nil, err
	}

	if len(resp.Docs) == 0 {
		return nil, errors.New("no quotes available")
	}
	return resp.Docs, err
}
