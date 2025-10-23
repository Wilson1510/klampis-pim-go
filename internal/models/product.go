package models

import (
	"github.com/Wilson1510/klampis-pim-go/pkg/utils"
	"gorm.io/gorm"
)

type Product struct {
	Base
	Name        string `gorm:"not null;type:varchar(150)" json:"name"`
	Slug        string `gorm:"uniqueIndex;not null;type:varchar(170)" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	CategoryID  uint   `gorm:"not null;index" json:"category_id"`

	// Relationship with Category
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	// Relationship with SKUs
	Skus []Sku `gorm:"foreignKey:ProductID" json:"skus,omitempty"`

	// Polymorphic relationship with Images
	Images []Image `gorm:"polymorphic:Imageable" json:"images,omitempty"`
}

// BeforeCreate is a GORM hook that runs before creating a record
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	return utils.GenerateModelSlug(p, tx)
}

// BeforeUpdate is a GORM hook that runs before updating a record
func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	return utils.GenerateModelSlug(p, tx)
}

// SlugModel interface implementation for Product

// GetName returns the name field for slug generation
func (p *Product) GetName() string {
	return p.Name
}

// GetSlug returns the current slug
func (p *Product) GetSlug() string {
	return p.Slug
}

// SetSlug sets the slug field
func (p *Product) SetSlug(slug string) {
	p.Slug = slug
}

// GetID returns the ID for database operations
func (p *Product) GetID() uint {
	return p.ID
}

// GetTableName returns the table name for database operations
func (p *Product) GetTableName() string {
	return "products"
}
