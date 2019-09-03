package main

import (
	"fmt"
	"log"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

func main() {
	if os.Getenv("ELASTICSEARCH_URL") == "" {
		log.Fatalf("ELASTICSEARCH_URL= is not set")
	}

	c, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to dial elasticsearch cluster: %s", err)
	}

	resp, err := c.Info()
	if err != nil {
		log.Fatalf("failed to query elasticsearch: %s", err)
	}

	fmt.Println(resp)
}
