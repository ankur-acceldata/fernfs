package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RespondWithSuccess sends a success response
func RespondWithSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// RespondWithError sends an error response
func RespondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Status: "error",
		Error:  message,
	})
}

// RespondWithValidationError sends a validation error response
func RespondWithValidationError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusBadRequest, message)
}

// RespondWithServerError sends a server error response
func RespondWithServerError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusInternalServerError, message)
}

// RespondWithNotFound sends a not found error response
func RespondWithNotFound(c *gin.Context, message string) {
	RespondWithError(c, http.StatusNotFound, message)
} 