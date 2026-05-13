package config

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/redis/go-redis/v9"

	"github.com/giladrozner/book_service/pkg/activity"
	"github.com/giladrozner/book_service/pkg/book"
)

var (
	BookRepo     book.Repository
	ActivityRepo activity.Repository
)

func Setup() error {
	cfg := Load()

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ESURL},
	})
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisURL,
		ReadTimeout:  cfg.RedisReadTimeout(),
		WriteTimeout: cfg.RedisWriteTimeout(),
	})

	BookRepo = book.NewESRepository(esClient, cfg.ESTimeout())
	ActivityRepo = activity.NewRedisRepository(redisClient)

	return nil
}
