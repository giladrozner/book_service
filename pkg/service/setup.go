package service

import (
	"github.com/giladrozner/book_service/pkg/activity"
	"github.com/giladrozner/book_service/pkg/book"
	"github.com/giladrozner/book_service/pkg/config"
)

var (
	BookRepo     book.Repository
	ActivityRepo activity.Repository
)

func Setup() {
	BookRepo = book.NewESRepository(config.ESClient, config.Load().ESTimeout())
	ActivityRepo = activity.NewRedisRepository(config.RedisClient)
}
