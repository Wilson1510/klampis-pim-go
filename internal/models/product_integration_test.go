//go:build integration
// +build integration

package models_test

import (
	"testing"

	"github.com/Wilson1510/klampis-pim-go/internal/models"
	"github.com/Wilson1510/klampis-pim-go/internal/testutil"
	"gorm.io/gorm"
)

func TestProductCreate_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create a category for testing
	category := models.Category{
		Name:        "Electronics",
		Description: "Electronic devices",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	testCases := []struct {
		name         string
		product      models.Product
		expectError  bool
		checkSlug    bool
		expectedSlug string
	}{
		{
			name: "Create valid product",
			product: models.Product{
				Name:        "Laptop Computer",
				Description: "High performance laptop",
				CategoryID:  category.ID,
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "laptop-computer",
		},
		{
			name: "Create product with special characters",
			product: models.Product{
				Name:        "Phone & Accessories",
				Description: "Mobile phone with accessories",
				CategoryID:  category.ID,
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "phone-accessories",
		},
		{
			name: "Create product without description",
			product: models.Product{
				Name:       "Tablet",
				CategoryID: category.ID,
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "tablet",
		},
		{
			name: "Create product with duplicate name",
			product: models.Product{
				Name:        "Laptop Computer",
				Description: "Another laptop",
				CategoryID:  category.ID,
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "laptop-computer-1",
		},
		{
			name: "Create product without name",
			product: models.Product{
				Description: "No name product",
				CategoryID:  category.ID,
			},
			expectError: false, // golang will automatically set the name to ""
		},
		{
			name: "Create product without category",
			product: models.Product{
				Name:        "Invalid Product",
				Description: "No category",
			},
			expectError: true,
		},
		{
			name: "Create product with invalid category ID",
			product: models.Product{
				Name:        "Invalid Product",
				Description: "Invalid category",
				CategoryID:  99999,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.product.CreatedBy = testUser.ID
			tc.product.UpdatedBy = testUser.ID
			result := db.Create(&tc.product)

			if tc.expectError {
				if result.Error == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if result.Error != nil {
					t.Errorf("Expected no error but got: %v", result.Error)
				}
				if tc.product.ID == 0 {
					t.Error("Expected product ID to be set after creation")
				}
				if tc.product.CreatedAt.IsZero() {
					t.Error("Expected CreatedAt to be set after creation")
				}
				if tc.checkSlug && tc.product.Slug != tc.expectedSlug {
					t.Errorf("Expected slug '%s', got '%s'", tc.expectedSlug, tc.product.Slug)
				}
			}
		})
	}
}

func TestProductRead_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create a test category
	category := models.Category{
		Name:        "Electronics",
		Description: "Electronic devices",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	// Create a test product
	product := models.Product{
		Name:        "Test Laptop",
		Description: "A test laptop",
		CategoryID:  category.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	t.Run("Read product by ID", func(t *testing.T) {
		var foundProduct models.Product
		result := db.First(&foundProduct, product.ID)

		if result.Error != nil {
			t.Errorf("Expected to find product but got error: %v", result.Error)
		}
		if foundProduct.Name != product.Name {
			t.Errorf("Expected name '%s', got '%s'", product.Name, foundProduct.Name)
		}
		if foundProduct.Slug != product.Slug {
			t.Errorf("Expected slug '%s', got '%s'", product.Slug, foundProduct.Slug)
		}
		if foundProduct.CategoryID != product.CategoryID {
			t.Errorf("Expected category ID %d, got %d", product.CategoryID, foundProduct.CategoryID)
		}
	})

	t.Run("Read product by slug", func(t *testing.T) {
		var foundProduct models.Product
		result := db.Where("slug = ?", product.Slug).First(&foundProduct)

		if result.Error != nil {
			t.Errorf("Expected to find product but got error: %v", result.Error)
		}
		if foundProduct.ID != product.ID {
			t.Errorf("Expected ID %d, got %d", product.ID, foundProduct.ID)
		}
	})

	t.Run("Read product with category relationship", func(t *testing.T) {
		var foundProduct models.Product
		result := db.Preload("Category").First(&foundProduct, product.ID)

		if result.Error != nil {
			t.Errorf("Expected to find product but got error: %v", result.Error)
		}
		if foundProduct.Category == nil {
			t.Error("Expected category to be loaded")
		} else {
			if foundProduct.Category.ID != category.ID {
				t.Errorf("Expected category ID %d, got %d", category.ID, foundProduct.Category.ID)
			}
			if foundProduct.Category.Name != category.Name {
				t.Errorf("Expected category name '%s', got '%s'", category.Name, foundProduct.Category.Name)
			}
		}
	})

	t.Run("Read non-existent product", func(t *testing.T) {
		var foundProduct models.Product
		result := db.First(&foundProduct, 99999)

		if result.Error != gorm.ErrRecordNotFound {
			t.Errorf("Expected ErrRecordNotFound but got: %v", result.Error)
		}
	})
}

func TestProductUpdate_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test categories
	category1 := models.Category{
		Name:        "Electronics",
		Description: "Electronic devices",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&category1).Error; err != nil {
		t.Fatalf("Failed to create test category1: %v", err)
	}

	category2 := models.Category{
		Name:        "Books",
		Description: "All kinds of books",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&category2).Error; err != nil {
		t.Fatalf("Failed to create test category2: %v", err)
	}

	// Create a test product
	product := models.Product{
		Name:        "Update Product",
		Description: "Original description",
		CategoryID:  category1.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	t.Run("Update product name and verify slug update", func(t *testing.T) {
		product.Name = "Updated Product Name"
		result := db.Save(&product)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedProduct models.Product
		db.First(&updatedProduct, product.ID)
		if updatedProduct.Name != "Updated Product Name" {
			t.Errorf("Expected name 'Updated Product Name', got '%s'", updatedProduct.Name)
		}
		if updatedProduct.Slug != "updated-product-name" {
			t.Errorf("Expected slug to be updated to 'updated-product-name', got '%s'", updatedProduct.Slug)
		}
	})

	t.Run("Update product description", func(t *testing.T) {
		product.Description = "Updated description"
		result := db.Save(&product)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedProduct models.Product
		db.First(&updatedProduct, product.ID)
		if updatedProduct.Description != "Updated description" {
			t.Errorf("Expected description 'Updated description', got '%s'", updatedProduct.Description)
		}
	})

	t.Run("Update product category", func(t *testing.T) {
		product.CategoryID = category2.ID
		result := db.Save(&product)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedProduct models.Product
		db.Preload("Category").First(&updatedProduct, product.ID)
		if updatedProduct.CategoryID != category2.ID {
			t.Errorf("Expected category ID %d, got %d", category2.ID, updatedProduct.CategoryID)
		}
		if updatedProduct.Category.Name != category2.Name {
			t.Errorf("Expected category name '%s', got '%s'", category2.Name, updatedProduct.Category.Name)
		}
	})

	t.Run("Update product with invalid category ID", func(t *testing.T) {
		product.CategoryID = 99999
		result := db.Save(&product)

		if result.Error == nil {
			t.Error("Expected error for invalid category ID but got none")
		}
	})
}

func TestProductDelete_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create a test category
	category := models.Category{
		Name:        "Electronics",
		Description: "Electronic devices",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	// Create a test product
	product := models.Product{
		Name:        "Delete Product",
		Description: "To be deleted",
		CategoryID:  category.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	t.Run("Soft delete product", func(t *testing.T) {
		result := db.Delete(&product)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify soft delete - should not be found in normal query
		var foundProduct models.Product
		result = db.First(&foundProduct, product.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected product to be soft deleted (not found in normal query)")
		}

		// Verify soft delete - should be found with Unscoped
		result = db.Unscoped().First(&foundProduct, product.ID)
		if result.Error != nil {
			t.Errorf("Expected to find soft deleted product with Unscoped but got error: %v", result.Error)
		}
		if foundProduct.DeletedAt.Time.IsZero() {
			t.Error("Expected DeletedAt to be set after soft delete")
		}
	})

	t.Run("Permanent delete product", func(t *testing.T) {
		// Create another product
		anotherProduct := models.Product{
			Name:        "Permanent Delete",
			Description: "To be permanently deleted",
			CategoryID:  category.ID,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&anotherProduct)

		// Permanently delete
		result := db.Unscoped().Delete(&anotherProduct)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify permanent delete - should not be found even with Unscoped
		var foundProduct models.Product
		result = db.Unscoped().First(&foundProduct, anotherProduct.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected product to be permanently deleted")
		}
	})
}

func TestProductQuery_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test categories
	category1 := models.Category{Name: "Electronics", Description: "Electronic devices", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	category2 := models.Category{Name: "Books", Description: "All kinds of books", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&category1)
	db.Create(&category2)

	// Create multiple test products
	products := []models.Product{
		{Name: "Laptop", Description: "High performance laptop", CategoryID: category1.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
		{Name: "Phone", Description: "Smartphone", CategoryID: category1.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
		{Name: "Novel", Description: "Fiction novel", CategoryID: category2.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
	}

	for _, prod := range products {
		if err := db.Create(&prod).Error; err != nil {
			t.Fatalf("Failed to create test product: %v", err)
		}
	}

	t.Run("Find all products", func(t *testing.T) {
		var allProducts []models.Product
		result := db.Find(&allProducts)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(allProducts) != 3 {
			t.Errorf("Expected 3 products, got %d", len(allProducts))
		}
	})

	t.Run("Find products by category", func(t *testing.T) {
		var electronicsProducts []models.Product
		result := db.Where("category_id = ?", category1.ID).Find(&electronicsProducts)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(electronicsProducts) != 2 {
			t.Errorf("Expected 2 electronics products, got %d", len(electronicsProducts))
		}
	})

	t.Run("Count products", func(t *testing.T) {
		var count int64
		result := db.Model(&models.Product{}).Count(&count)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if count != 3 {
			t.Errorf("Expected count 3, got %d", count)
		}
	})

	t.Run("Find products with pagination", func(t *testing.T) {
		var paginatedProducts []models.Product
		result := db.Limit(2).Offset(0).Find(&paginatedProducts)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(paginatedProducts) != 2 {
			t.Errorf("Expected 2 products in first page, got %d", len(paginatedProducts))
		}
	})

	t.Run("Search products by name", func(t *testing.T) {
		var foundProducts []models.Product
		result := db.Where("name LIKE ?", "%Laptop%").Find(&foundProducts)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(foundProducts) != 1 {
			t.Errorf("Expected 1 product, got %d", len(foundProducts))
		}
		if len(foundProducts) > 0 && foundProducts[0].Name != "Laptop" {
			t.Errorf("Expected product name 'Laptop', got '%s'", foundProducts[0].Name)
		}
	})
}

func TestProductRelationships_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create a category
	category := models.Category{
		Name:        "Electronics",
		Description: "Electronic devices",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	// Create a product
	product := models.Product{
		Name:        "Laptop",
		Description: "High performance laptop",
		CategoryID:  category.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	t.Run("Category has products", func(t *testing.T) {
		var foundCategory models.Category
		result := db.Preload("Products").First(&foundCategory, category.ID)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(foundCategory.Products) != 1 {
			t.Errorf("Expected 1 product in category, got %d", len(foundCategory.Products))
		}
		if len(foundCategory.Products) > 0 && foundCategory.Products[0].ID != product.ID {
			t.Errorf("Expected product ID %d, got %d", product.ID, foundCategory.Products[0].ID)
		}
	})

	t.Run("Product belongs to category", func(t *testing.T) {
		var foundProduct models.Product
		result := db.Preload("Category").First(&foundProduct, product.ID)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if foundProduct.Category == nil {
			t.Error("Expected category to be loaded")
		} else if foundProduct.Category.ID != category.ID {
			t.Errorf("Expected category ID %d, got %d", category.ID, foundProduct.Category.ID)
		}
	})
}

func TestProductSlugUniqueness_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create a test category
	category := models.Category{
		Name:        "Electronics",
		Description: "Electronic devices",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	t.Run("Create products with same name generates unique slugs", func(t *testing.T) {
		product1 := models.Product{Name: "Test Product", CategoryID: category.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
		product2 := models.Product{Name: "Test Product", CategoryID: category.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
		product3 := models.Product{Name: "Test Product", CategoryID: category.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}

		if err := db.Create(&product1).Error; err != nil {
			t.Fatalf("Failed to create product1: %v", err)
		}
		if err := db.Create(&product2).Error; err != nil {
			t.Fatalf("Failed to create product2: %v", err)
		}
		if err := db.Create(&product3).Error; err != nil {
			t.Fatalf("Failed to create product3: %v", err)
		}

		// Verify slugs are unique
		if product1.Slug != "test-product" {
			t.Errorf("Expected first slug 'test-product', got '%s'", product1.Slug)
		}
		if product2.Slug != "test-product-1" {
			t.Errorf("Expected second slug 'test-product-1', got '%s'", product2.Slug)
		}
		if product3.Slug != "test-product-2" {
			t.Errorf("Expected third slug 'test-product-2', got '%s'", product3.Slug)
		}
	})
}
