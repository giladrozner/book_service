package service

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/giladrozner/book_service/pkg/service/middleware/activity_middleware"
)

func Routes(router *gin.Engine) {
	Setup()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	// Tracked routes — wrapped in the activity-recording middleware
	tracked := router.Group("/")
	tracked.Use(activity_middleware.Middleware(ActivityRepo))
	{
		tracked.POST("/books", CreateBook)
		tracked.GET("/books/:id", GetBook)
		tracked.PUT("/books/:id", UpdateBook)
		tracked.DELETE("/books/:id", DeleteBook)
		tracked.GET("/search", Search)
		tracked.GET("/store", StoreStats)
	}

	// /activity is registered on the bare router — NOT tracked
	router.GET("/activity", Activity)
}
