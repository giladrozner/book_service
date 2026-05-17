package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func NewClient(url string, readTimeout, writeTimeout time.Duration) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         url,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})
}
