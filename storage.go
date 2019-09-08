package main

import (
	"context"
	"encoding/json"
	"fmt"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	esapi "github.com/elastic/go-elasticsearch/v7/esapi"
)

// SearchOptions define the options to be passed to Elasticsearch API seach request
type SearchOptions struct {
	// From is the number of documents to skip before returning the result
	From int
	// Size is the number of documents to return in result
	Size int
	// Sort is a list of fields to sort by followied by sort direction, i.e. ["field1:asc", "field2:desc"]
	Sort []string
}

// ElasticsearchStorage implements access to the Elasticsearch cluster
type ElasticsearchStorage struct {
	es *elasticsearch.Client
}

// NewElasticsearchStorage initializes a new instance of ElasticsearchStorage
func NewElasticsearchStorage(c *elasticsearch.Client) *ElasticsearchStorage {
	return &ElasticsearchStorage{es: c}
}

// Search queries the Elasticsearch cluster and returns a list of JSON documents
// matching the search query.
func (st *ElasticsearchStorage) Search(ctx context.Context, query string, opts SearchOptions) ([]json.RawMessage, error) {
	req := []func(*esapi.SearchRequest){
		st.es.Search.WithContext(ctx),
		st.es.Search.WithQuery(query),
	}

	if opts.From > 0 {
		req = append(req, st.es.Search.WithFrom(opts.From))
	}

	if opts.Size > 0 {
		req = append(req, st.es.Search.WithSize(opts.Size))
	}

	if len(opts.Sort) > 0 {
		req = append(req, st.es.Search.WithSort(opts.Sort...))
	}

	resp, err := st.es.Search(req...)
	if err != nil {
		return nil, fmt.Errorf("failed to query elasticsearch: %s", err)
	}
	defer resp.Body.Close()

	var searchResults struct {
		Hits struct {
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&searchResults); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %s", err)
	}

	var results []json.RawMessage
	for _, res := range searchResults.Hits.Hits {
		results = append(results, res.Source)
	}

	return results, nil
}
