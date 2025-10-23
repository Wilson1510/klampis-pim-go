//go:build integration
// +build integration

package models_test

import (
	"testing"

	"github.com/Wilson1510/klampis-pim-go/internal/models"
	"github.com/Wilson1510/klampis-pim-go/internal/testutil"
	"gorm.io/gorm"
)

func TestSkuAttributeValueCreate_Integration(t *testing.T) {
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

	// Create test data: Category, Product, SKU, and Attributes
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

	sku := models.Sku{
		Name:        "Laptop - 16GB",
		SkuNumber:   "LAP-16GB-001",
		Price:       999.99,
		ProductID:   product.ID,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	db.Create(&sku)

	ramAttr := models.Attribute{
		Name:     "RAM",
		Code:     "ram",
		DataType: models.DataTypeNumber,
		UOM:      "GB",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	db.Create(&ramAttr)

	colorAttr := models.Attribute{
		Name:     "Color",
		Code:     "color",
		DataType: models.DataTypeText,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	db.Create(&colorAttr)

	wirelessAttr := models.Attribute{
		Name:     "Is Wireless",
		Code:     "is_wireless",
		DataType: models.DataTypeBoolean,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	db.Create(&wirelessAttr)

	testCases := []struct {
		name        string
		attrValue   models.SkuAttributeValue
		expectError bool
	}{
		{
			name: "Create valid NUMBER attribute value",
			attrValue: models.SkuAttributeValue{
				SkuID:       sku.ID,
				AttributeID: ramAttr.ID,
				Value:       "16",
			},
			expectError: false,
		},
		{
			name: "Create valid TEXT attribute value",
			attrValue: models.SkuAttributeValue{
				SkuID:       sku.ID,
				AttributeID: colorAttr.ID,
				Value:       "Silver",
			},
			expectError: false,
		},
		{
			name: "Create valid BOOLEAN attribute value",
			attrValue: models.SkuAttributeValue{
				SkuID:       sku.ID,
				AttributeID: wirelessAttr.ID,
				Value:       "true",
			},
			expectError: false,
		},
		{
			name: "Create attribute value with invalid NUMBER",
			attrValue: models.SkuAttributeValue{
				SkuID:       sku.ID,
				AttributeID: ramAttr.ID,
				Value:       "not a number",
			},
			expectError: true,
		},
		{
			name: "Create attribute value with invalid BOOLEAN",
			attrValue: models.SkuAttributeValue{
				SkuID:       sku.ID,
				AttributeID: wirelessAttr.ID,
				Value:       "not a boolean",
			},
			expectError: true,
		},
		{
			name: "Create attribute value with non-existent attribute",
			attrValue: models.SkuAttributeValue{
				SkuID:       sku.ID,
				AttributeID: 99999,
				Value:       "test",
			},
			expectError: true,
		},
		{
			name: "Create attribute value with non-existent SKU",
			attrValue: models.SkuAttributeValue{
				SkuID:       99999,
				AttributeID: colorAttr.ID,
				Value:       "Red",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.attrValue.CreatedBy = testUser.ID
			tc.attrValue.UpdatedBy = testUser.ID
			result := db.Create(&tc.attrValue)

			if tc.expectError {
				if result.Error == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if result.Error != nil {
					t.Errorf("Expected no error but got: %v", result.Error)
				}
				if tc.attrValue.ID == 0 {
					t.Error("Expected attribute value ID to be set after creation")
				}
			}
		})
	}
}

func TestSkuAttributeValueRead_Integration(t *testing.T) {
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

	// Create test data
	category := models.Category{
		Name: "Electronics",
		Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&category)

	product := models.Product{
		Name:       "Laptop",
		CategoryID: category.ID,
		Base:       models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&product)

	sku := models.Sku{
		Name:      "Laptop - 16GB",
		SkuNumber: "LAP-16GB-001",
		Price:     999.99,
		ProductID: product.ID,
		Base:      models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&sku)

	ramAttr := models.Attribute{
		Name:     "RAM",
		Code:     "ram",
		DataType: models.DataTypeNumber,
		UOM:      "GB",
		Base:     models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&ramAttr)

	attrValue := models.SkuAttributeValue{
		SkuID:       sku.ID,
		AttributeID: ramAttr.ID,
		Value:       "16",
		CreatedBy:   testUser.ID,
		UpdatedBy:   testUser.ID,
	}
	if err := db.Create(&attrValue).Error; err != nil {
		t.Fatalf("Failed to create test attribute value: %v", err)
	}

	t.Run("Read attribute value by ID", func(t *testing.T) {
		var found models.SkuAttributeValue
		result := db.First(&found, attrValue.ID)

		if result.Error != nil {
			t.Errorf("Expected to find attribute value but got error: %v", result.Error)
		}
		if found.SkuID != attrValue.SkuID {
			t.Errorf("Expected SkuID %d, got %d", attrValue.SkuID, found.SkuID)
		}
		if found.AttributeID != attrValue.AttributeID {
			t.Errorf("Expected AttributeID %d, got %d", attrValue.AttributeID, found.AttributeID)
		}
		if found.Value != attrValue.Value {
			t.Errorf("Expected Value '%s', got '%s'", attrValue.Value, found.Value)
		}
	})

	t.Run("Read attribute value with SKU relationship", func(t *testing.T) {
		var found models.SkuAttributeValue
		result := db.Preload("Sku").First(&found, attrValue.ID)

		if result.Error != nil {
			t.Errorf("Expected to find attribute value but got error: %v", result.Error)
		}
		if found.Sku == nil {
			t.Error("Expected SKU to be loaded")
		} else {
			if found.Sku.ID != sku.ID {
				t.Errorf("Expected SKU ID %d, got %d", sku.ID, found.Sku.ID)
			}
		}
	})

	t.Run("Read attribute value with Attribute relationship", func(t *testing.T) {
		var found models.SkuAttributeValue
		result := db.Preload("Attribute").First(&found, attrValue.ID)

		if result.Error != nil {
			t.Errorf("Expected to find attribute value but got error: %v", result.Error)
		}
		if found.Attribute == nil {
			t.Error("Expected Attribute to be loaded")
		} else {
			if found.Attribute.ID != ramAttr.ID {
				t.Errorf("Expected Attribute ID %d, got %d", ramAttr.ID, found.Attribute.ID)
			}
		}
	})

	t.Run("Read attribute value with all relationships", func(t *testing.T) {
		var found models.SkuAttributeValue
		result := db.Preload("Sku").Preload("Attribute").First(&found, attrValue.ID)

		if result.Error != nil {
			t.Errorf("Expected to find attribute value but got error: %v", result.Error)
		}
		if found.Sku == nil {
			t.Error("Expected SKU to be loaded")
		}
		if found.Attribute == nil {
			t.Error("Expected Attribute to be loaded")
		}
	})

	t.Run("Read non-existent attribute value", func(t *testing.T) {
		var found models.SkuAttributeValue
		result := db.First(&found, 99999)

		if result.Error != gorm.ErrRecordNotFound {
			t.Errorf("Expected ErrRecordNotFound but got: %v", result.Error)
		}
	})
}

func TestSkuAttributeValueUpdate_Integration(t *testing.T) {
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

	// Create test data
	category := models.Category{
		Name: "Electronics",
		Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&category)

	product := models.Product{
		Name:       "Laptop",
		CategoryID: category.ID,
		Base:       models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&product)

	sku := models.Sku{
		Name:      "Laptop - 16GB",
		SkuNumber: "LAP-16GB-001",
		Price:     999.99,
		ProductID: product.ID,
		Base:      models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&sku)

	ramAttr := models.Attribute{
		Name:     "RAM",
		Code:     "ram",
		DataType: models.DataTypeNumber,
		UOM:      "GB",
		Base:     models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&ramAttr)

	attrValue := models.SkuAttributeValue{
		SkuID:       sku.ID,
		AttributeID: ramAttr.ID,
		Value:       "16",
		CreatedBy:   testUser.ID,
		UpdatedBy:   testUser.ID,
	}
	db.Create(&attrValue)

	t.Run("Update attribute value with valid value", func(t *testing.T) {
		attrValue.Value = "32"
		result := db.Save(&attrValue)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updated models.SkuAttributeValue
		db.First(&updated, attrValue.ID)
		if updated.Value != "32" {
			t.Errorf("Expected value '32', got '%s'", updated.Value)
		}
	})

	t.Run("Update attribute value with invalid value", func(t *testing.T) {
		attrValue.Value = "invalid number"
		result := db.Save(&attrValue)

		if result.Error == nil {
			t.Error("Expected error for invalid value but got none")
		}
	})

	t.Run("Update sequence", func(t *testing.T) {
		// Reset to valid value first
		attrValue.Value = "16"
		db.Save(&attrValue)

		attrValue.Sequence = 5
		result := db.Save(&attrValue)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updated models.SkuAttributeValue
		db.First(&updated, attrValue.ID)
		if updated.Sequence != 5 {
			t.Errorf("Expected sequence 5, got %d", updated.Sequence)
		}
	})
}

func TestSkuAttributeValueDelete_Integration(t *testing.T) {
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

	// Create test data
	category := models.Category{
		Name: "Electronics",
		Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&category)

	product := models.Product{
		Name:       "Laptop",
		CategoryID: category.ID,
		Base:       models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&product)

	sku := models.Sku{
		Name:      "Laptop - 16GB",
		SkuNumber: "LAP-16GB-001",
		Price:     999.99,
		ProductID: product.ID,
		Base:      models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&sku)

	ramAttr := models.Attribute{
		Name:     "RAM",
		Code:     "ram",
		DataType: models.DataTypeNumber,
		UOM:      "GB",
		Base:     models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&ramAttr)

	attrValue := models.SkuAttributeValue{
		SkuID:       sku.ID,
		AttributeID: ramAttr.ID,
		Value:       "16",
		CreatedBy:   testUser.ID,
		UpdatedBy:   testUser.ID,
	}
	db.Create(&attrValue)

	t.Run("Soft delete attribute value", func(t *testing.T) {
		result := db.Delete(&attrValue)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify soft delete
		var found models.SkuAttributeValue
		result = db.First(&found, attrValue.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected attribute value to be soft deleted")
		}

		// Verify with Unscoped
		result = db.Unscoped().First(&found, attrValue.ID)
		if result.Error != nil {
			t.Errorf("Expected to find soft deleted attribute value with Unscoped but got error: %v", result.Error)
		}
	})

	t.Run("Permanent delete attribute value", func(t *testing.T) {
		// Create another attribute value
		another := models.SkuAttributeValue{
			SkuID:       sku.ID,
			AttributeID: ramAttr.ID,
			Value:       "32",
			CreatedBy:   testUser.ID,
			UpdatedBy:   testUser.ID,
		}
		db.Create(&another)

		// Permanently delete
		result := db.Unscoped().Delete(&another)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify permanent delete
		var found models.SkuAttributeValue
		result = db.Unscoped().First(&found, another.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected attribute value to be permanently deleted")
		}
	})
}

func TestSkuAttributeValueQuery_Integration(t *testing.T) {
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

	// Create test data
	category := models.Category{
		Name: "Electronics",
		Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&category)

	product := models.Product{
		Name:       "Laptop",
		CategoryID: category.ID,
		Base:       models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&product)

	sku1 := models.Sku{
		Name:      "Laptop - 16GB",
		SkuNumber: "LAP-16GB-001",
		Price:     999.99,
		ProductID: product.ID,
		Base:      models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&sku1)

	sku2 := models.Sku{
		Name:      "Laptop - 32GB",
		SkuNumber: "LAP-32GB-001",
		Price:     1299.99,
		ProductID: product.ID,
		Base:      models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&sku2)

	ramAttr := models.Attribute{
		Name:     "RAM",
		Code:     "ram",
		DataType: models.DataTypeNumber,
		UOM:      "GB",
		Base:     models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&ramAttr)

	colorAttr := models.Attribute{
		Name:     "Color",
		Code:     "color",
		DataType: models.DataTypeText,
		Base:     models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&colorAttr)

	// Create attribute values
	attrValues := []models.SkuAttributeValue{
		{SkuID: sku1.ID, AttributeID: ramAttr.ID, Value: "16", CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
		{SkuID: sku1.ID, AttributeID: colorAttr.ID, Value: "Silver", CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
		{SkuID: sku2.ID, AttributeID: ramAttr.ID, Value: "32", CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
		{SkuID: sku2.ID, AttributeID: colorAttr.ID, Value: "Black", CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}

	for _, av := range attrValues {
		if err := db.Create(&av).Error; err != nil {
			t.Fatalf("Failed to create test attribute value: %v", err)
		}
	}

	t.Run("Find all attribute values", func(t *testing.T) {
		var all []models.SkuAttributeValue
		result := db.Find(&all)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(all) != 4 {
			t.Errorf("Expected 4 attribute values, got %d", len(all))
		}
	})

	t.Run("Find attribute values by SKU", func(t *testing.T) {
		var values []models.SkuAttributeValue
		result := db.Where("sku_id = ?", sku1.ID).Find(&values)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(values) != 2 {
			t.Errorf("Expected 2 attribute values for sku1, got %d", len(values))
		}
	})

	t.Run("Find attribute values by Attribute", func(t *testing.T) {
		var values []models.SkuAttributeValue
		result := db.Where("attribute_id = ?", ramAttr.ID).Find(&values)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(values) != 2 {
			t.Errorf("Expected 2 attribute values for RAM attribute, got %d", len(values))
		}
	})

	t.Run("Count attribute values", func(t *testing.T) {
		var count int64
		result := db.Model(&models.SkuAttributeValue{}).Count(&count)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if count != 4 {
			t.Errorf("Expected count 4, got %d", count)
		}
	})

	t.Run("Find attribute values with relationships", func(t *testing.T) {
		var values []models.SkuAttributeValue
		result := db.Preload("Sku").Preload("Attribute").Find(&values)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		for _, av := range values {
			if av.Sku == nil {
				t.Error("Expected SKU to be loaded")
			}
			if av.Attribute == nil {
				t.Error("Expected Attribute to be loaded")
			}
		}
	})
}

func TestSkuAttributeValueHelperMethods_Integration(t *testing.T) {
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

	// Create test data
	category := models.Category{
		Name: "Electronics",
		Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&category)

	product := models.Product{
		Name:       "Laptop",
		CategoryID: category.ID,
		Base:       models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&product)

	sku := models.Sku{
		Name:      "Laptop - 16GB",
		SkuNumber: "LAP-16GB-001",
		Price:     999.99,
		ProductID: product.ID,
		Base:      models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&sku)

	ramAttr := models.Attribute{
		Name:     "RAM",
		Code:     "ram",
		DataType: models.DataTypeNumber,
		UOM:      "GB",
		Base:     models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&ramAttr)

	t.Run("GetParsedValue returns correct type", func(t *testing.T) {
		attrValue := models.SkuAttributeValue{
			SkuID:       sku.ID,
			AttributeID: ramAttr.ID,
			Value:       "16",
			CreatedBy:   testUser.ID,
			UpdatedBy:   testUser.ID,
		}
		db.Create(&attrValue)

		parsed, err := attrValue.GetParsedValue(db)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if floatVal, ok := parsed.(float64); !ok || floatVal != 16.0 {
			t.Errorf("Expected float64(16.0), got %v", parsed)
		}
	})

	t.Run("GetDisplayValue includes UOM", func(t *testing.T) {
		attrValue := models.SkuAttributeValue{
			SkuID:       sku.ID,
			AttributeID: ramAttr.ID,
			Value:       "32",
			CreatedBy:   testUser.ID,
			UpdatedBy:   testUser.ID,
		}
		db.Create(&attrValue)

		display, err := attrValue.GetDisplayValue(db)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if display != "32 GB" {
			t.Errorf("Expected '32 GB', got '%s'", display)
		}
	})

	t.Run("SetValue converts and validates", func(t *testing.T) {
		attrValue := models.SkuAttributeValue{
			SkuID:       sku.ID,
			AttributeID: ramAttr.ID,
			CreatedBy:   testUser.ID,
			UpdatedBy:   testUser.ID,
		}

		err := attrValue.SetValue(64, db)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if attrValue.Value != "64" {
			t.Errorf("Expected value '64', got '%s'", attrValue.Value)
		}
	})
}

func TestSkuWithAttributeValues_Integration(t *testing.T) {
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

	// Create test data
	category := models.Category{
		Name: "Electronics",
		Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&category)

	product := models.Product{
		Name:       "Laptop",
		CategoryID: category.ID,
		Base:       models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&product)

	sku := models.Sku{
		Name:      "Laptop - 16GB",
		SkuNumber: "LAP-16GB-001",
		Price:     999.99,
		ProductID: product.ID,
		Base:      models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&sku)

	// Create attributes
	ramAttr := models.Attribute{
		Name:     "RAM",
		Code:     "ram",
		DataType: models.DataTypeNumber,
		UOM:      "GB",
		Base:     models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&ramAttr)

	colorAttr := models.Attribute{
		Name:     "Color",
		Code:     "color",
		DataType: models.DataTypeText,
		Base:     models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}
	db.Create(&colorAttr)

	// Create attribute values
	attrValues := []models.SkuAttributeValue{
		{SkuID: sku.ID, AttributeID: ramAttr.ID, Value: "16", Sequence: 1, CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
		{SkuID: sku.ID, AttributeID: colorAttr.ID, Value: "Silver", Sequence: 2, CreatedBy: testUser.ID, UpdatedBy: testUser.ID},
	}

	for _, av := range attrValues {
		if err := db.Create(&av).Error; err != nil {
			t.Fatalf("Failed to create attribute value: %v", err)
		}
	}

	t.Run("Load SKU with attribute values", func(t *testing.T) {
		var loadedSku models.Sku
		result := db.Preload("AttributeValues").First(&loadedSku, sku.ID)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(loadedSku.AttributeValues) != 2 {
			t.Errorf("Expected 2 attribute values, got %d", len(loadedSku.AttributeValues))
		}
	})

	t.Run("Load SKU with attribute values and attributes", func(t *testing.T) {
		var loadedSku models.Sku
		result := db.Preload("AttributeValues.Attribute").First(&loadedSku, sku.ID)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(loadedSku.AttributeValues) != 2 {
			t.Errorf("Expected 2 attribute values, got %d", len(loadedSku.AttributeValues))
		}
		for _, av := range loadedSku.AttributeValues {
			if av.Attribute == nil {
				t.Error("Expected Attribute to be loaded")
			}
		}
	})
}

