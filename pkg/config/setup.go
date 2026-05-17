package config

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/redis/go-redis/v9"

	elasticconnector "github.com/giladrozner/book_service/pkg/connectors/elastic"
	redisconnector "github.com/giladrozner/book_service/pkg/connectors/redis"
)

var (
	ESClient    *elasticsearch.Client
	RedisClient *redis.Client
)

func Setup() error {
	cfg := Load()

	esClient, err := elasticconnector.NewClient(cfg.ESURL)
	if err != nil {
		return err
	}

	ESClient = esClient
	RedisClient = redisconnector.NewClient(cfg.RedisURL, cfg.RedisReadTimeout(), cfg.RedisWriteTimeout())

	return nil
}
