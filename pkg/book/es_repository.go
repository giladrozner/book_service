package book

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"

	"github.com/giladrozner/book_service/pkg/config"
)

type ESRepository struct {
	client  *elasticsearch.Client
	timeout time.Duration
}

func NewESRepository(client *elasticsearch.Client, timeout time.Duration) *ESRepository {
	return &ESRepository{client: client, timeout: timeout}
}

var _ Repository = (*ESRepository)(nil)

func (r *ESRepository) Add(ctx context.Context, b Book) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	body, err := json.Marshal(b)
	if err != nil {
		return "", fmt.Errorf("marshal book: %w", err)
	}
	res, err := esapi.IndexRequest{
		Index: config.BooksIndex,
		Body:  bytes.NewReader(body),
	}.Do(ctx, r.client)
	if err != nil {
		return "", fmt.Errorf("add book: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return "", esError(res)
	}
	var resp struct {
		ID string `json:"_id"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", fmt.Errorf("parse add response: %w", err)
	}
	return resp.ID, nil
}

func (r *ESRepository) Get(ctx context.Context, id string) (Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := esapi.GetRequest{Index: config.BooksIndex, DocumentID: id}.Do(ctx, r.client)
	if err != nil {
		return Book{}, fmt.Errorf("get book %s: %w", id, err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return Book{}, esError(res)
	}
	var resp struct {
		Source Book `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return Book{}, fmt.Errorf("parse get response: %w", err)
	}
	return resp.Source, nil
}

func (r *ESRepository) UpdateTitle(ctx context.Context, id, newTitle string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	body, err := json.Marshal(map[string]any{"doc": map[string]any{"title": newTitle}})
	if err != nil {
		return fmt.Errorf("marshal update: %w", err)
	}
	res, err := esapi.UpdateRequest{
		Index:      config.BooksIndex,
		DocumentID: id,
		Body:       bytes.NewReader(body),
	}.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("update title %s: %w", id, err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return esError(res)
	}
	return nil
}

func (r *ESRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := esapi.DeleteRequest{Index: config.BooksIndex, DocumentID: id}.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("delete book %s: %w", id, err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return esError(res)
	}
	return nil
}

func (r *ESRepository) Search(ctx context.Context, c SearchCriteria) ([]BookWithID, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	must := []map[string]any{}
	if c.Title != nil {
		must = append(must, map[string]any{"match": map[string]any{"title": *c.Title}})
	}
	if c.AuthorName != nil {
		must = append(must, map[string]any{"match": map[string]any{"author_name": *c.AuthorName}})
	}
	if c.PriceMin != nil || c.PriceMax != nil {
		rng := map[string]any{}
		if c.PriceMin != nil {
			rng[config.ESRangeGTE] = *c.PriceMin
		}
		if c.PriceMax != nil {
			rng[config.ESRangeLTE] = *c.PriceMax
		}
		must = append(must, map[string]any{"range": map[string]any{"price": rng}})
	}
	if len(must) == 0 {
		must = append(must, map[string]any{"match_all": map[string]any{}})
	}
	body, err := json.Marshal(map[string]any{
		"query": map[string]any{"bool": map[string]any{"must": must}},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal search: %w", err)
	}
	res, err := esapi.SearchRequest{
		Index: []string{config.BooksIndex},
		Body:  bytes.NewReader(body),
	}.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("search books: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, esError(res)
	}
	var resp struct {
		Hits struct {
			Hits []struct {
				ID     string `json:"_id"`
				Source Book   `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("parse search response: %w", err)
	}
	out := make([]BookWithID, 0, len(resp.Hits.Hits))
	for _, h := range resp.Hits.Hits {
		out = append(out, BookWithID{ID: h.ID, Book: h.Source})
	}
	return out, nil
}

func (r *ESRepository) StoreStats(ctx context.Context) (int, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := esapi.SearchRequest{
		Index: []string{config.BooksIndex},
		Body:  bytes.NewReader([]byte(`{"size":0,"aggs":{"distinct_authors":{"cardinality":{"field":"author_name.keyword"}}}}`)),
	}.Do(ctx, r.client)
	if err != nil {
		return 0, 0, fmt.Errorf("store stats: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return 0, 0, esError(res)
	}
	var resp struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
		} `json:"hits"`
		Aggregations struct {
			DistinctAuthors struct {
				Value int `json:"value"`
			} `json:"distinct_authors"`
		} `json:"aggregations"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return 0, 0, fmt.Errorf("parse stats response: %w", err)
	}
	return resp.Hits.Total.Value, resp.Aggregations.DistinctAuthors.Value, nil
}

// esError wraps an Elasticsearch error response into a Go error.
// The format "Error %d: ..." allows ParseElasticsearchErrorCode to extract the status code.
func esError(res *esapi.Response) error {
	body, _ := io.ReadAll(res.Body)
	return fmt.Errorf("Error %d: %s", res.StatusCode, string(body))
}
