package utils

import (
	"gorm.io/gorm"
)

// SlugModel defines interface for models that have slug functionality
type SlugModel interface {
	GetName() string
	GetSlug() string
	SetSlug(slug string)
	GetID() uint
	GetTableName() string
}

// GenerateModelSlug generates slug for any model that implements SlugModel interface
func GenerateModelSlug(model SlugModel, tx *gorm.DB) error {
	// Only generate slug if it's empty or if name has changed
	if model.GetSlug() == "" || ShouldRegenerateModelSlug(model, tx) {
		baseSlug := GenerateSlug(model.GetName())
		if baseSlug == "" {
			return nil // Skip if name is empty
		}

		// Check for existing slugs to ensure uniqueness
		existingSlugs, err := getExistingSlugs(model, tx)
		if err != nil {
			return err
		}

		// Generate unique slug
		uniqueSlug := GenerateUniqueSlug(baseSlug, existingSlugs)
		model.SetSlug(uniqueSlug)
	}

	return nil
}

// ShouldRegenerateModelSlug checks if slug should be regenerated based on name changes
func ShouldRegenerateModelSlug(model SlugModel, tx *gorm.DB) bool {
	if model.GetID() == 0 {
		return true // New record
	}

	// For updates, check if name has changed
	var result struct {
		Name string `json:"name"`
	}

	err := tx.Table(model.GetTableName()).
		Select("name").
		Where("id = ?", model.GetID()).
		First(&result).Error

	if err != nil {
		return true // If we can't find original, regenerate
	}

	return result.Name != model.GetName()
}

// getExistingSlugs retrieves existing slugs for uniqueness check
func getExistingSlugs(model SlugModel, tx *gorm.DB) ([]string, error) {
	var results []struct {
		Slug string `json:"slug"`
	}

	query := tx.Table(model.GetTableName()).Select("slug")
	if model.GetID() != 0 {
		query = query.Where("id != ?", model.GetID())
	}

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	var slugs []string
	for _, result := range results {
		slugs = append(slugs, result.Slug)
	}

	return slugs, nil
}
