package activity_middleware

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/giladrozner/book_service/pkg/activity"
	"github.com/giladrozner/book_service/pkg/config"
)

func Middleware(repo activity.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // run the actual handler first

		username := c.Query(config.QueryParamUsername)
		if username == "" {
			return
		}
		route := c.FullPath() // gives "/books/:id" not "/books/abc"
		if route == "" {
			return
		}
		method := c.Request.Method

		if err := repo.Record(c.Request.Context(), username, activity.Action{
			Method: method,
			Route:  route,
		}); err != nil {
			log.Printf("activity record error: %v", err)
			// don't fail the request — recording is best-effort
		}
	}
}
