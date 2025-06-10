package responses

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Success sends a successful response with data
func Success(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	c.JSON(http.StatusOK, response)
}

// SuccessWithStatus sends a successful response with custom status code
func SuccessWithStatus(c *gin.Context, status int, message string, data interface{}) {
	response := SuccessResponse{
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	c.JSON(status, response)
}

// Created sends a 201 Created response
func Created(c *gin.Context, message string, data interface{}) {
	SuccessWithStatus(c, http.StatusCreated, message, data)
}

// Accepted sends a 202 Accepted response
func Accepted(c *gin.Context, message string, data interface{}) {
	SuccessWithStatus(c, http.StatusAccepted, message, data)
}

// NoContent sends a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// OK sends a simple 200 OK response with just a message
func OK(c *gin.Context, message string) {
	Success(c, message, nil)
}

// Data sends a successful response with only data (no message)
func Data(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data":      data,
		"timestamp": time.Now().UTC(),
	})
}

// Message sends a successful response with only a message (no data)
func Message(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"message":   message,
		"timestamp": time.Now().UTC(),
	})
}

// JSON sends a custom JSON response
func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// WithHeaders sends a response with custom headers
func WithHeaders(c *gin.Context, status int, headers map[string]string, data interface{}) {
	for key, value := range headers {
		c.Header(key, value)
	}
	c.JSON(status, data)
}

// WithRequestID sends a response including the request ID if available
func WithRequestID(c *gin.Context, status int, data interface{}) {
	// Try to get request ID from context
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			c.Header("X-Request-ID", id)
		}
	}
	c.JSON(status, data)
}

// File sends a file download response
func File(c *gin.Context, filepath, filename string) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filepath)
}

// Download sends a file as attachment
func Download(c *gin.Context, filepath, filename string) {
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.File(filepath)
}

// Redirect sends a redirect response
func Redirect(c *gin.Context, location string) {
	c.Redirect(http.StatusFound, location)
}

// PermanentRedirect sends a permanent redirect response
func PermanentRedirect(c *gin.Context, location string) {
	c.Redirect(http.StatusMovedPermanently, location)
}
