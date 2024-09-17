package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
)

var ES *elasticsearch.Client

func InitES() {
	var err error

	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		log.Fatal("ELASTICSEARCH_URL is not set")
	}

	log.Printf("Attempting to connect to Elasticsearch at %s", esURL)

	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
	}
	ES, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := ES.Info()
	if err != nil {
		log.Fatalf("Error getting response from Elasticsearch: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error response from Elasticsearch: %s", res.String())
	}

	log.Println("Connected to Elasticsearch")
}
