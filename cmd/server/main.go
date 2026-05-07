package main

import (
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"

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

	bookRepo := book.NewESRepository(esClient, cfg.ESTimeout())
	bookHandler := book.NewHandler(bookRepo)

	engine := gin.Default()
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})
	bookHandler.Register(engine)

	if err := engine.Run(":8080"); err != nil {
		panic(err)
	}
}
