package models

import (
	"github.com/Wilson1510/klampis-pim-go/pkg/utils"
	"gorm.io/gorm"
)

type Sku struct {
	Base
	Name        string  `gorm:"not null;type:varchar(200)" json:"name"`
	Slug        string  `gorm:"uniqueIndex;not null;type:varchar(220)" json:"slug"`
	Description string  `gorm:"type:text" json:"description"`
	SkuNumber   string  `gorm:"uniqueIndex;not null;type:varchar(50)" json:"sku_number"`
	Price       float64 `gorm:"not null;type:decimal(15,2)" json:"price"`
	ProductID   uint    `gorm:"not null;index" json:"product_id"`

	// Relationship with Product
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`

	// Polymorphic relationship with Images
	Images []Image `gorm:"polymorphic:Imageable;polymorphicValue:skus" json:"images,omitempty"`
}

// BeforeCreate is a GORM hook that runs before creating a record
func (s *Sku) BeforeCreate(tx *gorm.DB) error {
	return utils.GenerateModelSlug(s, tx)
}

// BeforeUpdate is a GORM hook that runs before updating a record
func (s *Sku) BeforeUpdate(tx *gorm.DB) error {
	return utils.GenerateModelSlug(s, tx)
}

// SlugModel interface implementation for Sku

// GetName returns the name field for slug generation
func (s *Sku) GetName() string {
	return s.Name
}

// GetSlug returns the current slug
func (s *Sku) GetSlug() string {
	return s.Slug
}

// SetSlug sets the slug field
func (s *Sku) SetSlug(slug string) {
	s.Slug = slug
}

// GetID returns the ID for database operations
func (s *Sku) GetID() uint {
	return s.ID
}

// GetTableName returns the table name for database operations
func (s *Sku) GetTableName() string {
	return "skus"
}
