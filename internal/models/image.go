package models

import (
	"gorm.io/gorm"
)

type Image struct {
	Base
	File          string `gorm:"not null;type:varchar(255)" json:"file"`
	Title         string `gorm:"type:varchar(200)" json:"title"`
	IsPrimary     bool   `gorm:"default:false" json:"is_primary"`
	ImageableID   uint   `gorm:"not null;index" json:"imageable_id"`
	ImageableType string `gorm:"not null;type:varchar(50);index" json:"imageable_type"`
}

// TableName overrides the table name used by Image to `images`
func (Image) TableName() string {
	return "images"
}

// BeforeCreate hook to ensure only one primary image per imageable
func (i *Image) BeforeCreate(tx *gorm.DB) error {
	if i.IsPrimary {
		// Set other images of the same imageable to non-primary
		return tx.Model(&Image{}).
			Where("imageable_id = ? AND imageable_type = ? AND is_primary = ?",
				i.ImageableID, i.ImageableType, true).
			Update("is_primary", false).Error
	}
	return nil
}

// BeforeUpdate hook to ensure only one primary image per imageable
func (i *Image) BeforeUpdate(tx *gorm.DB) error {
	if i.IsPrimary {
		// Set other images of the same imageable to non-primary
		return tx.Model(&Image{}).
			Where("imageable_id = ? AND imageable_type = ? AND id != ? AND is_primary = ?",
				i.ImageableID, i.ImageableType, i.ID, true).
			Update("is_primary", false).Error
	}
	return nil
}
