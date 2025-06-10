# Middleware Documentation

This guide covers all middleware components available in microservice-commons, from basic request handling to advanced authentication and monitoring.

## Table of Contents

1. [Overview](#overview)
2. [CORS Middleware](#cors-middleware)
3. [Authentication Middleware](#authentication-middleware)
4. [Health Check Middleware](#health-check-middleware)
5. [Logging Middleware](#logging-middleware)
6. [Recovery Middleware](#recovery-middleware)
7. [Request ID Middleware](#request-id-middleware)
8. [Custom Middleware](#custom-middleware)
9. [Middleware Ordering](#middleware-ordering)
10. [Best Practices](#best-practices)

## Overview

Middleware in microservice-commons provides reusable functionality that runs before or after your route handlers. All middleware is designed to work seamlessly with Gin and follows consistent patterns.

### Basic Usage

```go
import (
    "github.com/JorgeSaicoski/microservice-commons/middleware"
    "github.com/gin-gonic/gin"
)

func setupMiddleware(router *gin.Engine) {
    // Basic middleware stack
    router.Use(middleware.DefaultRecoveryMiddleware())
    router.Use(middleware.DefaultLoggingMiddleware())
    router.Use(middleware.DefaultCORSMiddleware())
    router.Use(middleware.DefaultRequestIDMiddleware())
}
```

## CORS Middleware

Cross-Origin Resource Sharing (CORS) middleware handles browser security policies for cross-origin requests.

### Quick Start

```go
// Default CORS (allows localhost:3000)
router.Use(middleware.DefaultCORSMiddleware())

// Custom origins
origins := []string{"https://app.company.com", "https://admin.company.com"}
router.Use(middleware.CustomCORSMiddleware(origins))

// Development (allows all origins)
router.Use(middleware.DevelopmentCORSMiddleware())

// Production (strict settings)
origins := []string{"https://app.company.com"}
router.Use(middleware.ProductionCORSMiddleware(origins))
```

### Advanced Configuration

```go
corsConfig := middleware.CORSConfig{
    AllowedOrigins: []string{"https://app.company.com"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders: []string{
        "Origin",
        "Content-Type",
        "Authorization",
        "X-Requested-With",
    },
    ExposedHeaders: []string{
        "Content-Length",
        "X-Request-ID",
    },
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}

router.Use(middleware.NewCORSMiddleware(corsConfig))
```

### Environment-Specific CORS

```go
func setupCORS(router *gin.Engine, cfg *config.Config) {
    if cfg.IsDevelopment() {
        // Permissive CORS for development
        router.Use(middleware.DevelopmentCORSMiddleware())
    } else {
        // Strict CORS for production
        router.Use(middleware.CustomCORSMiddleware(cfg.AllowedOrigins))
    }
}
```

### CORS Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `AllowedOrigins` | List of allowed origins | `["http://localhost:3000"]` |
| `AllowedMethods` | HTTP methods allowed | `["GET", "POST", "PUT", "DELETE", "OPTIONS"]` |
| `AllowedHeaders` | Headers allowed in requests | Standard headers + Authorization |
| `ExposedHeaders` | Headers exposed to client | `["Content-Length", "X-Request-ID"]` |
| `AllowCredentials` | Allow credentials (cookies) | `true` |
| `MaxAge` | Preflight cache duration | `12 hours` |

## Authentication Middleware

Flexible authentication middleware supporting multiple authentication methods.

### Basic Authentication

```go
// Bearer token authentication
tokenValidator := func(token string) (map[string]interface{}, error) {
    // Validate token (JWT, database lookup, etc.)
    if isValidToken(token) {
        return map[string]interface{}{
            "user_id": "123",
            "username": "john",
            "roles": []string{"user"},
        }, nil
    }
    return nil, middleware.ErrInvalidToken
}

router.Use(middleware.RequireAuth(tokenValidator))
```

### Custom Authentication Configuration

```go
authConfig := middleware.AuthConfig{
    SkipPaths: []string{"/health", "/public"},
    TokenExtractor: middleware.BearerTokenExtractor,
    TokenValidator: tokenValidator,
    ErrorHandler: func(c *gin.Context, err error) {
        responses.Unauthorized(c, "Authentication failed")
    },
}

router.Use(middleware.NewAuthMiddleware(authConfig))
```

### Token Extractors

```go
// Bearer token from Authorization header
middleware.BearerTokenExtractor(c)

// API key from X-API-Key header
middleware.APIKeyExtractor(c)

// Token from query parameter
tokenExtractor := middleware.QueryTokenExtractor("token")
```

### Authentication Methods

#### 1. JWT/Bearer Token Authentication

```go
protected := router.Group("/api")
protected.Use(middleware.RequireAuth(validateJWTToken))
{
    protected.GET("/profile", getProfile)
    protected.POST("/posts", createPost)
}
```

#### 2. API Key Authentication

```go
apiKeys := map[string]string{
    "key123": "user_1",
    "key456": "user_2",
}

router.Use(middleware.APIKeyAuth(apiKeys))
```

#### 3. Basic Authentication

```go
credentials := gin.Accounts{
    "admin": "secret",
    "user":  "password",
}

router.Use(middleware.BasicAuth(credentials))
```

#### 4. Optional Authentication

```go
// Don't fail if no token provided
router.Use(middleware.OptionalAuth(tokenValidator))

func handler(c *gin.Context) {
    userID, exists := middleware.GetUserID(c)
    if exists {
        // User is authenticated
    } else {
        // Anonymous user
    }
}
```

### Role-Based Access Control

```go
// Require specific role
adminRoutes := router.Group("/admin")
adminRoutes.Use(middleware.RequireAuth(tokenValidator))
adminRoutes.Use(middleware.RequireRole("admin"))
{
    adminRoutes.GET("/users", listAllUsers)
    adminRoutes.DELETE("/users/:id", deleteUser)
}

// Require any of multiple roles
moderatorRoutes := router.Group("/moderate")
moderatorRoutes.Use(middleware.RequireAuth(tokenValidator))
moderatorRoutes.Use(middleware.RequireAnyRole("admin", "moderator"))
{
    moderatorRoutes.POST("/approve/:id", approveContent)
}
```

### Working with Authentication Context

```go
func protectedHandler(c *gin.Context) {
    // Get user information from context
    userID, exists := middleware.GetUserID(c)
    if !exists {
        responses.Unauthorized(c, "User not found")
        return
    }

    // Get user roles
    roles, exists := middleware.GetUserRoles(c)
    if exists {
        fmt.Printf("User roles: %v\n", roles)
    }

    // Check specific role
    if middleware.HasRole(c, "admin") {
        // Admin-specific logic
    }

    responses.Success(c, "Protected resource", gin.H{
        "user_id": userID,
        "roles":   roles,
    })
}
```

### Authentication Error Handling

```go
authConfig := middleware.AuthConfig{
    ErrorHandler: func(c *gin.Context, err error) {
        switch err {
        case middleware.ErrMissingToken:
            responses.Unauthorized(c, "Authorization token required")
        case middleware.ErrInvalidToken:
            responses.Unauthorized(c, "Invalid token")
        case middleware.ErrExpiredToken:
            responses.Unauthorized(c, "Token expired")
        case middleware.ErrInsufficientPermissions:
            responses.Forbidden(c, "Insufficient permissions")
        default:
            responses.Unauthorized(c, "Authentication failed")
        }
    },
}
```

## Health Check Middleware

Comprehensive health monitoring for your service and its dependencies.

### Basic Health Checks

```go
// Simple health endpoint
router.Use(middleware.SimpleHealthMiddleware("my-service", "1.0.0"))
// Creates GET /health endpoint
```

### Advanced Health Configuration

```go
healthConfig := middleware.DefaultHealthConfig("my-service", "1.0.0")

// Add database health checker
healthConfig.AddHealthChecker("database", middleware.DatabaseHealthChecker(func() error {
    return db.Ping()
}))

// Add external service health checker
healthConfig.AddHealthChecker("redis", middleware.ExternalServiceHealthChecker(
    "redis", "http://redis:6379/ping", 5*time.Second,
))

// Add memory health checker (max 512MB)
healthConfig.AddHealthChecker("memory", middleware.MemoryHealthChecker(512))

// Add disk health checker (max 80% usage)
healthConfig.AddHealthChecker("disk", middleware.DiskHealthChecker("/", 80))

router.Use(middleware.HealthMiddleware(healthConfig))
```

### Health Check Endpoints

The health middleware automatically creates multiple endpoints:

| Endpoint | Purpose | Use Case |
|----------|---------|----------|
| `/health` | Basic health status | General monitoring |
| `/ready` | Readiness probe | Kubernetes readiness |
| `/live` | Liveness probe | Kubernetes liveness |

### Custom Health Checkers

```go
// Custom health checker for your service
customHealthChecker := func() middleware.HealthCheck {
    // Your custom health logic
    if isServiceHealthy() {
        return middleware.HealthCheck{
            Name:     "custom_service",
            Status:   middleware.HealthStatusHealthy,
            Duration: time.Since(start),
            Metadata: map[string]interface{}{
                "version": "1.0.0",
                "uptime":  getUptime(),
            },
        }
    }
    
    return middleware.HealthCheck{
        Name:     "custom_service",
        Status:   middleware.HealthStatusUnhealthy,
        Message:  "Service is not responding",
        Duration: time.Since(start),
    }
}

healthConfig.AddHealthChecker("custom", customHealthChecker)
```

### Health Check Response Format

```json
{
  "status": "healthy",
  "timestamp": "2023-12-25T15:30:45Z",
  "service": "my-service",
  "version": "1.0.0",
  "uptime": "2h30m15s",
  "checks": {
    "database": {
      "status": "healthy",
      "duration": "5ms",
      "metadata": {
        "response_time_ms": 5
      }
    },
    "memory": {
      "status": "healthy",
      "metadata": {
        "current_memory_mb": 45,
        "max_memory_mb": 512,
        "usage_percent": 8
      }
    }
  }
}
```

## Logging Middleware

Flexible request logging with multiple output formats.

### Basic Logging

```go
// Default logging (colored, readable)
router.Use(middleware.DefaultLoggingMiddleware())

// Detailed logging (includes user agent, referer)
router.Use(middleware.DetailedLoggingMiddleware())

// Production logging (JSON format)
router.Use(middleware.ProductionLoggingMiddleware())

// Silent logging (errors only)
router.Use(middleware.SilentLoggingMiddleware())
```

### Custom Logging Configuration

```go
loggingConfig := middleware.LoggingConfig{
    Level:     middleware.LogLevelInfo,
    SkipPaths: []string{"/health", "/metrics"},
    CustomFormat: func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("[%s] %s %s %d %s\n",
            param.TimeStamp.Format("15:04:05"),
            param.Method,
            param.Path,
            param.StatusCode,
            param.Latency,
        )
    },
}

router.Use(middleware.NewLoggingMiddleware(loggingConfig))
```

### Log Levels

```go
// Debug: Log all requests
router.Use(middleware.RequestLogger(middleware.LogLevelDebug))

// Info: Log successful requests and errors
router.Use(middleware.RequestLogger(middleware.LogLevelInfo))

// Warn: Log warnings and errors only
router.Use(middleware.RequestLogger(middleware.LogLevelWarn))

// Error: Log errors only
router.Use(middleware.RequestLogger(middleware.LogLevelError))
```

### Log Formats

#### Default Format (Development)
```
[GIN] 2023/12/25 - 15:04:05 | 200 |     2.547ms |       127.0.0.1 | GET     "/api/users"
```

#### Production Format (JSON)
```json
{"time":"2023-12-25T15:04:05Z","method":"GET","path":"/api/users","status":200,"latency":"2.547ms","ip":"127.0.0.1","user_agent":"curl/7.68.0","error":""}
```

#### Detailed Format
```
[2023/12/25 15:04:05] "GET /api/users HTTP/1.1" 200 2.547ms "curl/7.68.0" "" 127.0.0.1
```

## Recovery Middleware

Panic recovery with proper error handling and logging.

### Basic Recovery

```go
// Default recovery (logs stack trace)
router.Use(middleware.DefaultRecoveryMiddleware())

// Production recovery (no stack trace)
router.Use(middleware.ProductionRecoveryMiddleware())

// Development recovery (detailed errors)
router.Use(middleware.DevelopmentRecoveryMiddleware())

// Silent recovery (no logging)
router.Use(middleware.SilentRecoveryMiddleware())
```

### Custom Recovery

```go
recoveryConfig := middleware.RecoveryConfig{
    EnableStackTrace: true,
    SkipPaths:       []string{"/health"},
    CustomHandler: func(c *gin.Context, recovered interface{}) {
        requestID := middleware.MustGetRequestID(c)
        
        // Log the panic
        log.Printf("Panic recovered [%s]: %v", requestID, recovered)
        
        // Return error response
        responses.InternalError(c, "Something went wrong")
    },
}

router.Use(middleware.NewRecoveryMiddleware(recoveryConfig))
```

### Custom Panic Handlers

```go
// JSON error response handler
jsonHandler := middleware.JSONErrorPanicHandler(true) // includeDetails = true

// Custom panic handler
customHandler := func(c *gin.Context, recovered interface{}, requestID string) {
    // Log panic with context
    log.Printf("Panic in request %s: %v", requestID, recovered)
    
    // Send notification to monitoring service
    sendPanicAlert(requestID, recovered)
    
    // Return appropriate error response
    if strings.Contains(fmt.Sprintf("%v", recovered), "database") {
        responses.ServiceUnavailable(c, "Database temporarily unavailable")
    } else {
        responses.InternalError(c, "Internal server error")
    }
}

router.Use(middleware.WithCustomPanicHandler(customHandler))
```

## Request ID Middleware

Automatic request tracking for debugging and monitoring.

### Basic Request ID

```go
// Default request ID (generates random hex)
router.Use(middleware.DefaultRequestIDMiddleware())

// Short request ID (4 bytes hex)
router.Use(middleware.ShortRequestIDMiddleware())

// UUID-style request ID
router.Use(middleware.UUIDRequestIDMiddleware())
```

### Custom Request ID

```go
// Custom generator
customGenerator := func() string {
    return fmt.Sprintf("req_%d_%s", time.Now().Unix(), generateRandomString(8))
}

router.Use(middleware.CustomRequestIDMiddleware(customGenerator))
```

### Using Request IDs

```go
func handler(c *gin.Context) {
    // Get request ID
    requestID, exists := middleware.GetRequestID(c)
    if exists {
        log.Printf("Processing request %s", requestID)
    }
    
    // Or use the safe version
    requestID = middleware.MustGetRequestID(c)
    
    // Request ID is automatically added to response headers as X-Request-ID
    responses.Success(c, "Request processed", gin.H{
        "request_id": requestID,
    })
}
```

### Request ID Configuration

```go
requestIDConfig := middleware.RequestIDConfig{
    HeaderName: "X-Trace-ID",        // Custom header name
    ContextKey: "trace_id",          // Custom context key
    Generator:  customGenerator,     // Custom ID generator
}

router.Use(middleware.NewRequestIDMiddleware(requestIDConfig))
```

## Custom Middleware

Creating your own middleware following microservice-commons patterns.

### Basic Custom Middleware

```go
func TimingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // Process request
        c.Next()
        
        // Calculate and log timing
        duration := time.Since(start)
        log.Printf("Request %s took %v", c.Request.URL.Path, duration)
        
        // Add timing header
        c.Header("X-Response-Time", duration.String())
    }
}
```

### Advanced Custom Middleware

```go
func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
    // Rate limiting logic
    limiter := make(map[string][]time.Time)
    mutex := sync.RWMutex{}
    
    return func(c *gin.Context) {
        clientIP := c.ClientIP()
        now := time.Now()
        
        mutex.Lock()
        requests := limiter[clientIP]
        
        // Clean old requests
        var validRequests []time.Time
        for _, reqTime := range requests {
            if now.Sub(reqTime) < window {
                validRequests = append(validRequests, reqTime)
            }
        }
        
        // Check rate limit
        if len(validRequests) >= maxRequests {
            mutex.Unlock()
            responses.TooManyRequests(c, "Rate limit exceeded")
            c.Abort()
            return
        }
        
        // Add current request
        validRequests = append(validRequests, now)
        limiter[clientIP] = validRequests
        mutex.Unlock()
        
        c.Next()
    }
}
```

### Middleware with Configuration

```go
type CacheConfig struct {
    TTL        time.Duration
    SkipPaths  []string
    CacheStore CacheStore
}

func NewCacheMiddleware(config CacheConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if path should be cached
        for _, path := range config.SkipPaths {
            if c.Request.URL.Path == path {
                c.Next()
                return
            }
        }
        
        // Cache logic here
        cacheKey := generateCacheKey(c.Request)
        
        if cached := config.CacheStore.Get(cacheKey); cached != nil {
            c.Data(http.StatusOK, "application/json", cached)
            return
        }
        
        // Capture response
        w := &responseWriter{ResponseWriter: c.Writer}
        c.Writer = w
        
        c.Next()
        
        // Store in cache if successful
        if w.status == http.StatusOK {
            config.CacheStore.Set(cacheKey, w.body, config.TTL)
        }
    }
}
```

## Middleware Ordering

The order of middleware execution is crucial for proper functionality.

### Recommended Order

```go
func setupMiddleware(router *gin.Engine, cfg *config.Config) {
    // 1. Recovery (should be first to catch all panics)
    router.Use(middleware.DefaultRecoveryMiddleware())
    
    // 2. Request ID (for tracking)
    router.Use(middleware.DefaultRequestIDMiddleware())
    
    // 3. Logging (after request ID is set)
    router.Use(middleware.DefaultLoggingMiddleware())
    
    // 4. CORS (before authentication)
    router.Use(middleware.CustomCORSMiddleware(cfg.AllowedOrigins))
    
    // 5. Health checks (before authentication)
    router.Use(middleware.SimpleHealthMiddleware(cfg.ServiceName, cfg.ServiceVersion))
    
    // 6. Custom middleware (rate limiting, caching, etc.)
    router.Use(RateLimitMiddleware(100, time.Minute))
    
    // 7. Authentication (after all infrastructure middleware)
    router.Use(middleware.RequireAuth(tokenValidator))
}
```

### Why Order Matters

```go
// ❌ Bad: Logging before request ID
router.Use(middleware.DefaultLoggingMiddleware())  // Won't have request ID
router.Use(middleware.DefaultRequestIDMiddleware())

// ✅ Good: Request ID before logging
router.Use(middleware.DefaultRequestIDMiddleware())
router.Use(middleware.DefaultLoggingMiddleware())  // Will include request ID

// ❌ Bad: Authentication before CORS
router.Use(middleware.RequireAuth(validator))      // Browser preflight fails
router.Use(middleware.DefaultCORSMiddleware())

// ✅ Good: CORS before authentication
router.Use(middleware.DefaultCORSMiddleware())
router.Use(middleware.RequireAuth(validator))     // Preflight works
```

### Conditional Middleware

```go
func setupMiddleware(router *gin.Engine, cfg *config.Config) {
    // Always include recovery and logging
    router.Use(middleware.DefaultRecoveryMiddleware())
    router.Use(middleware.DefaultRequestIDMiddleware())
    
    // Environment-specific middleware
    if cfg.IsDevelopment() {
        router.Use(middleware.DetailedLoggingMiddleware())
        router.Use(middleware.DevelopmentCORSMiddleware())
    } else {
        router.Use(middleware.ProductionLoggingMiddleware())
        router.Use(middleware.CustomCORSMiddleware(cfg.AllowedOrigins))
    }
    
    // Optional middleware based on configuration
    if cfg.EnableRateLimit {
        router.Use(RateLimitMiddleware(cfg.RateLimit, time.Minute))
    }
    
    if cfg.EnableAuth {
        router.Use(middleware.RequireAuth(tokenValidator))
    }
}
```

## Best Practices

### 1. Error Handling in Middleware

```go
func ValidateAPIKeyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        
        if apiKey == "" {
            responses.BadRequest(c, "API key required")
            c.Abort() // Important: stop processing
            return
        }
        
        if !isValidAPIKey(apiKey) {
            responses.Unauthorized(c, "Invalid API key")
            c.Abort()
            return
        }
        
        c.Set("api_key", apiKey)
        c.Next() // Continue to next middleware/handler
    }
}
```

### 2. Context Values

```go
// ✅ Good: Use typed constants for context keys
const (
    UserIDKey    = "user_id"
    RequestIDKey = "request_id"
    TenantIDKey  = "tenant_id"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := extractUserID(c)
        c.Set(UserIDKey, userID)
        c.Next()
    }
}

func GetUserID(c *gin.Context) (string, bool) {
    userID, exists := c.Get(UserIDKey)
    if !exists {
        return "", false
    }
    return userID.(string), true
}
```

### 3. Middleware Testing

```go
func TestAuthMiddleware(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    router := gin.New()
    router.Use(AuthMiddleware())
    router.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // Test without token
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/test", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusUnauthorized, w.Code)
    
    // Test with valid token
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/test", nil)
    req.Header.Set("Authorization", "Bearer valid-token")
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### 4. Performance Considerations

```go
// ✅ Good: Efficient middleware
func EfficientMiddleware() gin.HandlerFunc {
    // Pre-compute expensive operations outside the handler
    compiledRegex := regexp.MustCompile(`^/api/v\d+`)
    
    return func(c *gin.Context) {
        // Fast path for common cases
        if c.Request.Method == "OPTIONS" {
            c.Next()
            return
        }
        
        // Use the pre-compiled regex
        if compiledRegex.MatchString(c.Request.URL.Path) {
            // API-specific logic
        }
        
        c.Next()
    }
}
```

### 5. Graceful Degradation

```go
func ExternalServiceMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Try external service with timeout
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()
        
        result, err := callExternalService(ctx)
        if err != nil {
            // Log error but don't fail the request
            log.Printf("External service unavailable: %v", err)
            c.Set("external_service_available", false)
        } else {
            c.Set("external_service_result", result)
            c.Set("external_service_available", true)
        }
        
        c.Next()
    }
}
```

### 6. Middleware Documentation

```go
// UserContextMiddleware extracts user information from JWT token
// and adds it to the Gin context for use by handlers.
//
// Context keys set:
//   - "user_id": string - User's unique identifier
//   - "username": string - User's username
//   - "roles": []string - User's roles
//
// Requirements:
//   - Valid JWT token in Authorization header
//   - Token must contain 'sub', 'preferred_username', and 'roles' claims
//
// Usage:
//   router.Use(UserContextMiddleware(jwtValidator))
func UserContextMiddleware(validator TokenValidator) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Implementation...
    }
}
```

### 7. Monitoring and Metrics

```go
func MetricsMiddleware() gin.HandlerFunc {
    requestCount := prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        status := fmt.Sprintf("%d", c.Writer.Status())
        
        requestCount.WithLabelValues(
            c.Request.Method,
            c.Request.URL.Path,
            status,
        ).Inc()
        
        // Record response time
        c.Header("X-Response-Time", duration.String())
    }
}
```

## Troubleshooting

### Common Issues

1. **CORS preflight failures**
   ```go
   // Make sure CORS middleware comes before authentication
   router.Use(middleware.DefaultCORSMiddleware())
   router.Use(middleware.RequireAuth(validator))
   ```

2. **Missing request IDs in logs**
   ```go
   // Request ID middleware must come before logging
   router.Use(middleware.DefaultRequestIDMiddleware())
   router.Use(middleware.DefaultLoggingMiddleware())
   ```

3. **Authentication not working**
   ```go
   // Check middleware order and error handling
   router.Use(middleware.RequireAuth(func(token string) (map[string]interface{}, error) {
       // Make sure to return proper errors
       if token == "" {
           return nil, middleware.ErrMissingToken
       }
       // Validate token...
   }))
   ```

4. **Health checks failing**
   ```go
   // Make sure health checker functions don't panic
   healthConfig.AddHealthChecker("database", func() middleware.HealthCheck {
       defer func() {
           if r := recover(); r != nil {
               log.Printf("Health check panic: %v", r)
           }
       }()
       // Health check logic...
   })
   ```

### Debugging Middleware

```go
func DebugMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Printf("Before: %s %s", c.Request.Method, c.Request.URL.Path)
        
        c.Next()
        
        log.Printf("After: %s %s - Status: %d", 
            c.Request.Method, c.Request.URL.Path, c.Writer.Status())
    }
}
```

This comprehensive middleware documentation covers all aspects of using and creating middleware with microservice-commons. Each middleware component is designed to work together seamlessly while providing maximum flexibility for customization.