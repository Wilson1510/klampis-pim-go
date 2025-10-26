package response

import "time"

// CategoryResponse represents the basic category response
type CategoryResponse struct {
	ID          uint      `json:"id" example:"1"`
	Name        string    `json:"name" example:"Electronics"`
	Slug        string    `json:"slug" example:"electronics"`
	Description string    `json:"description" example:"All electronic products"`
	ParentID    *uint     `json:"parent_id" example:"1"`
	CreatedAt   time.Time `json:"created_at" example:"2025-10-17T10:30:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-10-17T10:30:00Z"`
}

// CategoryDetailResponse represents a category with its parent information
type CategoryDetailResponse struct {
	ID          uint              `json:"id" example:"1"`
	Name        string            `json:"name" example:"Laptops"`
	Slug        string            `json:"slug" example:"laptops"`
	Description string            `json:"description" example:"Laptop computers"`
	ParentID    *uint             `json:"parent_id" example:"1"`
	Parent      *CategoryResponse `json:"parent,omitempty"`
	CreatedAt   time.Time         `json:"created_at" example:"2025-10-17T10:30:00Z"`
	UpdatedAt   time.Time         `json:"updated_at" example:"2025-10-17T10:30:00Z"`
}

// CategoryWithChildrenResponse represents a category with its children (hierarchical)
type CategoryWithChildrenResponse struct {
	ID          uint                           `json:"id" example:"1"`
	Name        string                         `json:"name" example:"Electronics"`
	Slug        string                         `json:"slug" example:"electronics"`
	Description string                         `json:"description" example:"All electronic products"`
	ParentID    *uint                          `json:"parent_id" example:"null"`
	Children    []CategoryWithChildrenResponse `json:"children,omitempty"`
	CreatedAt   time.Time                      `json:"created_at" example:"2025-10-17T10:30:00Z"`
	UpdatedAt   time.Time                      `json:"updated_at" example:"2025-10-17T10:30:00Z"`
}

// CategoryTreeResponse represents the full category tree structure
type CategoryTreeResponse struct {
	Categories []CategoryWithChildrenResponse `json:"categories"`
}

// CategoryWithProductCountResponse represents a category with product count
type CategoryWithProductCountResponse struct {
	ID           uint      `json:"id" example:"1"`
	Name         string    `json:"name" example:"Electronics"`
	Slug         string    `json:"slug" example:"electronics"`
	Description  string    `json:"description" example:"All electronic products"`
	ParentID     *uint     `json:"parent_id" example:"1"`
	ProductCount int64     `json:"product_count" example:"45"`
	CreatedAt    time.Time `json:"created_at" example:"2025-10-17T10:30:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2025-10-17T10:30:00Z"`
}

// SimpleCategoryResponse represents minimal category info (for nested responses)
type SimpleCategoryResponse struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"Electronics"`
	Slug string `json:"slug" example:"electronics"`
}

