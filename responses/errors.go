package responses

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Error     string      `json:"error"`
	Code      string      `json:"code"`
	Details   string      `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Path      string      `json:"path,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Metadata  interface{} `json:"metadata,omitempty"`
}

// APIError represents a structured error
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Details    string
	Metadata   interface{}
}

func (e *APIError) Error() string {
	return e.Message
}

// Common error codes
const (
	ErrCodeBadRequest           = "bad_request"
	ErrCodeUnauthorized         = "unauthorized"
	ErrCodeForbidden            = "forbidden"
	ErrCodeNotFound             = "not_found"
	ErrCodeMethodNotAllowed     = "method_not_allowed"
	ErrCodeConflict             = "conflict"
	ErrCodeUnprocessableEntity  = "unprocessable_entity"
	ErrCodeTooManyRequests      = "too_many_requests"
	ErrCodeInternalError        = "internal_error"
	ErrCodeServiceUnavailable   = "service_unavailable"
	ErrCodeValidationFailed     = "validation_failed"
	ErrCodeDatabaseError        = "database_error"
	ErrCodeExternalServiceError = "external_service_error"
)

// Error sends an error response with the given status, code, and message
func Error(c *gin.Context, status int, code, message string) {
	response := ErrorResponse{
		Error:     message,
		Code:      code,
		Timestamp: time.Now().UTC(),
		Path:      c.Request.URL.Path,
	}

	// Add request ID if available
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			response.RequestID = id
		}
	}

	c.JSON(status, response)
}

// ErrorWithDetails sends an error response with additional details
func ErrorWithDetails(c *gin.Context, status int, code, message, details string) {
	response := ErrorResponse{
		Error:     message,
		Code:      code,
		Details:   details,
		Timestamp: time.Now().UTC(),
		Path:      c.Request.URL.Path,
	}

	// Add request ID if available
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			response.RequestID = id
		}
	}

	c.JSON(status, response)
}

// ErrorWithMetadata sends an error response with metadata
func ErrorWithMetadata(c *gin.Context, status int, code, message string, metadata interface{}) {
	response := ErrorResponse{
		Error:     message,
		Code:      code,
		Timestamp: time.Now().UTC(),
		Path:      c.Request.URL.Path,
		Metadata:  metadata,
	}

	// Add request ID if available
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			response.RequestID = id
		}
	}

	c.JSON(status, response)
}

// BadRequest sends a 400 Bad Request error
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, ErrCodeBadRequest, message)
}

// BadRequestWithDetails sends a 400 Bad Request error with details
func BadRequestWithDetails(c *gin.Context, message, details string) {
	ErrorWithDetails(c, http.StatusBadRequest, ErrCodeBadRequest, message, details)
}

// Unauthorized sends a 401 Unauthorized error
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, ErrCodeUnauthorized, message)
}

// Forbidden sends a 403 Forbidden error
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, ErrCodeForbidden, message)
}

// NotFound sends a 404 Not Found error
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, ErrCodeNotFound, message)
}

// MethodNotAllowed sends a 405 Method Not Allowed error
func MethodNotAllowed(c *gin.Context, message string) {
	Error(c, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, message)
}

// Conflict sends a 409 Conflict error
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, ErrCodeConflict, message)
}

// UnprocessableEntity sends a 422 Unprocessable Entity error
func UnprocessableEntity(c *gin.Context, message string) {
	Error(c, http.StatusUnprocessableEntity, ErrCodeUnprocessableEntity, message)
}

// TooManyRequests sends a 429 Too Many Requests error
func TooManyRequests(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, ErrCodeTooManyRequests, message)
}

// InternalError sends a 500 Internal Server Error
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, ErrCodeInternalError, message)
}

// ServiceUnavailable sends a 503 Service Unavailable error
func ServiceUnavailable(c *gin.Context, message string) {
	Error(c, http.StatusServiceUnavailable, ErrCodeServiceUnavailable, message)
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, message string, validationErrors interface{}) {
	ErrorWithMetadata(c, http.StatusBadRequest, ErrCodeValidationFailed, message, validationErrors)
}

// DatabaseError sends a database error response
func DatabaseError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, ErrCodeDatabaseError, message)
}

// ExternalServiceError sends an external service error response
func ExternalServiceError(c *gin.Context, message string) {
	Error(c, http.StatusBadGateway, ErrCodeExternalServiceError, message)
}

// HandleAPIError handles APIError types
func HandleAPIError(c *gin.Context, err *APIError) {
	response := ErrorResponse{
		Error:     err.Message,
		Code:      err.Code,
		Details:   err.Details,
		Timestamp: time.Now().UTC(),
		Path:      c.Request.URL.Path,
		Metadata:  err.Metadata,
	}

	// Add request ID if available
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			response.RequestID = id
		}
	}

	c.JSON(err.StatusCode, response)
}

// HandleError handles generic errors
func HandleError(c *gin.Context, err error) {
	if apiErr, ok := err.(*APIError); ok {
		HandleAPIError(c, apiErr)
		return
	}

	// Default to internal server error
	InternalError(c, err.Error())
}

// NewAPIError creates a new APIError
func NewAPIError(statusCode int, code, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// NewAPIErrorWithDetails creates a new APIError with details
func NewAPIErrorWithDetails(statusCode int, code, message, details string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    details,
	}
}

// NewAPIErrorWithMetadata creates a new APIError with metadata
func NewAPIErrorWithMetadata(statusCode int, code, message string, metadata interface{}) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Metadata:   metadata,
	}
}

// Common error constructors
func NewBadRequestError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, ErrCodeBadRequest, message)
}

func NewUnauthorizedError(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, ErrCodeUnauthorized, message)
}

func NewForbiddenError(message string) *APIError {
	return NewAPIError(http.StatusForbidden, ErrCodeForbidden, message)
}

func NewNotFoundError(message string) *APIError {
	return NewAPIError(http.StatusNotFound, ErrCodeNotFound, message)
}

func NewConflictError(message string) *APIError {
	return NewAPIError(http.StatusConflict, ErrCodeConflict, message)
}

func NewValidationError(message string, validationErrors interface{}) *APIError {
	return NewAPIErrorWithMetadata(http.StatusBadRequest, ErrCodeValidationFailed, message, validationErrors)
}

func NewInternalError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, ErrCodeInternalError, message)
}
