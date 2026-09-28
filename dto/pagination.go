package dto

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
}
