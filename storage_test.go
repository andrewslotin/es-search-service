package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElasticsearchStorage_Search(t *testing.T) {
	testCases := map[string]struct {
		Query              string
		Options            SearchOptions
		ExpectedParameters url.Values
	}{
		"default": {
			Query: "search term",
			ExpectedParameters: url.Values{
				"q": []string{"search term"},
			},
		},
		"with from": {
			Query: "search term",
			Options: SearchOptions{
				From: 11,
			},
			ExpectedParameters: url.Values{
				"q":    []string{"search term"},
				"from": []string{"11"},
			},
		},
		"with size": {
			Query: "search term",
			Options: SearchOptions{
				Size: 123,
			},
			ExpectedParameters: url.Values{
				"q":    []string{"search term"},
				"size": []string{"123"},
			},
		},
		"with sort": {
			Query: "search term",
			Options: SearchOptions{
				Sort: []string{"a:asc", "b:desc"},
			},
			ExpectedParameters: url.Values{
				"q":    []string{"search term"},
				"sort": []string{"a:asc,b:desc"},
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			node, mux, teardown := setupTS()
			defer teardown()

			var numRequests int
			mux.Handle("/_search", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				numRequests++

				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, testCase.ExpectedParameters, req.URL.Query())

				fd, err := os.Open("testdata/search_results.json")
				if err != nil {
					panic(err)
				}
				defer fd.Close()

				io.Copy(w, fd)
			}))

			c, err := elasticsearch.NewClient(elasticsearch.Config{
				Addresses: []string{node},
			})
			require.NoError(t, err)

			st := NewElasticsearchStorage(c)

			results, err := st.Search(context.Background(), testCase.Query, testCase.Options)
			require.NoError(t, err)

			require.Len(t, results, 2)
			assert.JSONEq(t, string(results[0]), `{"key": "value"}`)
			assert.JSONEq(t, string(results[1]), `{"answer": 42}`)

			assert.Equal(t, 1, numRequests)
		})
	}
}

func setupTS() (string, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)

	return ts.URL, mux, ts.Close
}
