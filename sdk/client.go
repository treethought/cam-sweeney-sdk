package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const DEFAULT_BASE_URL = "https://the-one-api.dev/v2"

type OneAPIClient struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

type ClientConfig struct {
	Client  *http.Client
	BaseURL string
	ApiKey  string
}

func NewReadOnly() OneAPIClient {
	return OneAPIClient{
		client:  http.DefaultClient,
		baseURL: DEFAULT_BASE_URL,
	}
}

func New(apiKey string) OneAPIClient {
	return OneAPIClient{
		client:  http.DefaultClient,
		baseURL: DEFAULT_BASE_URL,
		apiKey:  apiKey,
	}
}

func NewWithConfig(config ClientConfig) OneAPIClient {
	var c OneAPIClient
	if config.ApiKey == "" {
		c = NewReadOnly()
	} else {
		c = New(config.ApiKey)
	}
	if config.Client != nil {
		c.client = config.Client
	}
	if config.BaseURL != "" {
		c.baseURL = config.BaseURL
	}
	return c
}

func (c OneAPIClient) buildEndpoint(path string) string {
	url := fmt.Sprintf("%s/%s", c.baseURL, strings.Trim(path, "/"))
	return strings.TrimSuffix(url, "/")

}

func (c OneAPIClient) doRequest(path string, opts ...RequestOption) (*http.Response, error) {
	endpoint := c.buildEndpoint(path)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	for _, f := range opts {
		f(req)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c OneAPIClient) doRequestInto(path string, v interface{}, opts ...RequestOption) error {
	resp, err := c.doRequest(path, opts...)
	if err != nil {
		return err
	}

	var bodyCopy bytes.Buffer
	r := io.TeeReader(resp.Body, &bodyCopy)

	// use TeeReader to allow reading into
	// response and error structs
	// TODO: maybe simplify this, and be more efficient
	// by allowing error info in response structs

	apiErr := APIError{}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	// check for error
	if err := json.Unmarshal(data, &apiErr); err == nil && apiErr.Message != "" {
		return SDKError{"API Error", path, apiErr}
	}

	// now unmarshal into provided struct
	err = json.Unmarshal(bodyCopy.Bytes(), v)
	if err != nil {
		return err
	}
	return nil
}

func (c OneAPIClient) Books() BooksClient {
	return BooksClient{c: c}
}

func (c OneAPIClient) Movies() MoviesClient {
	return MoviesClient{c: c}
}
func (c OneAPIClient) Characters() CharactersClient {
	return CharactersClient{c: c}
}
func (c OneAPIClient) Quotes() QuotesClient {
	return QuotesClient{c: c}
}
func (c OneAPIClient) Chapters() ChapterClient {
	return ChapterClient{c: c}
}
