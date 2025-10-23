package models

import (
	"testing"
)

// TestValidateRole tests the validateRole method (Pure Unit Test)
func TestValidateRole(t *testing.T) {
	testCases := []struct {
		name        string
		role        UserRole
		expectError bool
	}{
		// Valid roles
		{
			name:        "Valid System Role",
			role:        RoleSystem,
			expectError: false,
		},
		{
			name:        "Valid Admin Role",
			role:        RoleAdmin,
			expectError: false,
		},
		{
			name:        "Valid User Role",
			role:        RoleUser,
			expectError: false,
		},
		{
			name:        "Valid System Role String",
			role:        "SYSTEM",
			expectError: false,
		},
		{
			name:        "Valid Admin Role String",
			role:        "ADMIN",
			expectError: false,
		},
		{
			name:        "Valid User Role String",
			role:        "USER",
			expectError: false,
		},
		// Invalid roles
		{
			name:        "Invalid Role - SUPERUSER",
			role:        UserRole("SUPERUSER"),
			expectError: true,
		},
		{
			name:        "Invalid Role - Empty String",
			role:        UserRole(""),
			expectError: true,
		},
		{
			name:        "Invalid Role - Lowercase admin",
			role:        UserRole("admin"),
			expectError: true,
		},
		{
			name:        "Invalid Role - Random String",
			role:        UserRole("INVALID_ROLE"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create user with the test role
			user := User{Role: tc.role}

			// Call validateRole method
			err := user.validateRole()

			// Assert the result
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for role '%s', but got nil", tc.role)
				}
				// Check that error message contains expected text
				if err != nil && err.Error() == "" {
					t.Errorf("Expected non-empty error message for role '%s'", tc.role)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for role '%s', but got: %v", tc.role, err)
				}
			}
		})
	}
}
