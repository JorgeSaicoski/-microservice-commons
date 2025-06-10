package types

import "time"

// APIResponse represents a standard API response structure
type APIResponse[T any] struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
	Data      T         `json:"data,omitempty"`
	Error     *APIError `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id,omitempty"`
}

// APIError represents a structured API error
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details string      `json:"details,omitempty"`
	Field   string      `json:"field,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

// SuccessResponse represents a successful API response
type SuccessResponse[T any] struct {
	Message   string    `json:"message"`
	Data      T         `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

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

// ValidationErrorResponse represents validation errors
type ValidationErrorResponse struct {
	ErrorResponse
	Errors []ValidationError `json:"errors"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// ListResponse represents a list response with metadata
type ListResponse[T any] struct {
	Data  []T              `json:"data"`
	Meta  ResponseMetadata `json:"meta"`
	Links ResponseLinks    `json:"links,omitempty"`
}

// ResponseMetadata represents response metadata
type ResponseMetadata struct {
	Total      int64     `json:"total,omitempty"`
	Count      int       `json:"count"`
	Page       int       `json:"page,omitempty"`
	PageSize   int       `json:"page_size,omitempty"`
	TotalPages int       `json:"total_pages,omitempty"`
	HasNext    bool      `json:"has_next,omitempty"`
	HasPrev    bool      `json:"has_prev,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// ResponseLinks represents HATEOAS links
type ResponseLinks struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime,omitempty"`
	Checks    map[string]HealthCheck `json:"checks,omitempty"`
}

// HealthCheck represents a single health check
type HealthCheck struct {
	Status      string                 `json:"status"`
	Message     string                 `json:"message,omitempty"`
	LastChecked time.Time              `json:"last_checked"`
	Duration    string                 `json:"duration"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         *UserInfo `json:"user,omitempty"`
}

// UserInfo represents basic user information
type UserInfo struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Name     string   `json:"name,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	Avatar   string   `json:"avatar,omitempty"`
}

// FileUploadResponse represents file upload response
type FileUploadResponse struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	URL      string `json:"url"`
	Path     string `json:"path,omitempty"`
	Checksum string `json:"checksum,omitempty"`
}

// BatchResponse represents a batch operation response
type BatchResponse[T any] struct {
	Successful []T                    `json:"successful"`
	Failed     []BatchError           `json:"failed"`
	Summary    BatchSummary           `json:"summary"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// BatchError represents a single error in batch operation
type BatchError struct {
	Index   int         `json:"index"`
	ID      interface{} `json:"id,omitempty"`
	Error   string      `json:"error"`
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

// BatchSummary represents summary of batch operation
type BatchSummary struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

// SearchResponse represents search results
type SearchResponse[T any] struct {
	Results     []T                `json:"results"`
	Total       int64              `json:"total"`
	Query       string             `json:"query"`
	Filters     map[string]string  `json:"filters,omitempty"`
	Facets      map[string]Facet   `json:"facets,omitempty"`
	Suggestions []string           `json:"suggestions,omitempty"`
	Pagination  PaginationResponse `json:"pagination"`
}

// Facet represents search facets
type Facet struct {
	Values []FacetValue `json:"values"`
}

// FacetValue represents a facet value with count
type FacetValue struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// ExportResponse represents data export response
type ExportResponse struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	Format      string                 `json:"format"`
	URL         string                 `json:"url,omitempty"`
	Progress    int                    `json:"progress"`
	TotalRows   int64                  `json:"total_rows"`
	CreatedAt   time.Time              `json:"created_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	ExpiresAt   time.Time              `json:"expires_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ImportResponse represents data import response
type ImportResponse struct {
	ID            string                 `json:"id"`
	Status        string                 `json:"status"`
	Progress      int                    `json:"progress"`
	TotalRows     int64                  `json:"total_rows"`
	ProcessedRows int64                  `json:"processed_rows"`
	SuccessRows   int64                  `json:"success_rows"`
	ErrorRows     int64                  `json:"error_rows"`
	Errors        []ImportError          `json:"errors,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// ImportError represents an import error
type ImportError struct {
	Row    int    `json:"row"`
	Column string `json:"column,omitempty"`
	Value  string `json:"value,omitempty"`
	Error  string `json:"error"`
	Code   string `json:"code,omitempty"`
}

// WebhookResponse represents webhook delivery response
type WebhookResponse struct {
	ID          string                  `json:"id"`
	Event       string                  `json:"event"`
	Status      string                  `json:"status"`
	Attempts    int                     `json:"attempts"`
	NextRetry   *time.Time              `json:"next_retry,omitempty"`
	LastAttempt time.Time               `json:"last_attempt"`
	Response    WebhookDeliveryResponse `json:"response,omitempty"`
}

// WebhookDeliveryResponse represents webhook delivery details
type WebhookDeliveryResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
	Duration   string            `json:"duration"`
	Error      string            `json:"error,omitempty"`
}

// NotificationResponse represents notification response
type NotificationResponse struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Read      bool                   `json:"read"`
	CreatedAt time.Time              `json:"created_at"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Actions   []NotificationAction   `json:"actions,omitempty"`
}

// NotificationAction represents an action that can be taken on a notification
type NotificationAction struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	URL   string `json:"url,omitempty"`
	Type  string `json:"type"` // button, link, etc.
}

// StatisticsResponse represents statistics/analytics response
type StatisticsResponse struct {
	Period    string                 `json:"period"`
	StartDate time.Time              `json:"start_date"`
	EndDate   time.Time              `json:"end_date"`
	Metrics   map[string]MetricValue `json:"metrics"`
	Charts    map[string]ChartData   `json:"charts,omitempty"`
	Summary   map[string]interface{} `json:"summary,omitempty"`
}

// MetricValue represents a single metric value
type MetricValue struct {
	Current  interface{} `json:"current"`
	Previous interface{} `json:"previous,omitempty"`
	Change   float64     `json:"change,omitempty"`
	Trend    string      `json:"trend,omitempty"` // up, down, stable
}

// ChartData represents chart data points
type ChartData struct {
	Labels []string      `json:"labels"`
	Data   []interface{} `json:"data"`
	Type   string        `json:"type"` // line, bar, pie, etc.
}
