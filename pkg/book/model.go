package book

type Book struct {
	Title          string  `json:"title"`
	AuthorName     string  `json:"author_name"`
	Price          float64 `json:"price"`
	EbookAvailable *bool   `json:"ebook_available,omitempty"`
	PublishDate    string  `json:"publish_date"`
}
