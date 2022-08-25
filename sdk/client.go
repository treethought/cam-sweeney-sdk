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

// OneAPIClient is the sdk's interface to The One API
type OneAPIClient struct {
	client         *http.Client
	apiKey         string
	baseURL        string
	persistentOpts []RequestOption
}

// ClientConfig provides config to override client behavior
type ClientConfig struct {
	// http client used to make requests
	Client *http.Client
	// base url of API
	BaseURL string
	// APIKey required for authenticated endpoints
	ApiKey string

	// Request Options to apply to all requests
	// note any options provided to methods will overwrite any duplicates
	PersistentOptions []RequestOption
}

// NewUnAuthenticated creates a new client without authorization
func NewUnAuthenticated() OneAPIClient {
	return OneAPIClient{
		client:         http.DefaultClient,
		baseURL:        DEFAULT_BASE_URL,
		persistentOpts: []RequestOption{},
	}
}

// NewReadOnly creates a client using provided apiKey for authorization
func New(apiKey string) OneAPIClient {
	return OneAPIClient{
		client:         http.DefaultClient,
		baseURL:        DEFAULT_BASE_URL,
		apiKey:         apiKey,
		persistentOpts: []RequestOption{},
	}
}

// NewWithConfig creates a client configured with the provided config
func NewWithConfig(config ClientConfig) OneAPIClient {
	var c OneAPIClient
	if config.ApiKey == "" {
		c = NewUnAuthenticated()
	} else {
		c = New(config.ApiKey)
	}
	if config.Client != nil {
		c.client = config.Client
	}
	if config.BaseURL != "" {
		c.baseURL = config.BaseURL
	}
	if len(config.PersistentOptions) > 0 {
		c.persistentOpts = config.PersistentOptions
	}
	return c
}

func (c OneAPIClient) buildEndpoint(path string) string {
	url := fmt.Sprintf("%s/%s", c.baseURL, strings.Trim(path, "/"))
	return strings.TrimSuffix(url, "/")

}

func (c OneAPIClient) appendOptsToAuth(opts ...RequestOption) []RequestOption {
	auth := []RequestOption{WithAPIKey(c.apiKey)}
	return append(auth, opts...)

}

func (c OneAPIClient) doRequest(path string, opts ...RequestOption) (*http.Response, error) {
	endpoint := c.buildEndpoint(path)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	// apply persistent opts first
	for _, f := range c.persistentOpts {
		f(req)
	}

	// apply per-method opts
	for _, f := range opts {
		f(req)
	}

	return c.client.Do(req)
}

func (c OneAPIClient) doRequestInto(path string, v interface{}, opts ...RequestOption) error {
	resp, err := c.doRequest(path, opts...)
	if err != nil {
		return SDKError{"HTTP Error", path, err}
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
		return SDKError{"Error reading response", path, err}
	}

	// check for error
	if err := json.Unmarshal(data, &apiErr); err == nil && apiErr.Message != "" {
		return SDKError{"API Error", path, apiErr}
	}

	// now unmarshal into provided struct
	err = json.Unmarshal(bodyCopy.Bytes(), v)
	if err != nil {
		return SDKError{"Deserialization Error", path, err}
	}
	return nil
}

// Books provides access to the /book namespace of resources
func (c OneAPIClient) Books() BooksClient {
	return BooksClient{c: c}
}

// Books provides access to the /movie namespace of resources
func (c OneAPIClient) Movies() MoviesClient {
	return MoviesClient{c: c}
}

// Characters provides access to the /character namespace of resources
func (c OneAPIClient) Characters() CharactersClient {
	return CharactersClient{c: c}
}

// Characters provides access to the /quote namespace of resources
func (c OneAPIClient) Quotes() QuotesClient {
	return QuotesClient{c: c}
}

// Characters provides access to the /chapter namespace of resources
func (c OneAPIClient) Chapters() ChapterClient {
	return ChapterClient{c: c}
}
