package main

import (
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/giladrozner/book_service/internal/activity"
	"github.com/giladrozner/book_service/internal/book"
	"github.com/giladrozner/book_service/internal/config"
)

func main() {
	cfg := config.Load()

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ESURL},
	})
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisURL,
		ReadTimeout:  cfg.RedisReadTimeout(),
		WriteTimeout: cfg.RedisWriteTimeout(),
	})

	bookRepo := book.NewESRepository(esClient, cfg.ESTimeout())
	bookHandler := book.NewHandler(bookRepo)

	activityRepo := activity.NewRedisRepository(redisClient)
	activityHandler := activity.NewHandler(activityRepo)

	engine := gin.Default()
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	// Tracked routes — wrapped in the activity-recording middleware
	tracked := engine.Group("/")
	tracked.Use(activity.Middleware(activityRepo))
	bookHandler.Register(tracked)

	// /activity is registered on the bare engine — NOT tracked
	activityHandler.Register(engine)

	if err := engine.Run(":8080"); err != nil {
		panic(err)
	}
}
