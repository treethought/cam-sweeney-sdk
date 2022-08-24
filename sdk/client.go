package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

const BASE_URL = "https://the-one-api.dev/v2"

type OneAPIClient struct {
	client *http.Client
	apiKey string
}

func New(apiKey string) OneAPIClient {
	return OneAPIClient{
		client: http.DefaultClient,
		apiKey: apiKey,
	}
}

func (c OneAPIClient) authHeaderValue() string {
	return fmt.Sprintf("Bearer: %s", c.apiKey)
}

func (c OneAPIClient) buildEndpoint(path string) string {
	return fmt.Sprintf("%s/%s", BASE_URL, strings.Trim(path, "/"))
	return filepath.Join(BASE_URL, path)
}

// ListBooks returns a list of all "Lord of the Rings" books
func (c OneAPIClient) ListBooks() (response ListBooksResponse, err error) {
	endpoint := c.buildEndpoint("/book")
	fmt.Println(endpoint)

	resp, err := c.client.Get(endpoint)
	if err != nil {
		return response, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	books := ListBooksResponse{}
	err = json.Unmarshal(data, &books)
	return books, err
}

