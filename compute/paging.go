package compute

import "fmt"

// PagedResult represents the common fields for all paged results from the compute API.
type PagedResult struct {
	// The current page number.
	PageNumber int `json:"pageNumber"`

	// The number of items in the current page of results.
	PageCount int `json:"pageCount"`

	// The total number of results that match the requested filter criteria (if any).
	TotalCount int `json:"totalCount"`

	// The maximum number of results per page.
	PageSize int `json:"pageSize"`
}

// IsEmpty determines whether the page contains no results.
func (page *PagedResult) IsEmpty() bool {
	return page.PageCount == 0
}

// NextPage creates a Paging for the next page of results.
func (page *PagedResult) NextPage() *Paging {
	return &Paging{
		PageNumber: page.PageNumber + 1,
		PageSize:   page.PageSize,
	}
}

// Paging contains the paging configuration for a compute API operation.
type Paging struct {
	PageNumber int
	PageSize   int
}

// DefaultPaging creates Paging with default settings (page 1, 50 records per page).
func DefaultPaging() *Paging {
	return &Paging{
		PageNumber: 1,
		PageSize:   50,
	}
}

// EnsurePaging always returns a paging configuration (if the supplied Paging is nil, it returns the default configuration).
func (paging *Paging) EnsurePaging() *Paging {
	if paging != nil {
		return paging
	}

	return DefaultPaging()
}

func (paging *Paging) ensureValidPageSize() {
	if paging.PageSize < 5 {
		paging.PageSize = 5
	}
}

func (paging *Paging) toQueryParameters() string {
	return fmt.Sprintf("pageNumber=%d&pageSize=%d", paging.PageNumber, paging.PageSize)
}

// First configures the Paging for the first page of results.
func (paging *Paging) First() {
	paging.ensureValidPageSize()

	paging.PageNumber = 1
}

// Next configures the Paging for the next page of results.
func (paging *Paging) Next() {
	paging.ensureValidPageSize()

	paging.PageNumber++
}
