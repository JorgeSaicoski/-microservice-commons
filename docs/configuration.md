# Configuration Guide

This guide covers all configuration options available in microservice-commons, from basic setup to advanced customization.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Environment Variables](#environment-variables)
3. [Configuration Structure](#configuration-structure)
4. [Database Configuration](#database-configuration)
5. [Keycloak Configuration](#keycloak-configuration)
6. [Environment-Specific Settings](#environment-specific-settings)
7. [Validation](#validation)
8. [Best Practices](#best-practices)

## Quick Start

The simplest way to configure your service:

```go
package main

import (
    "github.com/JorgeSaicoski/microservice-commons/config"
    "github.com/JorgeSaicoski/microservice-commons/server"
)

func main() {
    server := server.NewServer(server.ServerOptions{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        SetupRoutes:    setupRoutes,
        // Config is loaded automatically from environment
    })
    server.Start()
}
```

## Environment Variables

### Core Service Configuration

| Variable | Default | Description | Required |
|----------|---------|-------------|----------|
| `SERVICE_NAME` | `"microservice"` | Name of your service | ✅ |
| `SERVICE_VERSION` | `"1.0.0"` | Version of your service | ✅ |
| `PORT` | `"8000"` | Port to run the server on | ✅ |
| `ENVIRONMENT` | `"dev"` | Environment: `dev`, `staging`, `prod` | ❌ |
| `LOG_LEVEL` | `"info"` | Log level: `debug`, `info`, `warn`, `error` | ❌ |
| `ALLOWED_ORIGINS` | `"http://localhost:3000"` | Comma-separated CORS origins | ❌ |

### Database Configuration

| Variable | Default | Description | Required |
|----------|---------|-------------|----------|
| `POSTGRES_HOST` | `"localhost"` | Database host | ✅ |
| `POSTGRES_PORT` | `"5432"` | Database port | ✅ |
| `POSTGRES_USER` | `"postgres"` | Database username | ✅ |
| `POSTGRES_PASSWORD` | `"postgres"` | Database password | ✅ |
| `POSTGRES_DB` | `"defaultdb"` | Database name | ✅ |
| `POSTGRES_SSLMODE` | `"disable"` | SSL mode: `disable`, `require`, `verify-full` | ❌ |
| `POSTGRES_TIMEZONE` | `"UTC"` | Database timezone | ❌ |
| `POSTGRES_MAX_IDLE_CONNS` | `10` | Maximum idle connections | ❌ |
| `POSTGRES_MAX_OPEN_CONNS` | `100` | Maximum open connections | ❌ |
| `POSTGRES_LOG_LEVEL` | `"silent"` | DB log level: `silent`, `error`, `warn`, `info` | ❌ |

### Database Connection Retry

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_MAX_RETRIES` | `3` | Maximum connection retry attempts |
| `DB_RETRY_DELAY_SECONDS` | `30` | Delay between retry attempts |
| `DB_LOG_LEVEL` | `"silent"` | Database connection log level |

### Keycloak Authentication

| Variable | Default | Description | Required |
|----------|---------|-------------|----------|
| `KEYCLOAK_URL` | `""` | Keycloak server URL | ❌* |
| `KEYCLOAK_REALM` | `"master"` | Keycloak realm | ❌* |
| `KEYCLOAK_PUBLIC_KEY` | `""` | Base64 encoded public key | ❌* |
| `KEYCLOAK_REQUIRED_CLAIMS` | `"sub,preferred_username"` | Required JWT claims | ❌ |
| `KEYCLOAK_SKIP_PATHS` | `"/health,/metrics"` | Paths to skip auth | ❌ |
| `KEYCLOAK_KEY_REFRESH_INTERVAL` | `"1h"` | Key refresh interval | ❌ |
| `KEYCLOAK_HTTP_TIMEOUT` | `"10s"` | HTTP timeout for Keycloak | ❌ |

*Either `KEYCLOAK_PUBLIC_KEY` or both `KEYCLOAK_URL` + `KEYCLOAK_REALM` must be provided.

## Configuration Structure

### Loading Configuration

```go
// Load from environment variables
config := config.LoadFromEnv()

// Load with service info override
config := config.LoadWithServiceInfo("my-service", "2.0.0")

// Validate configuration
if err := config.Validate(); err != nil {
    log.Fatal("Invalid configuration:", err)
}
```

### Configuration Struct

```go
type Config struct {
    Port           string            // Server port
    AllowedOrigins []string          // CORS allowed origins
    DatabaseConfig DatabaseConfig   // Database settings
    KeycloakConfig KeycloakConfig   // Authentication settings
    LogLevel       string           // Logging level
    Environment    string           // Environment name
    ServiceName    string           // Service name
    ServiceVersion string           // Service version
}
```

### Environment Helpers

```go
config := config.LoadFromEnv()

// Check environment
if config.IsDevelopment() {
    // Development-specific logic
}

if config.IsProduction() {
    // Production-specific logic
}

if config.IsStaging() {
    // Staging-specific logic
}

// Get structured log level
logLevel := config.GetLogLevel() // Returns LogLevel enum
```

## Database Configuration

### Basic Database Setup

```go
// Get database config
dbConfig := config.LoadDatabaseConfig()

// Validate database config
if err := dbConfig.Validate(); err != nil {
    log.Fatal("Invalid database config:", err)
}

// Get connection string
connStr := dbConfig.ConnectionString()
// Output: "host=localhost port=5432 user=postgres password=secret dbname=mydb sslmode=disable TimeZone=UTC"
```

### Database Configuration Methods

```go
dbConfig := config.LoadDatabaseConfig()

// Check SSL configuration
if dbConfig.IsSSLEnabled() {
    fmt.Println("SSL is enabled")
}

// Get validated log level
logLevel := dbConfig.GetLogLevel() // Returns: silent, error, warn, info
```

### Connection Pool Configuration

Understanding connection pool settings:

```bash
# Connection Pool Settings
POSTGRES_MAX_IDLE_CONNS=10    # Connections kept open when idle
POSTGRES_MAX_OPEN_CONNS=100   # Maximum total connections

# Rule: MAX_IDLE <= MAX_OPEN
# Recommended: MAX_IDLE = 10-20% of MAX_OPEN for most cases
```

**Connection Pool Guidelines:**

- **Small services**: `MAX_IDLE=5`, `MAX_OPEN=25`
- **Medium services**: `MAX_IDLE=10`, `MAX_OPEN=50`
- **Large services**: `MAX_IDLE=20`, `MAX_OPEN=100`
- **High-traffic services**: `MAX_IDLE=50`, `MAX_OPEN=200`

## Keycloak Configuration

### Static Key Configuration

If you have a static public key:

```bash
KEYCLOAK_PUBLIC_KEY=LS0tLS1CRUdJTi... # Base64 encoded public key
KEYCLOAK_REQUIRED_CLAIMS=sub,preferred_username,email
KEYCLOAK_SKIP_PATHS=/health,/metrics,/public
```

### JWKS Endpoint Configuration

For dynamic key fetching:

```bash
KEYCLOAK_URL=https://keycloak.company.com
KEYCLOAK_REALM=production
KEYCLOAK_KEY_REFRESH_INTERVAL=1h
KEYCLOAK_HTTP_TIMEOUT=10s
```

### Keycloak Validation

```go
keycloakConfig := config.LoadKeycloakConfig()

// Validate configuration
if err := keycloakConfig.Validate(); err != nil {
    log.Fatal("Invalid Keycloak config:", err)
}

// Check configuration type
if keycloakConfig.HasStaticKey() {
    fmt.Println("Using static public key")
}

if keycloakConfig.HasJWKS() {
    jwksURL := keycloakConfig.GetJWKSURL()
    fmt.Printf("JWKS URL: %s\n", jwksURL)
}

// Check if path should skip authentication
if keycloakConfig.ShouldSkipPath("/health") {
    fmt.Println("Health endpoint skips auth")
}
```

## Environment-Specific Settings

### Development Environment

```bash
# .env.dev
ENVIRONMENT=dev
LOG_LEVEL=debug
POSTGRES_HOST=localhost
POSTGRES_LOG_LEVEL=info
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

### Staging Environment

```bash
# .env.staging
ENVIRONMENT=staging
LOG_LEVEL=info
POSTGRES_HOST=staging-db.company.com
POSTGRES_SSLMODE=require
ALLOWED_ORIGINS=https://staging.company.com
```

### Production Environment

```bash
# .env.prod
ENVIRONMENT=prod
LOG_LEVEL=warn
POSTGRES_HOST=prod-db.company.com
POSTGRES_SSLMODE=verify-full
POSTGRES_MAX_IDLE_CONNS=20
POSTGRES_MAX_OPEN_CONNS=200
ALLOWED_ORIGINS=https://app.company.com
```

## Validation

### Built-in Validation

```go
config := config.LoadFromEnv()

// Comprehensive validation
if err := config.Validate(); err != nil {
    fmt.Printf("Configuration errors: %v\n", err)
    // Output: Configuration validation failed:
    //   - PORT is required
    //   - DATABASE_HOST is required
}
```

### Custom Validation

```go
import "github.com/JorgeSaicoski/microservice-commons/config"

// Create custom validator
validator := config.NewValidator()

// Add custom validation rules
validator.ValidateRequired("CUSTOM_API_KEY", os.Getenv("CUSTOM_API_KEY"))
validator.ValidateURL("WEBHOOK_URL", os.Getenv("WEBHOOK_URL"))
validator.ValidatePort("REDIS_PORT", os.Getenv("REDIS_PORT"))

// Check for errors
if validator.HasErrors() {
    for _, err := range validator.GetErrors() {
        fmt.Printf("Validation error: %s\n", err.Error())
    }
}
```

### Validation Rules

```go
validator := config.NewValidator()

// Field validation
validator.ValidateRequired("FIELD", value)
validator.ValidatePort("PORT", "8080")
validator.ValidateURL("URL", "https://example.com")
validator.ValidateOneOf("ENV", "prod", []string{"dev", "staging", "prod"})
validator.ValidateMinMax("TIMEOUT", "30", 1, 300)

// Get all errors
if err := validator.Error(); err != nil {
    log.Fatal(err)
}
```

## Best Practices

### 1. Environment Variable Naming

```bash
# ✅ Good: Consistent, descriptive naming
SERVICE_NAME=user-service
POSTGRES_HOST=localhost
KEYCLOAK_URL=https://auth.company.com

# ❌ Bad: Inconsistent, unclear naming
service=user-service
db_host=localhost
auth_server=https://auth.company.com
```

### 2. Secret Management

```bash
# ✅ Good: Use secrets management
POSTGRES_PASSWORD_FILE=/run/secrets/db_password
KEYCLOAK_PUBLIC_KEY_FILE=/run/secrets/keycloak_key

# ❌ Bad: Hardcoded secrets in environment
POSTGRES_PASSWORD=hardcoded_password_123
```

### 3. Environment-Specific Configuration

```go
// ✅ Good: Environment-aware configuration
func setupLogging(cfg *config.Config) {
    if cfg.IsDevelopment() {
        // Pretty, colorful logs
        gin.SetMode(gin.DebugMode)
    } else {
        // Structured JSON logs
        gin.SetMode(gin.ReleaseMode)
    }
}
```

### 4. Configuration Validation

```go
// ✅ Good: Fail fast with clear errors
func main() {
    config := config.LoadFromEnv()
    
    if err := config.Validate(); err != nil {
        log.Fatalf("Configuration error: %v", err)
    }
    
    // Continue with valid configuration
    server := server.NewServer(...)
}
```

### 5. Documentation

```bash
# ✅ Good: Document all environment variables
# .env.example with comments
SERVICE_NAME=my-service          # Name of the microservice
PORT=8000                       # Port to run the service on
ENVIRONMENT=dev                 # Environment: dev, staging, prod

# Database configuration
POSTGRES_HOST=localhost         # PostgreSQL host
POSTGRES_PORT=5432             # PostgreSQL port
POSTGRES_DB=myservice_db       # Database name
```

### 6. Default Values

```go
// ✅ Good: Sensible defaults for optional settings
POSTGRES_MAX_IDLE_CONNS=10     // Good default for most services
LOG_LEVEL=info                 // Balanced logging level
POSTGRES_SSLMODE=disable       // Good for development

// ✅ Required for critical settings
SERVICE_NAME=                  // No default - must be specified
POSTGRES_PASSWORD=             // No default - must be specified
```

### 7. Configuration Testing

```go
func TestConfigurationLoading(t *testing.T) {
    // Set test environment
    os.Setenv("SERVICE_NAME", "test-service")
    os.Setenv("POSTGRES_PASSWORD", "test-password")
    defer func() {
        os.Unsetenv("SERVICE_NAME")
        os.Unsetenv("POSTGRES_PASSWORD")
    }()
    
    config := config.LoadFromEnv()
    
    if config.ServiceName != "test-service" {
        t.Errorf("Expected service name 'test-service', got %s", config.ServiceName)
    }
}
```

### 8. Configuration in Docker

```dockerfile
# Use build args for compile-time settings
ARG SERVICE_NAME=my-service
ARG SERVICE_VERSION=1.0.0

# Use environment variables for runtime settings
ENV PORT=8000
ENV LOG_LEVEL=info

# Use secrets for sensitive data
COPY --from=secrets /run/secrets/db_password /run/secrets/db_password
```

### 9. Configuration in Kubernetes

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-service-config
data:
  SERVICE_NAME: "my-service"
  ENVIRONMENT: "prod"
  LOG_LEVEL: "info"
  POSTGRES_HOST: "postgres-service"

---
apiVersion: v1
kind: Secret
metadata:
  name: my-service-secrets
data:
  POSTGRES_PASSWORD: <base64-encoded-password>
  KEYCLOAK_PUBLIC_KEY: <base64-encoded-key>
```

## Troubleshooting

### Common Configuration Issues

1. **"PORT is required" error**
   ```bash
   # Make sure PORT is set
   export PORT=8000
   ```

2. **Database connection fails**
   ```bash
   # Check all required database variables
   export POSTGRES_HOST=localhost
   export POSTGRES_PASSWORD=your_password
   ```

3. **CORS issues in development**
   ```bash
   # Add your frontend URL to allowed origins
   export ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
   ```

4. **Keycloak authentication not working**
   ```bash
   # Ensure either static key or JWKS is configured
   export KEYCLOAK_PUBLIC_KEY=your_base64_key
   # OR
   export KEYCLOAK_URL=https://keycloak.company.com
   export KEYCLOAK_REALM=your_realm
   ```

### Configuration Debugging

```go
config := config.LoadFromEnv()

// Print configuration (be careful with secrets!)
fmt.Printf("Service: %s v%s\n", config.ServiceName, config.ServiceVersion)
fmt.Printf("Environment: %s\n", config.Environment)
fmt.Printf("Port: %s\n", config.Port)
fmt.Printf("Database Host: %s\n", config.DatabaseConfig.Host)

// Don't print passwords in production!
if config.IsDevelopment() {
    fmt.Printf("Database Config: %+v\n", config.DatabaseConfig)
}
```