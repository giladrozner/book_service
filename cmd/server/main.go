package main

import (
	"context"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v7"

	"github.com/giladrozner/book_service/internal/book"
	"github.com/giladrozner/book_service/internal/config"
)

func main() {
	cfg := config.Load()

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ESURL},
	})
	if err != nil {
		panic(err)
	}

	repo := book.NewESRepository(esClient, cfg.ESTimeout())
	ctx := context.Background()

	id, err := repo.Add(ctx, book.Book{
		Title:       "Smoke Test",
		AuthorName:  "Test Author",
		Price:       1.23,
		PublishDate: time.Now(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("added id:", id)

	got, err := repo.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("got: %+v\n", got)

	if err := repo.Delete(ctx, id); err != nil {
		panic(err)
	}
	fmt.Println("deleted ok")
}
