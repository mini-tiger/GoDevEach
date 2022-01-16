package main

import (
	"log"

	"github.com/elastic/go-elasticsearch/v7"
)

func main() {
	es, err := elasticsearch.NewClient(
		elasticsearch.Config{Addresses: []string{"http://172.16.71.31:9200"}},
	)
	if err != nil {
		panic(err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
}
