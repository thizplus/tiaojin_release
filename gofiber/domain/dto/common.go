package dto

type ListRequest struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=1,max=100"`
}

type ListResponse struct {
	Total   int64 `json:"total"`
	HasMore bool  `json:"has_more"`
}

type PaginationMeta struct {
	Total       int64 `json:"total"`
	HasMore     bool  `json:"has_more"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	PageSize    int   `json:"page_size"`
	Limit       int   `json:"limit"`  // Alias for PageSize for backward compatibility
	Offset      int   `json:"offset"`
}

type SearchRequest struct {
	Search *string `json:"search,omitempty"`
	Offset int     `json:"offset" validate:"min=0"`
	Limit  int     `json:"limit" validate:"min=1,max=100"`
}

type StatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type FileUploadRequest struct {
	FileData []byte `json:"file_data" validate:"required"`
	FileName string `json:"file_name" validate:"required"`
	FileType string `json:"file_type" validate:"required"`
}

// FileResponse moved to file.go to avoid duplication

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
}