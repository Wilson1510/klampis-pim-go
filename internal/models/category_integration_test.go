//go:build integration
// +build integration

package models_test

import (
	"testing"

	"github.com/Wilson1510/klampis-pim-go/internal/models"
	"github.com/Wilson1510/klampis-pim-go/internal/testutil"
	"gorm.io/gorm"
)

func TestCategoryCreate_Integration(t *testing.T) {
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

	testCases := []struct {
		name         string
		category     models.Category
		expectError  bool
		checkSlug    bool
		expectedSlug string
	}{
		{
			name: "Create valid category",
			category: models.Category{
				Name:        "Electronics",
				Description: "Electronic devices and accessories",
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "electronics",
		},
		{
			name: "Create category with special characters in name",
			category: models.Category{
				Name:        "Home & Garden",
				Description: "Home and garden supplies",
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "home-garden",
		},
		{
			name: "Create category without description",
			category: models.Category{
				Name: "Books",
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "books",
		},
		{
			name: "Create category with duplicate name",
			category: models.Category{
				Name:        "Electronics",
				Description: "Another electronics category",
			},
			expectError:  false,
			checkSlug:    true,
			expectedSlug: "electronics-1", // Slug should be made unique
		},
		{
			name: "Create category without name",
			category: models.Category{
				Description: "No name category",
			},
			expectError: false, // golang will automatically set the name to ""
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.category.CreatedBy = testUser.ID
			tc.category.UpdatedBy = testUser.ID
			result := db.Create(&tc.category)

			if tc.expectError {
				if result.Error == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if result.Error != nil {
					t.Errorf("Expected no error but got: %v", result.Error)
				}
				if tc.category.ID == 0 {
					t.Error("Expected category ID to be set after creation")
				}
				if tc.category.CreatedAt.IsZero() {
					t.Error("Expected CreatedAt to be set after creation")
				}
				if tc.checkSlug && tc.category.Slug != tc.expectedSlug {
					t.Errorf("Expected slug '%s', got '%s'", tc.expectedSlug, tc.category.Slug)
				}
			}
		})
	}
}

func TestCategoryRead_Integration(t *testing.T) {
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
	parent := models.Category{
		Name:        "Parent Category",
		Description: "This is a parent category",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("Failed to create parent category: %v", err)
	}

	child := models.Category{
		Name:        "Child Category",
		Description: "This is a child category",
		ParentID:    &parent.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&child).Error; err != nil {
		t.Fatalf("Failed to create child category: %v", err)
	}

	t.Run("Read category by ID", func(t *testing.T) {
		var foundCategory models.Category
		result := db.First(&foundCategory, parent.ID)

		if result.Error != nil {
			t.Errorf("Expected to find category but got error: %v", result.Error)
		}
		if foundCategory.Name != parent.Name {
			t.Errorf("Expected name '%s', got '%s'", parent.Name, foundCategory.Name)
		}
		if foundCategory.Slug != parent.Slug {
			t.Errorf("Expected slug '%s', got '%s'", parent.Slug, foundCategory.Slug)
		}
	})

	t.Run("Read category by slug", func(t *testing.T) {
		var foundCategory models.Category
		result := db.Where("slug = ?", parent.Slug).First(&foundCategory)

		if result.Error != nil {
			t.Errorf("Expected to find category but got error: %v", result.Error)
		}
		if foundCategory.ID != parent.ID {
			t.Errorf("Expected ID %d, got %d", parent.ID, foundCategory.ID)
		}
	})

	t.Run("Read category with parent relationship", func(t *testing.T) {
		var foundCategory models.Category
		result := db.Preload("Parent").First(&foundCategory, child.ID)

		if result.Error != nil {
			t.Errorf("Expected to find category but got error: %v", result.Error)
		}
		if foundCategory.Parent == nil {
			t.Error("Expected parent to be loaded")
		} else if foundCategory.Parent.ID != parent.ID {
			t.Errorf("Expected parent ID %d, got %d", parent.ID, foundCategory.Parent.ID)
		}
	})

	t.Run("Read category with children relationship", func(t *testing.T) {
		var foundCategory models.Category
		result := db.Preload("Children").First(&foundCategory, parent.ID)

		if result.Error != nil {
			t.Errorf("Expected to find category but got error: %v", result.Error)
		}
		if len(foundCategory.Children) != 1 {
			t.Errorf("Expected 1 child, got %d", len(foundCategory.Children))
		}
		if len(foundCategory.Children) > 0 && foundCategory.Children[0].ID != child.ID {
			t.Errorf("Expected child ID %d, got %d", child.ID, foundCategory.Children[0].ID)
		}
	})

	t.Run("Read non-existent category", func(t *testing.T) {
		var foundCategory models.Category
		result := db.First(&foundCategory, 99999)

		if result.Error != gorm.ErrRecordNotFound {
			t.Errorf("Expected ErrRecordNotFound but got: %v", result.Error)
		}
	})
}

func TestCategoryUpdate_Integration(t *testing.T) {
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
		Name:        "Update Category",
		Description: "Original description",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	t.Run("Update category name and verify slug update", func(t *testing.T) {
		category.Name = "Updated Category Name"
		result := db.Save(&category)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedCategory models.Category
		db.First(&updatedCategory, category.ID)
		if updatedCategory.Name != "Updated Category Name" {
			t.Errorf("Expected name 'Updated Category Name', got '%s'", updatedCategory.Name)
		}
		if updatedCategory.Slug != "updated-category-name" {
			t.Errorf("Expected slug to be updated to 'updated-category-name', got '%s'", updatedCategory.Slug)
		}
	})

	t.Run("Update category description", func(t *testing.T) {
		category.Description = "Updated description"
		result := db.Save(&category)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedCategory models.Category
		db.First(&updatedCategory, category.ID)
		if updatedCategory.Description != "Updated description" {
			t.Errorf("Expected description 'Updated description', got '%s'", updatedCategory.Description)
		}
	})

	t.Run("Update category to have a parent", func(t *testing.T) {
		parent := models.Category{
			Name:        "Parent Category",
			Description: "Parent for update test",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&parent)

		category.ParentID = &parent.ID
		result := db.Save(&category)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedCategory models.Category
		db.Preload("Parent").First(&updatedCategory, category.ID)
		if updatedCategory.ParentID == nil {
			t.Error("Expected ParentID to be set")
		} else if *updatedCategory.ParentID != parent.ID {
			t.Errorf("Expected ParentID %d, got %d", parent.ID, *updatedCategory.ParentID)
		}
	})
}

func TestCategoryDelete_Integration(t *testing.T) {
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
		Name:        "Delete Category",
		Description: "To be deleted",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	t.Run("Soft delete category", func(t *testing.T) {
		result := db.Delete(&category)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify soft delete - should not be found in normal query
		var foundCategory models.Category
		result = db.First(&foundCategory, category.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected category to be soft deleted (not found in normal query)")
		}

		// Verify soft delete - should be found with Unscoped
		result = db.Unscoped().First(&foundCategory, category.ID)
		if result.Error != nil {
			t.Errorf("Expected to find soft deleted category with Unscoped but got error: %v", result.Error)
		}
		if foundCategory.DeletedAt.Time.IsZero() {
			t.Error("Expected DeletedAt to be set after soft delete")
		}
	})

	t.Run("Permanent delete category", func(t *testing.T) {
		// Create another category
		anotherCategory := models.Category{
			Name:        "Permanent Delete",
			Description: "To be permanently deleted",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&anotherCategory)

		// Permanently delete
		result := db.Unscoped().Delete(&anotherCategory)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify permanent delete - should not be found even with Unscoped
		var foundCategory models.Category
		result = db.Unscoped().First(&foundCategory, anotherCategory.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected category to be permanently deleted")
		}
	})
}

func TestCategoryQuery_Integration(t *testing.T) {
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

	// Create multiple test categories
	categories := []models.Category{
		{Name: "Electronics", Description: "Electronic devices", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
		{Name: "Books", Description: "All kinds of books", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
		{Name: "Clothing", Description: "Apparel and accessories", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
	}

	for _, cat := range categories {
		if err := db.Create(&cat).Error; err != nil {
			t.Fatalf("Failed to create test category: %v", err)
		}
	}

	t.Run("Find all categories", func(t *testing.T) {
		var allCategories []models.Category
		result := db.Find(&allCategories)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(allCategories) != 3 {
			t.Errorf("Expected 3 categories, got %d", len(allCategories))
		}
	})

	t.Run("Count categories", func(t *testing.T) {
		var count int64
		result := db.Model(&models.Category{}).Count(&count)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if count != 3 {
			t.Errorf("Expected count 3, got %d", count)
		}
	})

	t.Run("Find categories with pagination", func(t *testing.T) {
		var paginatedCategories []models.Category
		result := db.Limit(2).Offset(0).Find(&paginatedCategories)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(paginatedCategories) != 2 {
			t.Errorf("Expected 2 categories in first page, got %d", len(paginatedCategories))
		}
	})

	t.Run("Search categories by name", func(t *testing.T) {
		var foundCategories []models.Category
		result := db.Where("name LIKE ?", "%Book%").Find(&foundCategories)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(foundCategories) != 1 {
			t.Errorf("Expected 1 category, got %d", len(foundCategories))
		}
		if len(foundCategories) > 0 && foundCategories[0].Name != "Books" {
			t.Errorf("Expected category name 'Books', got '%s'", foundCategories[0].Name)
		}
	})
}

func TestCategoryHierarchy_Integration(t *testing.T) {
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

	// Create a category hierarchy
	parent := models.Category{
		Name:        "Electronics",
		Description: "All electronics",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&parent).Error; err != nil {
		t.Fatalf("Failed to create parent category: %v", err)
	}

	child1 := models.Category{
		Name:        "Computers",
		Description: "Computer equipment",
		ParentID:    &parent.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&child1).Error; err != nil {
		t.Fatalf("Failed to create child1 category: %v", err)
	}

	child2 := models.Category{
		Name:        "Phones",
		Description: "Mobile phones",
		ParentID:    &parent.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&child2).Error; err != nil {
		t.Fatalf("Failed to create child2 category: %v", err)
	}

	grandchild := models.Category{
		Name:        "Laptops",
		Description: "Laptop computers",
		ParentID:    &child1.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&grandchild).Error; err != nil {
		t.Fatalf("Failed to create grandchild category: %v", err)
	}

	t.Run("Load parent with all children", func(t *testing.T) {
		var foundParent models.Category
		result := db.Preload("Children").First(&foundParent, parent.ID)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(foundParent.Children) != 2 {
			t.Errorf("Expected 2 children, got %d", len(foundParent.Children))
		}
	})

	t.Run("Load child with parent", func(t *testing.T) {
		var foundChild models.Category
		result := db.Preload("Parent").First(&foundChild, child1.ID)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if foundChild.Parent == nil {
			t.Error("Expected parent to be loaded")
		} else if foundChild.Parent.ID != parent.ID {
			t.Errorf("Expected parent ID %d, got %d", parent.ID, foundChild.Parent.ID)
		}
	})

	t.Run("Load grandchild with parent hierarchy", func(t *testing.T) {
		var foundGrandchild models.Category
		result := db.Preload("Parent.Parent").First(&foundGrandchild, grandchild.ID)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if foundGrandchild.Parent == nil {
			t.Error("Expected parent to be loaded")
		} else {
			if foundGrandchild.Parent.ID != child1.ID {
				t.Errorf("Expected parent ID %d, got %d", child1.ID, foundGrandchild.Parent.ID)
			}
			if foundGrandchild.Parent.Parent == nil {
				t.Error("Expected grandparent to be loaded")
			} else if foundGrandchild.Parent.Parent.ID != parent.ID {
				t.Errorf("Expected grandparent ID %d, got %d", parent.ID, foundGrandchild.Parent.Parent.ID)
			}
		}
	})

	t.Run("Find all root categories (no parent)", func(t *testing.T) {
		var rootCategories []models.Category
		result := db.Where("parent_id IS NULL").Find(&rootCategories)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(rootCategories) != 1 {
			t.Errorf("Expected 1 root category, got %d", len(rootCategories))
		}
		if len(rootCategories) > 0 && rootCategories[0].ID != parent.ID {
			t.Errorf("Expected root category ID %d, got %d", parent.ID, rootCategories[0].ID)
		}
	})

	t.Run("Find all subcategories of a parent", func(t *testing.T) {
		var subcategories []models.Category
		result := db.Where("parent_id = ?", parent.ID).Find(&subcategories)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(subcategories) != 2 {
			t.Errorf("Expected 2 subcategories, got %d", len(subcategories))
		}
	})
}

func TestCategorySlugUniqueness_Integration(t *testing.T) {
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

	t.Run("Create categories with same name generates unique slugs", func(t *testing.T) {
		category1 := models.Category{Name: "Test Category", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
		category2 := models.Category{Name: "Test Category", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}
		category3 := models.Category{Name: "Test Category", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}}

		if err := db.Create(&category1).Error; err != nil {
			t.Fatalf("Failed to create category1: %v", err)
		}
		if err := db.Create(&category2).Error; err != nil {
			t.Fatalf("Failed to create category2: %v", err)
		}
		if err := db.Create(&category3).Error; err != nil {
			t.Fatalf("Failed to create category3: %v", err)
		}

		// Verify slugs are unique
		if category1.Slug != "test-category" {
			t.Errorf("Expected first slug 'test-category', got '%s'", category1.Slug)
		}
		if category2.Slug != "test-category-1" {
			t.Errorf("Expected second slug 'test-category-1', got '%s'", category2.Slug)
		}
		if category3.Slug != "test-category-2" {
			t.Errorf("Expected third slug 'test-category-2', got '%s'", category3.Slug)
		}
	})
}
