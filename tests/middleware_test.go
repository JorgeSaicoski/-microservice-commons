package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JorgeSaicoski/microservice-commons/middleware"
	"github.com/gin-gonic/gin"
)

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.DefaultRequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		requestID := middleware.MustGetRequestID(c)
		c.JSON(200, gin.H{"request_id": requestID})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Check that request ID header is present
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("Expected X-Request-ID header to be set")
	}

	// Check that response contains request ID
	if !strings.Contains(w.Body.String(), requestID) {
		t.Error("Expected response to contain request ID")
	}
}

func TestRequestIDWithExistingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.DefaultRequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		requestID := middleware.MustGetRequestID(c)
		c.JSON(200, gin.H{"request_id": requestID})
	})

	existingID := "existing-request-id"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", existingID)
	router.ServeHTTP(w, req)

	// Should use existing request ID
	requestID := w.Header().Get("X-Request-ID")
	if requestID != existingID {
		t.Errorf("Expected request ID %s, got %s", existingID, requestID)
	}
}

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.CustomCORSMiddleware([]string{"http://localhost:3000"}))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	// Check CORS headers
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("Expected Access-Control-Allow-Origin header")
	}
}

func TestHealthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.SimpleHealthMiddleware("test-service", "1.0.0"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Body.String(), "test-service") {
		t.Error("Expected response to contain service name")
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.DefaultRecoveryMiddleware())
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if !strings.Contains(w.Body.String(), "internal_error") {
		t.Error("Expected response to contain error code")
	}
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock token validator
	tokenValidator := func(token string) (map[string]interface{}, error) {
		if token == "valid-token" {
			return map[string]interface{}{
				"user_id": "123",
				"roles":   []string{"user"},
			}, nil
		}
		return nil, middleware.ErrInvalidToken
	}

	router := gin.New()
	router.Use(middleware.RequireAuth(tokenValidator))
	router.GET("/protected", func(c *gin.Context) {
		userID, _ := middleware.GetUserID(c)
		c.JSON(200, gin.H{"user_id": userID})
	})

	// Test with valid token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Test with invalid token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Test without token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
