package models

import (
	"gorm.io/gorm"
)

type Base struct {
	gorm.Model
	CreatedBy uint `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"created_by"`
	UpdatedBy uint `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"updated_by"`
	IsActive  bool `gorm:"default:false" json:"is_active"`
	Sequence  uint  `gorm:"default:0" json:"sequence"`
	
	// Foreign key relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy" json:"updated_by_user,omitempty"`
}