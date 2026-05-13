package activity

import "context"

type Action struct {
	Method string `json:"method"`
	Route  string `json:"route"`
}

type Repository interface {
	Record(ctx context.Context, username string, action Action) error
	Recent(ctx context.Context, username string, n int) ([]Action, error)
}
