package book

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("book not found")

type SearchCriteria struct {
	Title      *string
	AuthorName *string
	PriceMin   *float64
	PriceMax   *float64
}

type BookWithID struct {
	ID   string `json:"id"`
	Book Book   `json:"book"`
}

type Repository interface {
	Add(ctx context.Context, b Book) (string, error)
	Get(ctx context.Context, id string) (Book, error)
	UpdateTitle(ctx context.Context, id, newTitle string) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, criteria SearchCriteria) ([]BookWithID, error)
	StoreStats(ctx context.Context) (totalBooks int, distinctAuthors int, err error)
}
