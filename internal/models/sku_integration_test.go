//go:build integration
// +build integration

package models_test

import (
	"testing"

	"github.com/Wilson1510/klampis-pim-go/internal/models"
	"github.com/Wilson1510/klampis-pim-go/internal/testutil"
	"gorm.io/gorm"
)

func TestSkuCreate_Integration(t *testing.T) {
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

	// Create category and product for testing
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
		t.Fatalf("Failed to create test product: %v", err)
	}

	testCases := []struct {
		name         string
		sku          models.Sku
		expectError  bool
		checkSlug    bool
		expectedSlug string
	}{
		{
			name: "Create valid SKU",
			sku: models.Sku{
				Name:        "Laptop - 16GB RAM",
				Description: "16GB RAM variant",
				SkuNumber:   "LAP-16GB-001",
				Price:       999.99,
				ProductID:   product.ID,
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "laptop-16gb-ram",
		},
		{
			name: "Create SKU with special characters",
			sku: models.Sku{
				Name:        "Laptop - 32GB RAM & SSD",
				Description: "32GB RAM with SSD",
				SkuNumber:   "LAP-32GB-001",
				Price:       1299.99,
				ProductID:   product.ID,
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "laptop-32gb-ram-ssd",
		},
		{
			name: "Create SKU without description",
			sku: models.Sku{
				Name:      "Laptop - 8GB RAM",
				SkuNumber: "LAP-8GB-001",
				Price:     799.99,
				ProductID: product.ID,
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "laptop-8gb-ram",
		},
		{
			name: "Create SKU with duplicate name",
			sku: models.Sku{
				Name:        "Laptop - 16GB RAM",
				Description: "Another 16GB variant",
				SkuNumber:   "LAP-16GB-002",
				Price:       999.99,
				ProductID:   product.ID,
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "laptop-16gb-ram-1",
		},
		{
			name: "Create SKU with duplicate SKU number",
			sku: models.Sku{
				Name:        "Laptop - Duplicate SKU",
				Description: "Duplicate SKU number",
				SkuNumber:   "LAP-16GB-001", // Already exists
				Price:       999.99,
				ProductID:   product.ID,
			},
			expectError: true,
		},
		{
			name: "Create SKU without name",
			sku: models.Sku{
				Description: "No name SKU",
				SkuNumber:   "LAP-NONAME-001",
				Price:       999.99,
				ProductID:   product.ID,
			},
			expectError: false, // golang will automatically set the name to ""
		},
		{
			name: "Create SKU without SKU number",
			sku: models.Sku{
				Name:        "Laptop - No SKU",
				Description: "No SKU number",
				Price:       999.99,
				ProductID:   product.ID,
			},
			expectError: false, // golang will automatically set the name to ""
		},
		{
			name: "Create SKU without product",
			sku: models.Sku{
				Name:        "Invalid SKU",
				Description: "No product",
				SkuNumber:   "INVALID-001",
				Price:       999.99,
			},
			expectError: true,
		},
		{
			name: "Create SKU with invalid product ID",
			sku: models.Sku{
				Name:        "Invalid SKU",
				Description: "Invalid product",
				SkuNumber:   "INVALID-002",
				Price:       999.99,
				ProductID:   99999,
			},
			expectError: true,
		},
		{
			name: "Create SKU with zero price",
			sku: models.Sku{
				Name:        "Free SKU",
				Description: "Free item",
				SkuNumber:   "FREE-001",
				Price:       0.0,
				ProductID:   product.ID,
			},
			expectError: false,
		},
		{
			name: "Create SKU with negative price",
			sku: models.Sku{
				Name:        "Negative Price SKU",
				Description: "Negative price",
				SkuNumber:   "NEG-001",
				Price:       -10.0,
				ProductID:   product.ID,
			},
			expectError: false, // GORM doesn't validate this by default
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.sku.CreatedBy = testUser.ID
			tc.sku.UpdatedBy = testUser.ID
			result := db.Create(&tc.sku)

			if tc.expectError {
				if result.Error == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if result.Error != nil {
					t.Errorf("Expected no error but got: %v", result.Error)
				}
				if tc.sku.ID == 0 {
					t.Error("Expected SKU ID to be set after creation")
				}
				if tc.sku.CreatedAt.IsZero() {
					t.Error("Expected CreatedAt to be set after creation")
				}
				if tc.checkSlug && tc.sku.Slug != tc.expectedSlug {
					t.Errorf("Expected slug '%s', got '%s'", tc.expectedSlug, tc.sku.Slug)
				}
			}
		})
	}
}

func TestSkuRead_Integration(t *testing.T) {
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

	// Create category and product
	category := models.Category{
		Name:        "Electronics",
		Description: "Electronic devices",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	db.Create(&category)

	product := models.Product{
		Name:        "Laptop",
		Description: "High performance laptop",
		CategoryID:  category.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	db.Create(&product)

	// Create a test SKU
	sku := models.Sku{
		Name:        "Laptop - 16GB RAM",
		Description: "16GB RAM variant",
		SkuNumber:   "LAP-16GB-001",
		Price:       999.99,
		ProductID:   product.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&sku).Error; err != nil {
		t.Fatalf("Failed to create test SKU: %v", err)
	}

	t.Run("Read SKU by ID", func(t *testing.T) {
		var foundSku models.Sku
		result := db.First(&foundSku, sku.ID)

		if result.Error != nil {
			t.Errorf("Expected to find SKU but got error: %v", result.Error)
		}
		if foundSku.Name != sku.Name {
			t.Errorf("Expected name '%s', got '%s'", sku.Name, foundSku.Name)
		}
		if foundSku.Slug != sku.Slug {
			t.Errorf("Expected slug '%s', got '%s'", sku.Slug, foundSku.Slug)
		}
		if foundSku.SkuNumber != sku.SkuNumber {
			t.Errorf("Expected SKU number '%s', got '%s'", sku.SkuNumber, foundSku.SkuNumber)
		}
		if foundSku.Price != sku.Price {
			t.Errorf("Expected price %.2f, got %.2f", sku.Price, foundSku.Price)
		}
		if foundSku.ProductID != sku.ProductID {
			t.Errorf("Expected product ID %d, got %d", sku.ProductID, foundSku.ProductID)
		}
	})

	t.Run("Read SKU by slug", func(t *testing.T) {
		var foundSku models.Sku
		result := db.Where("slug = ?", sku.Slug).First(&foundSku)

		if result.Error != nil {
			t.Errorf("Expected to find SKU but got error: %v", result.Error)
		}
		if foundSku.ID != sku.ID {
			t.Errorf("Expected ID %d, got %d", sku.ID, foundSku.ID)
		}
	})

	t.Run("Read SKU by SKU number", func(t *testing.T) {
		var foundSku models.Sku
		result := db.Where("sku_number = ?", sku.SkuNumber).First(&foundSku)

		if result.Error != nil {
			t.Errorf("Expected to find SKU but got error: %v", result.Error)
		}
		if foundSku.ID != sku.ID {
			t.Errorf("Expected ID %d, got %d", sku.ID, foundSku.ID)
		}
	})

	t.Run("Read SKU with product relationship", func(t *testing.T) {
		var foundSku models.Sku
		result := db.Preload("Product").First(&foundSku, sku.ID)

		if result.Error != nil {
			t.Errorf("Expected to find SKU but got error: %v", result.Error)
		}
		if foundSku.Product == nil {
			t.Error("Expected product to be loaded")
		} else {
			if foundSku.Product.ID != product.ID {
				t.Errorf("Expected product ID %d, got %d", product.ID, foundSku.Product.ID)
			}
			if foundSku.Product.Name != product.Name {
				t.Errorf("Expected product name '%s', got '%s'", product.Name, foundSku.Product.Name)
			}
		}
	})

	t.Run("Read SKU with nested relationships", func(t *testing.T) {
		var foundSku models.Sku
		result := db.Preload("Product.Category").First(&foundSku, sku.ID)

		if result.Error != nil {
			t.Errorf("Expected to find SKU but got error: %v", result.Error)
		}
		if foundSku.Product == nil {
			t.Error("Expected product to be loaded")
		} else if foundSku.Product.Category == nil {
			t.Error("Expected category to be loaded")
		} else {
			if foundSku.Product.Category.ID != category.ID {
				t.Errorf("Expected category ID %d, got %d", category.ID, foundSku.Product.Category.ID)
			}
		}
	})

	t.Run("Read non-existent SKU", func(t *testing.T) {
		var foundSku models.Sku
		result := db.First(&foundSku, 99999)

		if result.Error != gorm.ErrRecordNotFound {
			t.Errorf("Expected ErrRecordNotFound but got: %v", result.Error)
		}
	})
}

func TestSkuUpdate_Integration(t *testing.T) {
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

	// Create category and products
	category := models.Category{Name: "Electronics", Description: "Electronic devices", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&category)

	product1 := models.Product{Name: "Laptop", Description: "Laptop computer", CategoryID: category.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&product1)

	product2 := models.Product{Name: "Desktop", Description: "Desktop computer", CategoryID: category.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&product2)

	// Create a test SKU
	sku := models.Sku{
		Name:        "Update SKU",
		Description: "Original description",
		SkuNumber:   "UPD-001",
		Price:       999.99,
		ProductID:   product1.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&sku).Error; err != nil {
		t.Fatalf("Failed to create test SKU: %v", err)
	}

	t.Run("Update SKU name and verify slug update", func(t *testing.T) {
		sku.Name = "Updated SKU Name"
		result := db.Save(&sku)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedSku models.Sku
		db.First(&updatedSku, sku.ID)
		if updatedSku.Name != "Updated SKU Name" {
			t.Errorf("Expected name 'Updated SKU Name', got '%s'", updatedSku.Name)
		}
		if updatedSku.Slug != "updated-sku-name" {
			t.Errorf("Expected slug to be updated to 'updated-sku-name', got '%s'", updatedSku.Slug)
		}
	})

	t.Run("Update SKU description", func(t *testing.T) {
		sku.Description = "Updated description"
		result := db.Save(&sku)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedSku models.Sku
		db.First(&updatedSku, sku.ID)
		if updatedSku.Description != "Updated description" {
			t.Errorf("Expected description 'Updated description', got '%s'", updatedSku.Description)
		}
	})

	t.Run("Update SKU price", func(t *testing.T) {
		sku.Price = 1299.99
		result := db.Save(&sku)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedSku models.Sku
		db.First(&updatedSku, sku.ID)
		if updatedSku.Price != 1299.99 {
			t.Errorf("Expected price 1299.99, got %.2f", updatedSku.Price)
		}
	})

	t.Run("Update SKU number", func(t *testing.T) {
		sku.SkuNumber = "UPD-002"
		result := db.Save(&sku)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedSku models.Sku
		db.First(&updatedSku, sku.ID)
		if updatedSku.SkuNumber != "UPD-002" {
			t.Errorf("Expected SKU number 'UPD-002', got '%s'", updatedSku.SkuNumber)
		}
	})

	t.Run("Update SKU product", func(t *testing.T) {
		sku.ProductID = product2.ID
		result := db.Save(&sku)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedSku models.Sku
		db.Preload("Product").First(&updatedSku, sku.ID)
		if updatedSku.ProductID != product2.ID {
			t.Errorf("Expected product ID %d, got %d", product2.ID, updatedSku.ProductID)
		}
		if updatedSku.Product.Name != product2.Name {
			t.Errorf("Expected product name '%s', got '%s'", product2.Name, updatedSku.Product.Name)
		}
	})

	t.Run("Update SKU with invalid product ID", func(t *testing.T) {
		sku.ProductID = 99999
		result := db.Save(&sku)

		if result.Error == nil {
			t.Error("Expected error for invalid product ID but got none")
		}
	})

	t.Run("Update SKU number to duplicate", func(t *testing.T) {
		// Create another SKU
		anotherSku := models.Sku{
			Name:        "Another SKU",
			Description: "Another SKU",
			SkuNumber:   "ANOTHER-001",
			Price:       799.99,
			ProductID:   product1.ID,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&anotherSku)

		// Try to update first SKU's number to duplicate
		// First reset to valid state
		sku.ProductID = product1.ID
		db.Save(&sku)

		sku.SkuNumber = "ANOTHER-001"
		result := db.Save(&sku)

		if result.Error == nil {
			t.Error("Expected error for duplicate SKU number but got none")
		}
	})
}

func TestSkuDelete_Integration(t *testing.T) {
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

	// Create category and product
	category := models.Category{Name: "Electronics", Description: "Electronic devices", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&category)

	product := models.Product{Name: "Laptop", Description: "Laptop computer", CategoryID: category.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&product)

	// Create a test SKU
	sku := models.Sku{
		Name:        "Delete SKU",
		Description: "To be deleted",
		SkuNumber:   "DEL-001",
		Price:       999.99,
		ProductID:   product.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&sku).Error; err != nil {
		t.Fatalf("Failed to create test SKU: %v", err)
	}

	t.Run("Soft delete SKU", func(t *testing.T) {
		result := db.Delete(&sku)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify soft delete - should not be found in normal query
		var foundSku models.Sku
		result = db.First(&foundSku, sku.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected SKU to be soft deleted (not found in normal query)")
		}

		// Verify soft delete - should be found with Unscoped
		result = db.Unscoped().First(&foundSku, sku.ID)
		if result.Error != nil {
			t.Errorf("Expected to find soft deleted SKU with Unscoped but got error: %v", result.Error)
		}
		if foundSku.DeletedAt.Time.IsZero() {
			t.Error("Expected DeletedAt to be set after soft delete")
		}
	})

	t.Run("Permanent delete SKU", func(t *testing.T) {
		// Create another SKU
		anotherSku := models.Sku{
			Name:        "Permanent Delete",
			Description: "To be permanently deleted",
			SkuNumber:   "PERM-001",
			Price:       799.99,
			ProductID:   product.ID,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&anotherSku)

		// Permanently delete
		result := db.Unscoped().Delete(&anotherSku)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify permanent delete - should not be found even with Unscoped
		var foundSku models.Sku
		result = db.Unscoped().First(&foundSku, anotherSku.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected SKU to be permanently deleted")
		}
	})
}

func TestSkuQuery_Integration(t *testing.T) {
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

	// Create category and products
	category := models.Category{Name: "Electronics", Description: "Electronic devices", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&category)

	product1 := models.Product{Name: "Laptop", Description: "Laptop computer", CategoryID: category.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	product2 := models.Product{Name: "Desktop", Description: "Desktop computer", CategoryID: category.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&product1)
	db.Create(&product2)

	// Create multiple test SKUs
	skus := []models.Sku{
		{Name: "Laptop - 8GB", SkuNumber: "LAP-8GB", Price: 799.99, ProductID: product1.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
		{Name: "Laptop - 16GB", SkuNumber: "LAP-16GB", Price: 999.99, ProductID: product1.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
		{Name: "Desktop - 32GB", SkuNumber: "DSK-32GB", Price: 1499.99, ProductID: product2.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
	}

	for _, sku := range skus {
		if err := db.Create(&sku).Error; err != nil {
			t.Fatalf("Failed to create test SKU: %v", err)
		}
	}

	t.Run("Find all SKUs", func(t *testing.T) {
		var allSkus []models.Sku
		result := db.Find(&allSkus)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(allSkus) != 3 {
			t.Errorf("Expected 3 SKUs, got %d", len(allSkus))
		}
	})

	t.Run("Find SKUs by product", func(t *testing.T) {
		var laptopSkus []models.Sku
		result := db.Where("product_id = ?", product1.ID).Find(&laptopSkus)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(laptopSkus) != 2 {
			t.Errorf("Expected 2 laptop SKUs, got %d", len(laptopSkus))
		}
	})

	t.Run("Find SKUs by price range", func(t *testing.T) {
		var affordableSkus []models.Sku
		result := db.Where("price <= ?", 1000.0).Find(&affordableSkus)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(affordableSkus) != 2 {
			t.Errorf("Expected 2 affordable SKUs, got %d", len(affordableSkus))
		}
	})

	t.Run("Count SKUs", func(t *testing.T) {
		var count int64
		result := db.Model(&models.Sku{}).Count(&count)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if count != 3 {
			t.Errorf("Expected count 3, got %d", count)
		}
	})

	t.Run("Find SKUs with pagination", func(t *testing.T) {
		var paginatedSkus []models.Sku
		result := db.Limit(2).Offset(0).Find(&paginatedSkus)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(paginatedSkus) != 2 {
			t.Errorf("Expected 2 SKUs in first page, got %d", len(paginatedSkus))
		}
	})

	t.Run("Search SKUs by name", func(t *testing.T) {
		var foundSkus []models.Sku
		result := db.Where("name LIKE ?", "%16GB%").Find(&foundSkus)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(foundSkus) != 1 {
			t.Errorf("Expected 1 SKU, got %d", len(foundSkus))
		}
		if len(foundSkus) > 0 && foundSkus[0].Name != "Laptop - 16GB" {
			t.Errorf("Expected SKU name 'Laptop - 16GB', got '%s'", foundSkus[0].Name)
		}
	})

	t.Run("Order SKUs by price", func(t *testing.T) {
		var orderedSkus []models.Sku
		result := db.Order("price ASC").Find(&orderedSkus)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(orderedSkus) != 3 {
			t.Errorf("Expected 3 SKUs, got %d", len(orderedSkus))
		}
		if len(orderedSkus) >= 2 {
			if orderedSkus[0].Price > orderedSkus[1].Price {
				t.Error("Expected SKUs to be ordered by price ascending")
			}
		}
	})
}

func TestSkuRelationships_Integration(t *testing.T) {
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

	// Create category
	category := models.Category{Name: "Electronics", Description: "Electronic devices", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&category)

	// Create product
	product := models.Product{Name: "Laptop", Description: "Laptop computer", CategoryID: category.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&product)

	// Create SKU
	sku := models.Sku{
		Name:        "Laptop - 16GB",
		Description: "16GB RAM variant",
		SkuNumber:   "LAP-16GB",
		Price:       999.99,
		ProductID:   product.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	db.Create(&sku)

	t.Run("Product has SKUs", func(t *testing.T) {
		var foundProduct models.Product
		result := db.Preload("Skus").First(&foundProduct, product.ID)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(foundProduct.Skus) != 1 {
			t.Errorf("Expected 1 SKU in product, got %d", len(foundProduct.Skus))
		}
		if len(foundProduct.Skus) > 0 && foundProduct.Skus[0].ID != sku.ID {
			t.Errorf("Expected SKU ID %d, got %d", sku.ID, foundProduct.Skus[0].ID)
		}
	})

	t.Run("SKU belongs to product", func(t *testing.T) {
		var foundSku models.Sku
		result := db.Preload("Product").First(&foundSku, sku.ID)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if foundSku.Product == nil {
			t.Error("Expected product to be loaded")
		} else if foundSku.Product.ID != product.ID {
			t.Errorf("Expected product ID %d, got %d", product.ID, foundSku.Product.ID)
		}
	})
}

func TestSkuSlugUniqueness_Integration(t *testing.T) {
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

	// Create category and product
	category := models.Category{Name: "Electronics", Description: "Electronic devices", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&category)

	product := models.Product{Name: "Laptop", Description: "Laptop computer", CategoryID: category.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
	db.Create(&product)

	t.Run("Create SKUs with same name generates unique slugs", func(t *testing.T) {
		sku1 := models.Sku{Name: "Test SKU", SkuNumber: "TST-001", Price: 100.0, ProductID: product.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
		sku2 := models.Sku{Name: "Test SKU", SkuNumber: "TST-002", Price: 100.0, ProductID: product.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
		sku3 := models.Sku{Name: "Test SKU", SkuNumber: "TST-003", Price: 100.0, ProductID: product.ID, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}

		if err := db.Create(&sku1).Error; err != nil {
			t.Fatalf("Failed to create sku1: %v", err)
		}
		if err := db.Create(&sku2).Error; err != nil {
			t.Fatalf("Failed to create sku2: %v", err)
		}
		if err := db.Create(&sku3).Error; err != nil {
			t.Fatalf("Failed to create sku3: %v", err)
		}

		// Verify slugs are unique
		if sku1.Slug != "test-sku" {
			t.Errorf("Expected first slug 'test-sku', got '%s'", sku1.Slug)
		}
		if sku2.Slug != "test-sku-1" {
			t.Errorf("Expected second slug 'test-sku-1', got '%s'", sku2.Slug)
		}
		if sku3.Slug != "test-sku-2" {
			t.Errorf("Expected third slug 'test-sku-2', got '%s'", sku3.Slug)
		}
	})
}
