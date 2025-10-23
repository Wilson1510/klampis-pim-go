package models

import (
	"github.com/Wilson1510/klampis-pim-go/pkg/utils"
	"gorm.io/gorm"
)

type Category struct {
	Base
	Name        string `gorm:"not null;type:varchar(100)" json:"name"`
	Slug        string `gorm:"uniqueIndex;not null;type:varchar(120)" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	ParentID    *uint  `gorm:"index" json:"parent_id"`

	// Self-referencing relationships
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	// Relationship with Products
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// BeforeCreate is a GORM hook that runs before creating a record
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	return utils.GenerateModelSlug(c, tx)
}

// BeforeUpdate is a GORM hook that runs before updating a record
func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	return utils.GenerateModelSlug(c, tx)
}

// SlugModel interface implementation for Category

// GetName returns the name field for slug generation
func (c *Category) GetName() string {
	return c.Name
}

// GetSlug returns the current slug
func (c *Category) GetSlug() string {
	return c.Slug
}

// SetSlug sets the slug field
func (c *Category) SetSlug(slug string) {
	c.Slug = slug
}

// GetID returns the ID for database operations
func (c *Category) GetID() uint {
	return c.ID
}

// GetTableName returns the table name for database operations
func (c *Category) GetTableName() string {
	return "categories"
}
