package models

import (
	"testing"
	"time"
)

// TestAttributeDataTypeValidation tests data type validation
func TestAttributeDataTypeValidation(t *testing.T) {
	testCases := []struct {
		name        string
		attribute   Attribute
		expectError bool
		description string
	}{
		{
			name: "Valid TEXT data type",
			attribute: Attribute{
				Name:     "Color",
				Code:     "color",
				DataType: DataTypeText,
			},
			expectError: false,
			description: "TEXT data type should be valid",
		},
		{
			name: "Valid NUMBER data type",
			attribute: Attribute{
				Name:     "RAM",
				Code:     "ram",
				DataType: DataTypeNumber,
				UOM:      "GB",
			},
			expectError: false,
			description: "NUMBER data type should be valid",
		},
		{
			name: "Valid BOOLEAN data type",
			attribute: Attribute{
				Name:     "In Stock",
				Code:     "in_stock",
				DataType: DataTypeBoolean,
			},
			expectError: false,
			description: "BOOLEAN data type should be valid",
		},
		{
			name: "Valid DATE data type",
			attribute: Attribute{
				Name:     "Warranty Expiry",
				Code:     "warranty_expiry",
				DataType: DataTypeDate,
			},
			expectError: false,
			description: "DATE data type should be valid",
		},
		{
			name: "Invalid data type",
			attribute: Attribute{
				Name:     "Test",
				Code:     "test",
				DataType: DataType("INVALID"),
			},
			expectError: true,
			description: "Invalid data type should fail",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.attribute.ValidateDataType()

			if tc.expectError && err == nil {
				t.Errorf("Expected error for %s, but got nil", tc.description)
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error for %s, but got: %v", tc.description, err)
			}
		})
	}
}

// TestAttributeParseValue tests value parsing for different data types
func TestAttributeParseValue(t *testing.T) {
	testCases := []struct {
		name        string
		attribute   Attribute
		valueStr    string
		expected    interface{}
		expectError bool
	}{
		{
			name:        "Parse text value",
			attribute:   Attribute{DataType: DataTypeText},
			valueStr:    "Black",
			expected:    "Black",
			expectError: false,
		},
		{
			name:        "Parse number value (integer)",
			attribute:   Attribute{DataType: DataTypeNumber},
			valueStr:    "16",
			expected:    float64(16),
			expectError: false,
		},
		{
			name:        "Parse number value (decimal)",
			attribute:   Attribute{DataType: DataTypeNumber},
			valueStr:    "15.99",
			expected:    float64(15.99),
			expectError: false,
		},
		{
			name:        "Parse boolean true",
			attribute:   Attribute{DataType: DataTypeBoolean},
			valueStr:    "true",
			expected:    true,
			expectError: false,
		},
		{
			name:        "Parse boolean false",
			attribute:   Attribute{DataType: DataTypeBoolean},
			valueStr:    "false",
			expected:    false,
			expectError: false,
		},
		{
			name:        "Parse date ISO format",
			attribute:   Attribute{DataType: DataTypeDate},
			valueStr:    "2023-12-31",
			expected:    time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "Parse invalid number",
			attribute:   Attribute{DataType: DataTypeNumber},
			valueStr:    "not-a-number",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Parse invalid boolean",
			attribute:   Attribute{DataType: DataTypeBoolean},
			valueStr:    "not-a-bool",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Parse invalid date",
			attribute:   Attribute{DataType: DataTypeDate},
			valueStr:    "not-a-date",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.attribute.ParseValue(tc.valueStr)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Expected %v (%T), but got %v (%T)", tc.expected, tc.expected, result, result)
				}
			}
		})
	}
}

// TestAttributeFormatValue tests value formatting for different data types
func TestAttributeFormatValue(t *testing.T) {
	testCases := []struct {
		name        string
		attribute   Attribute
		value       interface{}
		expected    string
		expectError bool
	}{
		{
			name:        "Format text value",
			attribute:   Attribute{DataType: DataTypeText},
			value:       "Black",
			expected:    "Black",
			expectError: false,
		},
		{
			name:        "Format number value (int)",
			attribute:   Attribute{DataType: DataTypeNumber},
			value:       int64(16),
			expected:    "16",
			expectError: false,
		},
		{
			name:        "Format number value (float)",
			attribute:   Attribute{DataType: DataTypeNumber},
			value:       15.99,
			expected:    "15.99",
			expectError: false,
		},
		{
			name:        "Format boolean true",
			attribute:   Attribute{DataType: DataTypeBoolean},
			value:       true,
			expected:    "true",
			expectError: false,
		},
		{
			name:        "Format boolean false",
			attribute:   Attribute{DataType: DataTypeBoolean},
			value:       false,
			expected:    "false",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.attribute.FormatValue(tc.value)

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

// TestAttributeValidateValue tests value validation
func TestAttributeValidateValue(t *testing.T) {
	testCases := []struct {
		name        string
		attribute   Attribute
		valueStr    string
		expectError bool
	}{
		{
			name:        "Valid text",
			attribute:   Attribute{DataType: DataTypeText},
			valueStr:    "Any text is valid",
			expectError: false,
		},
		{
			name:        "Valid number (integer)",
			attribute:   Attribute{DataType: DataTypeNumber},
			valueStr:    "123",
			expectError: false,
		},
		{
			name:        "Valid number (decimal)",
			attribute:   Attribute{DataType: DataTypeNumber},
			valueStr:    "123.45",
			expectError: false,
		},
		{
			name:        "Invalid number",
			attribute:   Attribute{DataType: DataTypeNumber},
			valueStr:    "abc",
			expectError: true,
		},
		{
			name:        "Valid boolean",
			attribute:   Attribute{DataType: DataTypeBoolean},
			valueStr:    "true",
			expectError: false,
		},
		{
			name:        "Invalid boolean",
			attribute:   Attribute{DataType: DataTypeBoolean},
			valueStr:    "yes",
			expectError: true,
		},
		{
			name:        "Valid date",
			attribute:   Attribute{DataType: DataTypeDate},
			valueStr:    "2023-12-31",
			expectError: false,
		},
		{
			name:        "Invalid date",
			attribute:   Attribute{DataType: DataTypeDate},
			valueStr:    "invalid-date",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.attribute.ValidateValue(tc.valueStr)

			if tc.expectError && err == nil {
				t.Errorf("Expected error for value '%s', but got nil", tc.valueStr)
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error for value '%s', but got: %v", tc.valueStr, err)
			}
		})
	}
}

// TestAttributeDataTypeConstants tests that constants have expected values
func TestAttributeDataTypeConstants(t *testing.T) {
	if DataTypeText != "TEXT" {
		t.Errorf("Expected DataTypeText to be 'TEXT', got '%s'", DataTypeText)
	}

	if DataTypeNumber != "NUMBER" {
		t.Errorf("Expected DataTypeNumber to be 'NUMBER', got '%s'", DataTypeNumber)
	}

	if DataTypeBoolean != "BOOLEAN" {
		t.Errorf("Expected DataTypeBoolean to be 'BOOLEAN', got '%s'", DataTypeBoolean)
	}

	if DataTypeDate != "DATE" {
		t.Errorf("Expected DataTypeDate to be 'DATE', got '%s'", DataTypeDate)
	}
}
