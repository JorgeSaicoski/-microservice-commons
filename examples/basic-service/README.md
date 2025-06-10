# Basic Service Example

This example demonstrates how to create a simple REST API service using microservice-commons.

## Features Demonstrated

- ✅ Server setup with microservice-commons
- ✅ Database connection and migration
- ✅ CRUD operations with standardized responses
- ✅ Error handling
- ✅ Health checks (automatic)
- ✅ CORS configuration (automatic)
- ✅ Graceful shutdown (automatic)

## Quick Start

1. **Setup environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your database settings
   ```

2. **Install dependencies**:
   ```bash
   go mod init basic-service-example
   go get github.com/JorgeSaicoski/microservice-commons
   go get github.com/gin-gonic/gin
   go get gorm.io/gorm
   ```

3. **Setup PostgreSQL database**:
   ```bash
   # Using Docker
   docker run --name postgres-basic \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=basic_tasks \
     -p 5432:5432 \
     -d postgres:15
   ```

4. **Run the service**:
   ```bash
   go run main.go
   ```

## API Endpoints

### Health Check
- `GET /health` - Basic health status
- `GET /health/detailed` - Detailed health information

### Tasks API
- `GET /api/v1/tasks` - List all tasks
- `POST /api/v1/tasks` - Create a new task
- `GET /api/v1/tasks/:id` - Get a specific task
- `PUT /api/v1/tasks/:id` - Update a task
- `DELETE /api/v1/tasks/:id` - Delete a task

## Example Requests

### Create a task
```bash
curl -X POST http://localhost:8000/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn microservice-commons",
    "description": "Build a simple REST API",
    "completed": false
  }'
```

### Get all tasks
```bash
curl http://localhost:8000/api/v1/tasks
```

### Update a task
```bash
curl -X PUT http://localhost:8000/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn microservice-commons",
    "description": "Build a simple REST API",
    "completed": true
  }'
```

## Response Format

All responses follow the standardized format:

### Success Response
```json
{
  "message": "Tasks retrieved successfully",
  "data": [...],
  "timestamp": "2023-12-25T15:30:45Z"
}
```

### Error Response
```json
{
  "error": "Task not found",
  "code": "not_found",
  "timestamp": "2023-12-25T15:30:45Z",
  "path": "/api/v1/tasks/999"
}
```

## What You Get for Free

By using microservice-commons, this example automatically includes:

1. **Health Checks**: `/health` endpoint with database connectivity
2. **CORS**: Configured for your allowed origins
3. **Error Handling**: Standardized error responses
4. **Graceful Shutdown**: Proper cleanup on SIGTERM/SIGINT
5. **Request Logging**: Automatic request/response logging
6. **Recovery**: Panic recovery with proper error responses
7. **Database Connection**: With retry logic and health monitoring

## Next Steps

- Check out the [advanced-service](../advanced-service/) example for more features
- See how to add authentication, pagination, and custom middleware