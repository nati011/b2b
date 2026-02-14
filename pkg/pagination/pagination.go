package pagination

import (
	"net/http"
	"strconv"
)

// PageRequest represents pagination parameters
type PageRequest struct {
	Page   int // 1-based page number
	Limit  int // items per page
	Offset int // calculated from page and limit
}

// NewPageRequest creates a new page request and computes the offset.
// Callers are responsible for validating page and limit.
func NewPageRequest(page, limit int) PageRequest {
	return PageRequest{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

// PageResult represents a paginated result with metadata
type PageResult[T any] struct {
	Items      []T  `json:"items"`
	Total      int  `json:"total"`
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// NewPageResult creates a new paginated result
func NewPageResult[T any](items []T, total int, pageReq PageRequest) PageResult[T] {
	totalPages := (total + pageReq.Limit - 1) / pageReq.Limit
	if totalPages == 0 {
		totalPages = 1
	}

	return PageResult[T]{
		Items:      items,
		Total:      total,
		Page:       pageReq.Page,
		Limit:      pageReq.Limit,
		TotalPages: totalPages,
		HasNext:    pageReq.Page < totalPages,
		HasPrev:    pageReq.Page > 1,
	}
}

// FromRequest parses pagination parameters from an HTTP request.
// It extracts 'page' and 'limit' query parameters with defaults:
// - page: defaults to 1 if missing or invalid
// - limit: defaults to 10 if missing or invalid, max 100
func FromRequest(r *http.Request) PageRequest {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return NewPageRequest(page, limit)
}
