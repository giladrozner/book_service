package config

import "time"

const (
	// Server
	ServerPort = ":8080"

	// Elasticsearch
	BooksIndex = "books"

	// Activity
	ActivityKeyPrefix = "activity:%s"
	ActivityMaxItems  = 3
	ActivityTTL       = 24 * time.Hour

	// Config defaults
	DefaultESURL               = "http://localhost:9200"
	DefaultRedisURL            = "localhost:6379"
	DefaultESTimeoutMs         = 5000
	DefaultRedisReadTimeoutMs  = 2000
	DefaultRedisWriteTimeoutMs = 2000

	// Error messages
	ErrMsgInternal         = "internal error"
	ErrMsgBookNotFound     = "book not found"
	ErrMsgUsernameRequired = "username required"

	// Query parameters
	QueryParamTitle      = "title"
	QueryParamAuthorName = "author_name"
	QueryParamPriceRange = "price_range"
	QueryParamUsername   = "username"
)
