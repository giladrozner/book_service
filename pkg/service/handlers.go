package service

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/giladrozner/book_service/pkg/book"
	"github.com/giladrozner/book_service/pkg/config"
	"github.com/giladrozner/book_service/pkg/consts"
)

type createBookRequest struct {
	Title          string  `json:"title" binding:"required"`
	AuthorName     string  `json:"author_name" binding:"required"`
	Price          float64 `json:"price"`
	EbookAvailable *bool   `json:"ebook_available"`
	PublishDate    string  `json:"publish_date" binding:"required"`
}

type updateTitleRequest struct {
	Title string `json:"title" binding:"required"`
}

func CreateBook(c *gin.Context) {
	var req createBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := config.BookRepo.Add(c.Request.Context(), book.Book{
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

func GetBook(c *gin.Context) {
	id := c.Param("id")
	b, err := config.BookRepo.Get(c.Request.Context(), id)
	if errors.Is(err, book.ErrNotFound) {
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

func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var req updateTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := config.BookRepo.UpdateTitle(c.Request.Context(), id, req.Title)
	if errors.Is(err, book.ErrNotFound) {
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

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	err := config.BookRepo.Delete(c.Request.Context(), id)
	if errors.Is(err, book.ErrNotFound) {
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

func Search(c *gin.Context) {
	var crit book.SearchCriteria
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
	results, err := config.BookRepo.Search(c.Request.Context(), crit)
	if err != nil {
		log.Printf("search error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func StoreStats(c *gin.Context) {
	books, authors, err := config.BookRepo.StoreStats(c.Request.Context())
	if err != nil {
		log.Printf("store error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books, "authors": authors})
}

func Activity(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}
	actions, err := config.ActivityRepo.Recent(c.Request.Context(), username, consts.ActivityMaxItems)
	if err != nil {
		log.Printf("recent error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"actions": actions})
}
