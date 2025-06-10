# Advanced Service Example

This example demonstrates advanced features of microservice-commons including authentication, pagination, filtering, custom middleware, and more.

## Features Demonstrated

- ✅ **Authentication & Authorization**: Mock JWT with role-based access
- ✅ **Pagination**: Built-in pagination support with filtering
- ✅ **Advanced Models**: Using types from microservice-commons
- ✅ **Custom Middleware**: Request ID tracking and logging
- ✅ **Health Monitoring**: Database and memory health checks
- ✅ **Search & Filtering**: Advanced query capabilities
- ✅ **JSONB Support**: Metadata storage in PostgreSQL
- ✅ **Role-based Routes**: Admin-only endpoints
- ✅ **Request Validation**: Input validation and error handling

## Project Structure

```
advanced-service/
├── main.go              # Main application with all features
├── .env.example         # Environment configuration
└── README.md           # This file
```

## Quick Start

1. **Setup environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

2. **Install dependencies**:
   ```bash
   go mod init advanced-service-example
   go get github.com/JorgeSaicoski/microservice-commons
   go get github.com/gin-gonic/gin
   go get gorm.io/gorm
   ```

3. **Setup PostgreSQL with JSONB support**:
   ```bash
   docker run --name postgres-advanced \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=advanced_projects \
     -p 5432:5432 \
     -d postgres:15
   ```

4. **Run the service**:
   ```bash
   go run main.go
   ```

## API Documentation

### Authentication

Most endpoints require authentication. Use these tokens for testing:

- **User Token**: `Bearer user-token`
- **Admin Token**: `Bearer admin-token`

### Endpoints

#### Public Endpoints
```bash
# Register a new user
POST /api/v1/users/register
{
  "username": "johndoe",
  "email": "john@example.com",
  "name": "John Doe"
}

# Login
POST /api/v1/users/login
{
  "email": "john@example.com",
  "password": "password123"
}
```

#### Protected User Endpoints
```bash
# Get users with pagination
GET /api/v1/users?page=1&page_size=10
Authorization: Bearer user-token

# Get specific user
GET /api/v1/users/1
Authorization: Bearer user-token

# Update user
PUT /api/v1/users/1
Authorization: Bearer user-token
{
  "name": "John Updated",
  "status": "active"
}
```

#### Project Endpoints
```bash
# Get projects with filtering and pagination
GET /api/v1/projects?page=1&page_size=5&status=active&priority=high&search=api
Authorization: Bearer user-token

# Create project
POST /api/v1/projects
Authorization: Bearer user-token
{
  "name": "New Project",
  "description": "Project description",
  "status": "active",
  "priority": "high",
  "tags": ["api", "microservice"]
}

# Search projects
GET /api/v1/projects/search?q=microservice
Authorization: Bearer user-token
```

#### Admin Endpoints
```bash
# Get service statistics
GET /api/v1/admin/stats
Authorization: Bearer admin-token

# Delete user (admin only)
DELETE /api/v1/admin/users/1
Authorization: Bearer admin-token
```

#### Health Endpoints
```bash
# Basic health check
GET /health

# Detailed health with database and memory monitoring
GET /health/detailed

# Readiness probe (Kubernetes)
GET /ready

# Liveness probe (Kubernetes)  
GET /live
```

## Advanced Features

### 1. Pagination with Filtering

```bash
# Get projects with multiple filters
GET /api/v1/projects?page=1&page_size=10&status=active&priority=high&search=api

Response:
{
  "data": [...],
  "total": 50,
  "page": 1,
  "page_size": 10,
  "total_pages": 5,
  "has_next": true,
  "has_prev": false,
  "timestamp": "2023-12-25T15:30:45Z"
}
```

### 2. Advanced Models with Types

```go
type Project struct {
    types.BaseModel                    // ID, CreatedAt, UpdatedAt, DeletedAt
    Name        string       `json:"name"`
    Status      types.Status `json:"status"`      // Enum: active, inactive, etc.
    Priority    types.Priority `json:"priority"`  // Enum: low, medium, high, critical
    Tags        types.Tags   `json:"tags"`        // String array
    Metadata    types.Metadata `json:"metadata"`  // JSONB for flexible data
}
```

### 3. Health Monitoring

The service includes comprehensive health monitoring:

```json
{
  "status": "healthy",
  "timestamp": "2023-12-25T15:30:45Z",
  "service": "advanced-project-service",
  "version": "2.0.0",
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

### 4. Request Tracking

Every request gets a unique ID for tracing:

```
X-Request-ID: a1b2c3