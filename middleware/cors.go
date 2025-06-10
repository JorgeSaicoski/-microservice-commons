package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration options
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// DefaultCORSConfig returns a CORS configuration with sensible defaults
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Cache-Control",
		},
		ExposedHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// NewCORSMiddleware creates a new CORS middleware with the given configuration
func NewCORSMiddleware(config CORSConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     config.AllowedOrigins,
		AllowMethods:     config.AllowedMethods,
		AllowHeaders:     config.AllowedHeaders,
		ExposeHeaders:    config.ExposedHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           config.MaxAge,
	}

	return cors.New(corsConfig)
}

// DefaultCORSMiddleware creates a CORS middleware with default configuration
func DefaultCORSMiddleware() gin.HandlerFunc {
	return NewCORSMiddleware(DefaultCORSConfig())
}

// CustomCORSMiddleware creates a CORS middleware with custom origins
func CustomCORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	config := DefaultCORSConfig()
	config.AllowedOrigins = allowedOrigins
	return NewCORSMiddleware(config)
}

// DevelopmentCORSMiddleware creates a permissive CORS middleware for development
func DevelopmentCORSMiddleware() gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(corsConfig)
}

// ProductionCORSMiddleware creates a strict CORS middleware for production
func ProductionCORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	config := CORSConfig{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposedHeaders: []string{
			"Content-Type",
			"X-Request-ID",
		},
		AllowCredentials: false, // More secure for production
		MaxAge:           1 * time.Hour,
	}

	return NewCORSMiddleware(config)
}
