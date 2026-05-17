package utilities

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/giladrozner/book_service/pkg/config"
)

// ParseElasticsearchErrorCode extracts the HTTP status code from an Elasticsearch
// error string and returns it as an HttpError. Falls back to 500 if parsing fails.
func ParseElasticsearchErrorCode(err error) config.HttpError {
	re := regexp.MustCompile(`Error (\d{3})`)
	matches := re.FindStringSubmatch(err.Error())
	if len(matches) < 2 {
		return config.HttpError{Code: 500}
	}
	code, convErr := strconv.Atoi(matches[1])
	if convErr != nil {
		return config.HttpError{Code: 500}
	}
	return config.HttpError{Code: code}
}

// GetAddBookValidationErrors parses field-level validation errors from a failed
// ShouldBindJSON call on a create-book request and returns a map of field → message.
func GetAddBookValidationErrors(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	errorMessages := make(map[string]string)

	if errors.As(err, &validationErrors) {
		for _, fieldErr := range validationErrors {
			switch fieldErr.Field() {
			case "Title":
				errorMessages[config.FieldTitle] = config.ErrTitleRequired
			case "AuthorName":
				errorMessages[config.FieldAuthorName] = config.ErrAuthorNameRequired
			case "Price":
				errorMessages[config.FieldPrice] = config.ErrPriceRequired
			case "PublishDate":
				errorMessages[config.FieldPublishDate] = config.ErrPublishDateRequired
			default:
				errorMessages[fieldErr.Field()] = config.ErrFieldRequired
			}
		}
	}
	return errorMessages
}

// GetUpdateBookValidationErrors parses field-level validation errors from a failed
// ShouldBindJSON call on an update-book request and returns a map of field → message.
func GetUpdateBookValidationErrors(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	errorMessages := make(map[string]string)

	if errors.As(err, &validationErrors) {
		for _, fieldErr := range validationErrors {
			switch fieldErr.Field() {
			case "Title":
				errorMessages[config.FieldTitle] = config.ErrTitleRequired
			default:
				errorMessages[fieldErr.Field()] = config.ErrFieldRequired
			}
		}
	}
	return errorMessages
}
