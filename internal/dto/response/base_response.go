package response

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   interface{} `json:"error" example:"null"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool        `json:"success" example:"false"`
	Data    interface{} `json:"data" example:"null"`
	Error   ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string      `json:"code" example:"VALIDATION_ERROR"`
	Message string      `json:"message" example:"Invalid input data"`
	Details interface{} `json:"details,omitempty"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page  int   `json:"page" example:"1"`
	Limit int   `json:"limit" example:"20"`
	Total int64 `json:"total" example:"157"`
	Pages int   `json:"pages" example:"8"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	}
}

// NewSuccessResponseWithMeta creates a new success response with metadata
func NewSuccessResponseWithMeta(data interface{}, meta interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
		Error:   nil,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code string, message string, details interface{}) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Data:    nil,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// NewPaginationMeta creates pagination metadata
func NewPaginationMeta(page, limit int, total int64) PaginationMeta {
	pages := int(total) / limit
	if int(total)%limit != 0 {
		pages++
	}

	return PaginationMeta{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: pages,
	}
}

