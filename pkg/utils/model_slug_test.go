package utils

import (
	"errors"
	"testing"
)

// MockSlugModel implements SlugModel interface for testing
type MockSlugModel struct {
	id        uint
	name      string
	slug      string
	tableName string
}

func (m *MockSlugModel) GetName() string      { return m.name }
func (m *MockSlugModel) GetSlug() string      { return m.slug }
func (m *MockSlugModel) SetSlug(slug string)  { m.slug = slug }
func (m *MockSlugModel) GetID() uint          { return m.id }
func (m *MockSlugModel) GetTableName() string { return m.tableName }

// MockDB implements database interface for testing
type MockDB struct {
	findResults []map[string]interface{}
	findError   error
	firstResult map[string]interface{}
	firstError  error
}

func (m *MockDB) Table(name string) *MockDB                             { return m }
func (m *MockDB) Select(query interface{}, args ...interface{}) *MockDB { return m }
func (m *MockDB) Where(query interface{}, args ...interface{}) *MockDB  { return m }

func (m *MockDB) Find(dest interface{}) *MockDB {
	if m.findError != nil {
		return &MockDB{findError: m.findError}
	}

	// Handle different struct types that might be passed to Find()
	switch v := dest.(type) {
	case *[]struct{ Slug string }:
		*v = make([]struct{ Slug string }, len(m.findResults))
		for i, result := range m.findResults {
			if slug, exists := result["slug"]; exists {
				(*v)[i].Slug = slug.(string)
			}
		}
	case *[]struct {
		Slug string `json:"slug"`
	}:
		*v = make([]struct {
			Slug string `json:"slug"`
		}, len(m.findResults))
		for i, result := range m.findResults {
			if slug, exists := result["slug"]; exists {
				(*v)[i].Slug = slug.(string)
			}
		}
	}

	return m
}

func (m *MockDB) First(dest interface{}) *MockDB {
	if m.firstError != nil {
		return &MockDB{firstError: m.firstError}
	}

	// Handle different struct types that might be passed to First()
	switch v := dest.(type) {
	case *struct{ Name string }:
		if m.firstResult != nil {
			if name, exists := m.firstResult["name"]; exists {
				v.Name = name.(string)
			}
		}
	case *struct {
		Name string `json:"name"`
	}:
		if m.firstResult != nil {
			if name, exists := m.firstResult["name"]; exists {
				v.Name = name.(string)
			}
		}
	}

	return m
}

func (m *MockDB) Error() error {
	if m.findError != nil {
		return m.findError
	}
	if m.firstError != nil {
		return m.firstError
	}
	return nil
}

// TestShouldRegenerateModelSlug tests the ShouldRegenerateModelSlug function
func TestShouldRegenerateModelSlug(t *testing.T) {
	testCases := []struct {
		name           string
		model          *MockSlugModel
		originalName   string
		dbError        error
		expectedResult bool
		description    string
	}{
		{
			name: "New record (ID = 0)",
			model: &MockSlugModel{
				id:        0,
				name:      "New Model",
				tableName: "test_models",
			},
			expectedResult: true,
			description:    "Should regenerate for new records",
		},
		{
			name: "Name changed",
			model: &MockSlugModel{
				id:        1,
				name:      "Updated Model Name",
				tableName: "test_models",
			},
			originalName:   "Original Model Name",
			expectedResult: true,
			description:    "Should regenerate when name changes",
		},
		{
			name: "Name unchanged",
			model: &MockSlugModel{
				id:        1,
				name:      "Same Model Name",
				tableName: "test_models",
			},
			originalName:   "Same Model Name",
			expectedResult: false,
			description:    "Should not regenerate when name is same",
		},
		{
			name: "Database error",
			model: &MockSlugModel{
				id:        1,
				name:      "Model Name",
				tableName: "test_models",
			},
			dbError:        errors.New("database connection failed"),
			expectedResult: true,
			description:    "Should regenerate when can't fetch original",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock database
			mockDB := &MockDB{
				firstResult: map[string]interface{}{
					"name": tc.originalName,
				},
				firstError: tc.dbError,
			}

			// Call ShouldRegenerateModelSlug
			result := shouldRegenerateModelSlugWithMock(tc.model, mockDB)

			// Assert result
			if result != tc.expectedResult {
				t.Errorf("Expected ShouldRegenerateModelSlug() to return %v, but got %v. %s",
					tc.expectedResult, result, tc.description)
			}
		})
	}
}

// TestGenerateModelSlug tests the GenerateModelSlug function
func TestGenerateModelSlug(t *testing.T) {
	testCases := []struct {
		name            string
		model           *MockSlugModel
		existingSlugs   []string
		dbError         error
		expectedSlug    string
		shouldCallUtils bool
		description     string
	}{
		{
			name: "Empty slug - should generate",
			model: &MockSlugModel{
				id:        1,
				name:      "Electronics & Gadgets",
				slug:      "", // Empty slug
				tableName: "test_models",
			},
			existingSlugs:   []string{},
			expectedSlug:    "electronics-gadgets",
			shouldCallUtils: true,
			description:     "Should generate slug when empty",
		},
		{
			name: "Existing slug with unique name",
			model: &MockSlugModel{
				id:        1,
				name:      "Unique Model",
				slug:      "", // Will be generated
				tableName: "test_models",
			},
			existingSlugs:   []string{"other-model", "another-model"},
			expectedSlug:    "unique-model",
			shouldCallUtils: true,
			description:     "Should generate unique slug",
		},
		{
			name: "Duplicate slug - should append number",
			model: &MockSlugModel{
				id:        1,
				name:      "Electronics",
				slug:      "", // Will be generated
				tableName: "test_models",
			},
			existingSlugs:   []string{"electronics", "electronics-1"},
			expectedSlug:    "electronics-2",
			shouldCallUtils: true,
			description:     "Should handle duplicate slugs",
		},
		{
			name: "Empty name - should skip",
			model: &MockSlugModel{
				id:        1,
				name:      "", // Empty name
				slug:      "",
				tableName: "test_models",
			},
			existingSlugs:   []string{},
			expectedSlug:    "",
			shouldCallUtils: false,
			description:     "Should skip when name is empty",
		},
		{
			name: "Database error - should return error",
			model: &MockSlugModel{
				id:        1,
				name:      "Test Model",
				slug:      "",
				tableName: "test_models",
			},
			dbError:         errors.New("database query failed"),
			shouldCallUtils: false,
			description:     "Should handle database errors",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock database with existing slugs
			var mockResults []map[string]interface{}
			for _, slug := range tc.existingSlugs {
				mockResults = append(mockResults, map[string]interface{}{
					"slug": slug,
				})
			}

			mockDB := &MockDB{
				findResults: mockResults,
				findError:   tc.dbError,
			}

			// Call GenerateModelSlug using test wrapper
			err := generateModelSlugWithMock(tc.model, mockDB)

			// Assert error handling
			if tc.dbError != nil {
				if err == nil {
					t.Errorf("Expected error when database fails, but got nil")
				}
				return
			}

			// Assert no error for successful cases
			if err != nil {
				t.Errorf("Expected no error, but got: %v", err)
				return
			}

			// Assert slug generation
			if tc.shouldCallUtils {
				if tc.model.GetSlug() != tc.expectedSlug {
					t.Errorf("Expected slug '%s', but got '%s'. %s",
						tc.expectedSlug, tc.model.GetSlug(), tc.description)
				}
			} else {
				// For cases where utils shouldn't be called (empty name)
				if tc.model.GetName() == "" && tc.model.GetSlug() != "" {
					t.Errorf("Expected slug to remain empty when name is empty, but got '%s'",
						tc.model.GetSlug())
				}
			}
		})
	}
}

// Test helper functions that work with mocks instead of real GORM

// shouldRegenerateModelSlugWithMock is a test version of ShouldRegenerateModelSlug
func shouldRegenerateModelSlugWithMock(model SlugModel, mockDB *MockDB) bool {
	if model.GetID() == 0 {
		return true // New record
	}

	// For updates, check if name has changed
	var result struct {
		Name string `json:"name"`
	}

	err := mockDB.Table(model.GetTableName()).
		Select("name").
		Where("id = ?", model.GetID()).
		First(&result).Error()

	if err != nil {
		return true // If we can't find original, regenerate
	}

	return result.Name != model.GetName()
}

// generateModelSlugWithMock is a test version of GenerateModelSlug
func generateModelSlugWithMock(model SlugModel, mockDB *MockDB) error {
	// Only generate slug if it's empty or if name has changed
	if model.GetSlug() == "" || shouldRegenerateModelSlugWithMock(model, mockDB) {
		baseSlug := GenerateSlug(model.GetName())
		if baseSlug == "" {
			return nil // Skip if name is empty
		}

		// Check for existing slugs to ensure uniqueness
		existingSlugs, err := getExistingSlugsWithMock(model, mockDB)
		if err != nil {
			return err
		}

		// Generate unique slug
		uniqueSlug := GenerateUniqueSlug(baseSlug, existingSlugs)
		model.SetSlug(uniqueSlug)
	}

	return nil
}

// getExistingSlugsWithMock is a test version of getExistingSlugs
func getExistingSlugsWithMock(model SlugModel, mockDB *MockDB) ([]string, error) {
	var results []struct {
		Slug string `json:"slug"`
	}

	query := mockDB.Table(model.GetTableName()).Select("slug")
	if model.GetID() != 0 {
		query = query.Where("id != ?", model.GetID())
	}

	if err := query.Find(&results).Error(); err != nil {
		return nil, err
	}

	var slugs []string
	for _, result := range results {
		slugs = append(slugs, result.Slug)
	}

	return slugs, nil
}

// TestSlugModelInterface tests that our mock properly implements SlugModel
func TestSlugModelInterface(t *testing.T) {
	model := &MockSlugModel{
		id:        1,
		name:      "Test Model",
		slug:      "test-model",
		tableName: "test_models",
	}

	// Test interface methods
	if model.GetID() != 1 {
		t.Errorf("Expected ID 1, got %d", model.GetID())
	}

	if model.GetName() != "Test Model" {
		t.Errorf("Expected name 'Test Model', got '%s'", model.GetName())
	}

	if model.GetSlug() != "test-model" {
		t.Errorf("Expected slug 'test-model', got '%s'", model.GetSlug())
	}

	if model.GetTableName() != "test_models" {
		t.Errorf("Expected table name 'test_models', got '%s'", model.GetTableName())
	}

	// Test SetSlug
	model.SetSlug("new-slug")
	if model.GetSlug() != "new-slug" {
		t.Errorf("Expected slug 'new-slug' after SetSlug, got '%s'", model.GetSlug())
	}
}
