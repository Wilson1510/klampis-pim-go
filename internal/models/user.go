package models

import (
	"fmt"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleSystem UserRole = "SYSTEM"
    RoleAdmin  UserRole = "ADMIN"
    RoleUser   UserRole = "USER"
)

type User struct {
	gorm.Model
	Username string   `gorm:"uniqueIndex;not null;type:varchar(50)" json:"username"`
	Password string   `gorm:"not null;type:varchar(255)" json:"-"`
	Name     string   `gorm:"not null;type:varchar(50)" json:"name"`
	Role     UserRole `gorm:"not null;type:varchar(20)" json:"role"`
}

// BeforeCreate is a GORM hook that runs before creating a record
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.validateRole()
}

// BeforeUpdate is a GORM hook that runs before updating a record
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.validateRole()
}

// validateRole checks if the role is valid
func (u *User) validateRole() error {
	validRoles := []UserRole{RoleSystem, RoleAdmin, RoleUser}
	
	for _, validRole := range validRoles {
		if u.Role == validRole {
			return nil
		}
	}
	
	return fmt.Errorf("invalid role: %s. Valid roles are: %v", u.Role, validRoles)
}
