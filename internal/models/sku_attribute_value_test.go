package models

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// MockDBForSkuAttr is a mock database for SkuAttributeValue testing
type MockDBForSkuAttr struct {
	attribute  *Attribute
	firstError error
}

// First simulates GORM First operation
func (m *MockDBForSkuAttr) First(dest interface{}, conds ...interface{}) *MockDBForSkuAttr {
	if m.firstError != nil {
		return &MockDBForSkuAttr{firstError: m.firstError}
	}

	if attr, ok := dest.(*Attribute); ok && m.attribute != nil {
		*attr = *m.attribute
	}

	return &MockDBForSkuAttr{}
}

// Error returns the mock error
func (m *MockDBForSkuAttr) Error() error {
	if m.firstError != nil {
		return m.firstError
	}
	return nil
}

// TestValidateValueWithMock tests validateValue method with mocked database
func TestValidateValueWithMock(t *testing.T) {
	testCases := []struct {
		name        string
		skuAttrVal  SkuAttributeValue
		attribute   *Attribute
		dbError     error
		expectError bool
		description string
	}{
		{
			name: "Valid TEXT value",
			skuAttrVal: SkuAttributeValue{
				SkuID:       1,
				AttributeID: 1,
				Value:       "Space Black",
			},
			attribute: &Attribute{
				Name:     "Color",
				Code:     "color",
				DataType: DataTypeText,
			},
			expectError: false,
			description: "TEXT value should be valid",
		},
		{
			name: "Valid NUMBER value",
			skuAttrVal: SkuAttributeValue{
				SkuID:       1,
				AttributeID: 2,
				Value:       "16",
			},
			attribute: &Attribute{
				Name:     "RAM",
				Code:     "ram",
				DataType: DataTypeNumber,
				UOM:      "GB",
			},
			expectError: false,
			description: "NUMBER value should be valid",
		},
		{
			name: "Valid BOOLEAN value",
			skuAttrVal: SkuAttributeValue{
				SkuID:       1,
				AttributeID: 3,
				Value:       "true",
			},
			attribute: &Attribute{
				Name:     "In Stock",
				Code:     "in_stock",
				DataType: DataTypeBoolean,
			},
			expectError: false,
			description: "BOOLEAN value should be valid",
		},
		{
			name: "Valid DATE value",
			skuAttrVal: SkuAttributeValue{
				SkuID:       1,
				AttributeID: 4,
				Value:       "2023-12-31",
			},
			attribute: &Attribute{
				Name:     "Warranty Expiry",
				Code:     "warranty_expiry",
				DataType: DataTypeDate,
			},
			expectError: false,
			description: "DATE value should be valid",
		},
		{
			name: "Invalid NUMBER value",
			skuAttrVal: SkuAttributeValue{
				SkuID:       1,
				AttributeID: 2,
				Value:       "not-a-number",
			},
			attribute: &Attribute{
				Name:     "RAM",
				Code:     "ram",
				DataType: DataTypeNumber,
			},
			expectError: true,
			description: "Invalid NUMBER value should fail",
		},
		{
			name: "Invalid BOOLEAN value",
			skuAttrVal: SkuAttributeValue{
				SkuID:       1,
				AttributeID: 3,
				Value:       "yes",
			},
			attribute: &Attribute{
				Name:     "In Stock",
				Code:     "in_stock",
				DataType: DataTypeBoolean,
			},
			expectError: true,
			description: "Invalid BOOLEAN value should fail",
		},
		{
			name: "Invalid DATE value",
			skuAttrVal: SkuAttributeValue{
				SkuID:       1,
				AttributeID: 4,
				Value:       "invalid-date",
			},
			attribute: &Attribute{
				Name:     "Warranty Expiry",
				Code:     "warranty_expiry",
				DataType: DataTypeDate,
			},
			expectError: true,
			description: "Invalid DATE value should fail",
		},
		{
			name: "Attribute not found",
			skuAttrVal: SkuAttributeValue{
				SkuID:       1,
				AttributeID: 999,
				Value:       "test",
			},
			attribute:   nil,
			dbError:     errors.New("record not found"),
			expectError: true,
			description: "Should fail when attribute not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock database
			mockDB := &MockDBForSkuAttr{
				attribute:  tc.attribute,
				firstError: tc.dbError,
			}

			// Call validateValue with test wrapper
			err := tc.skuAttrVal.testValidateValue(mockDB)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got nil", tc.description)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for %s, but got: %v", tc.description, err)
				}
			}
		})
	}
}

// TestGetParsedValueWithMock tests GetParsedValue method
func TestGetParsedValueWithMock(t *testing.T) {
	testCases := []struct {
		name        string
		skuAttrVal  SkuAttributeValue
		attribute   *Attribute
		expectedVal interface{}
		expectError bool
	}{
		{
			name: "Parse TEXT value",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 1,
				Value:       "Space Black",
			},
			attribute: &Attribute{
				DataType: DataTypeText,
			},
			expectedVal: "Space Black",
			expectError: false,
		},
		{
			name: "Parse NUMBER value",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 2,
				Value:       "16",
			},
			attribute: &Attribute{
				DataType: DataTypeNumber,
			},
			expectedVal: float64(16),
			expectError: false,
		},
		{
			name: "Parse BOOLEAN value",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 3,
				Value:       "true",
			},
			attribute: &Attribute{
				DataType: DataTypeBoolean,
			},
			expectedVal: true,
			expectError: false,
		},
		{
			name: "Attribute not found",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 999,
				Value:       "test",
			},
			attribute:   nil,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock database
			var mockDB *MockDBForSkuAttr
			if tc.attribute == nil {
				mockDB = &MockDBForSkuAttr{
					firstError: errors.New("record not found"),
				}
			} else {
				mockDB = &MockDBForSkuAttr{
					attribute: tc.attribute,
				}
			}

			// Call GetParsedValue with test wrapper
			result, err := tc.skuAttrVal.testGetParsedValue(mockDB)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
				if result != tc.expectedVal {
					t.Errorf("Expected %v (%T), but got %v (%T)", tc.expectedVal, tc.expectedVal, result, result)
				}
			}
		})
	}
}

// TestGetDisplayValueWithMock tests GetDisplayValue method
func TestGetDisplayValueWithMock(t *testing.T) {
	testCases := []struct {
		name        string
		skuAttrVal  SkuAttributeValue
		attribute   *Attribute
		expected    string
		expectError bool
	}{
		{
			name: "Display value without UOM",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 1,
				Value:       "Space Black",
			},
			attribute: &Attribute{
				Name:     "Color",
				DataType: DataTypeText,
				UOM:      "",
			},
			expected:    "Space Black",
			expectError: false,
		},
		{
			name: "Display value with UOM",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 2,
				Value:       "16",
			},
			attribute: &Attribute{
				Name:     "RAM",
				DataType: DataTypeNumber,
				UOM:      "GB",
			},
			expected:    "16 GB",
			expectError: false,
		},
		{
			name: "Display value with UOM (different unit)",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 3,
				Value:       "15.6",
			},
			attribute: &Attribute{
				Name:     "Screen Size",
				DataType: DataTypeNumber,
				UOM:      "inch",
			},
			expected:    "15.6 inch",
			expectError: false,
		},
		{
			name: "Attribute not found",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 999,
				Value:       "test",
			},
			attribute:   nil,
			expectError: true,
		},
		{
			name: "Display value with preloaded Attribute (with UOM)",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 2,
				Value:       "32",
				Attribute: &Attribute{
					Name:     "RAM",
					DataType: DataTypeNumber,
					UOM:      "GB",
				},
			},
			attribute:   nil, // Not used because Attribute is preloaded
			expected:    "32 GB",
			expectError: false,
		},
		{
			name: "Display value with preloaded Attribute (without UOM)",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 1,
				Value:       "Blue",
				Attribute: &Attribute{
					Name:     "Color",
					DataType: DataTypeText,
					UOM:      "",
				},
			},
			attribute:   nil, // Not used because Attribute is preloaded
			expected:    "Blue",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock database
			var mockDB *MockDBForSkuAttr
			if tc.attribute == nil {
				mockDB = &MockDBForSkuAttr{
					firstError: errors.New("record not found"),
				}
			} else {
				mockDB = &MockDBForSkuAttr{
					attribute: tc.attribute,
				}
			}

			// Call GetDisplayValue with test wrapper
			result, err := tc.skuAttrVal.testGetDisplayValue(mockDB)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Expected '%s', but got '%s'", tc.expected, result)
				}
			}
		})
	}
}

// TestSetValueWithMock tests SetValue method
func TestSetValueWithMock(t *testing.T) {
	testCases := []struct {
		name        string
		skuAttrVal  SkuAttributeValue
		attribute   *Attribute
		inputValue  interface{}
		expected    string
		expectError bool
	}{
		{
			name: "Set TEXT value",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 1,
			},
			attribute: &Attribute{
				DataType: DataTypeText,
			},
			inputValue:  "Space Black",
			expected:    "Space Black",
			expectError: false,
		},
		{
			name: "Set NUMBER value (int)",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 2,
			},
			attribute: &Attribute{
				DataType: DataTypeNumber,
			},
			inputValue:  int64(16),
			expected:    "16",
			expectError: false,
		},
		{
			name: "Set NUMBER value (float)",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 2,
			},
			attribute: &Attribute{
				DataType: DataTypeNumber,
			},
			inputValue:  15.99,
			expected:    "15.99",
			expectError: false,
		},
		{
			name: "Set BOOLEAN value",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 3,
			},
			attribute: &Attribute{
				DataType: DataTypeBoolean,
			},
			inputValue:  true,
			expected:    "true",
			expectError: false,
		},
		{
			name: "Set DATE value",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 4,
			},
			attribute: &Attribute{
				DataType: DataTypeDate,
			},
			inputValue:  time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			expected:    "2023-12-31",
			expectError: false,
		},
		{
			name: "Attribute not found",
			skuAttrVal: SkuAttributeValue{
				AttributeID: 999,
			},
			attribute:   nil,
			inputValue:  "test",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock database
			var mockDB *MockDBForSkuAttr
			if tc.attribute == nil {
				mockDB = &MockDBForSkuAttr{
					firstError: errors.New("record not found"),
				}
			} else {
				mockDB = &MockDBForSkuAttr{
					attribute: tc.attribute,
				}
			}

			// Call SetValue with test wrapper
			err := tc.skuAttrVal.testSetValue(tc.inputValue, mockDB)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
				if tc.skuAttrVal.Value != tc.expected {
					t.Errorf("Expected value '%s', but got '%s'", tc.expected, tc.skuAttrVal.Value)
				}
			}
		})
	}
}

// TestTableName tests the TableName method
func TestTableName(t *testing.T) {
	skuAttrVal := SkuAttributeValue{}
	tableName := skuAttrVal.TableName()

	if tableName != "sku_attribute_values" {
		t.Errorf("Expected table name 'sku_attribute_values', got '%s'", tableName)
	}
}

// Test wrapper functions for unit testing

// testValidateValue is a test wrapper for validateValue
func (sav *SkuAttributeValue) testValidateValue(mockDB *MockDBForSkuAttr) error {
	// Fetch the attribute to get its data type
	var attribute Attribute
	if err := mockDB.First(&attribute, sav.AttributeID).Error(); err != nil {
		return errors.New("attribute not found: " + err.Error())
	}

	// Validate value according to attribute's data type
	if err := attribute.ValidateValue(sav.Value); err != nil {
		return errors.New("invalid value for attribute '" + attribute.Name + "': " + err.Error())
	}

	return nil
}

// testGetParsedValue is a test wrapper for GetParsedValue
func (sav *SkuAttributeValue) testGetParsedValue(mockDB *MockDBForSkuAttr) (interface{}, error) {
	var attribute Attribute
	if err := mockDB.First(&attribute, sav.AttributeID).Error(); err != nil {
		return nil, errors.New("attribute not found: " + err.Error())
	}

	return attribute.ParseValue(sav.Value)
}

// testSetValue is a test wrapper for SetValue
func (sav *SkuAttributeValue) testSetValue(value interface{}, mockDB *MockDBForSkuAttr) error {
	var attribute Attribute
	if err := mockDB.First(&attribute, sav.AttributeID).Error(); err != nil {
		return errors.New("attribute not found: " + err.Error())
	}

	valueStr, err := attribute.FormatValue(value)
	if err != nil {
		return err
	}

	sav.Value = valueStr
	return nil
}

// testGetDisplayValue is a test wrapper for GetDisplayValue (hybrid approach)
func (sav *SkuAttributeValue) testGetDisplayValue(mockDB *MockDBForSkuAttr) (string, error) {
	var attribute Attribute

	// Check if Attribute is already preloaded
	if sav.Attribute != nil {
		attribute = *sav.Attribute
	} else {
		// Fetch attribute if not preloaded
		if err := mockDB.First(&attribute, sav.AttributeID).Error(); err != nil {
			return "", errors.New("attribute not found: " + err.Error())
		}
	}

	// Format display value with UOM if available
	if attribute.UOM != "" {
		return fmt.Sprintf("%s %s", sav.Value, attribute.UOM), nil
	}

	return sav.Value, nil
}

// TestSkuAttributeValueWithDifferentTypes tests comprehensive scenarios with all data types
func TestSkuAttributeValueWithDifferentTypes(t *testing.T) {
	testCases := []struct {
		name       string
		attribute  Attribute
		value      string
		shouldPass bool
	}{
		{
			name:       "TEXT attribute - any value",
			attribute:  Attribute{DataType: DataTypeText},
			value:      "Any text including special chars !@#$%",
			shouldPass: true,
		},
		{
			name:       "NUMBER attribute - positive integer",
			attribute:  Attribute{DataType: DataTypeNumber},
			value:      "100",
			shouldPass: true,
		},
		{
			name:       "NUMBER attribute - negative number",
			attribute:  Attribute{DataType: DataTypeNumber},
			value:      "-50",
			shouldPass: true,
		},
		{
			name:       "NUMBER attribute - decimal",
			attribute:  Attribute{DataType: DataTypeNumber},
			value:      "99.99",
			shouldPass: true,
		},
		{
			name:       "NUMBER attribute - invalid text",
			attribute:  Attribute{DataType: DataTypeNumber},
			value:      "abc123",
			shouldPass: false,
		},
		{
			name:       "BOOLEAN attribute - true",
			attribute:  Attribute{DataType: DataTypeBoolean},
			value:      "true",
			shouldPass: true,
		},
		{
			name:       "BOOLEAN attribute - false",
			attribute:  Attribute{DataType: DataTypeBoolean},
			value:      "false",
			shouldPass: true,
		},
		{
			name:       "BOOLEAN attribute - 1 (valid)",
			attribute:  Attribute{DataType: DataTypeBoolean},
			value:      "1",
			shouldPass: true,
		},
		{
			name:       "BOOLEAN attribute - 0 (valid)",
			attribute:  Attribute{DataType: DataTypeBoolean},
			value:      "0",
			shouldPass: true,
		},
		{
			name:       "BOOLEAN attribute - invalid",
			attribute:  Attribute{DataType: DataTypeBoolean},
			value:      "yes",
			shouldPass: false,
		},
		{
			name:       "DATE attribute - ISO format",
			attribute:  Attribute{DataType: DataTypeDate},
			value:      "2023-12-31",
			shouldPass: true,
		},
		{
			name:       "DATE attribute - ISO datetime",
			attribute:  Attribute{DataType: DataTypeDate},
			value:      "2023-12-31 23:59:59",
			shouldPass: true,
		},
		{
			name:       "DATE attribute - invalid format",
			attribute:  Attribute{DataType: DataTypeDate},
			value:      "31-12-2023",
			shouldPass: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test validation
			err := tc.attribute.ValidateValue(tc.value)

			if tc.shouldPass {
				if err != nil {
					t.Errorf("Expected validation to pass, but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}
