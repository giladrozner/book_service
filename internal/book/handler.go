package book

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// gin.IRouter accepts both *gin.Engine and *gin.RouterGroup, so this works
// whether we attach to the root engine or to a sub-group with middleware.
func (h *Handler) Register(r gin.IRouter) {
	r.POST("/books", h.create)
	r.GET("/books/:id", h.get)
	r.PUT("/books/:id", h.updateTitle)
	r.DELETE("/books/:id", h.delete)
	r.GET("/search", h.search)
	r.GET("/store", h.storeStats)
}

type createBookRequest struct {
	Title          string    `json:"title" binding:"required"`
	AuthorName     string    `json:"author_name" binding:"required"`
	Price          float64   `json:"price"`
	EbookAvailable *bool     `json:"ebook_available"`
	PublishDate    time.Time `json:"publish_date" binding:"required"`
}

func (h *Handler) create(c *gin.Context) {
	var req createBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.repo.Add(c.Request.Context(), Book{
		Title:          req.Title,
		AuthorName:     req.AuthorName,
		Price:          req.Price,
		EbookAvailable: req.EbookAvailable,
		PublishDate:    req.PublishDate,
	})
	if err != nil {
		log.Printf("create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) get(c *gin.Context) {
	id := c.Param("id")
	b, err := h.repo.Get(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	if err != nil {
		log.Printf("get error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, b)
}

type updateTitleRequest struct {
	Title string `json:"title" binding:"required"`
}

func (h *Handler) updateTitle(c *gin.Context) {
	id := c.Param("id")
	var req updateTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.repo.UpdateTitle(c.Request.Context(), id, req.Title)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	if err != nil {
		log.Printf("update error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	err := h.repo.Delete(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	if err != nil {
		log.Printf("delete error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) search(c *gin.Context) {
	var crit SearchCriteria
	if t := c.Query("title"); t != "" {
		crit.Title = &t
	}
	if a := c.Query("author_name"); a != "" {
		crit.AuthorName = &a
	}
	if pr := c.Query("price_range"); pr != "" {
		parts := strings.SplitN(pr, ",", 2)
		if len(parts) == 2 {
			if parts[0] != "" {
				if v, err := strconv.ParseFloat(parts[0], 64); err == nil {
					crit.PriceMin = &v
				}
			}
			if parts[1] != "" {
				if v, err := strconv.ParseFloat(parts[1], 64); err == nil {
					crit.PriceMax = &v
				}
			}
		}
	}
	results, err := h.repo.Search(c.Request.Context(), crit)
	if err != nil {
		log.Printf("search error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (h *Handler) storeStats(c *gin.Context) {
	books, authors, err := h.repo.StoreStats(c.Request.Context())
	if err != nil {
		log.Printf("store error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books, "authors": authors})
}
