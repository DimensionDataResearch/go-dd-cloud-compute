package compute

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

// NextPage creates a PagingInfo for the next page of results.
func (page *PagedResult) NextPage() *PagingInfo {
	return &PagingInfo{
		PageNumber: page.PageNumber + 1,
		PageSize:   page.PageSize,
	}
}

// PagingInfo contains the paging configuration for a compute API operation.
type PagingInfo struct {
	PageNumber int
	PageSize   int
}

// DefaultPaging creates PagingInfo with default settings (page 1, 50 records per page).
func DefaultPaging() *PagingInfo {
	return &PagingInfo{
		PageNumber: 1,
		PageSize:   50,
	}
}

// EnsurePaging always returns a paging configuration (if the supplied PagingInfo is nil, it returns the default configuration).
func EnsurePaging(paging *PagingInfo) *PagingInfo {
	if paging != nil {
		return paging
	}

	return DefaultPaging()
}

func (pagingInfo *PagingInfo) ensureValidPageSize() {
	if pagingInfo.PageSize < 5 {
		pagingInfo.PageSize = 5
	}
}

// First configures the PagingInfo for the first page of results.
func (pagingInfo *PagingInfo) First() {
	pagingInfo.ensureValidPageSize()

	pagingInfo.PageNumber = 1
}

// Next configures the PagingInfo for the next page of results.
func (pagingInfo *PagingInfo) Next() {
	pagingInfo.ensureValidPageSize()

	pagingInfo.PageNumber++
}
