package httpx

import (
	"net/http"
	"strconv"
)

// Page size limits shared by every paginated collection (api.md section 6).
const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// PageQuery is the page a collection request asks for: page starts at 1.
type PageQuery struct {
	Page     int
	PageSize int
}

// Offset is the number of records before the requested page.
func (q PageQuery) Offset() int { return (q.Page - 1) * q.PageSize }

// ParsePageQuery reads the page and page_size query parameters, defaulting to
// the first page of DefaultPageSize records.
func ParsePageQuery(r *http.Request) (PageQuery, error) {
	query := r.URL.Query()
	page, pageOK := parseBounded(query.Get("page"), 1, 1, int(^uint(0)>>1))
	size, sizeOK := parseBounded(query.Get("page_size"), DefaultPageSize, 1, MaxPageSize)
	var fields []FieldError
	if !pageOK {
		fields = append(fields, FieldError{Field: "page", Code: "invalid", Message: "page must be a whole number of at least 1"})
	}
	if !sizeOK {
		fields = append(fields, FieldError{Field: "page_size", Code: "invalid", Message: "page_size must be a whole number from 1 to 100"})
	}
	if len(fields) > 0 {
		return PageQuery{}, ErrValidation(fields...)
	}
	return PageQuery{Page: page, PageSize: size}, nil
}

func parseBounded(raw string, fallback, minimum, maximum int) (int, bool) {
	if raw == "" {
		return fallback, true
	}
	n, err := strconv.Atoi(raw)
	return n, err == nil && n >= minimum && n <= maximum
}

// PageResponse is the body of every paginated collection.
type PageResponse[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// NewPageResponse wraps one page of items; a nil slice is sent as [].
func NewPageResponse[T any](items []T, q PageQuery, totalItems int) PageResponse[T] {
	if items == nil {
		items = []T{}
	}
	return PageResponse[T]{
		Items:      items,
		Page:       q.Page,
		PageSize:   q.PageSize,
		TotalItems: totalItems,
		TotalPages: (totalItems + q.PageSize - 1) / q.PageSize,
	}
}
