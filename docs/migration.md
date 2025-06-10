# Migration Guide

This guide helps you migrate existing microservices to use microservice-commons, reducing boilerplate code and ensuring consistency across your services.

## Table of Contents

1. [Before You Start](#before-you-start)
2. [Assessment Checklist](#assessment-checklist)
3. [Step-by-Step Migration](#step-by-step-migration)
4. [Common Migration Scenarios](#common-migration-scenarios)
5. [Breaking Changes](#breaking-changes)
6. [Testing Your Migration](#testing-your-migration)
7. [Rollback Strategy](#rollback-strategy)
8. [Performance Considerations](#performance-considerations)

## Before You Start

### Prerequisites

- Go 1.19 or later
- Existing service using Gin framework
- PostgreSQL database (if using database features)
- Basic understanding of middleware concepts

### What You'll Gain

✅ **50+ lines of boilerplate code eliminated**  
✅ **Standardized error responses across services**  
✅ **Built-in health checks and monitoring**  
✅ **Consistent CORS and middleware setup**  
✅ **Automatic graceful shutdown**  
✅ **Request ID tracking**  
✅ **Database connection management with retry logic**

## Assessment Checklist

Before migrating, assess your current service:

### Current Architecture Assessment

```bash
# Check your current structure
├── main.go                    # Server setup
├── handlers/                  # Route handlers
├── models/                    # Database models
├── middleware/               # Custom middleware
├── config/                   # Configuration
└── utils/                    # Utility functions
```

### Compatibility Check

- [ ] Using Gin framework
- [ ] PostgreSQL database (optional)
- [ ] Environment-based configuration
- [ ] Standard HTTP status codes
- [ ] JSON API responses

### Current Code Patterns

Identify these patterns in your existing code:

```go
// Server setup boilerplate
router := gin.Default()
router.Use(cors.New(corsConfig))
router.GET("/health", healthHandler)

// Manual graceful shutdown
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

// Manual error responses
c.JSON(400, gin.H{"error": "Bad request"})

// Custom pagination logic
page := c.DefaultQuery("page", "1")
pageSize := c.DefaultQuery("page_size", "10")
```

## Step-by-Step Migration

### Phase 1: Add Dependency

1. **Add microservice-commons to your project:**

```bash
go get github.com/JorgeSaicoski/microservice-commons
```

2. **Update your go.mod:**

```go
module your-service

go 1.19

require (
    github.com/JorgeSaicoski/microservice-commons v0.1.0
    github.com/gin-gonic/gin v1.9.1
    // ... other dependencies
)
```

### Phase 2: Environment Configuration

1. **Create .env.example based on your current config:**

```bash
# Before: Custom environment handling
API_KEY=your-api-key
DB_HOST=localhost
CORS_ORIGINS=http://localhost:3000

# After: Standardized environment variables
SERVICE_NAME=your-service
SERVICE_VERSION=1.0.0
PORT=8000
ENVIRONMENT=dev

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your-password
POSTGRES_DB=your-database

# CORS
ALLOWED_ORIGINS=http://localhost:3000
```

2. **Replace custom config loading:**

```go
// Before: Custom config
type Config struct {
    APIKey      string
    DBHost      string
    CORSOrigins string
}

func loadConfig() *Config {
    return &Config{
        APIKey:      os.Getenv("API_KEY"),
        DBHost:      os.Getenv("DB_HOST"),
        CORSOrigins: os.Getenv("CORS_ORIGINS"),
    }
}

// After: Use microservice-commons config
import "github.com/JorgeSaicoski/microservice-commons/config"

func main() {
    cfg := config.LoadFromEnv()
    if err := cfg.Validate(); err != nil {
        log.Fatal("Configuration error:", err)
    }
}
```

### Phase 3: Server Setup Migration

1. **Replace manual server setup:**

```go
// Before: Manual server setup (40+ lines)
func main() {
    router := gin.Default()
    
    // CORS setup
    allowedOrigins := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
    router.Use(cors.New(cors.Config{
        AllowOrigins:     allowedOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })
    
    // Routes
    setupRoutes(router)
    
    // Manual graceful shutdown
    srv := &http.Server{
        Addr:    ":" + os.Getenv("PORT"),
        Handler: router,
    }
    
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
}

// After: microservice-commons (10 lines)
import "github.com/JorgeSaicoski/microservice-commons/server"

func main() {
    server := server.NewServer(server.ServerOptions{
        ServiceName:    "your-service",
        ServiceVersion: "1.0.0",
        SetupRoutes:    setupRoutes,
    })
    server.Start() // Includes graceful shutdown automatically
}
```

### Phase 4: Response Migration

1. **Replace manual response handling:**

```go
// Before: Manual responses
func getUsers(c *gin.Context) {
    users := getUsersFromDB()
    c.JSON(200, gin.H{
        "status": "success",
        "data":   users,
    })
}

func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{
            "status": "error",
            "error":  "Invalid input",
        })
        return
    }
    
    // Create user logic...
    
    c.JSON(201, gin.H{
        "status": "success",
        "data":   user,
    })
}

// After: Standardized responses
import "github.com/JorgeSaicoski/microservice-commons/responses"

func getUsers(c *gin.Context) {
    users := getUsersFromDB()
    responses.Success(c, "Users retrieved successfully", users)
}

func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        responses.BadRequest(c, "Invalid user data")
        return
    }
    
    // Create user logic...
    
    responses.Created(c, "User created successfully", user)
}
```

### Phase 5: Database Migration

1. **Replace manual database connection:**

```go
// Before: Manual database setup
func connectDB() *gorm.DB {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    return db
}

// After: microservice-commons database
import "github.com/JorgeSaicoski/microservice-commons/database"

func setupRoutes(router *gin.Engine, cfg *config.Config) {
    db, err := database.ConnectWithConfig(cfg.DatabaseConfig)
    if err != nil {
        panic("Failed to connect to database: " + err.Error())
    }
    
    // Auto-migrate models
    if err := database.QuickMigrate(db, &User{}, &Task{}); err != nil {
        panic("Failed to migrate: " + err.Error())
    }
    
    // Setup routes with db...
}
```

### Phase 6: Middleware Migration

1. **Replace custom middleware with standardized versions:**

```go
// Before: Custom middleware
func customCORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
        c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}

func requestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := generateRequestID()
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}

// After: Use built-in middleware
import "github.com/JorgeSaicoski/microservice-commons/middleware"

// In your server setup - automatically included
// But if you want custom configuration:
func setupCustomMiddleware(router *gin.Engine, cfg *config.Config) {
    router.Use(middleware.DefaultRequestIDMiddleware())
    router.Use(middleware.CustomCORSMiddleware(cfg.AllowedOrigins))
}
```

## Common Migration Scenarios

### Scenario 1: Existing go-todo-list Service

**Current Structure:**
```go
// cmd/server/main.go
func main() {
    db.ConnectDatabase()
    
    config := api.DefaultRouterConfig()
    if origins := getEnv("ALLOWED_ORIGINS", ""); origins != "" {
        config.AllowedOrigins = origins
    }
    
    taskRouter := api.NewTaskRouter(db.DB, config)
    taskRouter.RegisterRoutes()
    
    port := getEnv("PORT", "8000")
    taskRouter.Run(":" + port)
}
```

**After Migration:**
```go
// main.go
func main() {
    server := server.NewServer(server.ServerOptions{
        ServiceName:    "go-todo-list",
        ServiceVersion: "1.0.0",
        SetupRoutes:    setupRoutes,
    })
    server.Start()
}

func setupRoutes(router *gin.Engine, cfg *config.Config) {
    db := database.MustConnect(cfg.DatabaseConfig)
    
    api := router.Group("/api")
    {
        api.GET("/tasks", handlers.GetTasks)
        api.POST("/tasks", handlers.CreateTask)
        api.GET("/tasks/:id", handlers.GetTask)
        api.PUT("/tasks/:id", handlers.UpdateTask)
        api.DELETE("/tasks/:id", handlers.DeleteTask)
    }
}
```

### Scenario 2: Service with Custom Authentication

**Before:**
```go
func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Token required"})
            c.Abort()
            return
        }
        
        // Validate token...
        if !isValidToken(token) {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", getUserIDFromToken(token))
        c.Next()
    }
}
```

**After:**
```go
import "github.com/JorgeSaicoski/microservice-commons/middleware"

func setupAuth(router *gin.Engine) {
    tokenValidator := func(token string) (map[string]interface{}, error) {
        if !isValidToken(token) {
            return nil, middleware.ErrInvalidToken
        }
        
        return map[string]interface{}{
            "user_id": getUserIDFromToken(token),
        }, nil
    }
    
    protected := router.Group("/api")
    protected.Use(middleware.RequireAuth(tokenValidator))
    {
        protected.GET("/profile", getProfile)
        protected.POST("/tasks", createTask)
    }
}
```

### Scenario 3: Service with Pagination

**Before:**
```go
func getTasks(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 10
    }
    
    offset := (page - 1) * pageSize
    
    var tasks []Task
    var total int64
    
    db.Model(&Task{}).Count(&total)
    db.Offset(offset).Limit(pageSize).Find(&tasks)
    
    totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
    
    c.JSON(200, gin.H{
        "data":        tasks,
        "page":        page,
        "page_size":   pageSize,
        "total":       total,
        "total_pages": totalPages,
        "has_next":    page < totalPages,
        "has_prev":    page > 1,
    })
}
```

**After:**
```go
import "github.com/JorgeSaicoski/microservice-commons/responses"

func getTasks(c *gin.Context) {
    params := responses.GetPaginationParams(c)
    
    var tasks []Task
    var total int64
    
    db.Model(&Task{}).Count(&total)
    db.Offset(params.Offset).Limit(params.Limit).Find(&tasks)
    
    responses.Paginated(c, tasks, total, params.Page, params.PageSize)
}
```

## Breaking Changes

### Response Format Changes

**Before:**
```json
{
  "status": "success",
  "data": {...}
}
```

**After:**
```json
{
  "message": "Operation successful",
  "data": {...},
  "timestamp": "2023-12-25T15:30:45Z"
}
```

**Migration Strategy:**
```go
// Option 1: Update clients to handle new format
// Option 2: Create compatibility wrapper
func compatibilityWrapper(c *gin.Context, message string, data interface{}) {
    if c.GetHeader("API-Version") == "v1" {
        // Old format for backward compatibility
        c.JSON(200, gin.H{
            "status": "success",
            "data":   data,
        })
    } else {
        // New standardized format
        responses.Success(c, message, data)
    }
}
```

### Error Response Changes

**Before:**
```json
{
  "status": "error",
  "error": "Not found"
}
```

**After:**
```json
{
  "error": "Not found",
  "code": "not_found",
  "timestamp": "2023-12-25T15:30:45Z",
  "path": "/api/users/123"
}
```

### Environment Variable Changes

**Before:**
```bash
DB_HOST=localhost
CORS_ORIGINS=http://localhost:3000
```

**After:**
```bash
POSTGRES_HOST=localhost
ALLOWED_ORIGINS=http://localhost:3000
```

**Migration Script:**
```bash
#!/bin/bash
# migrate-env.sh

# Backup current .env
cp .env .env.backup

# Replace old variable names
sed -i 's/DB_HOST/POSTGRES_HOST/g' .env
sed -i 's/DB_PORT/POSTGRES_PORT/g' .env
sed -i 's/DB_USER/POSTGRES_USER/g' .env
sed -i 's/DB_PASSWORD/POSTGRES_PASSWORD/g' .env
sed -i 's/DB_NAME/POSTGRES_DB/g' .env
sed -i 's/CORS_ORIGINS/ALLOWED_ORIGINS/g' .env

# Add new required variables
echo "SERVICE_NAME=your-service" >> .env
echo "SERVICE_VERSION=1.0.0" >> .env
echo "ENVIRONMENT=dev" >> .env
```

## Testing Your Migration

### 1. Unit Tests Migration

```go
// Before: Manual test setup
func TestGetUsers(t *testing.T) {
    router := gin.New()
    router.GET("/users", getUsersHandler)
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/users", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "success", response["status"])
}

// After: Test with microservice-commons
func TestGetUsers(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    // Use server setup for consistent testing
    server := server.NewServer(server.ServerOptions{
        ServiceName:    "test-service",
        ServiceVersion: "1.0.0",
        SetupRoutes:    setupTestRoutes,
        Config:         getTestConfig(),
    })
    
    router := server.GetRouter()
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/users", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
    
    var response responses.SuccessResponse
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Contains(t, response.Message, "retrieved successfully")
    assert.NotNil(t, response.Data)
}
```

### 2. Integration Tests

```go
func TestFullServiceIntegration(t *testing.T) {
    // Setup test database
    testDB := setupTestDatabase()
    defer teardownTestDatabase(testDB)
    
    // Create test server
    server := server.NewServer(server.ServerOptions{
        ServiceName:    "integration-test",
        ServiceVersion: "1.0.0",
        SetupRoutes:    setupRoutes,
        Config:         getTestConfig(),
    })
    
    // Test health endpoint
    router := server.GetRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/health", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "healthy")
}
```

### 3. Migration Validation Script

```bash
#!/bin/bash
# validate-migration.sh

echo "🔍 Validating migration..."

# Check if service starts
echo "Starting service..."
timeout 10s go run main.go &
SERVICE_PID=$!

sleep 3

# Test health endpoint
if curl -f http://localhost:8000/health > /dev/null 2>&1; then
    echo "✅ Health endpoint working"
else
    echo "❌ Health endpoint failed"
    kill $SERVICE_PID
    exit 1
fi

# Test main API endpoint
if curl -f http://localhost:8000/api/users > /dev/null 2>&1; then
    echo "✅ Main API endpoint working"
else
    echo "❌ Main API endpoint failed"
    kill $SERVICE_PID
    exit 1
fi

# Check response format
RESPONSE=$(curl -s http://localhost:8000/api/users)
if echo "$RESPONSE" | jq -e '.timestamp' > /dev/null 2>&1; then
    echo "✅ New response format detected"
else
    echo "⚠️  Old response format still in use"
fi

kill $SERVICE_PID
echo "🎉 Migration validation complete"
```

## Rollback Strategy

### 1. Prepare Rollback Branch

```bash
# Before starting migration
git checkout -b migration-to-commons
git checkout -b rollback-branch

# Work on migration
git checkout migration-to-commons
# ... migration work ...

# If rollback needed
git checkout rollback-branch
```

### 2. Environment Variable Rollback

```bash
#!/bin/bash
# rollback-env.sh

# Restore original environment variables
sed -i 's/POSTGRES_HOST/DB_HOST/g' .env
sed -i 's/POSTGRES_PORT/DB_PORT/g' .env
sed -i 's/POSTGRES_USER/DB_USER/g' .env
sed -i 's/POSTGRES_PASSWORD/DB_PASSWORD/g' .env
sed -i 's/POSTGRES_DB/DB_NAME/g' .env
sed -i 's/ALLOWED_ORIGINS/CORS_ORIGINS/g' .env

# Remove new variables
sed -i '/SERVICE_NAME/d' .env
sed -i '/SERVICE_VERSION/d' .env
sed -i '/ENVIRONMENT/d' .env
```

### 3. Code Rollback

Keep your original files during migration:

```bash
# Backup original files
mkdir migration-backup
cp main.go migration-backup/
cp -r handlers/ migration-backup/
cp -r middleware/ migration-backup/

# If rollback needed
cp migration-backup/* .
```

### 4. Dependency Rollback

```bash
# Remove microservice-commons
go mod edit -droprequire github.com/JorgeSaicoski/microservice-commons
go mod tidy

# Restore original dependencies
git checkout HEAD~1 go.mod go.sum
go mod download
```

## Performance Considerations

### 1. Memory Usage

**Before Migration:**
```go
// Custom middleware creates new instances
func customMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        data := make(map[string]interface{}) // New allocation per request
        // ...
    }
}
```

**After Migration:**
```go
// microservice-commons reuses objects
// Built-in middleware is optimized for performance
router.Use(middleware.DefaultRequestIDMiddleware())
```

### 2. Response Time Impact

Measure response times before and after:

```bash
# Before migration
ab -n 1000 -c 10 http://localhost:8000/api/users

# After migration
ab -n 1000 -c 10 http://localhost:8000/api/users
```

Expected changes:
- **Health endpoints**: Slightly slower due to additional checks
- **Regular endpoints**: Similar or faster due to optimized middleware
- **Error responses**: Faster due to pre-built responses

### 3. Database Connection Improvements

**Before:**
```go
// Single connection attempt
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
    log.Fatal(err)
}
```

**After:**
```go
// Automatic retry with exponential backoff
db, err := database.ConnectWithConfig(cfg.DatabaseConfig)
// Includes connection pooling and health monitoring
```

Performance improvements:
- **Connection reliability**: Retry logic handles temporary failures
- **Connection pooling**: Better resource management
- **Health monitoring**: Proactive issue detection

## Deployment Considerations

### 1. Blue-Green Deployment

```yaml
# docker-compose.yml for blue-green deployment
version: '3.8'
services:
  app-blue:
    build: .
    environment:
      - SERVICE_NAME=my-service
      - SERVICE_VERSION=1.0.0
    ports:
      - "8000:8000"
    
  app-green:
    build: .
    environment:
      - SERVICE_NAME=my-service
      - SERVICE_VERSION=1.1.0
    ports:
      - "8001:8000"
```

### 2. Kubernetes Migration

```yaml
# k8s-migration.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-service-v2
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-service
      version: v2
  template:
    metadata:
      labels:
        app: my-service
        version: v2
    spec:
      containers:
      - name: my-service
        image: my-service:v2
        env:
        - name: SERVICE_NAME
          value: "my-service"
        - name: SERVICE_VERSION
          value: "2.0.0"
        ports:
        - containerPort: 8000
        livenessProbe:
          httpGet:
            path: /live
            port: 8000
        readinessProbe:
          httpGet:
            path: /ready
            port: 8000
```

### 3. Monitoring During Migration

```bash
# Monitor service metrics during migration
watch -n 1 'curl -s http://localhost:8000/health | jq'

# Monitor response times
while true; do
  curl -w "@curl-format.txt" -s -o /dev/null http://localhost:8000/api/users
  sleep 1
done

# curl-format.txt
time_namelookup:  %{time_namelookup}\n
time_connect:     %{time_connect}\n
time_appconnect:  %{time_appconnect}\n
time_pretransfer: %{time_pretransfer}\n
time_redirect:    %{time_redirect}\n
time_starttransfer: %{time_starttransfer}\n
time_total:       %{time_total}\n
```

## Migration Checklist

### Pre-Migration

- [ ] Current service is working correctly
- [ ] All tests are passing
- [ ] Environment variables documented
- [ ] API responses documented
- [ ] Performance baseline established
- [ ] Rollback plan prepared

### During Migration

- [ ] Dependencies added successfully
- [ ] Environment variables migrated
- [ ] Server setup migrated
- [ ] Responses standardized
- [ ] Database connection migrated
- [ ] Middleware migrated
- [ ] Tests updated

### Post-Migration

- [ ] All tests passing
- [ ] Health endpoints working
- [ ] API responses in new format
- [ ] Performance acceptable
- [ ] Error handling working
- [ ] Graceful shutdown working
- [ ] Documentation updated

### Validation

- [ ] Service starts correctly
- [ ] All endpoints respond
- [ ] Database connectivity working
- [ ] Error responses standardized
- [ ] Request ID tracking working
- [ ] CORS configuration correct
- [ ] Authentication working (if applicable)

## Common Issues and Solutions

### Issue 1: Import Path Conflicts

**Problem:**
```go
// Conflicts with existing config package
import "your-service/config"
import "github.com/JorgeSaicoski/microservice-commons/config"
```

**Solution:**
```go
// Use alias for microservice-commons
import (
    localConfig "your-service/config"
    mcConfig "github.com/JorgeSaicoski/microservice-commons/config"
)
```

### Issue 2: Response Format Breaking Clients

**Problem:** Existing clients expect old response format

**Solution:**
```go
// Create adapter for backward compatibility
func adaptOldResponse(c *gin.Context, message string, data interface{}) {
    apiVersion := c.GetHeader("API-Version")
    
    if apiVersion == "v1" {
        // Old format
        c.JSON(200, gin.H{
            "status": "success",
            "data":   data,
        })
    } else {
        // New format
        responses.Success(c, message, data)
    }
}
```

### Issue 3: Database Migration Issues

**Problem:** Existing database schema conflicts

**Solution:**
```go
// Skip auto-migration for existing tables
if !database.TableExists(db, "users") {
    database.QuickMigrate(db, &User{})
} else {
    log.Println("Users table exists, skipping migration")
}
```

### Issue 4: Environment Variable Confusion

**Problem:** Mixed old and new environment variables

**Solution:**
```go
// Create compatibility layer
func getCompatibleEnv() *config.Config {
    cfg := config.LoadFromEnv()
    
    // Fallback to old variable names
    if cfg.DatabaseConfig.Host == "localhost" {
        if oldHost := os.Getenv("DB_HOST"); oldHost != "" {
            cfg.DatabaseConfig.Host = oldHost
        }
    }
    
    return cfg
}
```

This comprehensive migration guide should help you successfully transition your existing microservices to use microservice-commons while minimizing risks and ensuring a smooth transition.