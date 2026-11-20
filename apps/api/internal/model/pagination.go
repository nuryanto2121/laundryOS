package model

// PaginationMeta contains page navigation metadata.
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse wraps a slice of data with pagination metadata.
type PaginatedResponse struct {
	Data any            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// CalculatePagination computes SQL offset and populates PaginationMeta.
func CalculatePagination(page, limit int, totalItems int64) (offset int, meta PaginationMeta) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}
	offset = (page - 1) * limit
	meta = PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
	return offset, meta
}
