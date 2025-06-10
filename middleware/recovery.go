package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// RecoveryConfig holds configuration for recovery middleware
type RecoveryConfig struct {
	EnableStackTrace bool
	CustomHandler    gin.RecoveryFunc
	SkipPaths        []string
}

// DefaultRecoveryConfig returns default recovery configuration
func DefaultRecoveryConfig() RecoveryConfig {
	return RecoveryConfig{
		EnableStackTrace: true,
		CustomHandler:    nil,
		SkipPaths:        []string{},
	}
}

// NewRecoveryMiddleware creates a new recovery middleware
func NewRecoveryMiddleware(config RecoveryConfig) gin.HandlerFunc {
	if config.CustomHandler != nil {
		return gin.CustomRecovery(config.CustomHandler)
	}

	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Check if path should be skipped
		for _, path := range config.SkipPaths {
			if c.Request.URL.Path == path {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		// Get request ID if available
		requestID := MustGetRequestID(c)

		// Log the panic
		logPanic(recovered, config.EnableStackTrace, requestID, c)

		// Return error response
		errorResponse := gin.H{
			"error": "Internal server error",
			"code":  "internal_error",
		}

		if requestID != "" {
			errorResponse["request_id"] = requestID
		}

		c.JSON(http.StatusInternalServerError, errorResponse)
	})
}

// DefaultRecoveryMiddleware creates a recovery middleware with default configuration
func DefaultRecoveryMiddleware() gin.HandlerFunc {
	return NewRecoveryMiddleware(DefaultRecoveryConfig())
}

// ProductionRecoveryMiddleware creates a recovery middleware for production use
func ProductionRecoveryMiddleware() gin.HandlerFunc {
	config := RecoveryConfig{
		EnableStackTrace: false, // Don't log stack traces in production
		CustomHandler: func(c *gin.Context, recovered interface{}) {
			requestID := MustGetRequestID(c)

			// Log error without stack trace
			fmt.Printf("[PANIC RECOVERY] Request ID: %s, Error: %v\n", requestID, recovered)

			// Return minimal error response
			errorResponse := gin.H{
				"error": "Internal server error",
				"code":  "internal_error",
			}

			if requestID != "" {
				errorResponse["request_id"] = requestID
			}

			c.JSON(http.StatusInternalServerError, errorResponse)
		},
	}

	return NewRecoveryMiddleware(config)
}

// DevelopmentRecoveryMiddleware creates a recovery middleware for development use
func DevelopmentRecoveryMiddleware() gin.HandlerFunc {
	config := RecoveryConfig{
		EnableStackTrace: true,
		CustomHandler: func(c *gin.Context, recovered interface{}) {
			requestID := MustGetRequestID(c)

			// Log detailed error with stack trace
			logPanic(recovered, true, requestID, c)

			// Return detailed error response for development
			errorResponse := gin.H{
				"error":   "Internal server error",
				"code":    "internal_error",
				"details": fmt.Sprintf("%v", recovered),
				"path":    c.Request.URL.Path,
				"method":  c.Request.Method,
			}

			if requestID != "" {
				errorResponse["request_id"] = requestID
			}

			c.JSON(http.StatusInternalServerError, errorResponse)
		},
	}

	return NewRecoveryMiddleware(config)
}

// SilentRecoveryMiddleware creates a recovery middleware that doesn't log
func SilentRecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := MustGetRequestID(c)

		errorResponse := gin.H{
			"error": "Internal server error",
			"code":  "internal_error",
		}

		if requestID != "" {
			errorResponse["request_id"] = requestID
		}

		c.JSON(http.StatusInternalServerError, errorResponse)
	})
}

// CustomErrorRecoveryMiddleware creates a recovery middleware with custom error responses
func CustomErrorRecoveryMiddleware(errorHandler func(c *gin.Context, err interface{}, requestID string)) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := MustGetRequestID(c)
		errorHandler(c, recovered, requestID)
	})
}

// logPanic logs panic information
func logPanic(recovered interface{}, enableStackTrace bool, requestID string, c *gin.Context) {
	fmt.Printf("[PANIC RECOVERY] Request ID: %s\n", requestID)
	fmt.Printf("[PANIC RECOVERY] URL: %s %s\n", c.Request.Method, c.Request.URL.Path)
	fmt.Printf("[PANIC RECOVERY] Error: %v\n", recovered)

	if enableStackTrace {
		fmt.Printf("[PANIC RECOVERY] Stack trace:\n%s\n", debug.Stack())
	}
}

// PanicHandler is a type for custom panic handling functions
type PanicHandler func(c *gin.Context, recovered interface{}, requestID string)

// WithCustomPanicHandler creates a recovery middleware with a custom panic handler
func WithCustomPanicHandler(handler PanicHandler) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := MustGetRequestID(c)
		handler(c, recovered, requestID)
	})
}

// JSONErrorPanicHandler returns a panic handler that responds with JSON errors
func JSONErrorPanicHandler(includeDetails bool) PanicHandler {
	return func(c *gin.Context, recovered interface{}, requestID string) {
		logPanic(recovered, true, requestID, c)

		errorResponse := gin.H{
			"error": "Internal server error",
			"code":  "internal_error",
		}

		if requestID != "" {
			errorResponse["request_id"] = requestID
		}

		if includeDetails {
			errorResponse["details"] = fmt.Sprintf("%v", recovered)
		}

		c.JSON(http.StatusInternalServerError, errorResponse)
	}
}
