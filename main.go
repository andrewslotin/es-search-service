package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

var args struct {
	NodesList   string
	ConnTimeout time.Duration
}

func main() {
	if timeout := os.Getenv("ELASTICSEARCH_CONN_TIMEOUT"); timeout != "" {
		dur, err := time.ParseDuration(timeout)
		if err != nil {
			log.Fatalf("invalid elasticsearch cluster connection timeout value: %s", timeout)
		}

		args.ConnTimeout = dur
	}

	flag.StringVar(&args.NodesList, "nodes", os.Getenv("ELASTICSEARCH_NODES"), "Comma-separated list of Elasticsearch cluster nodes")
	flag.DurationVar(&args.ConnTimeout, "timeout", args.ConnTimeout, "Elastisearch cluster connection timeout")
	flag.Parse()

	nodes := strings.Split(args.NodesList, ",")
	if len(nodes) == 0 {
		log.Fatal("there were no elasticsearch nodes provided, did you forget to populate ELASTICSEARCH_NODES=?")
	}

	ctx, _ := context.WithTimeout(context.Background(), args.ConnTimeout)
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
