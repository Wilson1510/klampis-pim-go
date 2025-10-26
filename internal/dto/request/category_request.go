package request

// CreateCategoryRequest represents the request body for creating a new category
type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100" example:"Electronics"`
	Description string  `json:"description" binding:"omitempty" example:"All electronic products"`
	ParentID    *uint   `json:"parent_id" binding:"omitempty" example:"1"`
}

// UpdateCategoryRequest represents the request body for updating an existing category
type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100" example:"Electronics"`
	Description *string `json:"description" binding:"omitempty" example:"All electronic products"`
	ParentID    *uint   `json:"parent_id" binding:"omitempty" example:"1"`
}

// GetCategoriesRequest represents query parameters for listing categories
type GetCategoriesRequest struct {
	PaginationRequest
	Name     string `form:"name" binding:"omitempty,max=100" example:"Electronics"`
	ParentID *uint  `form:"parent_id" binding:"omitempty" example:"1"`
	// Include root categories only (categories without parent)
	RootOnly bool `form:"root_only" binding:"omitempty" example:"false"`
}

