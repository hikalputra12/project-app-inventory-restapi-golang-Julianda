package utils

import (
	"app-inventory/dto"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

// ResponseSuccess sends a successful JSON response with Gin
func ResponseSuccess(c *gin.Context, code int, message string, data any) {
	c.JSON(code, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// ResponseError sends an error JSON response with Gin
func ResponseError(c *gin.Context, code int, message string, errs any) {
	c.JSON(code, Response{
		Status:  false,
		Message: message,
		Errors:  errs,
	})
}

// ResponseBadRequest sends a 400/error JSON response with Gin
func ResponseBadRequest(c *gin.Context, code int, message string, errs any) {
	ResponseError(c, code, message, errs)
}

// ResponsePagination sends a paginated JSON response with Gin
func ResponsePagination(c *gin.Context, code int, message string, data any, pagination dto.Pagination) {
	c.JSON(code, gin.H{
		"status":     true,
		"message":    message,
		"data":       data,
		"pagination": pagination,
	})
}
