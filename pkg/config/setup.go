package config

import (
	"github.com/giladrozner/book_service/pkg/activity"
	"github.com/giladrozner/book_service/pkg/book"
	"github.com/giladrozner/book_service/pkg/connectors/elastic"
	"github.com/giladrozner/book_service/pkg/connectors/redis"
)

var (
	BookRepo     book.Repository
	ActivityRepo activity.Repository
)

func Setup() error {
	cfg := Load()

	esClient, err := elastic.NewClient(cfg.ESURL)
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(cfg.RedisURL, cfg.RedisReadTimeout(), cfg.RedisWriteTimeout())

	BookRepo = book.NewESRepository(esClient, cfg.ESTimeout())
	ActivityRepo = activity.NewRedisRepository(redisClient)

	return nil
}
