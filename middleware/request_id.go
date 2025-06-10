package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"
	// RequestIDKey is the context key for request ID
	RequestIDKey = "request_id"
)

// RequestIDConfig holds configuration for request ID middleware
type RequestIDConfig struct {
	HeaderName string        // Custom header name (default: X-Request-ID)
	ContextKey string        // Custom context key (default: request_id)
	Generator  func() string // Custom ID generator
}

// DefaultRequestIDConfig returns default request ID configuration
func DefaultRequestIDConfig() RequestIDConfig {
	return RequestIDConfig{
		HeaderName: RequestIDHeader,
		ContextKey: RequestIDKey,
		Generator:  generateRequestID,
	}
}

// NewRequestIDMiddleware creates a new request ID middleware
func NewRequestIDMiddleware(config RequestIDConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists in header
		requestID := c.GetHeader(config.HeaderName)

		// Generate new ID if not present
		if requestID == "" {
			requestID = config.Generator()
		}

		// Set request ID in context
		c.Set(config.ContextKey, requestID)

		// Set request ID in response header
		c.Header(config.HeaderName, requestID)

		c.Next()
	}
}

// DefaultRequestIDMiddleware creates a request ID middleware with default configuration
func DefaultRequestIDMiddleware() gin.HandlerFunc {
	return NewRequestIDMiddleware(DefaultRequestIDConfig())
}

// CustomRequestIDMiddleware creates a request ID middleware with custom generator
func CustomRequestIDMiddleware(generator func() string) gin.HandlerFunc {
	config := DefaultRequestIDConfig()
	config.Generator = generator
	return NewRequestIDMiddleware(config)
}

// GetRequestID extracts request ID from Gin context
func GetRequestID(c *gin.Context) (string, bool) {
	requestID, exists := c.Get(RequestIDKey)
	if !exists {
		return "", false
	}

	if id, ok := requestID.(string); ok {
		return id, true
	}

	return "", false
}

// MustGetRequestID extracts request ID from context or returns empty string
func MustGetRequestID(c *gin.Context) string {
	if id, exists := GetRequestID(c); exists {
		return id
	}
	return ""
}

// generateRequestID generates a random request ID
func generateRequestID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("req_%d", getCurrentTimestamp())
	}
	return hex.EncodeToString(bytes)
}

// generateShortRequestID generates a shorter request ID (4 bytes)
func generateShortRequestID() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("req_%d", getCurrentTimestamp()%10000)
	}
	return hex.EncodeToString(bytes)
}

// generateUUIDRequestID generates a UUID-like request ID
func generateUUIDRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return generateRequestID()
	}

	// Set version (4) and variant bits
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}

// getCurrentTimestamp returns current timestamp in milliseconds
func getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// ShortRequestIDMiddleware creates middleware with shorter request IDs
func ShortRequestIDMiddleware() gin.HandlerFunc {
	return CustomRequestIDMiddleware(generateShortRequestID)
}

// UUIDRequestIDMiddleware creates middleware with UUID-style request IDs
func UUIDRequestIDMiddleware() gin.HandlerFunc {
	return CustomRequestIDMiddleware(generateUUIDRequestID)
}
