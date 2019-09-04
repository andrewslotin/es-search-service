package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchHandler(t *testing.T) {
	testCases := map[string]struct {
		Request       *http.Request
		SearchResult  []json.RawMessage
		ExpectedCode  int
		ExpectedBody  string
		ExpectedQuery string
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
		"missing query": {
			Request:      httptest.NewRequest(http.MethodGet, "/", nil),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: `{"status": "error", "code": 400, "error": "missing query parameter"}`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			m := &searcherMock{
				Results: testCase.SearchResult,
			}
			h := SearchHandler(m)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, testCase.Request)

			assert.Equal(t, testCase.ExpectedCode, rec.Code)
			assert.JSONEq(t, testCase.ExpectedBody, rec.Body.String())
			assert.Equal(t, testCase.ExpectedQuery, m.Query)
		})
	}
}

type searcherMock struct {
	Query   string
	Results []json.RawMessage
}

func (m *searcherMock) Search(ctx context.Context, query string) ([]json.RawMessage, error) {
	m.Query = query

	return m.Results, nil
}
