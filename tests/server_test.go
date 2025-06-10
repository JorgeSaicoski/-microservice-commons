package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JorgeSaicoski/microservice-commons/config"
	"github.com/JorgeSaicoski/microservice-commons/server"
	"github.com/gin-gonic/gin"
)

func TestNewServer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	options := server.ServerOptions{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		SetupRoutes: func(router *gin.Engine, cfg *config.Config) {
			router.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "test"})
			})
		},
	}

	srv := server.NewServer(options)

	if srv == nil {
		t.Error("Expected server to be created")
	}

	router := srv.GetRouter()
	if router == nil {
		t.Error("Expected router to be available")
	}
}

func TestServerOptions_Validate(t *testing.T) {
	tests := []struct {
		name        string
		options     server.ServerOptions
		expectError bool
	}{
		{
			name: "valid options",
			options: server.ServerOptions{
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
				SetupRoutes:    func(router *gin.Engine, cfg *config.Config) {},
			},
			expectError: false,
		},
		{
			name: "missing service name",
			options: server.ServerOptions{
				ServiceVersion: "1.0.0",
				SetupRoutes:    func(router *gin.Engine, cfg *config.Config) {},
			},
			expectError: true,
		},
		{
			name: "missing setup routes",
			options: server.ServerOptions{
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Validate()
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestHealthEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	options := server.ServerOptions{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		SetupRoutes:    func(router *gin.Engine, cfg *config.Config) {},
	}

	srv := server.NewServer(options)
	router := srv.GetRouter()

	// Test basic health endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Test detailed health endpoint
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/health/detailed", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestServerWithCustomConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	customConfig := &config.Config{
		Port:           "9999",
		ServiceName:    "custom-service",
		ServiceVersion: "2.0.0",
		AllowedOrigins: []string{"http://example.com"},
		DatabaseConfig: config.DatabaseConfig{
			Host:         "localhost",
			Port:         "5432",
			User:         "user",
			Password:     "pass",
			DatabaseName: "testdb",
		},
		KeycloakConfig: config.KeycloakConfig{
			PublicKeyBase64: "test-key",
		},
	}

	options := server.ServerOptions{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Config:         customConfig,
		SetupRoutes:    func(router *gin.Engine, cfg *config.Config) {},
	}

	srv := server.NewServer(options)
	config := srv.GetConfig()

	if config.ServiceName != "custom-service" {
		t.Errorf("Expected service name 'custom-service', got %s", config.ServiceName)
	}

	if config.Port != "9999" {
		t.Errorf("Expected port '9999', got %s", config.Port)
	}
}
