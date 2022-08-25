package sdk

import "fmt"

type chapterResponse struct {
	Docs []Chapter
}

type Chapter struct {
	ID   string `json:"_id,omitempty"`
	Name string `json:"chapterName,omitempty"`
}

type ChapterClient struct {
	c OneAPIClient
}

func (ch ChapterClient) List() ([]Chapter, error) {
	resp := chapterResponse{}
	err := ch.c.doRequestInto("/chapter", &resp, WithAPIKey(ch.c.apiKey))
	if err != nil {
		return nil, err
	}
	return resp.Docs, nil
}

func (ch ChapterClient) Get(id string) (Chapter, error) {
	path := fmt.Sprintf("/chapter/%s", id)
	resp := chapterResponse{}
	err := ch.c.doRequestInto(path, &resp, WithAPIKey(ch.c.apiKey))
	if err != nil {
		return Chapter{}, err
	}
	return resp.Docs[0], err
}
