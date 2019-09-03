package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

func main() {
	nodes := strings.Split(os.Getenv("ELASTICSEARCH_NODES"), ",")
	if len(nodes) == 0 {
		log.Fatal("there were no elasticsearch nodes provided, did you forget to populate ELASTICSEARCH_NODES=?")
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	c, err := DialElasticsearch(ctx, nodes)
	if err != nil {
		log.Fatalf("failed to connect to elasticsearch cluster: %s", err)
	}

	resp, err := c.Info()
	if err != nil {
		log.Fatalf("failed to query elasticsearch: %s", err)
	}

	fmt.Println(resp)
}

// DialElasticsearch establishes connection with Elasticsearch cluster and ensures that it's
// up and running. If there is a non-nil context provided, this function will keep retrying to
// connect to cluster in case of an error until the supplied context is done.
func DialElasticsearch(ctx context.Context, nodes []string) (*elasticsearch.Client, error) {
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: nodes,
	})
	if err != nil {
		return nil, err
	}

	_, err = c.Info()
	if err != nil && ctx == nil {
		// do not retry if there was no context provided for cancellation/timeout
		return nil, err
	}

	// return immediately if connection succeeded
	if err == nil {
		return c, nil
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// attempt to reach the cluster each 100ms until success or the context is cancelled
	for {
		select {
		case <-ticker.C:
			_, err := c.Info()
			if err == nil {
				return c, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
