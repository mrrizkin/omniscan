package types

type Response struct {
	Title   string          `json:"title"`
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Debug   string          `json:"debug,omitempty"`
	Data    interface{}     `json:"data"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}

type PaginationMeta struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	PageCount int `json:"page_count"`
}
