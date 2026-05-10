package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ESURL               string
	RedisURL            string
	ESTimeoutMs         int
	RedisReadTimeoutMs  int
	RedisWriteTimeoutMs int
}

func Load() Config {
	return Config{
		ESURL:               getenv("ES_URL", "http://localhost:9200"),
		RedisURL:            getenv("REDIS_URL", "localhost:6379"),
		ESTimeoutMs:         getenvInt("ES_TIMEOUT_MS", 5000),
		RedisReadTimeoutMs:  getenvInt("REDIS_READ_TIMEOUT_MS", 2000),
		RedisWriteTimeoutMs: getenvInt("REDIS_WRITE_TIMEOUT_MS", 2000),
	}
}

func (c Config) ESTimeout() time.Duration {
	return time.Duration(c.ESTimeoutMs) * time.Millisecond
}

func (c Config) RedisReadTimeout() time.Duration {
	return time.Duration(c.RedisReadTimeoutMs) * time.Millisecond
}

func (c Config) RedisWriteTimeout() time.Duration {
	return time.Duration(c.RedisWriteTimeoutMs) * time.Millisecond
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
