# Response Patterns Guide

This guide covers the standardized response patterns provided by microservice-commons, ensuring consistent API responses across all your services.

## Table of Contents

1. [Overview](#overview)
2. [Success Responses](#success-responses)
3. [Error Responses](#error-responses)
4. [Pagination Responses](#pagination-responses)
5. [File Responses](#file-responses)
6. [Custom Responses](#custom-responses)
7. [Response Standards](#response-standards)
8. [Best Practices](#best-practices)

## Overview

microservice-commons provides a comprehensive set of response helpers that ensure consistency across your microservices. All responses follow standardized formats with proper HTTP status codes, timestamps, and error codes.

### Basic Usage

```go
import (
    "github.com/JorgeSaicoski/microservice-commons/responses"
    "github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {
    // Success response
    responses.Success(c, "Operation completed", data)
    
    // Error response
    responses.BadRequest(c, "Invalid input")
    
    // Paginated response
    responses.Paginated(c, items, total, page, pageSize)
}
```

## Success Responses

### Basic Success Response

```go
func getUser(c *gin.Context) {
    user := User{ID: 1, Name: "John Doe"}
    responses.Success(c, "User retrieved successfully", user)
}
```

**Response:**
```json
{
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "name": "John Doe"
  },
  "timestamp": "2023-12-25T15:30:45Z"
}
```

### Success Response Variants

```go
// Success with data
responses.Success(c, "Task completed", task)

// Success with custom status
responses.SuccessWithStatus(c, http.StatusAccepted, "Request accepted", nil)

// Created response (201)
responses.Created(c, "User created successfully", user)

// Accepted response (202)
responses.Accepted(c, "Request queued for processing", queueInfo)

// No content response (204)
responses.NoContent(c)

// Simple OK message
responses.OK(c, "Operation successful")

// Data only (no message)
responses.Data(c, users)

// Message only (no data)
responses.Message(c, "Cache cleared successfully")
```

### Success Response Formats

#### Standard Success Format
```json
{
  "message": "string",
  "data": "any",
  "timestamp": "2023-12-25T15:30:45Z"
}
```

#### Data-Only Format
```json
{
  "data": [...],
  "timestamp": "2023-12-25T15:30:45Z"
}
```

#### Message-Only Format
```json
{
  "message": "Operation completed",
  "timestamp": "2023-12-25T15:30:45Z"
}
```

## Error Responses

### Standard Error Responses

```go
// 400 Bad Request
responses.BadRequest(c, "Invalid input data")

// 401 Unauthorized
responses.Unauthorized(c, "Authentication required")

// 403 Forbidden
responses.Forbidden(c, "Access denied")

// 404 Not Found
responses.NotFound(c, "Resource not found")

// 405 Method Not Allowed
responses.MethodNotAllowed(c, "POST method not allowed")

// 409 Conflict
responses.Conflict(c, "Resource already exists")

// 422 Unprocessable Entity
responses.UnprocessableEntity(c, "Validation failed")

// 429 Too Many Requests
responses.TooManyRequests(c, "Rate limit exceeded")

// 500 Internal Server Error
responses.InternalError(c, "Database connection failed")

// 503 Service Unavailable
responses.ServiceUnavailable(c, "Service temporarily unavailable")
```

### Error Response with Details

```go
func validateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        responses.BadRequestWithDetails(c, 
            "Invalid user data", 
            err.Error())
        return
    }
    
    // Validation logic...
}
```

### Error Response with Metadata

```go
func processPayment(c *gin.Context) {
    if err := chargeCard(cardInfo); err != nil {
        metadata := map[string]interface{}{
            "transaction_id": txnID,
            "error_code":     err.Code,
            "retry_after":    300, // seconds
        }
        
        responses.ErrorWithMetadata(c, 
            http.StatusPaymentRequired,
            "payment_failed",
            "Payment processing failed",
            metadata)
        return
    }
}
```

### Specialized Error Responses

```go
// Validation error with field details
validationErrors := []ValidationError{
    {Field: "email", Message: "Invalid email format"},
    {Field: "password", Message: "Password too weak"},
}
responses.ValidationError(c, "Validation failed", validationErrors)

// Database error
responses.DatabaseError(c, "Failed to save user")

// External service error
responses.ExternalServiceError(c, "Payment gateway unavailable")
```

### Error Response Format

```json
{
  "error": "Resource not found",
  "code": "not_found",
  "details": "User with ID 123 does not exist",
  "timestamp": "2023-12-25T15:30:45Z",
  "path": "/api/users/123",
  "request_id": "req_abc123",
  "metadata": {
    "suggested_action": "Check user ID and try again"
  }
}
```

### Standard Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | `bad_request` | Invalid request data |
| 401 | `unauthorized` | Authentication required |
| 403 | `forbidden` | Access denied |
| 404 | `not_found` | Resource not found |
| 405 | `method_not_allowed` | HTTP method not allowed |
| 409 | `conflict` | Resource conflict |
| 422 | `unprocessable_entity` | Validation failed |
| 429 | `too_many_requests` | Rate limit exceeded |
| 500 | `internal_error` | Internal server error |
| 502 | `external_service_error` | External service error |
| 503 | `service_unavailable` | Service unavailable |

### Custom Error Types

```go
// Create custom API errors
customError := responses.NewAPIError(
    http.StatusTeapot, 
    "coffee_not_available", 
    "I'm a teapot, cannot brew coffee")

// Handle the error
responses.HandleAPIError(c, customError)

// Or use convenience constructors
badRequestErr := responses.NewBadRequestError("Invalid email format")
notFoundErr := responses.NewNotFoundError("User not found")
validationErr := responses.NewValidationError("Validation failed", validationDetails)
```

## Pagination Responses

### Basic Pagination

```go
func getUsers(c *gin.Context) {
    // Get pagination parameters from query
    params := responses.GetPaginationParams(c)
    
    // Fetch data
    users, total := getUsersPaginated(params.Offset, params.Limit)
    
    // Return paginated response
    responses.Paginated(c, users, total, params.Page, params.PageSize)
}
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 10, max: 100)

**Response:**
```json
{
  "data": [...],
  "total": 150,
  "page": 1,
  "page_size": 10,
  "total_pages": 15,
  "has_next": true,
  "has_prev": false,
  "timestamp": "2023-12-25T15:30:45Z"
}
```

### Advanced Pagination

```go
func getProjects(c *gin.Context) {
    // Custom pagination defaults
    params := responses.GetPaginationParamsWithDefaults(c, 1, 20, 50)
    
    // Validate parameters
    if err := responses.ValidatePaginationParams(params.Page, params.PageSize); err != nil {
        responses.HandleError(c, err)
        return
    }
    
    // Fetch data with filtering
    projects, total := getProjectsFiltered(params, c.Query("status"))
    
    // Return with custom status
    responses.PaginatedWithStatus(c, http.StatusOK, projects, total, params.Page, params.PageSize)
}
```

### Pagination with Metadata

```go
func getTasksPaginated(c *gin.Context) {
    params := responses.GetPaginationParams(c)
    tasks, total := getTasksData(params)
    
    // Create pagination metadata
    meta := responses.CreatePaginationMeta(total, params.Page, params.PageSize)
    
    // Return with separate metadata
    responses.PaginatedWithMeta(c, tasks, meta)
}
```

### Cursor-Based Pagination

```go
func getTimelineEvents(c *gin.Context) {
    cursor, limit := responses.GetCursorParams(c)
    
    events, nextCursor, hasNext := getEventsFromCursor(cursor, limit)
    
    responses.CursorPaginated(c, events, nextCursor, "", hasNext, false)
}
```

**Response:**
```json
{
  "data": [...],
  "next_cursor": "eyJpZCI6MTIzLCJ0aW1lIjoiMjAyMy0xMi0yNSJ9",
  "prev_cursor": "",
  "has_next": true,
  "has_prev": false,
  "timestamp": "2023-12-25T15:30:45Z"
}
```

### Pagination with HATEOAS Links

```go
func getUsersWithLinks(c *gin.Context) {
    params := responses.GetPaginationParams(c)
    users, total := getUsersData(params)
    
    baseURL := "https://api.company.com/users"
    
    responses.PaginatedWithLinks(c, users, total, params.Page, params.PageSize, baseURL)
}
```

**Response:**
```json
{
  "data": [...],
  "meta": {
    "total": 150,
    "page": 2,
    "page_size": 10,
    "total_pages": 15,
    "has_next": true,
    "has_prev": true
  },
  "links": {
    "self": "https://api.company.com/users?page=2&page_size=10",
    "first": "https://api.company.com/users?page=1&page_size=10",
    "last": "https://api.company.com/users?page=15&page_size=10",
    "next": "https://api.company.com/users?page=3&page_size=10",
    "prev": "https://api.company.com/users?page=1&page_size=10"
  },
  "timestamp": "2023-12-25T15:30:45Z"
}
```

### Empty Results

```go
func getFilteredResults(c *gin.Context) {
    params := responses.GetPaginationParams(c)
    results := filterResults(params, c.Query("filter"))
    
    if len(results) == 0 {
        responses.EmptyPaginatedResponse(c, params.Page, params.PageSize)
        return
    }
    
    responses.Paginated(c, results, int64(len(results)), params.Page, params.PageSize)
}
```

## File Responses

### File Downloads

```go
func downloadFile(c *gin.Context) {
    fileID := c.Param("id")
    filePath := getFilePath(fileID)
    fileName := getFileName(fileID)
    
    // Download as attachment
    responses.Download(c, filePath, fileName)
}

func exportData(c *gin.Context) {
    data := generateCSVData()
    tempFile := createTempCSV(data)
    defer os.Remove(tempFile)
    
    // File response with custom headers
    responses.File(c, tempFile, "export.csv")
}
```

### Streaming Responses

```go
func streamLogs(c *gin.Context) {
    c.Header("Content-Type", "text/plain")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    
    // Stream data
    c.Stream(func(w io.Writer) bool {
        line := getNextLogLine()
        if line == "" {
            return false // End stream
        }
        
        fmt.Fprintf(w, "%s\n", line)
        return true // Continue streaming
    })
}
```

## Custom Responses

### Custom JSON Response

```go
func customResponse(c *gin.Context) {
    customData := map[string]interface{}{
        "status": "processing",
        "job_id": "job_123",
        "estimated_completion": time.Now().Add(5 * time.Minute),
    }
    
    responses.JSON(c, http.StatusAccepted, customData)
}
```

### Response with Custom Headers

```go
func responseWithHeaders(c *gin.Context) {
    headers := map[string]string{
        "X-API-Version":    "2.0",
        "X-Rate-Limit":     "100",
        "X-Rate-Remaining": "99",
    }
    
    data := getData()
    responses.WithHeaders(c, http.StatusOK, headers, data)
}
```

### Response with Request ID

```go
func trackedResponse(c *gin.Context) {
    data := processRequest()
    
    // Automatically includes X-Request-ID header
    responses.WithRequestID(c, http.StatusOK, data)
}
```

### Redirect Responses

```go
func redirectHandler(c *gin.Context) {
    // Temporary redirect
    responses.Redirect(c, "https://new-location.com")
    
    // Permanent redirect
    responses.PermanentRedirect(c, "https://new-permanent-location.com")
}
```

## Response Standards

### Consistent Timestamp Format

All responses include timestamps in RFC3339 format (UTC):

```json
{
  "timestamp": "2023-12-25T15:30:45Z"
}
```

### Request ID Tracking

When request ID middleware is enabled, all responses include request tracking:

```json
{
  "request_id": "req_abc123"
}
```

### Error Response Consistency

All error responses follow the same structure:

```json
{
  "error": "Human-readable error message",
  "code": "machine_readable_error_code",
  "details": "Additional details (optional)",
  "timestamp": "2023-12-25T15:30:45Z",
  "path": "/api/users/123",
  "request_id": "req_abc123"
}
```

### Success Response Consistency

All success responses include:

```json
{
  "message": "Human-readable success message",
  "data": "Response data (optional)",
  "timestamp": "2023-12-25T15:30:45Z"
}
```

## Best Practices

### 1. Meaningful Messages

```go
// ✅ Good: Descriptive messages
responses.Success(c, "User profile updated successfully", user)
responses.BadRequest(c, "Email address is required")
responses.NotFound(c, "Task with ID 123 not found")

// ❌ Bad: Generic messages
responses.Success(c, "OK", user)
responses.BadRequest(c, "Error")
responses.NotFound(c, "Not found")
```

### 2. Consistent Error Codes

```go
// ✅ Good: Use standard error codes
const (
    ErrUserNotFound     = "user_not_found"
    ErrInvalidEmail     = "invalid_email"
    ErrEmailAlreadyUsed = "email_already_used"
)

func createUser(c *gin.Context) {
    if !isValidEmail(email) {
        responses.Error(c, http.StatusBadRequest, ErrInvalidEmail, "Email format is invalid")
        return
    }
    
    if emailExists(email) {
        responses.Error(c, http.StatusConflict, ErrEmailAlreadyUsed, "Email is already registered")
        return
    }
}
```

### 3. Appropriate HTTP Status Codes

```go
func handleUserOperations(c *gin.Context) {
    switch c.Request.Method {
    case "GET":
        // 200 OK for successful retrieval
        responses.Success(c, "User found", user)
        
    case "POST":
        // 201 Created for successful creation
        responses.Created(c, "User created successfully", user)
        
    case "PUT":
        // 200 OK for successful update
        responses.Success(c, "User updated successfully", user)
        
    case "DELETE":
        // 204 No Content for successful deletion
        responses.NoContent(c)
    }
}
```

### 4. Validation Error Handling

```go
func validateAndCreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        // Parse binding errors into structured format
        var validationErrors []ValidationError
        
        if fieldErrors, ok := err.(validator.ValidationErrors); ok {
            for _, fieldError := range fieldErrors {
                validationErrors = append(validationErrors, ValidationError{
                    Field:   fieldError.Field(),
                    Value:   fmt.Sprintf("%v", fieldError.Value()),
                    Message: getValidationMessage(fieldError),
                })
            }
        }
        
        responses.ValidationError(c, "Validation failed", validationErrors)
        return
    }
    
    // Additional business logic validation
    if !isUniqueEmail(user.Email) {
        responses.Conflict(c, "Email address is already in use")
        return
    }
    
    // Create user...
    responses.Created(c, "User created successfully", user)
}
```

### 5. Pagination Best Practices

```go
func paginateResults(c *gin.Context) {
    // Validate pagination parameters early
    params := responses.GetPaginationParams(c)
    
    if err := responses.ValidatePaginationParams(params.Page, params.PageSize); err != nil {
        responses.HandleError(c, err)
        return
    }
    
    // Apply filters
    filters := extractFilters(c)
    
    // Get total count (for pagination metadata)
    total := countResults(filters)
    
    // Early return for empty results
    if total == 0 {
        responses.EmptyPaginatedResponse(c, params.Page, params.PageSize)
        return
    }
    
    // Fetch paginated data
    results := getResults(params.Offset, params.Limit, filters)
    
    responses.Paginated(c, results, total, params.Page, params.PageSize)
}
```

### 6. Error Context and Debugging

```go
func handleDatabaseOperation(c *gin.Context) {
    user, err := db.GetUser(userID)
    if err != nil {
        // Log detailed error for debugging
        log.Printf("Database error for user %s: %v", userID, err)
        
        // Return appropriate user-facing error
        if errors.Is(err, gorm.ErrRecordNotFound) {
            responses.NotFound(c, "User not found")
        } else {
            responses.DatabaseError(c, "Failed to retrieve user")
        }
        return
    }
    
    responses.Success(c, "User retrieved successfully", user)
}
```

### 7. Response Caching Headers

```go
func getCacheableData(c *gin.Context) {
    data := getStaticData()
    
    // Set caching headers
    c.Header("Cache-Control", "public, max-age=3600") // 1 hour
    c.Header("ETag", generateETag(data))
    
    // Check if client has cached version
    if c.GetHeader("If-None-Match") == generateETag(data) {
        c.Status(http.StatusNotModified)
        return
    }
    
    responses.Success(c, "Data retrieved", data)
}
```

### 8. Conditional Responses

```go
func conditionalResponse(c *gin.Context) {
    format := c.Query("format")
    data := getData()
    
    switch format {
    case "csv":
        c.Header("Content-Type", "text/csv")
        c.Header("Content-Disposition", "attachment; filename=data.csv")
        c.String(http.StatusOK, convertToCSV(data))
        
    case "xml":
        c.Header("Content-Type", "application/xml")
        c.XML(http.StatusOK, data)
        
    default:
        responses.Success(c, "Data retrieved", data)
    }
}
```

### 9. Rate Limiting Information

```go
func rateLimitedHandler(c *gin.Context) {
    // Get rate limit info
    limit, remaining, resetTime := getRateLimitInfo(c.ClientIP())
    
    // Add rate limit headers
    c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
    c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
    c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
    
    if remaining <= 0 {
        responses.TooManyRequests(c, "Rate limit exceeded")
        return
    }
    
    // Process request...
    responses.Success(c, "Request processed", data)
}
```

### 10. Async Operation Responses

```go
func startLongRunningJob(c *gin.Context) {
    jobID := uuid.New().String()
    
    // Start background job
    go processLongRunningJob(jobID)
    
    // Return immediate response with job tracking
    response := map[string]interface{}{
        "job_id": jobID,
        "status": "started",
        "check_url": fmt.Sprintf("/api/jobs/%s", jobID),
        "estimated_duration": "5-10 minutes",
    }
    
    responses.Accepted(c, "Job started successfully", response)
}

func checkJobStatus(c *gin.Context) {
    jobID := c.Param("id")
    status := getJobStatus(jobID)
    
    response := map[string]interface{}{
        "job_id": jobID,
        "status": status.Status,
        "progress": status.Progress,
        "created_at": status.CreatedAt,
    }
    
    switch status.Status {
    case "completed":
        response["result"] = status.Result
        response["completed_at"] = status.CompletedAt
        
    case "failed":
        response["error"] = status.Error
        response["failed_at"] = status.FailedAt
        
    case "running":
        response["estimated_completion"] = status.EstimatedCompletion
    }
    
    responses.Success(c, "Job status retrieved", response)
}
```

## Response Testing

### Testing Success Responses

```go
func TestSuccessResponse(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    router.GET("/test", func(c *gin.Context) {
        data := map[string]string{"test": "value"}
        responses.Success(c, "Test successful", data)
    })
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/test", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response responses.SuccessResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Test successful", response.Message)
    assert.NotNil(t, response.Data)
}
```

### Testing Error Responses

```go
func TestErrorResponse(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    router.GET("/error", func(c *gin.Context) {
        responses.BadRequest(c, "Invalid input")
    })
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/error", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusBadRequest, w.Code)
    
    var response responses.ErrorResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "Invalid input", response.Error)
    assert.Equal(t, "bad_request", response.Code)
}
```

### Testing Pagination

```go
func TestPaginationResponse(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    router.GET("/items", func(c *gin.Context) {
        items := []string{"item1", "item2", "item3"}
        responses.Paginated(c, items, 10, 1, 3)
    })
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/items?page=1&page_size=3", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response responses.PaginationResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, int64(10), response.Total)
    assert.Equal(t, 1, response.Page)
    assert.Equal(t, 3, response.PageSize)
    assert.True(t, response.HasNext)
    assert.False(t, response.HasPrev)
}
```

## Common Response Patterns

### API Health Status

```go
func healthStatus(c *gin.Context) {
    status := checkSystemHealth()
    
    if status.IsHealthy {
        responses.Success(c, "System is healthy", status)
    } else {
        responses.ServiceUnavailable(c, "System is experiencing issues")
    }
}
```

### Bulk Operations

```go
func bulkCreateUsers(c *gin.Context) {
    var users []User
    if err := c.ShouldBindJSON(&users); err != nil {
        responses.BadRequest(c, "Invalid user data")
        return
    }
    
    results := make([]map[string]interface{}, 0)
    errors := make([]map[string]interface{}, 0)
    
    for i, user := range users {
        if err := createUser(user); err != nil {
            errors = append(errors, map[string]interface{}{
                "index": i,
                "user":  user,
                "error": err.Error(),
            })
        } else {
            results = append(results, map[string]interface{}{
                "index": i,
                "user":  user,
            })
        }
    }
    
    response := map[string]interface{}{
        "successful": results,
        "failed":     errors,
        "summary": map[string]int{
            "total":      len(users),
            "successful": len(results),
            "failed":     len(errors),
        },
    }
    
    if len(errors) == 0 {
        responses.Success(c, "All users created successfully", response)
    } else if len(results) == 0 {
        responses.BadRequest(c, "No users were created", response)
    } else {
        responses.SuccessWithStatus(c, http.StatusPartialContent, 
            "Some users created successfully", response)
    }
}
```

### Search Results

```go
func searchItems(c *gin.Context) {
    query := c.Query("q")
    if query == "" {
        responses.BadRequest(c, "Search query is required")
        return
    }
    
    results := performSearch(query)
    
    response := map[string]interface{}{
        "query":        query,
        "total_results": len(results),
        "results":      results,
        "search_time":  "45ms",
    }
    
    if len(results) == 0 {
        responses.Success(c, "No results found", response)
    } else {
        responses.Success(c, fmt.Sprintf("Found %d results", len(results)), response)
    }
}
```

This comprehensive response patterns guide ensures that all your microservices maintain consistent, predictable API responses that are easy to consume and debug.