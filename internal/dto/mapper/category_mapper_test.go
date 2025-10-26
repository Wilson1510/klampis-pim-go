package mapper

import (
	"testing"
	"time"

	"github.com/Wilson1510/klampis-pim-go/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestToCategoryResponse(t *testing.T) {
	// Setup
	now := time.Now()
	parentID := uint(1)
	category := &models.Category{
		Base: models.Base{
			ID:        2,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        "Laptops",
		Slug:        "laptops",
		Description: "Laptop computers",
		ParentID:    &parentID,
	}

	// Execute
	response := ToCategoryResponse(category)

	// Assert
	assert.Equal(t, uint(2), response.ID)
	assert.Equal(t, "Laptops", response.Name)
	assert.Equal(t, "laptops", response.Slug)
	assert.Equal(t, "Laptop computers", response.Description)
	assert.Equal(t, &parentID, response.ParentID)
	assert.Equal(t, now, response.CreatedAt)
	assert.Equal(t, now, response.UpdatedAt)
}

func TestToCategoryDetailResponse(t *testing.T) {
	// Setup
	now := time.Now()
	parentID := uint(1)

	parent := &models.Category{
		Base: models.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        "Electronics",
		Slug:        "electronics",
		Description: "All electronics",
		ParentID:    nil,
	}

	category := &models.Category{
		Base: models.Base{
			ID:        2,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        "Laptops",
		Slug:        "laptops",
		Description: "Laptop computers",
		ParentID:    &parentID,
		Parent:      parent,
	}

	// Execute
	response := ToCategoryDetailResponse(category)

	// Assert
	assert.Equal(t, uint(2), response.ID)
	assert.Equal(t, "Laptops", response.Name)
	assert.NotNil(t, response.Parent)
	assert.Equal(t, uint(1), response.Parent.ID)
	assert.Equal(t, "Electronics", response.Parent.Name)
}

func TestToCategoryDetailResponseWithoutParent(t *testing.T) {
	// Setup
	now := time.Now()

	category := &models.Category{
		Base: models.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        "Electronics",
		Slug:        "electronics",
		Description: "All electronics",
		ParentID:    nil,
		Parent:      nil,
	}

	// Execute
	response := ToCategoryDetailResponse(category)

	// Assert
	assert.Equal(t, uint(1), response.ID)
	assert.Nil(t, response.Parent)
}

func TestToCategoryWithChildrenResponse(t *testing.T) {
	// Setup
	now := time.Now()
	parentID := uint(1)
	childParentID := uint(2)

	grandchild := models.Category{
		Base: models.Base{
			ID:        3,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        "Gaming Laptops",
		Slug:        "gaming-laptops",
		Description: "High-performance gaming laptops",
		ParentID:    &childParentID,
		Children:    []models.Category{},
	}

	child := models.Category{
		Base: models.Base{
			ID:        2,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        "Laptops",
		Slug:        "laptops",
		Description: "Laptop computers",
		ParentID:    &parentID,
		Children:    []models.Category{grandchild},
	}

	parent := &models.Category{
		Base: models.Base{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        "Electronics",
		Slug:        "electronics",
		Description: "All electronics",
		ParentID:    nil,
		Children:    []models.Category{child},
	}

	// Execute
	response := ToCategoryWithChildrenResponse(parent)

	// Assert
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "Electronics", response.Name)
	assert.Len(t, response.Children, 1)
	assert.Equal(t, "Laptops", response.Children[0].Name)
	assert.Len(t, response.Children[0].Children, 1)
	assert.Equal(t, "Gaming Laptops", response.Children[0].Children[0].Name)
}

func TestToSimpleCategoryResponse(t *testing.T) {
	// Setup
	category := &models.Category{
		Base: models.Base{
			ID: 1,
		},
		Name: "Electronics",
		Slug: "electronics",
	}

	// Execute
	response := ToSimpleCategoryResponse(category)

	// Assert
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "Electronics", response.Name)
	assert.Equal(t, "electronics", response.Slug)
}

func TestToCategoryResponseList(t *testing.T) {
	// Setup
	now := time.Now()
	categories := []models.Category{
		{
			Base: models.Base{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Name:        "Electronics",
			Slug:        "electronics",
			Description: "All electronics",
		},
		{
			Base: models.Base{
				ID:        2,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Name:        "Books",
			Slug:        "books",
			Description: "All books",
		},
	}

	// Execute
	responses := ToCategoryResponseList(categories)

	// Assert
	assert.Len(t, responses, 2)
	assert.Equal(t, "Electronics", responses[0].Name)
	assert.Equal(t, "Books", responses[1].Name)
}

func TestToCategoryTreeResponse(t *testing.T) {
	// Setup
	now := time.Now()
	parentID := uint(1)

	child := models.Category{
		Base: models.Base{
			ID:        2,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        "Laptops",
		Slug:        "laptops",
		Description: "Laptop computers",
		ParentID:    &parentID,
		Children:    []models.Category{},
	}

	categories := []models.Category{
		{
			Base: models.Base{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Name:        "Electronics",
			Slug:        "electronics",
			Description: "All electronics",
			ParentID:    nil,
			Children:    []models.Category{child},
		},
	}

	// Execute
	response := ToCategoryTreeResponse(categories)

	// Assert
	assert.Len(t, response.Categories, 1)
	assert.Equal(t, "Electronics", response.Categories[0].Name)
	assert.Len(t, response.Categories[0].Children, 1)
	assert.Equal(t, "Laptops", response.Categories[0].Children[0].Name)
}

