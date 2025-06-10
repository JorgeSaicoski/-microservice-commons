package server

import (
	"github.com/JorgeSaicoski/microservice-commons/config"
	"github.com/gin-gonic/gin"
)

// SetupRoutesFunc is a function type for route setup
type SetupRoutesFunc func(router *gin.Engine, cfg *config.Config)

// ServerOptions holds configuration options for the server
type ServerOptions struct {
	// Required fields
	ServiceName    string
	ServiceVersion string
	SetupRoutes    SetupRoutesFunc

	// Optional fields with defaults
	Config         *config.Config // If nil, will load from environment
	Port           string         // Override config port
	GinMode        string         // gin.ReleaseMode, gin.DebugMode, gin.TestMode
	DisableLogging bool           // Disable request logging middleware
	DisableCORS    bool           // Disable CORS middleware
	DisableHealth  bool           // Disable health check endpoints
	DisableRecover bool           // Disable recovery middleware

	// Advanced options
	CustomMiddleware []gin.HandlerFunc // Additional middleware to apply
	HealthPath       string            // Custom health check path (default: /health)
	MetricsPath      string            // Custom metrics path (default: /metrics)
}

// DefaultServerOptions returns ServerOptions with sensible defaults
func DefaultServerOptions() ServerOptions {
	return ServerOptions{
		GinMode:        gin.ReleaseMode,
		DisableLogging: false,
		DisableCORS:    false,
		DisableHealth:  false,
		DisableRecover: false,
		HealthPath:     "/health",
		MetricsPath:    "/metrics",
	}
}

// Validate validates the server options
func (o *ServerOptions) Validate() error {
	if o.ServiceName == "" {
		return &ServerError{
			Code:    "invalid_options",
			Message: "ServiceName is required",
		}
	}

	if o.ServiceVersion == "" {
		return &ServerError{
			Code:    "invalid_options",
			Message: "ServiceVersion is required",
		}
	}

	if o.SetupRoutes == nil {
		return &ServerError{
			Code:    "invalid_options",
			Message: "SetupRoutes function is required",
		}
	}

	return nil
}
