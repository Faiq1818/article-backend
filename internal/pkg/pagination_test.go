package pkg

import (
	"article/internal/models"
	"testing"
)

func TestPagination_Normalize(t *testing.T) {
	tests := []struct {
		name     string
		input    Pagination
		expected Pagination
	}{
		{
			name:     "valid values",
			input:    Pagination{Page: 2, Limit: 5},
			expected: Pagination{Page: 2, Limit: 5},
		},
		{
			name:     "page less than 1",
			input:    Pagination{Page: 0, Limit: 5},
			expected: Pagination{Page: 1, Limit: 5},
		},
		{
			name:     "limit less or equal 0",
			input:    Pagination{Page: 2, Limit: 0},
			expected: Pagination{Page: 2, Limit: 10},
		},
		{
			name:     "both invalid",
			input:    Pagination{Page: 0, Limit: -1},
			expected: Pagination{Page: 1, Limit: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.input
			p.Normalize()

			if p.Page != tt.expected.Page {
				t.Errorf("expected Page %d, got %d", tt.expected.Page, p.Page)
			}

			if p.Limit != tt.expected.Limit {
				t.Errorf("expected Limit %d, got %d", tt.expected.Limit, p.Limit)
			}
		})
	}
}

func TestPagination_MakeOffset(t *testing.T) {
	tests := []struct {
		name     string
		input    Pagination
		expected int
	}{
		{
			name:     "page 1",
			input:    Pagination{Page: 1, Limit: 10},
			expected: 0,
		},
		{
			name:     "page 2",
			input:    Pagination{Page: 2, Limit: 10},
			expected: 10,
		},
		{
			name:     "page 3",
			input:    Pagination{Page: 3, Limit: 5},
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := tt.input.MakeOffset()

			if offset != tt.expected {
				t.Errorf("expected offset %d, got %d", tt.expected, offset)
			}
		})
	}
}

func TestPagination_MakeMeta(t *testing.T) {
	tests := []struct {
		name     string
		input    Pagination
		total    int
		expected models.PaginationMeta
	}{
		{
			name:  "first page",
			input: Pagination{Page: 1, Limit: 10},
			total: 25,
			expected: models.PaginationMeta{
				CurrentPage: 1,
				Limit:       10,
				TotalItems:  25,
				TotalPages:  3,
				HasNext:     true,
				HasPrev:     false,
			},
		},
		{
			name:  "middle page",
			input: Pagination{Page: 2, Limit: 10},
			total: 25,
			expected: models.PaginationMeta{
				CurrentPage: 2,
				Limit:       10,
				TotalItems:  25,
				TotalPages:  3,
				HasNext:     true,
				HasPrev:     true,
			},
		},
		{
			name:  "last page",
			input: Pagination{Page: 3, Limit: 10},
			total: 25,
			expected: models.PaginationMeta{
				CurrentPage: 3,
				Limit:       10,
				TotalItems:  25,
				TotalPages:  3,
				HasNext:     false,
				HasPrev:     true,
			},
		},
		{
			name:  "single page",
			input: Pagination{Page: 1, Limit: 10},
			total: 5,
			expected: models.PaginationMeta{
				CurrentPage: 1,
				Limit:       10,
				TotalItems:  5,
				TotalPages:  1,
				HasNext:     false,
				HasPrev:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := tt.input.MakeMeta(tt.total)

			if meta.CurrentPage != tt.expected.CurrentPage {
				t.Errorf("expected CurrentPage %d, got %d", tt.expected.CurrentPage, meta.CurrentPage)
			}
			if meta.Limit != tt.expected.Limit {
				t.Errorf("expected Limit %d, got %d", tt.expected.Limit, meta.Limit)
			}
			if meta.TotalItems != tt.expected.TotalItems {
				t.Errorf("expected TotalItems %d, got %d", tt.expected.TotalItems, meta.TotalItems)
			}
			if meta.TotalPages != tt.expected.TotalPages {
				t.Errorf("expected TotalPages %d, got %d", tt.expected.TotalPages, meta.TotalPages)
			}
			if meta.HasNext != tt.expected.HasNext {
				t.Errorf("expected HasNext %v, got %v", tt.expected.HasNext, meta.HasNext)
			}
			if meta.HasPrev != tt.expected.HasPrev {
				t.Errorf("expected HasPrev %v, got %v", tt.expected.HasPrev, meta.HasPrev)
			}
		})
	}
}
