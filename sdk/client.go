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
