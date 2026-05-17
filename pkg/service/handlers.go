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
	id, err := BookRepo.Add(c.Request.Context(), book.Book{
		Title:          req.Title,
		AuthorName:     req.AuthorName,
		Price:          req.Price,
		EbookAvailable: req.EbookAvailable,
		PublishDate:    req.PublishDate,
	})
	if err != nil {
		log.Printf("create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.ErrMsgInternal})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func GetBook(c *gin.Context) {
	id := c.Param("id")
	b, err := BookRepo.Get(c.Request.Context(), id)
	if errors.Is(err, book.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": config.ErrMsgBookNotFound})
		return
	}
	if err != nil {
		log.Printf("get error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.ErrMsgInternal})
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
	err := BookRepo.UpdateTitle(c.Request.Context(), id, req.Title)
	if errors.Is(err, book.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": config.ErrMsgBookNotFound})
		return
	}
	if err != nil {
		log.Printf("update error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.ErrMsgInternal})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	err := BookRepo.Delete(c.Request.Context(), id)
	if errors.Is(err, book.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": config.ErrMsgBookNotFound})
		return
	}
	if err != nil {
		log.Printf("delete error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.ErrMsgInternal})
		return
	}
	c.Status(http.StatusNoContent)
}

func Search(c *gin.Context) {
	var crit book.SearchCriteria
	if t := c.Query(config.QueryParamTitle); t != "" {
		crit.Title = &t
	}
	if a := c.Query(config.QueryParamAuthorName); a != "" {
		crit.AuthorName = &a
	}
	if pr := c.Query(config.QueryParamPriceRange); pr != "" {
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
	results, err := BookRepo.Search(c.Request.Context(), crit)
	if err != nil {
		log.Printf("search error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.ErrMsgInternal})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func StoreStats(c *gin.Context) {
	books, authors, err := BookRepo.StoreStats(c.Request.Context())
	if err != nil {
		log.Printf("store error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.ErrMsgInternal})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books, "authors": authors})
}

func Activity(c *gin.Context) {
	username := c.Query(config.QueryParamUsername)
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.ErrMsgUsernameRequired})
		return
	}
	actions, err := ActivityRepo.Recent(c.Request.Context(), username, config.ActivityMaxItems)
	if err != nil {
		log.Printf("recent error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.ErrMsgInternal})
		return
	}
	c.JSON(http.StatusOK, gin.H{"actions": actions})
}
