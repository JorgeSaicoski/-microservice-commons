# microservice-commons Documentation

Welcome to the comprehensive documentation for microservice-commons - a Go library that eliminates boilerplate code and ensures consistency across your microservices architecture.

## 📚 Documentation Index

| Document | Description | Audience |
|----------|-------------|----------|
| [Configuration Guide](configuration.md) | Complete configuration reference | All Developers |
| [Middleware Documentation](middleware.md) | All middleware components and usage | Backend Developers |
| [Response Patterns Guide](responses.md) | Standardized API response patterns | API Developers |
| [Migration Guide](migration.md) | Step-by-step migration from existing services | Teams migrating existing services |

## 🚀 Quick Start

If you're new to microservice-commons, start here:

### 1. Installation

```bash
go get github.com/JorgeSaicoski/microservice-commons
```

### 2. Your First Service (5 minutes)

```go
package main

import (
    "github.com/JorgeSaicoski/microservice-commons/server"
    "github.com/JorgeSaicoski/microservice-commons/responses"
    "github.com/gin-gonic/gin"
)

func main() {
    server := server.NewServer(server.ServerOptions{
        ServiceName:    "hello-service",
        ServiceVersion: "1.0.0",
        SetupRoutes:    setupRoutes,
    })
    server.Start()
}

func setupRoutes(router *gin.Engine, cfg *config.Config) {
    api := router.Group("/api/v1")
    {
        api.GET("/hello", func(c *gin.Context) {
            responses.Success(c, "Hello from microservice-commons!", nil)
        })
    }
}
```

### 3. Run Your Service

```bash
# Set required environment variables
export SERVICE_NAME=hello-service
export PORT=8000

# Run the service
go run main.go
```

### 4. Test Your Service

```bash
# Health check (automatic)
curl http://localhost:8000/health

# Your API endpoint
curl http://localhost:8000/api/v1/hello
```

**That's it!** You now have a production-ready microservice with:
- ✅ Health checks
- ✅ CORS handling
- ✅ Graceful shutdown
- ✅ Request logging
- ✅ Standardized responses
- ✅ Error recovery

## 📖 Learning Path

### For New Projects

1. **[Configuration Guide](configuration.md)** - Set up your environment
2. **[Middleware Documentation](middleware.md)** - Add authentication, logging, etc.
3. **[Response Patterns Guide](responses.md)** - Standardize your API responses
4. **[Examples](../examples/)** - See working examples

### For Existing Services

1. **[Migration Guide](migration.md)** - Step-by-step migration process
2. **[Configuration Guide](configuration.md)** - Update your configuration
3. **[Response Patterns Guide](responses.md)** - Standardize existing responses

## 🏗️ Architecture Overview

microservice-commons is built around these core components:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Server      │    │   Middleware    │    │   Responses     │
│                 │    │                 │    │                 │
│ • Setup         │    │ • CORS          │    │ • Success       │
│ • Graceful      │    │ • Auth          │    │ • Errors        │
│   Shutdown      │    │ • Logging       │    │ • Pagination    │
│ • Health        │    │ • Recovery      │    │ • Validation    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
         ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
         │   Configuration │    │    Database     │    │     Types       │
         │                 │    │                 │    │                 │
         │ • Environment   │    │ • Connection    │    │ • Common Models │
         │ • Validation    │    │ • Health        │    │ • Enums         │
         │ • Multi-env     │    │ • Migration     │    │ • Pagination    │
         └─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🎯 What Problems Does This Solve?

### Before microservice-commons

```go
// 50+ lines of boilerplate per service
func main() {
    // Manual CORS setup
    router := gin.Default()
    allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
    origins := strings.Split(allowedOrigins, ",")
    router.Use(cors.New(cors.Config{
        AllowOrigins: origins,
        // ... 10 more lines of CORS config
    }))
    
    // Manual health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })
    
    // Manual server setup with graceful shutdown
    srv := &http.Server{
        Addr:    ":" + getEnv("PORT", "8000"),
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
    
    // ... 20+ more lines for graceful shutdown
}
```

### After microservice-commons

```go
// 10 lines total
func main() {
    server := server.NewServer(server.ServerOptions{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        SetupRoutes:    setupRoutes,
    })
    server.Start() // Everything else is handled automatically
}
```

## 🔧 Key Features

### 🚀 **Server Management**
- Automatic CORS configuration
- Built-in health check endpoints (`/health`, `/ready`, `/live`)
- Graceful shutdown handling
- Request logging and recovery
- Environment-based configuration

### 🔐 **Authentication & Security**
- JWT/Bearer token authentication
- API key authentication
- Role-based access control
- Request ID tracking
- CORS protection

### 📊 **Monitoring & Health**
- Database health checks
- Memory usage monitoring
- External service health checks
- Detailed health reporting
- Kubernetes-ready probes

### 📄 **Standardized Responses**
- Consistent success/error formats
- Built-in pagination support
- HTTP status code helpers
- Request tracking
- Validation error formatting

### 🗄️ **Database Integration**
- Connection pooling with retry logic
- Health monitoring and statistics
- Migration utilities with safety checks
- Integration with pgconnect library

## 🌟 Benefits

| Benefit | Description | Time Saved |
|---------|-------------|------------|
| **Reduced Boilerplate** | Eliminate 50+ lines of setup code per service | 2-3 hours per service |
| **Consistency** | Standardized patterns across all services | Ongoing maintenance |
| **Reliability** | Battle-tested middleware and error handling | Debugging time |
| **Developer Experience** | Clear APIs and comprehensive documentation | Learning curve |
| **Production Ready** | Health checks, monitoring, graceful shutdown | DevOps integration |

## 📊 Comparison with Manual Setup

| Feature | Manual Implementation | microservice-commons | Lines Saved |
|---------|----------------------|----------------------|-------------|
| Server Setup | 30+ lines | 5 lines | 25+ |
| CORS Configuration | 15+ lines | Automatic | 15+ |
| Health Checks | 10+ lines | Automatic | 10+ |
| Graceful Shutdown | 20+ lines | Automatic | 20+ |
| Error Responses | 5+ lines each | 1 line each | 4+ per endpoint |
| Pagination | 15+ lines | 3 lines | 12+ |
| **Total per service** | **~100 lines** | **~10 lines** | **~90 lines** |

## 🔄 Integration with Existing Tools

microservice-commons works seamlessly with:

- **[pgconnect](https://github.com/JorgeSaicoski/pgconnect)** - PostgreSQL repository pattern
- **[keycloak-auth](https://github.com/JorgeSaicoski/keycloak-auth)** - JWT authentication middleware
- **Gin Framework** - HTTP web framework
- **GORM** - Database ORM
- **Docker** - Containerization
- **Kubernetes** - Container orchestration

## 🚦 Getting Started Paths

Choose your path based on your situation:

### 🆕 New Project
```bash
# 1. Start with basic example
cd examples/basic-service
cp .env.example .env
go run main.go

# 2. Check advanced features
cd ../advanced-service
go run main.go
```

### 🔄 Migrating Existing Service
```bash
# 1. Read migration guide
open docs/migration.md

# 2. Follow step-by-step process
# 3. Use validation scripts
./validate-migration.sh
```

### 🏢 Enterprise Team
```bash
# 1. Review architecture decisions
open docs/configuration.md

# 2. Set up standards
# 3. Create team templates
# 4. Migrate services gradually
```

## 📞 Support and Resources

### Documentation
- 📖 **[Configuration Guide](configuration.md)** - Complete setup reference
- 🔧 **[Middleware Guide](middleware.md)** - All middleware components
- 📄 **[Response Patterns](responses.md)** - API response standards
- 🔄 **[Migration Guide](migration.md)** - Migrate existing services

### Examples
- 🎯 **[Basic Service](../examples/basic-service/)** - Simple CRUD API
- 🚀 **[Advanced Service](../examples/advanced-service/)** - Full-featured service

### Community
- 🐛 **Issues** - Report bugs and request features
- 💬 **Discussions** - Ask questions and share experiences
- 📚 **Wiki** - Community-driven documentation

## 🔄 Version Compatibility

| microservice-commons | Go Version | Dependencies |
|----------------------|------------|---------------|
| v1.x.x | Go 1.19+ | Gin v1.9+, GORM v1.25+ |

## 📈 Roadmap

### Current Version (v1.0)
- ✅ Core server setup
- ✅ Standard middleware
- ✅ Response patterns
- ✅ Database integration
- ✅ Health monitoring

### Future Versions
- 🔄 Metrics collection (Prometheus)
- 🔄 Distributed tracing
- 🔄 Rate limiting
- 🔄 Caching middleware
- 🔄 gRPC support

## 🤝 Contributing

We welcome contributions! Whether you're:
- 🐛 Reporting bugs
- 💡 Suggesting features
- 📝 Improving documentation
- 🔧 Contributing code

Please see our contributing guidelines and join our community.

---

**Ready to eliminate boilerplate and build better microservices?**

Start with the [Configuration Guide](configuration.md) or jump into the [examples](../examples/) to see microservice-commons in action!