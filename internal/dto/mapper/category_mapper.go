package mapper

import (
	"github.com/Wilson1510/klampis-pim-go/internal/dto/response"
	"github.com/Wilson1510/klampis-pim-go/internal/models"
)

// ToCategoryResponse converts a Category model to CategoryResponse DTO
func ToCategoryResponse(category *models.Category) response.CategoryResponse {
	return response.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// ToCategoryDetailResponse converts a Category model with parent to CategoryDetailResponse DTO
func ToCategoryDetailResponse(category *models.Category) response.CategoryDetailResponse {
	resp := response.CategoryDetailResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	// Include parent if it exists
	if category.Parent != nil {
		parentResp := ToCategoryResponse(category.Parent)
		resp.Parent = &parentResp
	}

	return resp
}

// ToCategoryWithChildrenResponse converts a Category model to CategoryWithChildrenResponse DTO (recursive)
func ToCategoryWithChildrenResponse(category *models.Category) response.CategoryWithChildrenResponse {
	resp := response.CategoryWithChildrenResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
		Children:    []response.CategoryWithChildrenResponse{},
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	// Recursively map children
	if len(category.Children) > 0 {
		for _, child := range category.Children {
			resp.Children = append(resp.Children, ToCategoryWithChildrenResponse(&child))
		}
	}

	return resp
}

// ToSimpleCategoryResponse converts a Category model to SimpleCategoryResponse DTO
func ToSimpleCategoryResponse(category *models.Category) response.SimpleCategoryResponse {
	return response.SimpleCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}
}

// ToCategoryResponseList converts a slice of Category models to a slice of CategoryResponse DTOs
func ToCategoryResponseList(categories []models.Category) []response.CategoryResponse {
	responses := make([]response.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = ToCategoryResponse(&category)
	}
	return responses
}

// ToCategoryDetailResponseList converts a slice of Category models to a slice of CategoryDetailResponse DTOs
func ToCategoryDetailResponseList(categories []models.Category) []response.CategoryDetailResponse {
	responses := make([]response.CategoryDetailResponse, len(categories))
	for i, category := range categories {
		responses[i] = ToCategoryDetailResponse(&category)
	}
	return responses
}

// ToCategoryTreeResponse converts root categories with children to CategoryTreeResponse
func ToCategoryTreeResponse(categories []models.Category) response.CategoryTreeResponse {
	treeCategories := make([]response.CategoryWithChildrenResponse, len(categories))
	for i, category := range categories {
		treeCategories[i] = ToCategoryWithChildrenResponse(&category)
	}

	return response.CategoryTreeResponse{
		Categories: treeCategories,
	}
}

