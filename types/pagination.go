package types

// PaginationRequest represents pagination parameters from client
type PaginationRequest struct {
	Page      int    `json:"page" form:"page" binding:"min=1"`
	PageSize  int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order" binding:"oneof=asc desc"`
}

// PaginationResponse represents paginated response metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// CursorPaginationRequest represents cursor-based pagination parameters
type CursorPaginationRequest struct {
	Cursor    string `json:"cursor" form:"cursor"`
	Limit     int    `json:"limit" form:"limit" binding:"min=1,max=100"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order" binding:"oneof=asc desc"`
}

// CursorPaginationResponse represents cursor-based pagination metadata
type CursorPaginationResponse struct {
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
	Limit      int    `json:"limit"`
}

// PaginatedResult represents a paginated result set
type PaginatedResult[T any] struct {
	Data       []T                `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

// CursorPaginatedResult represents a cursor-paginated result set
type CursorPaginatedResult[T any] struct {
	Data       []T                      `json:"data"`
	Pagination CursorPaginationResponse `json:"pagination"`
}

// FilterOptions represents common filtering options
type FilterOptions struct {
	Search   string            `json:"search" form:"search"`
	Status   string            `json:"status" form:"status"`
	Category string            `json:"category" form:"category"`
	Tags     []string          `json:"tags" form:"tags"`
	DateFrom string            `json:"date_from" form:"date_from"`
	DateTo   string            `json:"date_to" form:"date_to"`
	UserID   *uint             `json:"user_id" form:"user_id"`
	Custom   map[string]string `json:"custom" form:"custom"`
}

// QueryOptions combines pagination, sorting, and filtering
type QueryOptions struct {
	PaginationRequest
	FilterOptions
}

// CursorQueryOptions combines cursor pagination, sorting, and filtering
type CursorQueryOptions struct {
	CursorPaginationRequest
	FilterOptions
}

// DefaultPaginationRequest returns default pagination parameters
func DefaultPaginationRequest() PaginationRequest {
	return PaginationRequest{
		Page:      1,
		PageSize:  10,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// DefaultCursorPaginationRequest returns default cursor pagination parameters
func DefaultCursorPaginationRequest() CursorPaginationRequest {
	return CursorPaginationRequest{
		Limit:     10,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// Validate validates pagination request
func (p *PaginationRequest) Validate() error {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	if p.SortOrder != "asc" && p.SortOrder != "desc" {
		p.SortOrder = "desc"
	}
	return nil
}

// Validate validates cursor pagination request
func (c *CursorPaginationRequest) Validate() error {
	if c.Limit < 1 {
		c.Limit = 10
	}
	if c.Limit > 100 {
		c.Limit = 100
	}
	if c.SortOrder != "asc" && c.SortOrder != "desc" {
		c.SortOrder = "desc"
	}
	return nil
}

// GetOffset calculates the offset for database queries
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit returns the page size as limit
func (p *PaginationRequest) GetLimit() int {
	return p.PageSize
}
