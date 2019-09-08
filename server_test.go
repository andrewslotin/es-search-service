package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrewslotin/es-search-service/storage"

	"github.com/stretchr/testify/assert"
)

func TestSearchHandler(t *testing.T) {
	testCases := map[string]struct {
		Request       *http.Request
		SearchResult  []json.RawMessage
		ExpectedCode  int
		ExpectedBody  string
		ExpectedQuery string
		ExpectedOpts  storage.SearchOptions
	}{
		"with results": {
			Request: httptest.NewRequest(http.MethodGet, "/?q=search+term", nil),
			SearchResult: []json.RawMessage{
				json.RawMessage(`{"key": "value"}`),
				json.RawMessage(`{"answer": 42}`),
			},
			ExpectedCode:  http.StatusOK,
			ExpectedBody:  `{"status": "success", "results": [{"key": "value"}, {"answer": 42}]}`,
			ExpectedQuery: "search term",
		},
		"with empty results": {
			Request:       httptest.NewRequest(http.MethodGet, "/?q=search+term", nil),
			ExpectedCode:  http.StatusOK,
			ExpectedBody:  `{"status": "success", "results": []}`,
			ExpectedQuery: "search term",
		},
		"with pagination": {
			Request:       httptest.NewRequest(http.MethodGet, "/?q=search+term&from=11&size=123", nil),
			ExpectedCode:  http.StatusOK,
			ExpectedBody:  `{"status": "success", "results": []}`,
			ExpectedQuery: "search term",
			ExpectedOpts:  storage.SearchOptions{From: 11, Size: 123},
		},
		"with sort": {
			Request:       httptest.NewRequest(http.MethodGet, "/?q=search+term&sort=a:asc&sort=b:desc", nil),
			ExpectedCode:  http.StatusOK,
			ExpectedBody:  `{"status": "success", "results": []}`,
			ExpectedQuery: "search term",
			ExpectedOpts:  storage.SearchOptions{Sort: []string{"a:asc", "b:desc"}},
		},
		"with filter": {
			Request:       httptest.NewRequest(http.MethodGet, "/?q=search+term&filter=a:1+OR+b:2+and+c:3", nil),
			ExpectedCode:  http.StatusOK,
			ExpectedBody:  `{"status": "success", "results": []}`,
			ExpectedQuery: "search term",
			ExpectedOpts:  storage.SearchOptions{Filter: "a:1 OR b:2 and c:3"},
		},
		"missing query": {
			Request:      httptest.NewRequest(http.MethodGet, "/", nil),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: `{"status": "error", "code": 400, "error": "missing query parameter"}`,
		},
		"malformed from": {
			Request:      httptest.NewRequest(http.MethodGet, "/?q=search+term&from=abc&size=123", nil),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: `{"status": "error", "code": 400, "error": "malformed from parameter"}`,
		},
		"negative from": {
			Request:      httptest.NewRequest(http.MethodGet, "/?q=search+term&from=-1&size=123", nil),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: `{"status": "error", "code": 400, "error": "malformed from parameter"}`,
		},
		"malformed size": {
			Request:      httptest.NewRequest(http.MethodGet, "/?q=search+term&from=11&size=abc", nil),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: `{"status": "error", "code": 400, "error": "malformed size parameter"}`,
		},
		"negative size": {
			Request:      httptest.NewRequest(http.MethodGet, "/?q=search+term&from=11&size=-1", nil),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: `{"status": "error", "code": 400, "error": "malformed size parameter"}`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			m := &searcherMock{
				Results: testCase.SearchResult,
			}
			h := SearchHandler(m)
			rec := httptest.NewRecorder()

			h(rec, AuthenticatedRequest{
				Request:  testCase.Request,
				Username: "test1",
			})

			assert.Equal(t, testCase.ExpectedCode, rec.Code)
			assert.JSONEq(t, testCase.ExpectedBody, rec.Body.String())
			assert.Equal(t, testCase.ExpectedQuery, m.Query)
			assert.Equal(t, testCase.ExpectedOpts, m.Opts)
		})
	}
}

type searcherMock struct {
	Query   string
	Opts    storage.SearchOptions
	Results []json.RawMessage
}

func (m *searcherMock) Search(ctx context.Context, query string, opts storage.SearchOptions) ([]json.RawMessage, error) {
	m.Query = query
	m.Opts = opts

	return m.Results, nil
}
