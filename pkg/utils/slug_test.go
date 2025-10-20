package utils

import (
	"testing"
)

// TestGenerateSlug tests the GenerateSlug function
func TestGenerateSlug(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple text",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "Text with special characters",
			input:    "Electronics & Gadgets",
			expected: "electronics-gadgets",
		},
		{
			name:     "Text with numbers",
			input:    "iPhone 15 Pro Max",
			expected: "iphone-15-pro-max",
		},
		{
			name:     "Text with multiple spaces",
			input:    "Smart   Phones    &   Accessories",
			expected: "smart-phones-accessories",
		},
		{
			name:     "Text with underscores and hyphens",
			input:    "Test_Product-Name",
			expected: "test-product-name",
		},
		{
			name:     "Text with leading/trailing spaces",
			input:    "  Trimmed Text  ",
			expected: "trimmed-text",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only special characters",
			input:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "Mixed case with accents",
			input:    "Café & Restaurant",
			expected: "cafe-restaurant",
		},
		{
			name:     "Long text with multiple hyphens",
			input:    "This---is---a---very---long---product---name",
			expected: "this-is-a-very-long-product-name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateSlug(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', but got '%s'", tc.expected, result)
			}
		})
	}
}

// TestGenerateUniqueSlug tests the GenerateUniqueSlug function
func TestGenerateUniqueSlug(t *testing.T) {
	testCases := []struct {
		name          string
		baseSlug      string
		existingSlugs []string
		expected      string
	}{
		{
			name:          "Unique slug",
			baseSlug:      "electronics",
			existingSlugs: []string{"computers", "phones", "tablets"},
			expected:      "electronics",
		},
		{
			name:          "Duplicate slug - should append number",
			baseSlug:      "electronics",
			existingSlugs: []string{"electronics", "computers", "phones"},
			expected:      "electronics-1",
		},
		{
			name:          "Multiple duplicates",
			baseSlug:      "electronics",
			existingSlugs: []string{"electronics", "electronics-1", "electronics-2"},
			expected:      "electronics-3",
		},
		{
			name:          "Empty base slug",
			baseSlug:      "",
			existingSlugs: []string{"test"},
			expected:      "",
		},
		{
			name:          "No existing slugs",
			baseSlug:      "new-product",
			existingSlugs: []string{},
			expected:      "new-product",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateUniqueSlug(tc.baseSlug, tc.existingSlugs)
			if result != tc.expected {
				t.Errorf("Expected '%s', but got '%s'", tc.expected, result)
			}
		})
	}
}

// TestContains tests the contains helper function
func TestContains(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "Item exists",
			slice:    []string{"apple", "banana", "orange"},
			item:     "banana",
			expected: true,
		},
		{
			name:     "Item does not exist",
			slice:    []string{"apple", "banana", "orange"},
			item:     "grape",
			expected: false,
		},
		{
			name:     "Empty slice",
			slice:    []string{},
			item:     "apple",
			expected: false,
		},
		{
			name:     "Empty item",
			slice:    []string{"apple", "banana"},
			item:     "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := contains(tc.slice, tc.item)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
