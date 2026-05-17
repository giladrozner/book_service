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
	ErrMsgUsernameRequired = "username required"

	// Query parameters
	QueryParamTitle      = "title"
	QueryParamAuthorName = "author_name"
	QueryParamPriceRange = "price_range"
	QueryParamUsername   = "username"

	// URL parameters
	ParamID = "id"

	// Elasticsearch range operators
	ESRangeGTE = "gte"
	ESRangeLTE = "lte"

	// Book field names (used in validation errors and ES queries)
	FieldTitle       = "title"
	FieldAuthorName  = "author_name"
	FieldPrice       = "price"
	FieldPublishDate = "publish_date"

	// Field-level validation error messages
	ErrTitleRequired       = "title is a required field."
	ErrAuthorNameRequired  = "author name is a required field."
	ErrPriceRequired       = "price is a required field and must be a number."
	ErrPublishDateRequired = "publish date is a required field."
	ErrFieldRequired       = "this field is required."

	// Response field keys
	ErrorField           = "error"
	MessageField         = "message"
	ActionsField         = "actions"
	BookCountField       = "book_count"
	DistinctAuthorsField = "distinct_authors"

	// Success messages
	SuccessBookUpdated = "book updated successfully"
	SuccessBookDeleted = "book deleted successfully"
)

// ElasticsearchErrorMap maps ES HTTP status codes to human-readable messages.
var ElasticsearchErrorMap = map[int]string{
	400: "Bad Request: The request was invalid or improperly formatted.",
	401: "Unauthorized: Authentication required or failed.",
	403: "Forbidden: You do not have permission to perform this action.",
	404: "Not Found: The requested resource could not be found.",
	408: "Request Timeout: The request took too long to process.",
	409: "Conflict: A version conflict occurred, possibly due to concurrent updates.",
	429: "Too Many Requests: Elasticsearch is throttling due to high load.",
	500: "Internal Server Error: An error occurred within Elasticsearch.",
	503: "Service Unavailable: The cluster is unavailable.",
}
