package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	if err := engine.Run(":8080"); err != nil {
		panic(err)
	}
}
