package models

import (
	"testing"
)

// TestCategorySlugModelInterface tests that Category properly implements SlugModel interface
func TestCategorySlugModelInterface(t *testing.T) {
	category := Category{
		Name: "Test Category",
		Slug: "test-category",
	}
	// Manually set ID since gorm.Model uses lowercase fields
	category.ID = 1

	// Test SlugModel interface implementation
	if category.GetID() != 1 {
		t.Errorf("Expected ID 1, got %d", category.GetID())
	}
	
	if category.GetName() != "Test Category" {
		t.Errorf("Expected name 'Test Category', got '%s'", category.GetName())
	}
	
	if category.GetSlug() != "test-category" {
		t.Errorf("Expected slug 'test-category', got '%s'", category.GetSlug())
	}
	
	if category.GetTableName() != "categories" {
		t.Errorf("Expected table name 'categories', got '%s'", category.GetTableName())
	}
	
	// Test SetSlug
	category.SetSlug("new-category-slug")
	if category.GetSlug() != "new-category-slug" {
		t.Errorf("Expected slug 'new-category-slug' after SetSlug, got '%s'", category.GetSlug())
	}
}
