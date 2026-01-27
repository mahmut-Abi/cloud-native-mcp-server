package common

// PaginationOptions defines pagination parameters for list operations
type PaginationOptions struct {
	Limit         int    // Maximum items per page
	Page          int    // Current page number (for page-based pagination)
	Continue      string // Continuation token for next page (for token-based pagination)
	FieldSelector string // Field selector for filtering
	LabelSelector string // Label selector for filtering
}

// PaginationInfo represents pagination metadata for API responses
// Supports both page-based and token-based pagination
type PaginationInfo struct {
	// Page-based fields
	CurrentPage     int   `json:"currentPage,omitempty"`
	PerPage         int   `json:"perPage,omitempty"`
	TotalCount      int64 `json:"totalCount,omitempty"`
	TotalPages      int   `json:"totalPages,omitempty"`
	HasNextPage     bool  `json:"hasNextPage,omitempty"`
	HasPreviousPage bool  `json:"hasPreviousPage,omitempty"`

	// Token-based fields
	ContinueToken   string `json:"continueToken,omitempty"`
	RemainingCount  int64  `json:"remainingCount,omitempty"`
	CurrentPageSize int64  `json:"currentPageSize,omitempty"`
	HasMore         bool   `json:"hasMore,omitempty"`
}

// ValidatePaginationOptions validates pagination parameters
func ValidatePaginationOptions(opts *PaginationOptions, maxLimit int) error {
	if opts == nil {
		return nil // nil options are valid
	}

	if opts.Limit < 0 {
		return &ValidationError{Field: "limit", Message: "limit must be non-negative"}
	}

	if maxLimit > 0 && opts.Limit > maxLimit {
		return &ValidationError{
			Field:   "limit",
			Message: "limit exceeds maximum allowed value",
			Value:   opts.Limit,
			Max:     maxLimit,
		}
	}

	if opts.Page < 0 {
		return &ValidationError{Field: "page", Message: "page must be non-negative"}
	}

	return nil
}

// DefaultPaginationOptions returns default pagination options
func DefaultPaginationOptions() *PaginationOptions {
	return &PaginationOptions{
		Limit: 100,
		Page:  1,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
	Max     int         `json:"max,omitempty"`
}

func (e *ValidationError) Error() string {
	if e.Max > 0 {
		return e.Message + " (max: " + string(rune(e.Max)) + ")"
	}
	return e.Message
}