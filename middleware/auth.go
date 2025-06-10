package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthConfig holds configuration for authentication middleware
type AuthConfig struct {
	SkipPaths      []string
	TokenExtractor func(*gin.Context) (string, error)
	TokenValidator func(string) (map[string]interface{}, error)
	ErrorHandler   func(*gin.Context, error)
}

// AuthError represents authentication errors
type AuthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AuthError) Error() string {
	return e.Message
}

// Common auth error types
var (
	ErrMissingToken = &AuthError{
		Code:    "missing_token",
		Message: "Authorization token is required",
	}
	ErrInvalidToken = &AuthError{
		Code:    "invalid_token",
		Message: "Invalid authorization token",
	}
	ErrExpiredToken = &AuthError{
		Code:    "expired_token",
		Message: "Authorization token has expired",
	}
	ErrInsufficientPermissions = &AuthError{
		Code:    "insufficient_permissions",
		Message: "Insufficient permissions for this operation",
	}
)

// DefaultAuthConfig returns default authentication configuration
func DefaultAuthConfig() AuthConfig {
	return AuthConfig{
		SkipPaths:      []string{"/health", "/metrics"},
		TokenExtractor: BearerTokenExtractor,
		ErrorHandler:   DefaultAuthErrorHandler,
	}
}

// NewAuthMiddleware creates an authentication middleware with the given configuration
func NewAuthMiddleware(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path should be skipped
		if shouldSkipAuth(c.Request.URL.Path, config.SkipPaths) {
			c.Next()
			return
		}

		// Extract token
		token, err := config.TokenExtractor(c)
		if err != nil {
			config.ErrorHandler(c, err)
			return
		}

		// Validate token
		if config.TokenValidator != nil {
			claims, err := config.TokenValidator(token)
			if err != nil {
				config.ErrorHandler(c, err)
				return
			}

			// Store claims in context
			for key, value := range claims {
				c.Set(key, value)
			}
		}

		c.Next()
	}
}

// BearerTokenExtractor extracts Bearer token from Authorization header
func BearerTokenExtractor(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", ErrMissingToken
	}

	// Check if it starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", &AuthError{
			Code:    "invalid_auth_format",
			Message: "Authorization header must use Bearer format",
		}
	}

	// Extract the token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", ErrMissingToken
	}

	return token, nil
}

// APIKeyExtractor extracts API key from X-API-Key header
func APIKeyExtractor(c *gin.Context) (string, error) {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		return "", &AuthError{
			Code:    "missing_api_key",
			Message: "API key is required",
		}
	}
	return apiKey, nil
}

// QueryTokenExtractor extracts token from query parameter
func QueryTokenExtractor(paramName string) func(*gin.Context) (string, error) {
	return func(c *gin.Context) (string, error) {
		token := c.Query(paramName)
		if token == "" {
			return "", &AuthError{
				Code:    "missing_token",
				Message: "Token parameter is required",
			}
		}
		return token, nil
	}
}

// DefaultAuthErrorHandler handles authentication errors
func DefaultAuthErrorHandler(c *gin.Context, err error) {
	if authErr, ok := err.(*AuthError); ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": authErr.Message,
			"code":  authErr.Code,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication failed",
			"code":  "auth_failed",
		})
	}
	c.Abort()
}

// shouldSkipAuth checks if authentication should be skipped for the given path
func shouldSkipAuth(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// RequireAuth creates a simple authentication middleware that requires a valid token
func RequireAuth(tokenValidator func(string) (map[string]interface{}, error)) gin.HandlerFunc {
	config := DefaultAuthConfig()
	config.TokenValidator = tokenValidator
	return NewAuthMiddleware(config)
}

// OptionalAuth creates an authentication middleware that doesn't fail if no token is provided
func OptionalAuth(tokenValidator func(string) (map[string]interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to extract token
		token, err := BearerTokenExtractor(c)
		if err != nil {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		// Validate token if provided
		if tokenValidator != nil {
			claims, err := tokenValidator(token)
			if err != nil {
				// Invalid token, but don't fail the request
				c.Next()
				return
			}

			// Store claims in context
			for key, value := range claims {
				c.Set(key, value)
			}
		}

		c.Next()
	}
}

// APIKeyAuth creates an API key authentication middleware
func APIKeyAuth(validAPIKeys map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			DefaultAuthErrorHandler(c, &AuthError{
				Code:    "missing_api_key",
				Message: "API key is required",
			})
			return
		}

		// Check if API key is valid
		if userID, exists := validAPIKeys[apiKey]; exists {
			c.Set("user_id", userID)
			c.Set("auth_method", "api_key")
			c.Next()
		} else {
			DefaultAuthErrorHandler(c, &AuthError{
				Code:    "invalid_api_key",
				Message: "Invalid API key",
			})
		}
	}
}

// BasicAuth creates a basic authentication middleware
func BasicAuth(validCredentials map[string]string) gin.HandlerFunc {
	return gin.BasicAuth(validCredentials)
}

// RequireRole creates a middleware that requires a specific role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user has the required role
		userRoles, exists := c.Get("roles")
		if !exists {
			DefaultAuthErrorHandler(c, ErrInsufficientPermissions)
			return
		}

		// Convert to string slice
		roles, ok := userRoles.([]string)
		if !ok {
			DefaultAuthErrorHandler(c, ErrInsufficientPermissions)
			return
		}

		// Check if required role is present
		for _, userRole := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		DefaultAuthErrorHandler(c, ErrInsufficientPermissions)
	}
}

// RequireAnyRole creates a middleware that requires any of the specified roles
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("roles")
		if !exists {
			DefaultAuthErrorHandler(c, ErrInsufficientPermissions)
			return
		}

		userRolesList, ok := userRoles.([]string)
		if !ok {
			DefaultAuthErrorHandler(c, ErrInsufficientPermissions)
			return
		}

		// Check if any required role is present
		for _, requiredRole := range roles {
			for _, userRole := range userRolesList {
				if userRole == requiredRole {
					c.Next()
					return
				}
			}
		}

		DefaultAuthErrorHandler(c, ErrInsufficientPermissions)
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	if id, ok := userID.(string); ok {
		return id, true
	}

	return "", false
}

// GetUserRoles extracts user roles from context
func GetUserRoles(c *gin.Context) ([]string, bool) {
	roles, exists := c.Get("roles")
	if !exists {
		return nil, false
	}

	if rolesList, ok := roles.([]string); ok {
		return rolesList, true
	}

	return nil, false
}

// HasRole checks if the user has a specific role
func HasRole(c *gin.Context, role string) bool {
	roles, exists := GetUserRoles(c)
	if !exists {
		return false
	}

	for _, userRole := range roles {
		if userRole == role {
			return true
		}
	}

	return false
}
