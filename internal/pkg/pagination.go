package pkg

import "article/internal/models"

// Pagination represents the parameters used to navigate through paginated datasets.
type Pagination struct {
	Page  int
	Limit int
}

// Normalize ensures that the pagination parameters are within valid bounds.
// It mutates the struct fields in-place.
//
// Behavior:
// - If Page < 1, it will be set to 1 (1-based indexing).
// - If Limit <= 0, it will be set to a default value of 10.
//
// This method must be called before using Offset or building pagination metadata
// to avoid invalid calculations such as negative offsets or division by zero.
func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
}

// MakeOffset calculates the database offset based on the given page and limit.
//
// Formula:
//
//	offset = (page - 1) * limit
func (p *Pagination) MakeOffset() int {
	offset := (p.Page - 1) * p.Limit
	return offset
}

// MakeMeta computes metadata for the current result set based on the total record count.
// It calculates the total number of pages and determines the existence of adjacent pages.
//
// Parameters:
//   - total: The total number of records available in the dataset.
func (p *Pagination) MakeMeta(total int) models.PaginationMeta {
	totalPages := (total + p.Limit - 1) / p.Limit

	hasNext := p.Page < totalPages
	hasPrev := p.Page > 1

	meta := models.PaginationMeta{
		CurrentPage: p.Page,
		Limit:       p.Limit,
		TotalItems:  total,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}

	return meta
}
