package responses

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PaginationResponse represents a paginated API response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
	Timestamp  time.Time   `json:"timestamp"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// PaginationParams represents pagination parameters from request
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Offset   int `json:"offset"`
	Limit    int `json:"limit"`
}

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Paginated sends a paginated response
func Paginated(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
		Timestamp:  time.Now().UTC(),
	}

	c.JSON(http.StatusOK, response)
}

// PaginatedWithMeta sends a paginated response with separate metadata
func PaginatedWithMeta(c *gin.Context, data interface{}, meta PaginationMeta) {
	response := gin.H{
		"data":      data,
		"meta":      meta,
		"timestamp": time.Now().UTC(),
	}

	c.JSON(http.StatusOK, response)
}

// PaginatedWithStatus sends a paginated response with custom status
func PaginatedWithStatus(c *gin.Context, status int, data interface{}, total int64, page, pageSize int) {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
		Timestamp:  time.Now().UTC(),
	}

	c.JSON(status, response)
}

// GetPaginationParams extracts pagination parameters from query string
func GetPaginationParams(c *gin.Context) PaginationParams {
	page := getIntParam(c, "page", DefaultPage)
	pageSize := getIntParam(c, "page_size", DefaultPageSize)

	// Validate and constrain parameters
	if page < 1 {
		page = DefaultPage
	}

	if pageSize < 1 {
		pageSize = DefaultPageSize
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	offset := (page - 1) * pageSize

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
		Limit:    pageSize,
	}
}

// GetPaginationParamsWithDefaults extracts pagination parameters with custom defaults
func GetPaginationParamsWithDefaults(c *gin.Context, defaultPage, defaultPageSize, maxPageSize int) PaginationParams {
	page := getIntParam(c, "page", defaultPage)
	pageSize := getIntParam(c, "page_size", defaultPageSize)

	// Validate and constrain parameters
	if page < 1 {
		page = defaultPage
	}

	if pageSize < 1 {
		pageSize = defaultPageSize
	}

	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	offset := (page - 1) * pageSize

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
		Limit:    pageSize,
	}
}

// CreatePaginationMeta creates pagination metadata
func CreatePaginationMeta(total int64, page, pageSize int) PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return PaginationMeta{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// getIntParam safely extracts integer parameter from query string
func getIntParam(c *gin.Context, key string, defaultValue int) int {
	param := c.Query(key)
	if param == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}

	return value
}

// ValidatePaginationParams validates pagination parameters
func ValidatePaginationParams(page, pageSize int) error {
	if page < 1 {
		return NewBadRequestError("Page must be greater than 0")
	}

	if pageSize < 1 {
		return NewBadRequestError("Page size must be greater than 0")
	}

	if pageSize > MaxPageSize {
		return NewBadRequestError("Page size cannot exceed " + strconv.Itoa(MaxPageSize))
	}

	return nil
}

// CalculateOffset calculates the offset for database queries
func CalculateOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}

// CalculateTotalPages calculates total pages from total records and page size
func CalculateTotalPages(total int64, pageSize int) int {
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

// CursorPaginationResponse represents a cursor-based pagination response
type CursorPaginationResponse struct {
	Data       interface{} `json:"data"`
	NextCursor string      `json:"next_cursor,omitempty"`
	PrevCursor string      `json:"prev_cursor,omitempty"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
	Timestamp  time.Time   `json:"timestamp"`
}

// CursorPaginated sends a cursor-based paginated response
func CursorPaginated(c *gin.Context, data interface{}, nextCursor, prevCursor string, hasNext, hasPrev bool) {
	response := CursorPaginationResponse{
		Data:       data,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
		Timestamp:  time.Now().UTC(),
	}

	c.JSON(http.StatusOK, response)
}

// GetCursorParams extracts cursor pagination parameters
func GetCursorParams(c *gin.Context) (cursor string, limit int) {
	cursor = c.Query("cursor")
	limit = getIntParam(c, "limit", DefaultPageSize)

	if limit > MaxPageSize {
		limit = MaxPageSize
	}

	return cursor, limit
}

// EmptyPaginatedResponse sends an empty paginated response
func EmptyPaginatedResponse(c *gin.Context, page, pageSize int) {
	Paginated(c, []interface{}{}, 0, page, pageSize)
}

// PaginationLinks represents pagination links for HATEOAS
type PaginationLinks struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}

// PaginatedWithLinks sends a paginated response with HATEOAS links
func PaginatedWithLinks(c *gin.Context, data interface{}, total int64, page, pageSize int, baseURL string) {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	links := PaginationLinks{
		Self:  buildURL(baseURL, page, pageSize),
		First: buildURL(baseURL, 1, pageSize),
		Last:  buildURL(baseURL, totalPages, pageSize),
	}

	if page < totalPages {
		links.Next = buildURL(baseURL, page+1, pageSize)
	}

	if page > 1 {
		links.Prev = buildURL(baseURL, page-1, pageSize)
	}

	response := gin.H{
		"data":      data,
		"meta":      CreatePaginationMeta(total, page, pageSize),
		"links":     links,
		"timestamp": time.Now().UTC(),
	}

	c.JSON(http.StatusOK, response)
}

// buildURL builds pagination URL
func buildURL(baseURL string, page, pageSize int) string {
	return baseURL + "?page=" + strconv.Itoa(page) + "&page_size=" + strconv.Itoa(pageSize)
}
