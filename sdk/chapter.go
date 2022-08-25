package sdk

import "fmt"

type chapterResponse struct {
	paginatedResponse
	Docs []Chapter
}

// Chapter represents a single book chapter
type Chapter struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"chapterName,omitempty"`
	Book string `json:"book,omitempty"`
}

// ChapterClient provides methods for interacting with chapter resources
type ChapterClient struct {
	c OneAPIClient
}

// List provides all chapters across all books
func (ch ChapterClient) List(opts ...RequestOption) ([]Chapter, error) {
	resp := chapterResponse{}
	opts = ch.c.appendOptsToAuth(opts...)
	err := ch.c.doRequestInto("/chapter", &resp, opts...)
	if err != nil {
		return nil, err
	}
	return resp.Docs, nil
}

// Get returns a single chapter by ID
func (ch ChapterClient) Get(id string, opts ...RequestOption) (Chapter, error) {
	path := fmt.Sprintf("/chapter/%s", id)
	resp := chapterResponse{}

	opts = ch.c.appendOptsToAuth(opts...)
	err := ch.c.doRequestInto(path, &resp, opts...)
	if err != nil {
		return Chapter{}, err
	}
	return resp.Docs[0], err
}
