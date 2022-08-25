package sdk

import "fmt"

type characterResponse struct {
	paginatedResponse
	Docs []Character
}

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

type CharactersClient struct {
	c OneAPIClient
}

// ListBooks returns a list of all "Lord of the Rings" books
func (ch CharactersClient) List() ([]Character, error) {
	resp := characterResponse{}
	err := ch.c.doRequestInto("/character", &resp, WithAPIKey(ch.c.apiKey))
	if err != nil {
		return nil, err
	}
	return resp.Docs, err
}

func (ch CharactersClient) Get(id string) (Character, error) {
	path := fmt.Sprintf("/character/%s", id)
	resp := characterResponse{}
	err := ch.c.doRequestInto(path, &resp, WithAPIKey(ch.c.apiKey))
	if err != nil {
		return Character{}, err
	}
	return resp.Docs[0], err
}
