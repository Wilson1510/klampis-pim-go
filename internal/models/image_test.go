package models

import (
	"testing"
)

// TestImageTableName tests the custom table name
func TestImageTableName(t *testing.T) {
	image := Image{}
	if image.TableName() != "images" {
		t.Errorf("Expected table name 'images', got '%s'", image.TableName())
	}
}
