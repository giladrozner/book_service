package service

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/giladrozner/book_service/pkg/book"
	"github.com/giladrozner/book_service/pkg/config"
	"github.com/giladrozner/book_service/pkg/utilities"
)

type createBookRequest struct {
	Title          string  `json:"title" binding:"required"`
	AuthorName     string  `json:"author_name" binding:"required"`
	Price          float64 `json:"price" binding:"required"`
	EbookAvailable *bool   `json:"ebook_available"`
	PublishDate    string  `json:"publish_date" binding:"required"`
}

type updateTitleRequest struct {
	Title string `json:"title" binding:"required"`
}

func CreateBook(c *gin.Context) {
	var req createBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorField: utilities.GetAddBookValidationErrors(err)})
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
		httpErr := utilities.ParseElasticsearchErrorCode(err)
		c.JSON(httpErr.Code, gin.H{config.ErrorField: config.ElasticsearchErrorMap[httpErr.Code]})
		return
	}
	c.JSON(http.StatusCreated, gin.H{config.ParamID: id})
}

func GetBook(c *gin.Context) {
	id := c.Param(config.ParamID)
	b, err := BookRepo.Get(c.Request.Context(), id)
	if err != nil {

		httpErr := utilities.ParseElasticsearchErrorCode(err)
		c.JSON(httpErr.Code, gin.H{config.ErrorField: config.ElasticsearchErrorMap[httpErr.Code]})
		return
	}
	c.JSON(http.StatusOK, b)
}

func UpdateBook(c *gin.Context) {
	id := c.Param(config.ParamID)
	var req updateTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorField: utilities.GetUpdateBookValidationErrors(err)})
		return
	}
	if err := BookRepo.UpdateTitle(c.Request.Context(), id, req.Title); err != nil {
		httpErr := utilities.ParseElasticsearchErrorCode(err)
		c.JSON(httpErr.Code, gin.H{config.ErrorField: config.ElasticsearchErrorMap[httpErr.Code]})
		return
	}
	c.JSON(http.StatusOK, gin.H{config.MessageField: config.SuccessBookUpdated})
}

func DeleteBook(c *gin.Context) {
	id := c.Param(config.ParamID)
	if err := BookRepo.Delete(c.Request.Context(), id); err != nil {
		httpErr := utilities.ParseElasticsearchErrorCode(err)
		c.JSON(httpErr.Code, gin.H{config.ErrorField: config.ElasticsearchErrorMap[httpErr.Code]})
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
		httpErr := utilities.ParseElasticsearchErrorCode(err)
		c.JSON(httpErr.Code, gin.H{config.ErrorField: config.ElasticsearchErrorMap[httpErr.Code]})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func StoreStats(c *gin.Context) {
	books, authors, err := BookRepo.StoreStats(c.Request.Context())
	if err != nil {
		httpErr := utilities.ParseElasticsearchErrorCode(err)
		c.JSON(httpErr.Code, gin.H{config.ErrorField: config.ElasticsearchErrorMap[httpErr.Code]})
		return
	}
	c.JSON(http.StatusOK, gin.H{config.BookCountField: books, config.DistinctAuthorsField: authors})
}

func Activity(c *gin.Context) {
	username := c.Query(config.QueryParamUsername)
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorField: config.ErrMsgUsernameRequired})
		return
	}
	actions, err := ActivityRepo.Recent(c.Request.Context(), username, config.ActivityMaxItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorField: config.ErrMsgInternal})
		return
	}
	c.JSON(http.StatusOK, gin.H{config.ActionsField: actions})
}
