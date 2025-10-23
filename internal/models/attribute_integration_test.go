//go:build integration
// +build integration

package models_test

import (
	"testing"

	"github.com/Wilson1510/klampis-pim-go/internal/models"
	"github.com/Wilson1510/klampis-pim-go/internal/testutil"
	"gorm.io/gorm"
)

func TestAttributeCreate_Integration(t *testing.T) {
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
		name        string
		attribute   models.Attribute
		expectError bool
	}{
		{
			name: "Create valid TEXT attribute",
			attribute: models.Attribute{
				Name:     "Color",
				Code:     "color",
				DataType: models.DataTypeText,
			},
			expectError: false,
		},
		{
			name: "Create valid NUMBER attribute with UOM",
			attribute: models.Attribute{
				Name:     "RAM",
				Code:     "ram",
				DataType: models.DataTypeNumber,
				UOM:      "GB",
			},
			expectError: false,
		},
		{
			name: "Create valid BOOLEAN attribute",
			attribute: models.Attribute{
				Name:     "Is Wireless",
				Code:     "is_wireless",
				DataType: models.DataTypeBoolean,
			},
			expectError: false,
		},
		{
			name: "Create valid DATE attribute",
			attribute: models.Attribute{
				Name:     "Warranty Expiry",
				Code:     "warranty_expiry",
				DataType: models.DataTypeDate,
			},
			expectError: false,
		},
		{
			name: "Create attribute with invalid data type",
			attribute: models.Attribute{
				Name:     "Invalid",
				Code:     "invalid",
				DataType: models.DataType("INVALID"),
			},
			expectError: true,
		},
		{
			name: "Create attribute with duplicate code",
			attribute: models.Attribute{
				Name:     "Color Duplicate",
				Code:     "color", // Already exists
				DataType: models.DataTypeText,
			},
			expectError: true,
		},
		{
			name: "Create attribute without name",
			attribute: models.Attribute{
				Code:     "no_name",
				DataType: models.DataTypeText,
			},
			expectError: false, // golang will automatically set the name to ""
		},
		{
			name: "Create attribute without code",
			attribute: models.Attribute{
				Name:     "No Code",
				DataType: models.DataTypeText,
			},
			expectError: false, // golang will automatically set the code to ""
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.attribute.CreatedBy = testUser.ID
			tc.attribute.UpdatedBy = testUser.ID
			result := db.Create(&tc.attribute)

			if tc.expectError {
				if result.Error == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if result.Error != nil {
					t.Errorf("Expected no error but got: %v", result.Error)
				}
				if tc.attribute.ID == 0 {
					t.Error("Expected attribute ID to be set after creation")
				}
				if tc.attribute.CreatedAt.IsZero() {
					t.Error("Expected CreatedAt to be set after creation")
				}
			}
		})
	}
}

func TestAttributeRead_Integration(t *testing.T) {
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

	// Create test attributes
	attributes := []models.Attribute{
		{
			Name:     "RAM",
			Code:     "ram",
			DataType: models.DataTypeNumber,
			UOM:      "GB",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		},
		{
			Name:     "Color",
			Code:     "color",
			DataType: models.DataTypeText,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		},
	}

	for i := range attributes {
		if err := db.Create(&attributes[i]).Error; err != nil {
			t.Fatalf("Failed to create test attribute: %v", err)
		}
	}

	t.Run("Read attribute by ID", func(t *testing.T) {
		var foundAttribute models.Attribute
		result := db.First(&foundAttribute, attributes[0].ID)

		if result.Error != nil {
			t.Errorf("Expected to find attribute but got error: %v", result.Error)
		}
		if foundAttribute.Name != attributes[0].Name {
			t.Errorf("Expected name '%s', got '%s'", attributes[0].Name, foundAttribute.Name)
		}
		if foundAttribute.Code != attributes[0].Code {
			t.Errorf("Expected code '%s', got '%s'", attributes[0].Code, foundAttribute.Code)
		}
		if foundAttribute.DataType != attributes[0].DataType {
			t.Errorf("Expected data type '%s', got '%s'", attributes[0].DataType, foundAttribute.DataType)
		}
		if foundAttribute.UOM != attributes[0].UOM {
			t.Errorf("Expected UOM '%s', got '%s'", attributes[0].UOM, foundAttribute.UOM)
		}
	})

	t.Run("Read attribute by code", func(t *testing.T) {
		var foundAttribute models.Attribute
		result := db.Where("code = ?", attributes[0].Code).First(&foundAttribute)

		if result.Error != nil {
			t.Errorf("Expected to find attribute but got error: %v", result.Error)
		}
		if foundAttribute.ID != attributes[0].ID {
			t.Errorf("Expected ID %d, got %d", attributes[0].ID, foundAttribute.ID)
		}
	})

	t.Run("Read non-existent attribute", func(t *testing.T) {
		var foundAttribute models.Attribute
		result := db.First(&foundAttribute, 99999)

		if result.Error != gorm.ErrRecordNotFound {
			t.Errorf("Expected ErrRecordNotFound but got: %v", result.Error)
		}
	})
}

func TestAttributeUpdate_Integration(t *testing.T) {
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

	// Create a test attribute
	attribute := models.Attribute{
		Name:     "Screen Size",
		Code:     "screen_size",
		DataType: models.DataTypeNumber,
		UOM:      "inch",
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&attribute).Error; err != nil {
		t.Fatalf("Failed to create test attribute: %v", err)
	}

	t.Run("Update attribute name", func(t *testing.T) {
		attribute.Name = "Display Size"
		result := db.Save(&attribute)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedAttribute models.Attribute
		db.First(&updatedAttribute, attribute.ID)
		if updatedAttribute.Name != "Display Size" {
			t.Errorf("Expected name 'Display Size', got '%s'", updatedAttribute.Name)
		}
	})

	t.Run("Update attribute UOM", func(t *testing.T) {
		attribute.UOM = "inches"
		result := db.Save(&attribute)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedAttribute models.Attribute
		db.First(&updatedAttribute, attribute.ID)
		if updatedAttribute.UOM != "inches" {
			t.Errorf("Expected UOM 'inches', got '%s'", updatedAttribute.UOM)
		}
	})

	t.Run("Update attribute data type to valid type", func(t *testing.T) {
		attribute.DataType = models.DataTypeText
		result := db.Save(&attribute)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedAttribute models.Attribute
		db.First(&updatedAttribute, attribute.ID)
		if updatedAttribute.DataType != models.DataTypeText {
			t.Errorf("Expected data type TEXT, got '%s'", updatedAttribute.DataType)
		}
	})

	t.Run("Update attribute data type to invalid type", func(t *testing.T) {
		attribute.DataType = models.DataType("INVALID")
		result := db.Save(&attribute)

		if result.Error == nil {
			t.Error("Expected error for invalid data type but got none")
		}
	})
}

func TestAttributeDelete_Integration(t *testing.T) {
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

	// Create a test attribute
	attribute := models.Attribute{
		Name:     "Delete Attribute",
		Code:     "delete_attr",
		DataType: models.DataTypeText,
		Base: models.Base{
			CreatedBy: testUser.ID,
			UpdatedBy: testUser.ID,
		},
	}
	if err := db.Create(&attribute).Error; err != nil {
		t.Fatalf("Failed to create test attribute: %v", err)
	}

	t.Run("Soft delete attribute", func(t *testing.T) {
		result := db.Delete(&attribute)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify soft delete - should not be found in normal query
		var foundAttribute models.Attribute
		result = db.First(&foundAttribute, attribute.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected attribute to be soft deleted (not found in normal query)")
		}

		// Verify soft delete - should be found with Unscoped
		result = db.Unscoped().First(&foundAttribute, attribute.ID)
		if result.Error != nil {
			t.Errorf("Expected to find soft deleted attribute with Unscoped but got error: %v", result.Error)
		}
		if foundAttribute.DeletedAt.Time.IsZero() {
			t.Error("Expected DeletedAt to be set after soft delete")
		}
	})

	t.Run("Permanent delete attribute", func(t *testing.T) {
		// Create another attribute
		anotherAttribute := models.Attribute{
			Name:     "Permanent Delete",
			Code:     "perm_delete",
			DataType: models.DataTypeText,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&anotherAttribute)

		// Permanently delete
		result := db.Unscoped().Delete(&anotherAttribute)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify permanent delete - should not be found even with Unscoped
		var foundAttribute models.Attribute
		result = db.Unscoped().First(&foundAttribute, anotherAttribute.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected attribute to be permanently deleted")
		}
	})
}

func TestAttributeQuery_Integration(t *testing.T) {
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

	// Create multiple test attributes
	attributes := []models.Attribute{
		{Name: "RAM", Code: "ram", DataType: models.DataTypeNumber, UOM: "GB", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
		{Name: "Storage", Code: "storage", DataType: models.DataTypeNumber, UOM: "GB", Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
		{Name: "Color", Code: "color", DataType: models.DataTypeText, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
		{Name: "Is Wireless", Code: "is_wireless", DataType: models.DataTypeBoolean, Base: models.Base{CreatedBy: testUser.ID, UpdatedBy: testUser.ID}},
	}

	for _, attr := range attributes {
		if err := db.Create(&attr).Error; err != nil {
			t.Fatalf("Failed to create test attribute: %v", err)
		}
	}

	t.Run("Find all attributes", func(t *testing.T) {
		var allAttributes []models.Attribute
		result := db.Find(&allAttributes)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(allAttributes) != 4 {
			t.Errorf("Expected 4 attributes, got %d", len(allAttributes))
		}
	})

	t.Run("Find attributes by data type", func(t *testing.T) {
		var numberAttributes []models.Attribute
		result := db.Where("data_type = ?", models.DataTypeNumber).Find(&numberAttributes)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(numberAttributes) != 2 {
			t.Errorf("Expected 2 number attributes, got %d", len(numberAttributes))
		}
		for _, attr := range numberAttributes {
			if attr.DataType != models.DataTypeNumber {
				t.Errorf("Expected data type NUMBER, got %s", attr.DataType)
			}
		}
	})

	t.Run("Count attributes", func(t *testing.T) {
		var count int64
		result := db.Model(&models.Attribute{}).Count(&count)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if count != 4 {
			t.Errorf("Expected count 4, got %d", count)
		}
	})

	t.Run("Find attributes with pagination", func(t *testing.T) {
		var paginatedAttributes []models.Attribute
		result := db.Limit(2).Offset(0).Find(&paginatedAttributes)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(paginatedAttributes) != 2 {
			t.Errorf("Expected 2 attributes in first page, got %d", len(paginatedAttributes))
		}
	})

	t.Run("Search attributes by name", func(t *testing.T) {
		var foundAttributes []models.Attribute
		result := db.Where("name LIKE ?", "%RAM%").Find(&foundAttributes)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(foundAttributes) != 1 {
			t.Errorf("Expected 1 attribute, got %d", len(foundAttributes))
		}
		if len(foundAttributes) > 0 && foundAttributes[0].Name != "RAM" {
			t.Errorf("Expected attribute name 'RAM', got '%s'", foundAttributes[0].Name)
		}
	})
}

func TestAttributeValueParsing_Integration(t *testing.T) {
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

	t.Run("Parse TEXT value", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "Color",
			Code:     "color",
			DataType: models.DataTypeText,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		value, err := attribute.ParseValue("Red")
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if value != "Red" {
			t.Errorf("Expected 'Red', got '%v'", value)
		}
	})

	t.Run("Parse NUMBER value", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "RAM",
			Code:     "ram",
			DataType: models.DataTypeNumber,
			UOM:      "GB",
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		value, err := attribute.ParseValue("16")
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if floatVal, ok := value.(float64); !ok || floatVal != 16.0 {
			t.Errorf("Expected float64(16.0), got %v", value)
		}
	})

	t.Run("Parse BOOLEAN value", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "Is Wireless",
			Code:     "is_wireless",
			DataType: models.DataTypeBoolean,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		value, err := attribute.ParseValue("true")
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if boolVal, ok := value.(bool); !ok || !boolVal {
			t.Errorf("Expected bool(true), got %v", value)
		}
	})

	t.Run("Parse DATE value", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "Warranty Expiry",
			Code:     "warranty_expiry",
			DataType: models.DataTypeDate,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		value, err := attribute.ParseValue("2024-12-31")
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if _, ok := value.(interface{}); !ok {
			t.Errorf("Expected time.Time value, got %T", value)
		}
	})

	t.Run("Parse invalid NUMBER value", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "RAM",
			Code:     "ram2",
			DataType: models.DataTypeNumber,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		_, err := attribute.ParseValue("invalid")
		if err == nil {
			t.Error("Expected error for invalid number but got none")
		}
	})
}

func TestAttributeValueFormatting_Integration(t *testing.T) {
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

	t.Run("Format TEXT value", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "Color",
			Code:     "color_fmt",
			DataType: models.DataTypeText,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		formatted, err := attribute.FormatValue("Blue")
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if formatted != "Blue" {
			t.Errorf("Expected 'Blue', got '%s'", formatted)
		}
	})

	t.Run("Format NUMBER value (int)", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "RAM",
			Code:     "ram_fmt",
			DataType: models.DataTypeNumber,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		formatted, err := attribute.FormatValue(16)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if formatted != "16" {
			t.Errorf("Expected '16', got '%s'", formatted)
		}
	})

	t.Run("Format NUMBER value (float)", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "Price",
			Code:     "price_fmt",
			DataType: models.DataTypeNumber,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		formatted, err := attribute.FormatValue(99.99)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if formatted != "99.99" {
			t.Errorf("Expected '99.99', got '%s'", formatted)
		}
	})

	t.Run("Format BOOLEAN value", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "Is Wireless",
			Code:     "is_wireless_fmt",
			DataType: models.DataTypeBoolean,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		formatted, err := attribute.FormatValue(true)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if formatted != "true" {
			t.Errorf("Expected 'true', got '%s'", formatted)
		}
	})
}

func TestAttributeValidateValue_Integration(t *testing.T) {
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

	t.Run("Validate valid NUMBER value", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "RAM",
			Code:     "ram_validate",
			DataType: models.DataTypeNumber,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		err := attribute.ValidateValue("16")
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})

	t.Run("Validate invalid NUMBER value", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "RAM",
			Code:     "ram_validate2",
			DataType: models.DataTypeNumber,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		err := attribute.ValidateValue("not a number")
		if err == nil {
			t.Error("Expected error for invalid number but got none")
		}
	})

	t.Run("Validate valid BOOLEAN value", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "Is Wireless",
			Code:     "is_wireless_validate",
			DataType: models.DataTypeBoolean,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		err := attribute.ValidateValue("true")
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	})

	t.Run("Validate invalid BOOLEAN value", func(t *testing.T) {
		attribute := models.Attribute{
			Name:     "Is Wireless",
			Code:     "is_wireless_validate2",
			DataType: models.DataTypeBoolean,
			Base: models.Base{
				CreatedBy: testUser.ID,
				UpdatedBy: testUser.ID,
			},
		}
		db.Create(&attribute)

		err := attribute.ValidateValue("not a boolean")
		if err == nil {
			t.Error("Expected error for invalid boolean but got none")
		}
	})
}

