package elastic

import (
	"github.com/elastic/go-elasticsearch/v7"
)

func NewClient(url string) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})
}
