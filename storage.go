package main

import (
	"context"
	"encoding/json"
	"fmt"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

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
func (st *ElasticsearchStorage) Search(ctx context.Context, query string) ([]json.RawMessage, error) {
	resp, err := st.es.Search(
		st.es.Search.WithContext(ctx),
		st.es.Search.WithQuery(query),
	)
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
