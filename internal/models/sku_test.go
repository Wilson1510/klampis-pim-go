package models

import (
	"testing"
)

// TestSkuSlugModelInterface tests that Sku properly implements SlugModel interface
func TestSkuSlugModelInterface(t *testing.T) {
	sku := Sku{
		Name: "Test SKU",
		Slug: "test-sku",
	}
	// Manually set ID since gorm.Model uses lowercase fields
	sku.ID = 1

	// Test SlugModel interface implementation
	if sku.GetID() != 1 {
		t.Errorf("Expected ID 1, got %d", sku.GetID())
	}
	
	if sku.GetName() != "Test SKU" {
		t.Errorf("Expected name 'Test SKU', got '%s'", sku.GetName())
	}
	
	if sku.GetSlug() != "test-sku" {
		t.Errorf("Expected slug 'test-sku', got '%s'", sku.GetSlug())
	}
	
	if sku.GetTableName() != "skus" {
		t.Errorf("Expected table name 'skus', got '%s'", sku.GetTableName())
	}
	
	// Test SetSlug
	sku.SetSlug("new-sku-slug")
	if sku.GetSlug() != "new-sku-slug" {
		t.Errorf("Expected slug 'new-sku-slug' after SetSlug, got '%s'", sku.GetSlug())
	}
}
