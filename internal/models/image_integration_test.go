//go:build integration
// +build integration

package models_test

import (
	"testing"

	"github.com/Wilson1510/klampis-pim-go/internal/models"
	"github.com/Wilson1510/klampis-pim-go/internal/testutil"
	"gorm.io/gorm"
)

func TestImageCreate_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy (required for Base model)
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	testCases := []struct {
		name         string
		image        models.Image
		expectError  bool
		errorMsg     string
		validateFunc func(*testing.T, *models.Image)
	}{
		{
			name: "Create valid image with all fields",
			image: models.Image{
				File:          "/uploads/product-image-1.jpg",
				Title:         "Product Main Image",
				IsPrimary:     true,
				ImageableID:   1,
				ImageableType: "Product",
				Base: models.Base{
					CreatedBy: testUser.ID,
					UpdatedBy: testUser.ID,
				},
			},
			expectError: false,
			validateFunc: func(t *testing.T, img *models.Image) {
				if img.ID == 0 {
					t.Error("Expected image ID to be set after creation")
				}
				if img.File != "/uploads/product-image-1.jpg" {
					t.Errorf("Expected file '/uploads/product-image-1.jpg', got '%s'", img.File)
				}
				if img.Title != "Product Main Image" {
					t.Errorf("Expected title 'Product Main Image', got '%s'", img.Title)
				}
				if !img.IsPrimary {
					t.Error("Expected IsPrimary to be true")
				}
			},
		},
		{
			name: "Create image without title (nullable field)",
			image: models.Image{
				File:          "/uploads/product-image-2.jpg",
				IsPrimary:     false,
				ImageableID:   1,
				ImageableType: "Product",
				Base: models.Base{
					CreatedBy: testUser.ID,
					UpdatedBy: testUser.ID,
				},
			},
			expectError: false,
			validateFunc: func(t *testing.T, img *models.Image) {
				if img.Title != "" {
					t.Errorf("Expected empty title, got '%s'", img.Title)
				}
			},
		},
		{
			name: "Create image with IsPrimary not set (should default to false)",
			image: models.Image{
				File:          "/uploads/sku-image-1.jpg",
				Title:         "SKU Image",
				ImageableID:   10,
				ImageableType: "Sku",
				Base: models.Base{
					CreatedBy: testUser.ID,
					UpdatedBy: testUser.ID,
				},
			},
			expectError: false,
			validateFunc: func(t *testing.T, img *models.Image) {
				if img.IsPrimary {
					t.Error("Expected IsPrimary to default to false")
				}
			},
		},
		{
			name: "Create image for Product (polymorphic relationship)",
			image: models.Image{
				File:          "/uploads/product-banner.jpg",
				Title:         "Product Banner",
				ImageableID:   5,
				ImageableType: "Product",
				Base: models.Base{
					CreatedBy: testUser.ID,
					UpdatedBy: testUser.ID,
				},
			},
			expectError: false,
			validateFunc: func(t *testing.T, img *models.Image) {
				if img.ImageableType != "Product" {
					t.Errorf("Expected ImageableType 'Product', got '%s'", img.ImageableType)
				}
				if img.ImageableID != 5 {
					t.Errorf("Expected ImageableID 5, got %d", img.ImageableID)
				}
			},
		},
		{
			name: "Create image for Sku (polymorphic relationship)",
			image: models.Image{
				File:          "/uploads/sku-variant.jpg",
				Title:         "SKU Variant Image",
				ImageableID:   15,
				ImageableType: "Sku",
				Base: models.Base{
					CreatedBy: testUser.ID,
					UpdatedBy: testUser.ID,
				},
			},
			expectError: false,
			validateFunc: func(t *testing.T, img *models.Image) {
				if img.ImageableType != "Sku" {
					t.Errorf("Expected ImageableType 'Sku', got '%s'", img.ImageableType)
				}
			},
		},
		{
			name: "Create image with empty File (zero value for NOT NULL field)",
			image: models.Image{
				File:          "", // Zero value for string
				Title:         "Test Image",
				ImageableID:   1,
				ImageableType: "Product",
				Base: models.Base{
					CreatedBy: testUser.ID,
					UpdatedBy: testUser.ID,
				},
			},
			expectError: false, // NOT NULL allows zero values
			validateFunc: func(t *testing.T, img *models.Image) {
				if img.File != "" {
					t.Errorf("Expected empty file, got '%s'", img.File)
				}
			},
		},
		{
			name: "Create image with zero ImageableID (zero value for NOT NULL field)",
			image: models.Image{
				File:          "/uploads/test.jpg",
				Title:         "Test Image",
				ImageableID:   0, // Zero value for uint
				ImageableType: "Product",
				Base: models.Base{
					CreatedBy: testUser.ID,
					UpdatedBy: testUser.ID,
				},
			},
			expectError: false, // NOT NULL allows zero values
			validateFunc: func(t *testing.T, img *models.Image) {
				if img.ImageableID != 0 {
					t.Errorf("Expected ImageableID 0, got %d", img.ImageableID)
				}
			},
		},
		{
			name: "Create image with empty ImageableType (zero value for NOT NULL field)",
			image: models.Image{
				File:          "/uploads/test.jpg",
				Title:         "Test Image",
				ImageableID:   1,
				ImageableType: "", // Zero value for string
				Base: models.Base{
					CreatedBy: testUser.ID,
					UpdatedBy: testUser.ID,
				},
			},
			expectError: false, // NOT NULL allows zero values
			validateFunc: func(t *testing.T, img *models.Image) {
				if img.ImageableType != "" {
					t.Errorf("Expected empty ImageableType, got '%s'", img.ImageableType)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := db.Create(&tc.image)

			if tc.expectError {
				if result.Error == nil {
					t.Error("Expected error but got none")
				}
				if tc.errorMsg != "" && result.Error != nil {
					// Optionally check error message
				}
			} else {
				if result.Error != nil {
					t.Errorf("Expected no error but got: %v", result.Error)
				}
				if tc.validateFunc != nil {
					tc.validateFunc(t, &tc.image)
				}
			}
		})
	}
}

func TestImagePrimaryConstraint_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy (required for Base model)
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("Create first primary image for Product", func(t *testing.T) {
		img := models.Image{
			File:          "/uploads/product-1-primary.jpg",
			Title:         "Product 1 Primary",
			IsPrimary:     true,
			ImageableID:   1,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		result := db.Create(&img)
		if result.Error != nil {
			t.Errorf("Expected no error, got: %v", result.Error)
		}

		var images []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", 1, "Product").Find(&images)
		
		if len(images) != 1 {
			t.Errorf("Expected 1 image, got %d", len(images))
		}
		if !images[0].IsPrimary {
			t.Error("Expected image to be primary")
		}
	})

	t.Run("Create second primary image - should set first to non-primary", func(t *testing.T) {
		// Create first primary image
		first := models.Image{
			File:          "/uploads/product-2-first.jpg",
			Title:         "First Primary",
			IsPrimary:     true,
			ImageableID:   2,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&first)

		// Create second primary image
		second := models.Image{
			File:          "/uploads/product-2-second.jpg",
			Title:         "Second Primary",
			IsPrimary:     true,
			ImageableID:   2,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		result := db.Create(&second)
		if result.Error != nil {
			t.Errorf("Expected no error, got: %v", result.Error)
		}

		// Verify only one primary image exists
		var images []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", 2, "Product").Find(&images)
		
		if len(images) != 2 {
			t.Errorf("Expected 2 images, got %d", len(images))
		}

		primaryCount := 0
		for _, img := range images {
			if img.IsPrimary {
				primaryCount++
				if img.Title != "Second Primary" {
					t.Error("Expected new image to be primary")
				}
			}
		}

		if primaryCount != 1 {
			t.Errorf("Expected 1 primary image, got %d", primaryCount)
		}
	})

	t.Run("Multiple non-primary images are allowed", func(t *testing.T) {
		// Create multiple non-primary images
		images := []models.Image{
			{
				File:          "/uploads/product-3-img1.jpg",
				IsPrimary:     false,
				ImageableID:   3,
				ImageableType: "Product",
				Base: models.Base{
					CreatedBy: testUser.ID,
					UpdatedBy: testUser.ID,
				},
			},
			{
				File:          "/uploads/product-3-img2.jpg",
				IsPrimary:     false,
				ImageableID:   3,
				ImageableType: "Product",
				Base: models.Base{
					CreatedBy: testUser.ID,
					UpdatedBy: testUser.ID,
				},
			},
			{
				File:          "/uploads/product-3-img3.jpg",
				IsPrimary:     false,
				ImageableID:   3,
				ImageableType: "Product",
				Base: models.Base{
					CreatedBy: testUser.ID,
					UpdatedBy: testUser.ID,
				},
			},
		}

		for _, img := range images {
			result := db.Create(&img)
			if result.Error != nil {
				t.Errorf("Expected no error, got: %v", result.Error)
			}
		}

		var allImages []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", 3, "Product").Find(&allImages)
		
		if len(allImages) != 3 {
			t.Errorf("Expected 3 images, got %d", len(allImages))
		}

		for _, img := range allImages {
			if img.IsPrimary {
				t.Error("Expected all images to be non-primary")
			}
		}
	})

	t.Run("Primary image for different imageable should not affect others", func(t *testing.T) {
		// Create primary image for Product 4
		img4 := models.Image{
			File:          "/uploads/product-4-primary.jpg",
			IsPrimary:     true,
			ImageableID:   4,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&img4)

		// Create primary image for Product 5
		img5 := models.Image{
			File:          "/uploads/product-5-primary.jpg",
			IsPrimary:     true,
			ImageableID:   5,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&img5)

		// Verify both have primary images
		var product4Images []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", 4, "Product").Find(&product4Images)
		if !product4Images[0].IsPrimary {
			t.Error("Product 4 should still have primary image")
		}

		var product5Images []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", 5, "Product").Find(&product5Images)
		if !product5Images[0].IsPrimary {
			t.Error("Product 5 should have primary image")
		}
	})

	t.Run("Primary image for different ImageableType should not affect others", func(t *testing.T) {
		// Create primary image for Product 6
		productImg := models.Image{
			File:          "/uploads/product-6-primary.jpg",
			IsPrimary:     true,
			ImageableID:   6,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&productImg)

		// Create primary image for Sku 6 (same ID, different type)
		skuImg := models.Image{
			File:          "/uploads/sku-6-primary.jpg",
			IsPrimary:     true,
			ImageableID:   6,
			ImageableType: "Sku",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&skuImg)

		// Verify both have primary images
		var productImages []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", 6, "Product").Find(&productImages)
		if !productImages[0].IsPrimary {
			t.Error("Product 6 should still have primary image")
		}

		var skuImages []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", 6, "Sku").Find(&skuImages)
		if !skuImages[0].IsPrimary {
			t.Error("Sku 6 should have primary image")
		}
	})
}

func TestImageUpdate_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy (required for Base model)
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("Update non-primary to primary", func(t *testing.T) {
		img := models.Image{
			File:          "/uploads/product-10-img1.jpg",
			IsPrimary:     false,
			ImageableID:   10,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&img)

		// Update to primary
		img.IsPrimary = true
		result := db.Save(&img)
		if result.Error != nil {
			t.Errorf("Expected no error, got: %v", result.Error)
		}

		// Verify update
		var updated models.Image
		db.First(&updated, img.ID)
		if !updated.IsPrimary {
			t.Error("Expected image to be primary after update")
		}
	})

	t.Run("Update to primary when another is primary - should demote old primary", func(t *testing.T) {
		// Create first primary
		first := models.Image{
			File:          "/uploads/product-11-first.jpg",
			IsPrimary:     true,
			ImageableID:   11,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&first)

		// Create second non-primary
		second := models.Image{
			File:          "/uploads/product-11-second.jpg",
			IsPrimary:     false,
			ImageableID:   11,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&second)

		// Update second to primary
		second.IsPrimary = true
		result := db.Save(&second)
		if result.Error != nil {
			t.Errorf("Expected no error, got: %v", result.Error)
		}

		// Verify only one primary exists
		var images []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", 11, "Product").Find(&images)

		primaryCount := 0
		for _, img := range images {
			if img.IsPrimary {
				primaryCount++
				if img.ID != second.ID {
					t.Error("Expected updated image to be primary")
				}
			}
		}

		if primaryCount != 1 {
			t.Errorf("Expected 1 primary image, got %d", primaryCount)
		}
	})

	t.Run("Update file path and title", func(t *testing.T) {
		img := models.Image{
			File:          "/uploads/old-path.jpg",
			Title:         "Original Title",
			ImageableID:   12,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&img)

		// Update fields
		img.File = "/uploads/new-path.jpg"
		img.Title = "Updated Title"
		result := db.Save(&img)
		if result.Error != nil {
			t.Errorf("Expected no error, got: %v", result.Error)
		}

		// Verify updates
		var updated models.Image
		db.First(&updated, img.ID)
		if updated.File != "/uploads/new-path.jpg" {
			t.Errorf("Expected file '/uploads/new-path.jpg', got '%s'", updated.File)
		}
		if updated.Title != "Updated Title" {
			t.Errorf("Expected title 'Updated Title', got '%s'", updated.Title)
		}
	})
}

func TestImageQuery_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy (required for Base model)
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Setup test data
	db.Create(&models.Image{
		File:          "/uploads/product-20-img1.jpg",
		Title:         "Image 1",
		IsPrimary:     true,
		ImageableID:   20,
		ImageableType: "Product",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	})
	db.Create(&models.Image{
		File:          "/uploads/product-20-img2.jpg",
		Title:         "Image 2",
		IsPrimary:     false,
		ImageableID:   20,
		ImageableType: "Product",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	})
	db.Create(&models.Image{
		File:          "/uploads/sku-20-img1.jpg",
		Title:         "SKU Image",
		IsPrimary:     true,
		ImageableID:   20,
		ImageableType: "Sku",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	})

	t.Run("Find all images for a Product", func(t *testing.T) {
		var images []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", 20, "Product").Find(&images)
		if len(images) != 2 {
			t.Errorf("Expected 2 images, got %d", len(images))
		}
	})

	t.Run("Find primary image for a Product", func(t *testing.T) {
		var image models.Image
		err := db.Where("imageable_id = ? AND imageable_type = ? AND is_primary = ?",
			20, "Product", true).First(&image).Error
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if image.Title != "Image 1" {
			t.Errorf("Expected title 'Image 1', got '%s'", image.Title)
		}
	})

	t.Run("Find all images for a Sku", func(t *testing.T) {
		var images []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", 20, "Sku").Find(&images)
		if len(images) != 1 {
			t.Errorf("Expected 1 image, got %d", len(images))
		}
		if images[0].Title != "SKU Image" {
			t.Errorf("Expected title 'SKU Image', got '%s'", images[0].Title)
		}
	})

	t.Run("Count images per imageable", func(t *testing.T) {
		var count int64
		db.Model(&models.Image{}).
			Where("imageable_id = ? AND imageable_type = ?", 20, "Product").
			Count(&count)
		if count != 2 {
			t.Errorf("Expected count 2, got %d", count)
		}
	})
}

func TestImageDeletion_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user for CreatedBy/UpdatedBy (required for Base model)
	testUser := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(&testUser).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("Delete non-primary image", func(t *testing.T) {
		img := models.Image{
			File:          "/uploads/delete-test-1.jpg",
			IsPrimary:     false,
			ImageableID:   30,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&img)

		result := db.Delete(&img)
		if result.Error != nil {
			t.Errorf("Expected no error, got: %v", result.Error)
		}
		if result.RowsAffected != 1 {
			t.Errorf("Expected 1 row affected, got %d", result.RowsAffected)
		}

		// Verify deletion
		var found models.Image
		err := db.First(&found, img.ID).Error
		if err == nil {
			t.Error("Expected error (not found), but image still exists")
		}
	})

	t.Run("Delete primary image", func(t *testing.T) {
		primary := models.Image{
			File:          "/uploads/delete-test-2-primary.jpg",
			IsPrimary:     true,
			ImageableID:   31,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&primary)

		result := db.Delete(&primary)
		if result.Error != nil {
			t.Errorf("Expected no error, got: %v", result.Error)
		}

		// Verify deletion
		var found models.Image
		err := db.First(&found, primary.ID).Error
		if err == nil {
			t.Error("Expected error (not found), but image still exists")
		}
	})

	t.Run("Soft delete sets DeletedAt", func(t *testing.T) {
		// Create another image to test soft delete
		softDeleteImg := models.Image{
			File:          "/uploads/soft-delete-test.jpg",
			IsPrimary:     false,
			ImageableID:   32,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&softDeleteImg)

		// Soft delete
		result := db.Delete(&softDeleteImg)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify soft delete - should not be found in normal query
		var foundImage models.Image
		result = db.First(&foundImage, softDeleteImg.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected image to be soft deleted (not found in normal query)")
		}

		// Verify soft delete - should be found with Unscoped
		result = db.Unscoped().First(&foundImage, softDeleteImg.ID)
		if result.Error != nil {
			t.Errorf("Expected to find soft deleted image with Unscoped but got error: %v", result.Error)
		}
		if foundImage.DeletedAt.Time.IsZero() {
			t.Error("Expected DeletedAt to be set after soft delete")
		}
	})

	t.Run("Permanent delete image", func(t *testing.T) {
		// Create another image
		permanentDeleteImg := models.Image{
			File:          "/uploads/permanent-delete.jpg",
			IsPrimary:     false,
			ImageableID:   33,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&permanentDeleteImg)

		// Permanently delete
		result := db.Unscoped().Delete(&permanentDeleteImg)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify permanent delete - should not be found even with Unscoped
		var foundImage models.Image
		result = db.Unscoped().First(&foundImage, permanentDeleteImg.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected image to be permanently deleted")
		}
	})
}

func TestImageRelationships_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create test user
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

	// Create product
	product := models.Product{
		Name:        "Laptop Computer",
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

	// Create SKU
	sku := models.Sku{
		Name:      "Laptop 16GB RAM",
		SkuNumber: "LAP-16GB",
		Price:     999.00,
		ProductID: product.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&sku).Error; err != nil {
		t.Fatalf("Failed to create sku: %v", err)
	}

	// Create images for product
	productImages := []models.Image{
		{
			File:          "/uploads/laptop-front.jpg",
			Title:         "Front View",
			IsPrimary:     true,
			ImageableID:   product.ID,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		},
		{
			File:          "/uploads/laptop-side.jpg",
			Title:         "Side View",
			IsPrimary:     false,
			ImageableID:   product.ID,
			ImageableType: "Product",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		},
	}
	for _, img := range productImages {
		if err := db.Create(&img).Error; err != nil {
			t.Fatalf("Failed to create product image: %v", err)
		}
	}

	// Create images for SKU
	skuImages := []models.Image{
		{
			File:          "/uploads/laptop-16gb-front.jpg",
			Title:         "16GB Front",
			IsPrimary:     true,
			ImageableID:   sku.ID,
			ImageableType: "Sku",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		},
		{
			File:          "/uploads/laptop-16gb-back.jpg",
			Title:         "16GB Back",
			IsPrimary:     false,
			ImageableID:   sku.ID,
			ImageableType: "Sku",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		},
	}
	for _, img := range skuImages {
		if err := db.Create(&img).Error; err != nil {
			t.Fatalf("Failed to create sku image: %v", err)
		}
	}

	t.Run("Product has images", func(t *testing.T) {
		var foundImages []models.Image
		result := db.Where("imageable_id = ? AND imageable_type = ?", product.ID, "Product").Find(&foundImages)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(foundImages) != 2 {
			t.Errorf("Expected 2 images for product, got %d", len(foundImages))
		}
	})

	t.Run("Sku has images", func(t *testing.T) {
		var foundImages []models.Image
		result := db.Where("imageable_id = ? AND imageable_type = ?", sku.ID, "Sku").Find(&foundImages)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(foundImages) != 2 {
			t.Errorf("Expected 2 images for sku, got %d", len(foundImages))
		}
	})

	t.Run("Find primary image for product", func(t *testing.T) {
		var primaryImage models.Image
		result := db.Where("imageable_id = ? AND imageable_type = ? AND is_primary = ?",
			product.ID, "Product", true).First(&primaryImage)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if primaryImage.Title != "Front View" {
			t.Errorf("Expected title 'Front View', got '%s'", primaryImage.Title)
		}
	})

	t.Run("Find primary image for sku", func(t *testing.T) {
		var primaryImage models.Image
		result := db.Where("imageable_id = ? AND imageable_type = ? AND is_primary = ?",
			sku.ID, "Sku", true).First(&primaryImage)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if primaryImage.Title != "16GB Front" {
			t.Errorf("Expected title '16GB Front', got '%s'", primaryImage.Title)
		}
	})

	t.Run("Polymorphic - Product and Sku images are separate", func(t *testing.T) {
		// Count all Product images
		var productCount int64
		db.Model(&models.Image{}).Where("imageable_type = ?", "Product").Count(&productCount)
		if productCount != 2 {
			t.Errorf("Expected 2 Product images, got %d", productCount)
		}

		// Count all Sku images
		var skuCount int64
		db.Model(&models.Image{}).Where("imageable_type = ?", "Sku").Count(&skuCount)
		if skuCount != 2 {
			t.Errorf("Expected 2 Sku images, got %d", skuCount)
		}

		// Verify images are properly isolated
		var productImgs []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", product.ID, "Product").Find(&productImgs)
		for _, img := range productImgs {
			if img.ImageableType != "Product" {
				t.Errorf("Expected ImageableType 'Product', got '%s'", img.ImageableType)
			}
		}

		var skuImgs []models.Image
		db.Where("imageable_id = ? AND imageable_type = ?", sku.ID, "Sku").Find(&skuImgs)
		for _, img := range skuImgs {
			if img.ImageableType != "Sku" {
				t.Errorf("Expected ImageableType 'Sku', got '%s'", img.ImageableType)
			}
		}
	})
}
