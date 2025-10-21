package models

import (
	"fmt"
	"gorm.io/gorm"
)

type SkuAttributeValue struct {
	gorm.Model
	SkuID       uint   `gorm:"not null;index:idx_sku_attr" json:"sku_id"`
	AttributeID uint   `gorm:"not null;index:idx_sku_attr" json:"attribute_id"`
	Value       string `gorm:"type:text;not null" json:"value"`
	CreatedBy   uint   `gorm:"" json:"created_by"` // Optional: untuk audit trail
	UpdatedBy   uint   `gorm:"" json:"updated_by"` // Optional: untuk audit trail
	Sequence    int    `gorm:"default:0" json:"sequence"` // Untuk ordering attributes display
	
	// Relationships
	Sku           *Sku       `gorm:"foreignKey:SkuID" json:"sku,omitempty"`
	Attribute     *Attribute `gorm:"foreignKey:AttributeID" json:"attribute,omitempty"`
	CreatedByUser *User      `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	UpdatedByUser *User      `gorm:"foreignKey:UpdatedBy" json:"updated_by_user,omitempty"`
}

// BeforeCreate GORM hook for validation
func (sav *SkuAttributeValue) BeforeCreate(tx *gorm.DB) error {
	return sav.validateValue(tx)
}

// BeforeUpdate GORM hook for validation
func (sav *SkuAttributeValue) BeforeUpdate(tx *gorm.DB) error {
	return sav.validateValue(tx)
}

// validateValue validates the value against the attribute's data type
func (sav *SkuAttributeValue) validateValue(tx *gorm.DB) error {
	// Fetch the attribute to get its data type
	var attribute Attribute
	if err := tx.First(&attribute, sav.AttributeID).Error; err != nil {
		return fmt.Errorf("attribute not found: %w", err)
	}
	
	// Validate value according to attribute's data type
	if err := attribute.ValidateValue(sav.Value); err != nil {
		return fmt.Errorf("invalid value for attribute '%s' (type: %s): %w", 
			attribute.Name, attribute.DataType, err)
	}
	
	return nil
}

// GetParsedValue returns the value parsed according to the attribute's data type
func (sav *SkuAttributeValue) GetParsedValue(tx *gorm.DB) (interface{}, error) {
	var attribute Attribute
	if err := tx.First(&attribute, sav.AttributeID).Error; err != nil {
		return nil, fmt.Errorf("attribute not found: %w", err)
	}
	
	return attribute.ParseValue(sav.Value)
}

// SetValue sets the value with automatic type conversion
func (sav *SkuAttributeValue) SetValue(value interface{}, tx *gorm.DB) error {
	var attribute Attribute
	if err := tx.First(&attribute, sav.AttributeID).Error; err != nil {
		return fmt.Errorf("attribute not found: %w", err)
	}
	
	valueStr, err := attribute.FormatValue(value)
	if err != nil {
		return err
	}
	
	sav.Value = valueStr
	return nil
}

// GetDisplayValue returns a formatted display value with UOM if applicable
func (sav *SkuAttributeValue) GetDisplayValue(tx *gorm.DB) (string, error) {
	var attribute Attribute
	// Check if Attribute is already preloaded
	if sav.Attribute != nil {
		attribute = *sav.Attribute
	} else {
		// Fetch attribute if not preloaded
		if err := tx.First(&attribute, sav.AttributeID).Error; err != nil {
			return "", fmt.Errorf("attribute not found: %w", err)
		}
	}
	
	// Format display value with UOM if available
	if attribute.UOM != "" {
		return fmt.Sprintf("%s %s", sav.Value, attribute.UOM), nil
	}
	
	return sav.Value, nil
}

// TableName specifies the table name for SkuAttributeValue
func (SkuAttributeValue) TableName() string {
	return "sku_attribute_values"
}

// Example usage:
// 
// // Create attribute
// ramAttr := Attribute{
//     Name:     "RAM",
//     Code:     "ram",
//     DataType: DataTypeNumber,
//     UOM:      "GB",
// }
//
// // Create SKU attribute value
// skuAttrValue := SkuAttributeValue{
//     SkuID:       1,
//     AttributeID: ramAttr.ID,
//     Value:       "16",  // Stored as string
//     Sequence:    1,
// }
//
// // Get parsed value
// parsedValue, _ := skuAttrValue.GetParsedValue(db) // Returns: int64(16)
//
// // Get display value
// displayValue, _ := skuAttrValue.GetDisplayValue(db) // Returns: "16 GB"
