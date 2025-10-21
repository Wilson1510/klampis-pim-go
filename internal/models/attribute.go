package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// DataType constants for attribute data types
type DataType string

const (
	DataTypeText    DataType = "TEXT"
	DataTypeNumber  DataType = "NUMBER"
	DataTypeBoolean DataType = "BOOLEAN"
	DataTypeDate    DataType = "DATE"
)

type Attribute struct {
	Base
	Name     string   `gorm:"not null;type:varchar(50)" json:"name"`
	Code     string   `gorm:"uniqueIndex;not null;type:varchar(70)" json:"code"`
	DataType DataType `gorm:"not null;type:varchar(20)" json:"data_type"`
	UOM      string   `gorm:"type:varchar(15)" json:"uom"` // Unit of measurement: GB, inch, GHz, years, etc.
	
	// Relationships
	SkuAttributeValues []SkuAttributeValue `gorm:"foreignKey:AttributeID" json:"sku_attribute_values,omitempty"`
}

// ValidateDataType validates if the data type is valid
func (a *Attribute) ValidateDataType() error {
	validTypes := []DataType{DataTypeText, DataTypeNumber, DataTypeBoolean, DataTypeDate}
	
	for _, validType := range validTypes {
		if a.DataType == validType {
			return nil
		}
	}
	
	return fmt.Errorf("invalid data type: %s. Valid types are: %v", a.DataType, validTypes)
}

// BeforeCreate GORM hook
func (a *Attribute) BeforeCreate(tx interface{}) error {
	return a.ValidateDataType()
}

// BeforeUpdate GORM hook
func (a *Attribute) BeforeUpdate(tx interface{}) error {
	return a.ValidateDataType()
}

// GetTableName returns the table name for database operations
func (a *Attribute) GetTableName() string {
	return "attributes"
}

// Helper methods for working with different data types

// ParseValue parses the string value according to the attribute's data type
func (a *Attribute) ParseValue(valueStr string) (interface{}, error) {
	switch a.DataType {
	case DataTypeText:
		return valueStr, nil
		
	case DataTypeNumber:
		return strconv.ParseFloat(valueStr, 64)
		
	case DataTypeBoolean:
		return strconv.ParseBool(valueStr)
		
	case DataTypeDate:
		// Try multiple date formats
		formats := []string{
			"2006-01-02",                // ISO date: 2023-12-31
			"2006-01-02 15:04:05",       // ISO datetime: 2023-12-31 23:59:59
			time.RFC3339,                 // RFC3339: 2023-12-31T23:59:59Z
			"02/01/2006",                 // DD/MM/YYYY: 31/12/2023
			"01/02/2006",                 // MM/DD/YYYY: 12/31/2023
		}
		
		var lastErr error
		for _, format := range formats {
			if t, err := time.Parse(format, valueStr); err == nil {
				return t, nil
			} else {
				lastErr = err
			}
		}
		return nil, fmt.Errorf("invalid date format: %w", lastErr)
		
	default:
		return nil, fmt.Errorf("unsupported data type: %s", a.DataType)
	}
}

// FormatValue formats a value to string for storage
func (a *Attribute) FormatValue(value interface{}) (string, error) {
	switch a.DataType {
	case DataTypeText:
		if str, ok := value.(string); ok {
			return str, nil
		}
		return fmt.Sprintf("%v", value), nil
		
	case DataTypeNumber:
		switch v := value.(type) {
		case int:
			return fmt.Sprintf("%d", v), nil
		case int32:
			return fmt.Sprintf("%d", v), nil
		case int64:
			return fmt.Sprintf("%d", v), nil
		case float32:
			return fmt.Sprintf("%.2f", v), nil
		case float64:
			return fmt.Sprintf("%.2f", v), nil
		default:
			return "", fmt.Errorf("invalid number value: %v", value)
		}
		
	case DataTypeBoolean:
		if b, ok := value.(bool); ok {
			return strconv.FormatBool(b), nil
		}
		return "", fmt.Errorf("invalid boolean value: %v", value)
		
	case DataTypeDate:
		switch v := value.(type) {
		case time.Time:
			return v.Format("2006-01-02"), nil
		case string:
			// Validate it's a valid date string
			if _, err := time.Parse("2006-01-02", v); err == nil {
				return v, nil
			}
			return "", fmt.Errorf("invalid date string format, use YYYY-MM-DD: %s", v)
		default:
			return "", fmt.Errorf("invalid date value: %v", value)
		}
		
	default:
		return "", fmt.Errorf("unsupported data type: %s", a.DataType)
	}
}

// ValidateValue validates if a value is valid for this attribute's data type
func (a *Attribute) ValidateValue(valueStr string) error {
	_, err := a.ParseValue(valueStr)
	return err
}

// MarshalJSON custom JSON marshaling
func (a *Attribute) MarshalJSON() ([]byte, error) {
	type Alias Attribute
	return json.Marshal(&struct {
		*Alias
		DataType string `json:"data_type"`
	}{
		Alias:    (*Alias)(a),
		DataType: string(a.DataType),
	})
}
