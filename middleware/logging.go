package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LogLevel represents different logging levels
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LoggingConfig holds configuration for logging middleware
type LoggingConfig struct {
	Level        LogLevel
	SkipPaths    []string
	CustomFormat gin.LogFormatter
}

// DefaultLoggingConfig returns default logging configuration
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:     LogLevelInfo,
		SkipPaths: []string{"/health", "/metrics"},
	}
}

// NewLoggingMiddleware creates a new logging middleware with the given configuration
func NewLoggingMiddleware(config LoggingConfig) gin.HandlerFunc {
	// Use custom format if provided, otherwise use default
	if config.CustomFormat != nil {
		return gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: config.CustomFormat,
			SkipPaths: config.SkipPaths,
		})
	}

	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: defaultLogFormat,
		SkipPaths: config.SkipPaths,
	})
}

// DefaultLoggingMiddleware creates a logging middleware with default configuration
func DefaultLoggingMiddleware() gin.HandlerFunc {
	return NewLoggingMiddleware(DefaultLoggingConfig())
}

// DetailedLoggingMiddleware creates a logging middleware with detailed output
func DetailedLoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: detailedLogFormat,
		SkipPaths: []string{"/health"},
	})
}

// ProductionLoggingMiddleware creates a logging middleware optimized for production
func ProductionLoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: productionLogFormat,
		SkipPaths: []string{"/health", "/metrics"},
	})
}

// SilentLoggingMiddleware creates a logging middleware that only logs errors
func SilentLoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			// Only log errors and slow requests
			if param.StatusCode >= 400 || param.Latency > time.Second {
				return fmt.Sprintf("[%s] %s %s %d %s %s\n",
					param.TimeStamp.Format("2006/01/02 - 15:04:05"),
					param.Method,
					param.Path,
					param.StatusCode,
					param.Latency,
					param.ErrorMessage,
				)
			}
			return ""
		},
		SkipPaths: []string{"/health", "/metrics"},
	})
}

// defaultLogFormat provides a clean, readable log format
func defaultLogFormat(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

// detailedLogFormat provides detailed logging with request/response info
func detailedLogFormat(param gin.LogFormatterParams) string {
	return fmt.Sprintf(`[%s] "%s %s %s" %d %s "%s" "%s" %s %s
`,
		param.TimeStamp.Format("2006/01/02 15:04:05"),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.Request.Referer(),
		param.ClientIP,
		param.ErrorMessage,
	)
}

// productionLogFormat provides JSON-like structured logging for production
func productionLogFormat(param gin.LogFormatterParams) string {
	return fmt.Sprintf(`{"time":"%s","method":"%s","path":"%s","status":%d,"latency":"%s","ip":"%s","user_agent":"%s","error":"%s"}
`,
		param.TimeStamp.Format(time.RFC3339),
		param.Method,
		param.Path,
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

// RequestLogger logs individual request details
func RequestLogger(level LogLevel) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request details
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		// Log based on level and status
		switch {
		case statusCode >= 500:
			fmt.Printf("[ERROR] %s | %3d | %13v | %15s | %-7s %s\n",
				time.Now().Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
			)
		case statusCode >= 400:
			fmt.Printf("[WARN] %s | %3d | %13v | %15s | %-7s %s\n",
				time.Now().Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
			)
		case level == LogLevelDebug:
			fmt.Printf("[DEBUG] %s | %3d | %13v | %15s | %-7s %s\n",
				time.Now().Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
			)
		case level == LogLevelInfo && statusCode < 400:
			fmt.Printf("[INFO] %s | %3d | %13v | %15s | %-7s %s\n",
				time.Now().Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
			)
		}
	}
}
