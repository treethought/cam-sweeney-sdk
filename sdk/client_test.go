package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWithConfig(t *testing.T) {
	type args struct {
		config ClientConfig
	}
	httpClient := &http.Client{Timeout: time.Hour * 3}
	confWithClient := ClientConfig{Client: httpClient}
	confWithAPIKey := ClientConfig{ApiKey: "123"}
	configWithBaseURL := ClientConfig{BaseURL: "https://example.com"}

	tests := []struct {
		name   string
		config ClientConfig
		want   OneAPIClient
	}{
		{"with client", confWithClient, OneAPIClient{baseURL: DEFAULT_BASE_URL, apiKey: "", client: httpClient}},
		{"with apiKey", confWithAPIKey, OneAPIClient{baseURL: DEFAULT_BASE_URL, apiKey: "123", client: http.DefaultClient}},
		{"with apiKey", configWithBaseURL, OneAPIClient{baseURL: "https://example.com", apiKey: "", client: http.DefaultClient}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWithConfig(tt.config)
			assert.Equal(t, got.client.Timeout, tt.want.client.Timeout)
			assert.Equal(t, got.baseURL, tt.want.baseURL)
			assert.Equal(t, got.apiKey, tt.want.apiKey)
		})
	}
}

func TestOneAPIClient_buildEndpoint(t *testing.T) {
	client := NewReadOnly()
	type args struct {
		path string
	}
	tests := []struct {
		name string
		path string
		want string
	}{
		// TODO: Add test cases.
		{"empty", "", DEFAULT_BASE_URL},
		{"bare slash", "/", DEFAULT_BASE_URL},
		{"leading slash", "/book", fmt.Sprintf("%s/book", DEFAULT_BASE_URL)},
		{"missing leading slash", "book", fmt.Sprintf("%s/book", DEFAULT_BASE_URL)},
		{"trailing slash", "/book/", fmt.Sprintf("%s/book", DEFAULT_BASE_URL)},
		{"nested", "/book/id/chapter", fmt.Sprintf("%s/book/id/chapter", DEFAULT_BASE_URL)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.buildEndpoint(tt.path); got != tt.want {
				t.Errorf("OneAPIClient.buildEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOneAPIClient_doRequest(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			resp := APIError{Success: false, Message: "sample error message"}
			data, err := json.Marshal(resp)
			assert.Nil(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
			return
		}
		data, err := io.ReadAll(r.Body)
		assert.Nil(err)
		for k, v := range r.Header {
			fmt.Println("setting ", k, v)
			w.Header().Set(k, v[0])
		}
		w.Write(data)

	}))
	defer server.Close()

	client := NewWithConfig(ClientConfig{BaseURL: server.URL})

	type args struct {
		path string
		opts []RequestOption
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantPath   string
		wantHeader map[string]string
	}{
		{"with error", args{path: "/error"}, true, "/error", nil},
		{"no options", args{path: "/book"}, false, "/book", nil},
		{"with limit", args{path: "/book", opts: []RequestOption{WithLimit(3)}}, false, "/book?limit=1", nil},
		{"with offset", args{path: "/book", opts: []RequestOption{WithOffset(2)}}, false, "/book?offset=2", nil},
		{"with page", args{path: "/book", opts: []RequestOption{WithPage(4)}}, false, "/book?page=4", nil},
		{"with sort asc", args{path: "/book", opts: []RequestOption{WithSort("name", "asc")}}, false, "/book?sort=name:asc", nil},
		{"with sort dsc", args{path: "/book", opts: []RequestOption{WithSort("character", "dsc")}}, false, "/book?sort=character:dsc", nil},
		{"with apikey", args{path: "/book", opts: []RequestOption{WithAPIKey("123")}}, false, "/book", map[string]string{"Authorization": "Bearer 123"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.doRequest(tt.args.path, tt.args.opts...)
			if tt.wantErr {
				expectedErr := APIError{Success: false, Message: "sample error message"}

				data, err := io.ReadAll(got.Body)
				assert.Nil(err)

				gotErr := APIError{}
				err = json.Unmarshal(data, &gotErr)
				assert.Nil(err)

				assert.Equal(expectedErr, gotErr)

				return
			}
			assert.Nil(err)
			assert.Equal(tt.args.path, got.Request.URL.Path)

			fmt.Println(got.Header)

			for k, v := range tt.wantHeader {
				assert.Equal(v, got.Header.Get(k))
			}

		})
	}
}

