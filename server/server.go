package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JorgeSaicoski/microservice-commons/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server with configuration
type Server struct {
	router   *gin.Engine
	server   *http.Server
	config   *config.Config
	options  ServerOptions
	shutdown *ShutdownManager
}

// ServerError represents server-related errors
type ServerError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ServerError) Error() string {
	return e.Message
}

// NewServer creates a new server instance with the given options
func NewServer(options ServerOptions) *Server {
	// Validate options
	if err := options.Validate(); err != nil {
		panic(fmt.Sprintf("Invalid server options: %v", err))
	}

	// Load configuration if not provided
	cfg := options.Config
	if cfg == nil {
		cfg = config.LoadWithServiceInfo(options.ServiceName, options.ServiceVersion)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		panic(fmt.Sprintf("Invalid configuration: %v", err))
	}

	// Set Gin mode
	if options.GinMode != "" {
		gin.SetMode(options.GinMode)
	} else if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Create server instance
	server := &Server{
		router:  router,
		config:  cfg,
		options: options,
	}

	// Setup middleware
	server.setupMiddleware()

	// Setup default routes
	server.setupDefaultRoutes()

	// Setup user routes
	options.SetupRoutes(router, cfg)

	// Create HTTP server
	port := options.Port
	if port == "" {
		port = cfg.Port
	}

	server.server = &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Setup graceful shutdown
	server.shutdown = setupGracefulShutdown(server.server)

	return server
}

// setupMiddleware configures the middleware stack
func (s *Server) setupMiddleware() {
	// Recovery middleware (unless disabled)
	if !s.options.DisableRecover {
		s.router.Use(gin.Recovery())
	}

	// Logging middleware (unless disabled)
	if !s.options.DisableLogging {
		s.router.Use(gin.Logger())
	}

	// CORS middleware (unless disabled)
	if !s.options.DisableCORS {
		s.setupCORS()
	}

	// Custom middleware
	for _, middleware := range s.options.CustomMiddleware {
		s.router.Use(middleware)
	}
}

// setupCORS configures CORS middleware
func (s *Server) setupCORS() {
	corsConfig := cors.Config{
		AllowOrigins:     s.config.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	s.router.Use(cors.New(corsConfig))
}

// setupDefaultRoutes sets up default routes like health checks
func (s *Server) setupDefaultRoutes() {
	if !s.options.DisableHealth {
		s.setupHealthRoutes()
	}
}

// setupHealthRoutes sets up health check endpoints
func (s *Server) setupHealthRoutes() {
	healthPath := s.options.HealthPath
	if healthPath == "" {
		healthPath = "/health"
	}

	s.router.GET(healthPath, s.healthCheckHandler)

	// Detailed health check
	s.router.GET(healthPath+"/detailed", s.detailedHealthHandler)
}

// healthCheckHandler handles basic health checks
func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   s.config.ServiceName,
		"version":   s.config.ServiceVersion,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// detailedHealthHandler handles detailed health checks
func (s *Server) detailedHealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"service":     s.config.ServiceName,
		"version":     s.config.ServiceVersion,
		"environment": s.config.Environment,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"uptime":      time.Since(time.Now()).String(), // This would need to be tracked properly
		"config": gin.H{
			"port":        s.config.Port,
			"log_level":   s.config.LogLevel,
			"environment": s.config.Environment,
		},
	})
}

// Start starts the server with graceful shutdown handling
func (s *Server) Start() error {
	fmt.Printf("Starting %s v%s\n", s.config.ServiceName, s.config.ServiceVersion)
	fmt.Printf("Environment: %s\n", s.config.Environment)
	fmt.Printf("Log Level: %s\n", s.config.LogLevel)

	return s.shutdown.StartWithGracefulShutdown()
}

// Stop stops the server gracefully
func (s *Server) Stop() error {
	return s.shutdown.WaitForShutdown()
}

// ForceStop forces immediate server shutdown
func (s *Server) ForceStop() error {
	return s.shutdown.ForceShutdown()
}

// GetRouter returns the Gin router (useful for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// GetConfig returns the server configuration
func (s *Server) GetConfig() *config.Config {
	return s.config
}

// GetHTTPServer returns the underlying HTTP server
func (s *Server) GetHTTPServer() *http.Server {
	return s.server
}
