package request

// PaginationRequest represents common pagination parameters
type PaginationRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1" example:"1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100" example:"20"`
	SortField string `form:"sort_field" binding:"omitempty" example:"created_at"`
	OrderRule string `form:"order_rule" binding:"omitempty,oneof=asc desc" example:"asc"`
}

// GetPage returns the page number with default value of 1
func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetLimit returns the limit with default value of 20 and max of 100
func (p *PaginationRequest) GetLimit() int {
	if p.Limit <= 0 {
		return 20
	}
	if p.Limit > 100 {
		return 100
	}
	return p.Limit
}

// GetOffset calculates the offset for database queries
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

// GetSortField returns the sort field with default value
func (p *PaginationRequest) GetSortField() string {
	if p.SortField == "" {
		return "created_at"
	}
	return p.SortField
}

// GetOrderRule returns the order rule (asc/desc) with default value of asc
func (p *PaginationRequest) GetOrderRule() string {
	if p.OrderRule == "" {
		return "asc"
	}
	return p.OrderRule
}

