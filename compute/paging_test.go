package compute

import (
	"testing"
)

// Verify that a PagedResult correctly determines that a page is not the last page.
func TestPagedResult_IsLastPage_PageSize_2_FirstPage_2of6Records(test *testing.T) {
	expect := expect(test)

	result := &PagedResult{
		PageNumber: 1,
		PageSize:   2,
		PageCount:  2,
		TotalCount: 6,
	}

	expect.IsFalse("PagedResult.IsLastPage",
		result.IsLastPage(),
	)
}

// Verify that a PagedResult correctly determines that a page is not the last page.
func TestPagedResult_IsLastPage_PageSize_2_SecondLastPage_3of6Records(test *testing.T) {
	expect := expect(test)

	result := &PagedResult{
		PageNumber: 2,
		PageSize:   2,
		PageCount:  1,
		TotalCount: 6,
	}

	expect.IsFalse("PagedResult.IsLastPage",
		result.IsLastPage(),
	)
}

// Verify that a PagedResult correctly determines that a page is the last page.
func TestPagedResult_IsLastPage_PageSize_2_LastPage_6of6Records(test *testing.T) {
	expect := expect(test)

	result := &PagedResult{
		PageNumber: 3,
		PageSize:   2,
		PageCount:  2,
		TotalCount: 6,
	}

	expect.IsTrue("PagedResult.IsLastPage",
		result.IsLastPage(),
	)
}
