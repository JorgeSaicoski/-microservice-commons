# microservice-commons

A comprehensive Go library providing common utilities and patterns for microservices architecture. Designed to eliminate boilerplate code and ensure consistency across services in the Personal Manager ecosystem.

## 🚀 Quick Start

```bash
go get github.com/JorgeSaicoski/microservice-commons
```

**Before microservice-commons (50+ lines of boilerplate):**
```go
// Your old main.go
func main() {
    router := gin.Default()
    
    // Setup CORS manually
    allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
    origins := strings.Split(allowedOrigins, ",")
    router.Use(cors.New(cors.Config{
        AllowOrigins: origins,
        // ... 10 more lines of CORS config
    }))
    
    // Setup health check manually
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })
    
    // Manual server setup with graceful shutdown
    // ... 20+ more lines
}
```

**After microservice-commons (10 lines):**
```go
import "github.com/JorgeSaicoski/microservice-commons/server"

func main() {
    server := server.NewServer(server.ServerOptions{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        SetupRoutes:    setupRoutes,
    })
    server.Start() // Includes graceful shutdown automatically
}
```

## 📁 Project Structure

```
microservice-commons/
├── config/                 # Configuration management
│   ├── config.go          # Main configuration struct and loader
│   ├── database.go        # Database-specific configuration
│   ├── keycloak.go        # Keycloak authentication configuration
│   └── validation.go      # Configuration validation utilities
├── 
├── middleware/             # Gin middleware components
│   ├── cors.go           # CORS middleware setup
│   ├── logging.go        # Logging middleware and utilities
│   ├── health.go         # Health check middleware
│   ├── auth.go           # Authentication middleware helpers
│   ├── recovery.go       # Recovery middleware customizations
│   └── request_id.go     # Request ID middleware
├── 
├── server/                # Server setup and lifecycle
│   ├── server.go         # Main server struct and setup
│   ├── graceful.go       # Graceful shutdown utilities
│   └── options.go        # Server configuration options
├── 
├── responses/             # Standardized API responses
│   ├── responses.go      # Standard response helpers
│   ├── errors.go         # Error response types and helpers
│   └── pagination.go     # Pagination response helpers
├── 
├── database/              # Database utilities
│   ├── connection.go     # Database connection with retry logic
│   ├── migration.go      # Migration utilities and helpers
│   └── health.go         # Database health check utilities
├── 
├── utils/                 # Generic utility functions
│   ├── env.go            # Environment variable utilities
│   ├── strings.go        # String manipulation utilities
│   ├── time.go           # Time utilities
│   └── validation.go     # Common validation functions
├── 
├── types/                 # Shared types and structs
│   ├── common.go         # Common types across services
│   ├── pagination.go     # Pagination types
│   └── responses.go      # Response types
├── 
├── examples/              # Working examples
│   ├── basic-service/    # Simple service setup example
│   └── advanced-service/ # Advanced features example
├── 
├── docs/                  # Documentation
│   ├── configuration.md  # Configuration guide
│   ├── middleware.md     # Middleware documentation
│   ├── responses.md      # Response patterns guide
│   └── migration.md      # Migration from existing services
└── 
└── tests/                 # Test coverage
    ├── config_test.go
    ├── middleware_test.go
    ├── server_test.go
    └── responses_test.go
```

## 🎯 Key Features

### ✅ Environment-Based Configuration
- Automatic environment variable loading
- Validation with helpful error messages
- Support for dev/staging/prod environments
- Database, Keycloak, and service configuration

### ✅ Standardized Server Setup
- Automatic CORS configuration
- Built-in health check endpoints
- Graceful shutdown handling
- Request logging and recovery

### ✅ Database Management
- Connection pooling with retry logic
- Health monitoring and statistics
- Migration utilities with safety checks
- Integration with existing pgconnect library

### ✅ Consistent API Responses
- Standardized success/error formats
- Pagination support
- HTTP status code helpers
- JSON response utilities

### ✅ Middleware Collection
- Authentication helpers
- Request ID tracking
- Custom logging formats
- Recovery with error reporting

## 📖 Usage Examples

### Basic Service Setup

```go
package main

import (
    "github.com/JorgeSaicoski/microservice-commons/server"
    "github.com/JorgeSaicoski/microservice-commons/config"
    "github.com/JorgeSaicoski/microservice-commons/database"
    "github.com/gin-gonic/gin"
)

func main() {
    server := server.NewServer(server.ServerOptions{
        ServiceName:    "task-service",
        ServiceVersion: "1.0.0",
        SetupRoutes:    setupRoutes,
    })
    server.Start()
}

func setupRoutes(router *gin.Engine, cfg *config.Config) {
    // Connect to database
    db, err := database.ConnectWithConfig(cfg.DatabaseConfig)
    if err != nil {
        panic(err)
    }
    
    // Setup your routes
    api := router.Group("/api")
    {
        api.GET("/tasks", getTasksHandler)
        api.POST("/tasks", createTaskHandler)
    }
}
```

### Database Connection with Health Monitoring

```go
import (
    "github.com/JorgeSaicoski/microservice-commons/database"
    "github.com/JorgeSaicoski/microservice-commons/config"
)

func setupDatabase() {
    cfg := config.LoadFromEnv()
    
    // Connect with automatic retry
    db, err := database.ConnectWithConfig(cfg.DatabaseConfig)
    if err != nil {
        panic(err)
    }
    
    // Check health
    if !database.QuickHealthCheck(db) {
        panic("Database health check failed")
    }
    
    // Run migrations
    if err := database.QuickMigrate(db, &Task{}, &User{}); err != nil {
        panic(err)
    }
}
```

### Standardized API Responses

```go
import (
    "github.com/JorgeSaicoski/microservice-commons/responses"
    "github.com/gin-gonic/gin"
)

func getTasksHandler(c *gin.Context) {
    tasks := []Task{} // Your task data
    
    // Paginated response
    responses.Paginated(c, tasks, 100, 1, 10)
}

func createTaskHandler(c *gin.Context) {
    task := Task{} // Created task
    
    // Success response
    responses.Success(c, "Task created successfully", task)
}

func errorHandler(c *gin.Context) {
    // Standardized error responses
    responses.BadRequest(c, "Invalid task data")
    responses.NotFound(c, "Task not found")
    responses.InternalError(c, "Database connection failed")
}
```

### Configuration Management

```go
import "github.com/JorgeSaicoski/microservice-commons/config"

func main() {
    // Load configuration
    cfg := config.LoadFromEnv()
    
    // Validate configuration
    if err := cfg.Validate(); err != nil {
        panic(err)
    }
    
    // Environment-specific behavior
    if cfg.IsDevelopment() {
        fmt.Println("Running in development mode")
    }
    
    // Access configuration
    fmt.Printf("Service: %s v%s\n", cfg.ServiceName, cfg.ServiceVersion)
    fmt.Printf("Database: %s:%s\n", cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Port)
}
```

## 🔧 Environment Variables

### Core Configuration
```bash
# Service Configuration
SERVICE_NAME=my-service
SERVICE_VERSION=1.0.0
PORT=8000
ENVIRONMENT=dev                    # dev, staging, prod
LOG_LEVEL=info                     # debug, info, warn, error
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080

# Database Configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=yourpassword
POSTGRES_DB=myservice_db
POSTGRES_SSLMODE=disable
POSTGRES_TIMEZONE=UTC
POSTGRES_MAX_IDLE_CONNS=10
POSTGRES_MAX_OPEN_CONNS=100
POSTGRES_LOG_LEVEL=silent         # silent, error, warn, info

# Keycloak Configuration
KEYCLOAK_URL=http://localhost:8080/keycloak
KEYCLOAK_REALM=master
KEYCLOAK_PUBLIC_KEY=               # Base64 encoded public key (optional)
KEYCLOAK_REQUIRED_CLAIMS=sub,preferred_username
KEYCLOAK_SKIP_PATHS=/health,/metrics
KEYCLOAK_KEY_REFRESH_INTERVAL=1h
KEYCLOAK_HTTP_TIMEOUT=10s

# Database Connection Retry
DB_MAX_RETRIES=3
DB_RETRY_DELAY_SECONDS=30
DB_LOG_LEVEL=silent
```

### Environment-Specific Examples

**Development (.env.dev):**
```bash
ENVIRONMENT=dev
LOG_LEVEL=debug
POSTGRES_HOST=localhost
ALLOWED_ORIGINS=http://localhost:3000
```

**Production (.env.prod):**
```bash
ENVIRONMENT=prod
LOG_LEVEL=warn
POSTGRES_HOST=prod-db.company.com
POSTGRES_SSLMODE=require
ALLOWED_ORIGINS=https://app.company.com
```

## 🗄️ Database Features

### Connection Management
- **Automatic retry logic** - Handles database startup delays
- **Connection pooling** - Configurable idle/max connections
- **Health monitoring** - Real-time connection statistics
- **Graceful handling** - Proper connection cleanup

### Connection Pool Stats
Each service maintains its own connection pool:

```json
{
  "connections": {
    "open": 8,        // Total connections established
    "in_use": 3,      // Currently executing queries  
    "idle": 5,        // Ready for reuse
    "max_open": 50    // Configuration limit
  }
}
```

**Understanding the numbers:**
- **Open**: Total TCP connections to database
- **In Use**: Connections actively processing queries
- **Idle**: Connections waiting in pool for reuse
- **Max Open**: Hard limit to prevent overwhelming database

**Multiple users, shared connections:**
- 100 users can share 10 connections
- Connections are borrowed per query, not per user
- Users only wait if all connections are busy with slow queries

### Migration Support
```go
// Safe migration with options
migrator := database.NewMigrator(db, database.MigrationOptions{
    DropTables:    false,  // NEVER true in production
    CreateIndexes: true,   // Add performance indexes
    Verbose:       true,   // Log migration details
})

err := migrator.AddModels(&Task{}, &User{}).Migrate()
```

### Health Monitoring
```go
// Quick health check
healthy := database.QuickHealthCheck(db)

// Detailed health information
status, details := database.DetailedHealthCheck(db)
```

## 🔐 Authentication Integration

### Keycloak Setup
Works with your existing keycloak-auth library:

```go
import (
    keycloakauth "github.com/JorgeSaicoski/keycloak-auth"
    "github.com/JorgeSaicoski/microservice-commons/config"
)

func setupAuth(router *gin.Engine, cfg *config.Config) {
    // Convert to keycloak-auth config
    authConfig := keycloakauth.Config{
        KeycloakURL:        cfg.KeycloakConfig.URL,
        Realm:              cfg.KeycloakConfig.Realm,
        PublicKeyBase64:    cfg.KeycloakConfig.PublicKeyBase64,
        RequiredClaims:     cfg.KeycloakConfig.RequiredClaims,
        SkipPaths:          cfg.KeycloakConfig.SkipPaths,
        KeyRefreshInterval: cfg.KeycloakConfig.KeyRefreshInterval,
        HTTPTimeout:        cfg.KeycloakConfig.HTTPTimeout,
    }
    
    router.Use(keycloakauth.SimpleAuthMiddleware(authConfig))
}
```

## 📊 API Response Standards

### Success Responses
```json
{
  "message": "Operation completed successfully",
  "data": { /* your data */ }
}
```

### Error Responses
```json
{
  "error": "Validation failed",
  "code": "bad_request",
  "details": "Title field is required"
}
```

### Paginated Responses
```json
{
  "data": [/* your items */],
  "total": 150,
  "page": 1,
  "pageSize": 10,
  "totalPages": 15
}
```

### Standard Error Codes
- `bad_request` - Invalid input data
- `unauthorized` - Authentication required
- `forbidden` - Insufficient permissions
- `not_found` - Resource doesn't exist
- `internal_error` - Server-side error

## 🔄 Migration from Existing Services

### Before (Current go-todo-list/cmd/server/main.go)
```go
func main() {
    // Connect to the database
    db.ConnectDatabase()

    // Get router config from environment
    config := api.DefaultRouterConfig()
    if origins := getEnv("ALLOWED_ORIGINS", ""); origins != "" {
        config.AllowedOrigins = origins
    }

    // Create router with full configuration
    taskRouter := api.NewTaskRouter(db.DB, config)
    taskRouter.RegisterRoutes()

    // Start the server
    port := getEnv("PORT", "8000")
    taskRouter.Run(":" + port)
}
```

### After (Using microservice-commons)
```go
import (
    "github.com/JorgeSaicoski/microservice-commons/server"
    "github.com/JorgeSaicoski/microservice-commons/database"
)

func main() {
    server := server.NewServer(server.ServerOptions{
        ServiceName:    "go-todo-list",
        ServiceVersion: "1.0.0",
        SetupRoutes:    setupRoutes,
    })
    server.Start()
}

func setupRoutes(router *gin.Engine, cfg *config.Config) {
    // Connect to database
    db := database.MustConnect(cfg.DatabaseConfig)
    
    // Setup your existing routes
    api := router.Group("/api")
    {
        api.GET("/tasks", handlers.GetTasks)
        api.POST("/tasks", handlers.CreateTask)
    }
}
```

### Migration Steps
1. **Replace environment handling** - Remove custom `getEnv` functions
2. **Remove CORS setup** - Handled automatically by server
3. **Remove health check** - Included by default
4. **Simplify database connection** - Use database.MustConnect()
5. **Update response format** - Use responses.Success(), responses.Error()
6. **Remove graceful shutdown** - Handled by server.Start()

## 🧪 Testing

### Unit Tests
```go
func TestConfigValidation(t *testing.T) {
    cfg := &config.Config{
        Port: "",  // Invalid
    }
    
    err := cfg.Validate()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "PORT is required")
}
```

### Integration Tests
```go
func TestServerStartup(t *testing.T) {
    server := server.NewServer(server.ServerOptions{
        ServiceName: "test-service",
        SetupRoutes: func(router *gin.Engine, cfg *config.Config) {
            router.GET("/test", func(c *gin.Context) {
                c.JSON(200, gin.H{"status": "ok"})
            })
        },
    })
    
    // Test router without starting server
    router := server.GetRouter()
    // ... test your routes
}
```

## 🚀 Deployment Considerations

### Docker Environment
```dockerfile
# Set environment variables in Dockerfile
ENV SERVICE_NAME=my-service
ENV SERVICE_VERSION=1.0.0
ENV ENVIRONMENT=prod
ENV LOG_LEVEL=warn
```

### Kubernetes ConfigMap
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-service-config
data:
  SERVICE_NAME: "my-service"
  ENVIRONMENT: "prod"
  POSTGRES_HOST: "postgres-service"
  KEYCLOAK_URL: "http://keycloak-service:8080/keycloak"
```

### Health Check Endpoints
Every service automatically gets:
- `GET /health` - Basic health status
- Health checks include database connectivity
- Connection pool statistics available via health details

### Monitoring Integration
```go
// Custom health check data
func setupRoutes(router *gin.Engine, cfg *config.Config) {
    db := database.MustConnect(cfg.DatabaseConfig)
    
    // Add custom health data
    router.GET("/health/detailed", func(c *gin.Context) {
        status, details := database.DetailedHealthCheck(db)
        c.JSON(200, gin.H{
            "service": cfg.ServiceName,
            "version": cfg.ServiceVersion,
            "database": status,
            "details": details,
        })
    })
}
```

## 🛠️ Development Workflow

### 1. Create New Service
```bash
mkdir my-new-service
cd my-new-service
go mod init github.com/company/my-new-service
go get github.com/JorgeSaicoski/microservice-commons
```

### 2. Basic Service Structure
```go
// main.go
package main

import "github.com/JorgeSaicoski/microservice-commons/server"

func main() {
    server := server.NewServer(server.ServerOptions{
        ServiceName:    "my-new-service",
        ServiceVersion: "1.0.0",
        SetupRoutes:    setupRoutes,
    })
    server.Start()
}

// routes.go  
func setupRoutes(router *gin.Engine, cfg *config.Config) {
    // Your service routes here
}
```

### 3. Add Database Models
```go
// models.go
type MyModel struct {
    gorm.Model
    Name string `json:"name"`
}

// In setupRoutes:
db := database.MustConnect(cfg.DatabaseConfig)
database.QuickMigrate(db, &MyModel{})
```

### 4. Add Handlers
```go
import "github.com/JorgeSaicoski/microservice-commons/responses"

func createHandler(c *gin.Context) {
    var model MyModel
    if err := c.ShouldBindJSON(&model); err != nil {
        responses.BadRequest(c, "Invalid data")
        return
    }
    
    // Save model...
    
    responses.Success(c, "Created successfully", model)
}
```

## 🤝 Contributing

This library is part of the Personal Manager ecosystem. When adding features:

1. **Keep it generic** - Features should be useful across multiple services
2. **Maintain backward compatibility** - Existing services depend on this
3. **Add tests** - All new features need test coverage
4. **Update examples** - Keep documentation current
5. **Follow patterns** - Consistency with existing code style

### Adding New Middleware
```go
// middleware/custom.go
func NewCustomMiddleware(options CustomOptions) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Your middleware logic
        c.Next()
    }
}
```

### Adding New Response Types
```go
// responses/custom.go
func CustomResponse(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, CustomFormat{
        Data: data,
        // Custom format
    })
}
```

## 📚 Related Libraries

This library works with the Personal Manager ecosystem:

- **[pgconnect](https://github.com/JorgeSaicoski/pgconnect)** - PostgreSQL repository pattern
- **[keycloak-auth](https://github.com/JorgeSaicoski/keycloak-auth)** - JWT authentication middleware
- **[go-todo-list](https://github.com/JorgeSaicoski/go-todo-list)** - Task management service
- **[go-project-manager](https://github.com/JorgeSaicoski/go-project-manager)** - Project core service

## 📄 License

MIT License - See [LICENSE](LICENSE) for details.

---

## 🔍 Quick Reference

### Import Patterns
```go
import (
    "github.com/JorgeSaicoski/microservice-commons/server"
    "github.com/JorgeSaicoski/microservice-commons/config"
    "github.com/JorgeSaicoski/microservice-commons/database"
    "github.com/JorgeSaicoski/microservice-commons/responses"
    "github.com/JorgeSaicoski/microservice-commons/utils"
)
```

### Common Functions
```go
// Server
server := server.NewServer(options)
server.Start()

// Config
cfg := config.LoadFromEnv()
cfg.Validate()

// Database
db := database.MustConnect(cfg.DatabaseConfig)
database.QuickMigrate(db, models...)
healthy := database.QuickHealthCheck(db)

// Responses
responses.Success(c, "message", data)
responses.BadRequest(c, "error message")
responses.Paginated(c, items, total, page, pageSize)

// Utils
value := utils.GetEnv("KEY", "default")
intValue := utils.GetEnvInt("KEY", 0)
boolValue := utils.GetEnvBool("KEY", false)
```

This library eliminates ~50 lines of boilerplate per service and ensures consistency across the entire Personal Manager microservices ecosystem.