package book

import "time"

type Book struct {
	Title          string    `json:"title"`
	AuthorName     string    `json:"author_name"`
	Price          float64   `json:"price"`
	EbookAvailable *bool     `json:"ebook_available,omitempty"`
	PublishDate    time.Time `json:"publish_date"`
}
