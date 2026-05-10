package activity

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Register(r gin.IRouter) {
	r.GET("/activity", h.recent)
}

func (h *Handler) recent(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}
	actions, err := h.repo.Recent(c.Request.Context(), username, 3)
	if err != nil {
		log.Printf("recent error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"actions": actions})
}
