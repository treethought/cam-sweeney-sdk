package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
	client := NewUnAuthenticated()
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
		{"with match", args{path: "/character", opts: []RequestOption{WithFilterMatch("name", "frodo")}}, false, "/character?name=frodo", nil},
		{"with match negate", args{path: "/character", opts: []RequestOption{WithFilterNegate("name", "frodo")}}, false, "/character?name!=frodo", nil},
		{"with match negate and multiple", args{path: "/character", opts: []RequestOption{WithLimit(2), WithFilterNegate("name", "frodo")}}, false, "/character?limit=2&name!=frodo", nil},
		{"with include", args{path: "/character", opts: []RequestOption{WithFilterInclude("race", "hobit,human")}}, false, "/character?name=hobit,human", nil},
		{"with exclude", args{path: "/character", opts: []RequestOption{WithFilterExclude("race", "hobit,human")}}, false, "/character?name!=hobit,human", nil},
		{"with regex include", args{path: "/character", opts: []RequestOption{WithRegexInclude("name", "/foot/i")}}, false, "/character?name=/foot/i", nil},
		{"with regex exclude", args{path: "/character", opts: []RequestOption{WithRegexExclude("name", "/foot/i")}}, false, "/character?name!=/foot/i,human", nil},
		{"with less than", args{path: "/movie", opts: []RequestOption{WithComparison("RuntimeInMinutes", "<", 180)}}, false, "/movie?RuntimeInMinutes<180", nil},
		{"with greater than", args{path: "/movie", opts: []RequestOption{WithComparison("RuntimeInMinutes", ">", 180)}}, false, "/movie?RuntimeInMinutes>180", nil},
		{"with greater than or equal", args{path: "/movie", opts: []RequestOption{WithComparison("RuntimeInMinutes", ">=", 180)}}, false, "/movie?RuntimeInMinutes>=180", nil},
		{"with comparison and multiple ", args{path: "/movie", opts: []RequestOption{WithLimit(2), WithComparison("RuntimeInMinutes", ">=", 180)}}, false, "/movie?limit=2&RuntimeInMinutes>=180", nil},
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
			assert.Contains(got.Request.URL.String(), tt.args.path)
			// assert.Equal(tt.args.path, got.Request.URL.Path)
			fmt.Println(got.Request.URL.String())

			for k, v := range tt.wantHeader {
				assert.Equal(v, got.Header.Get(k))
			}

		})
	}
}

func TestOneAPIClient_doRequestInto(t *testing.T) {
	assert := assert.New(t)
	var resp interface{}

	apiErr := APIError{Success: false, Message: "sample error message"}

	wantErrResp := SDKError{message: "API Error", endpoint: "/error", err: apiErr}
	wantBookResp := booksResponse{Docs: []Book{{ID: "123", Name: "Sample Book"}}}
	wantMovieResp := moviesResponse{Docs: []Movie{{ID: "123", Name: "Sample Movie", RuntimeInMinutes: 122}}}
	wantChapterResp := chapterResponse{Docs: []Chapter{{ID: "123", Name: "Sample Chapter", Book: "smaple book"}}}
	wantCharacterResp := characterResponse{Docs: []Character{{ID: "123", Name: "Frodo", Height: "4 ft", Race: "Hobbit"}}}
	wantQuoteResp := quoteResponse{Docs: []Quote{{ID: "123", Character: "Frodo", Dialog: "smaple dialog"}}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/error":
			resp = apiErr
		case "/book":
			resp = wantBookResp
		case "/movie":
			resp = wantMovieResp
		case "/chapter":
			resp = wantChapterResp
		case "/character":
			resp = wantCharacterResp
		case "/quote":
			resp = wantQuoteResp
		}
		data, err := json.Marshal(resp)
		assert.Nil(err)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		w.Header().Set("Content-Type", "application/json")

	}))
	defer server.Close()

	client := NewWithConfig(ClientConfig{BaseURL: server.URL})

	type args struct {
		path string
		v    interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantResp interface{}
	}{
		{"error message", args{path: "/error", v: make(map[string]interface{})}, true, wantErrResp},
		{"book resp", args{path: "/book", v: booksResponse{}}, false, wantBookResp},
		{"movie resp", args{path: "/movie", v: moviesResponse{}}, false, wantMovieResp},
		{"chapter resp", args{path: "/chapter", v: chapterResponse{}}, false, wantChapterResp},
		{"character resp", args{path: "/character", v: characterResponse{}}, false, wantCharacterResp},
		{"quote resp", args{path: "/quote", v: quoteResponse{}}, false, wantQuoteResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var err error
			switch tt.args.path {
			case "/error":
				err = client.doRequestInto(tt.args.path, make(map[string]string))
				assert.NotNil(err)
				assert.ErrorIs(err, wantErrResp)
			case "/book":
				v := booksResponse{}
				err = client.doRequestInto(tt.args.path, &v)
				assert.Nil(err)
				assert.Equal(wantBookResp, v)
			case "/movie":
				v := moviesResponse{}
				err = client.doRequestInto(tt.args.path, &v)
				assert.Nil(err)
				assert.Equal(wantMovieResp, v)
			case "/chapter":
				v := chapterResponse{}
				err = client.doRequestInto(tt.args.path, &v)
				assert.Nil(err)
				assert.Equal(wantChapterResp, v)
			case "/character":
				v := characterResponse{}
				err = client.doRequestInto(tt.args.path, &v)
				assert.Nil(err)
				assert.Equal(wantCharacterResp, v)
			case "/quote":
				v := quoteResponse{}
				err = client.doRequestInto(tt.args.path, &v)
				assert.Nil(err)
				assert.Equal(wantQuoteResp, v)
			}

		})
	}
}
