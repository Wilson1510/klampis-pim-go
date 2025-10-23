package models

import (
	"testing"
)

// TestProductSlugModelInterface tests that Product properly implements SlugModel interface
func TestProductSlugModelInterface(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Slug: "test-product",
	}
	// Manually set ID since gorm.Model uses lowercase fields
	product.ID = 1

	// Test SlugModel interface implementation
	if product.GetID() != 1 {
		t.Errorf("Expected ID 1, got %d", product.GetID())
	}

	if product.GetName() != "Test Product" {
		t.Errorf("Expected name 'Test Product', got '%s'", product.GetName())
	}

	if product.GetSlug() != "test-product" {
		t.Errorf("Expected slug 'test-product', got '%s'", product.GetSlug())
	}

	if product.GetTableName() != "products" {
		t.Errorf("Expected table name 'products', got '%s'", product.GetTableName())
	}

	// Test SetSlug
	product.SetSlug("new-product-slug")
	if product.GetSlug() != "new-product-slug" {
		t.Errorf("Expected slug 'new-product-slug' after SetSlug, got '%s'", product.GetSlug())
	}
}
