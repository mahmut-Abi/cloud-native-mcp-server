package kubernetes

import (
	"fmt"
)

// PaginationOptions defines pagination parameters for list operations
type PaginationOptions struct {
	Limit         int    // Maximum items per page
	Continue      string // Continuation token for next page
	FieldSelector string // Field selector for filtering
	LabelSelector string // Label selector for filtering
}

// PaginationResult represents paginated results
type PaginationResult struct {
	Items           []interface{} // Items in current page
	ContinueToken   string        // Token to fetch next page
	RemainingCount  int64         // Remaining items count
	CurrentPageSize int           // Size of current page
	TotalCount      int64         // Total items available
}

// PaginationInfo represents pagination metadata for API responses
type PaginationInfo struct {
	ContinueToken   string `json:"continueToken"`
	RemainingCount  int64  `json:"remainingCount"`
	CurrentPageSize int64  `json:"currentPageSize"`
	HasMore         bool   `json:"hasMore"`
}

// ValidatePaginationOptions validates pagination parameters
func ValidatePaginationOptions(opts *PaginationOptions) error {
	if opts == nil {
		return fmt.Errorf("pagination options cannot be nil")
	}

	if opts.Limit < 1 || opts.Limit > 500 {
		return fmt.Errorf("limit must be between 1 and 500, got %d", opts.Limit)
	}

	return nil
}

// DefaultPaginationOptions returns default pagination options
func DefaultPaginationOptions() *PaginationOptions {
	return &PaginationOptions{
		Limit: 100,
	}
}
