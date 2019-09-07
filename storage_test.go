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
	node, mux, teardown := setupTS()
	defer teardown()

	var numRequests int
	mux.Handle("/_search", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		numRequests++

		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, url.Values{
			"q": []string{"search term"},
		}, req.URL.Query())

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

	results, err := st.Search(context.Background(), "search term", SearchOptions{})
	require.NoError(t, err)

	require.Len(t, results, 2)
	assert.JSONEq(t, string(results[0]), `{"key": "value"}`)
	assert.JSONEq(t, string(results[1]), `{"answer": 42}`)
}

func TestElasticsearchStorage_Search_WithFrom(t *testing.T) {
	node, mux, teardown := setupTS()
	defer teardown()

	var numRequests int
	mux.Handle("/_search", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		numRequests++

		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, url.Values{
			"q":    []string{"search term"},
			"from": []string{"11"},
		}, req.URL.Query())

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

	results, err := st.Search(context.Background(), "search term", SearchOptions{
		From: 11,
	})
	require.NoError(t, err)

	require.Len(t, results, 2)
	assert.JSONEq(t, string(results[0]), `{"key": "value"}`)
	assert.JSONEq(t, string(results[1]), `{"answer": 42}`)
}

func TestElasticsearchStorage_Search_WithPageSize(t *testing.T) {
	node, mux, teardown := setupTS()
	defer teardown()

	var numRequests int
	mux.Handle("/_search", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		numRequests++

		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, url.Values{
			"q":    []string{"search term"},
			"size": []string{"123"},
		}, req.URL.Query())

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

	results, err := st.Search(context.Background(), "search term", SearchOptions{
		Size: 123,
	})
	require.NoError(t, err)

	require.Len(t, results, 2)
	assert.JSONEq(t, string(results[0]), `{"key": "value"}`)
	assert.JSONEq(t, string(results[1]), `{"answer": 42}`)
}

func setupTS() (string, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)

	return ts.URL, mux, ts.Close
}
