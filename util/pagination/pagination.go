package pagination

import "math"

// Pagination 用於定義分頁資訊的結構
type Pagination struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
}

func NewPagination(currentPage, pageSize int, totalRecords int64) Pagination {
	return Pagination{
		CurrentPage:  currentPage,
		PageSize:     pageSize,
		TotalRecords: totalRecords,
		TotalPages:   int(math.Ceil(float64(totalRecords) / float64(pageSize))),
	}
}
