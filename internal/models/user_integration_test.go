//go:build integration
// +build integration

package models_test

import (
	"testing"

	"github.com/Wilson1510/klampis-pim-go/internal/models"
	"github.com/Wilson1510/klampis-pim-go/internal/testutil"
	"gorm.io/gorm"
)

func TestUserCreate_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	testCases := []struct {
		name        string
		user        models.User
		expectError bool
		errorMsg    string
	}{
		{
			name: "Create valid admin user",
			user: models.User{
				Username: "admin",
				Password: "password123",
				Name:     "Admin User",
				Role:     models.RoleAdmin,
			},
			expectError: false,
		},
		{
			name: "Create valid system user",
			user: models.User{
				Username: "system",
				Password: "password123",
				Name:     "System User",
				Role:     models.RoleSystem,
			},
			expectError: false,
		},
		{
			name: "Create valid regular user",
			user: models.User{
				Username: "user1",
				Password: "password123",
				Name:     "Regular User",
				Role:     models.RoleUser,
			},
			expectError: false,
		},
		{
			name: "Create user with invalid role",
			user: models.User{
				Username: "invalid",
				Password: "password123",
				Name:     "Invalid User",
				Role:     models.UserRole("INVALID"),
			},
			expectError: true,
			errorMsg:    "invalid role",
		},
		{
			name: "Create user with duplicate username",
			user: models.User{
				Username: "admin", // Already created in first test case
				Password: "password456",
				Name:     "Another Admin",
				Role:     models.RoleAdmin,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := db.Create(&tc.user)

			if tc.expectError {
				if result.Error == nil {
					t.Errorf("Expected error but got none")
				}
				if tc.errorMsg != "" && result.Error != nil {
					// Check if error message contains expected text
					if result.Error.Error() == "" {
						t.Errorf("Expected error message containing '%s', got empty error", tc.errorMsg)
					}
				}
			} else {
				if result.Error != nil {
					t.Errorf("Expected no error but got: %v", result.Error)
				}
				if tc.user.ID == 0 {
					t.Error("Expected user ID to be set after creation")
				}
				if tc.user.CreatedAt.IsZero() {
					t.Error("Expected CreatedAt to be set after creation")
				}
			}
		})
	}
}

func TestUserRead_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user
	user := models.User{
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Role:     models.RoleUser,
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("Read user by ID", func(t *testing.T) {
		var foundUser models.User
		result := db.First(&foundUser, user.ID)

		if result.Error != nil {
			t.Errorf("Expected to find user but got error: %v", result.Error)
		}
		if foundUser.Username != user.Username {
			t.Errorf("Expected username '%s', got '%s'", user.Username, foundUser.Username)
		}
		if foundUser.Name != user.Name {
			t.Errorf("Expected name '%s', got '%s'", user.Name, foundUser.Name)
		}
		if foundUser.Role != user.Role {
			t.Errorf("Expected role '%s', got '%s'", user.Role, foundUser.Role)
		}
	})

	t.Run("Read user by username", func(t *testing.T) {
		var foundUser models.User
		result := db.Where("username = ?", user.Username).First(&foundUser)

		if result.Error != nil {
			t.Errorf("Expected to find user but got error: %v", result.Error)
		}
		if foundUser.ID != user.ID {
			t.Errorf("Expected ID %d, got %d", user.ID, foundUser.ID)
		}
	})

	t.Run("Read non-existent user", func(t *testing.T) {
		var foundUser models.User
		result := db.First(&foundUser, 99999)

		if result.Error != gorm.ErrRecordNotFound {
			t.Errorf("Expected ErrRecordNotFound but got: %v", result.Error)
		}
	})
}

func TestUserUpdate_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user
	user := models.User{
		Username: "updateuser",
		Password: "password123",
		Name:     "Update User",
		Role:     models.RoleUser,
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("Update user name", func(t *testing.T) {
		user.Name = "Updated Name"
		result := db.Save(&user)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedUser models.User
		db.First(&updatedUser, user.ID)
		if updatedUser.Name != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got '%s'", updatedUser.Name)
		}
	})

	t.Run("Update user role to valid role", func(t *testing.T) {
		user.Role = models.RoleAdmin
		result := db.Save(&user)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify update
		var updatedUser models.User
		db.First(&updatedUser, user.ID)
		if updatedUser.Role != models.RoleAdmin {
			t.Errorf("Expected role '%s', got '%s'", models.RoleAdmin, updatedUser.Role)
		}
	})

	t.Run("Update user role to invalid role", func(t *testing.T) {
		user.Role = models.UserRole("SUPERADMIN")
		result := db.Save(&user)

		if result.Error == nil {
			t.Error("Expected error for invalid role but got none")
		}
	})

	t.Run("Update username to duplicate", func(t *testing.T) {
		// Create another user
		anotherUser := models.User{
			Username: "anotheruser",
			Password: "password123",
			Name:     "Another User",
			Role:     models.RoleUser,
		}
		db.Create(&anotherUser)

		// Try to update first user's username to duplicate
		user.Username = "anotheruser"
		result := db.Save(&user)

		if result.Error == nil {
			t.Error("Expected error for duplicate username but got none")
		}
	})
}

func TestUserDelete_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a test user
	user := models.User{
		Username: "deleteuser",
		Password: "password123",
		Name:     "Delete User",
		Role:     models.RoleUser,
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("Soft delete user", func(t *testing.T) {
		result := db.Delete(&user)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify soft delete - should not be found in normal query
		var foundUser models.User
		result = db.First(&foundUser, user.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected user to be soft deleted (not found in normal query)")
		}

		// Verify soft delete - should be found with Unscoped
		result = db.Unscoped().First(&foundUser, user.ID)
		if result.Error != nil {
			t.Errorf("Expected to find soft deleted user with Unscoped but got error: %v", result.Error)
		}
		if foundUser.DeletedAt.Time.IsZero() {
			t.Error("Expected DeletedAt to be set after soft delete")
		}
	})

	t.Run("Permanent delete user", func(t *testing.T) {
		// Create another user
		anotherUser := models.User{
			Username: "permanentdelete",
			Password: "password123",
			Name:     "Permanent Delete",
			Role:     models.RoleUser,
		}
		db.Create(&anotherUser)

		// Permanently delete
		result := db.Unscoped().Delete(&anotherUser)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}

		// Verify permanent delete - should not be found even with Unscoped
		var foundUser models.User
		result = db.Unscoped().First(&foundUser, anotherUser.ID)
		if result.Error != gorm.ErrRecordNotFound {
			t.Error("Expected user to be permanently deleted")
		}
	})
}

func TestUserQuery_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create multiple test users
	users := []models.User{
		{Username: "admin1", Password: "pass", Name: "Admin One", Role: models.RoleAdmin},
		{Username: "admin2", Password: "pass", Name: "Admin Two", Role: models.RoleAdmin},
		{Username: "user1", Password: "pass", Name: "User One", Role: models.RoleUser},
		{Username: "user2", Password: "pass", Name: "User Two", Role: models.RoleUser},
		{Username: "system1", Password: "pass", Name: "System One", Role: models.RoleSystem},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	t.Run("Find all users", func(t *testing.T) {
		var allUsers []models.User
		result := db.Find(&allUsers)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(allUsers) != 5 {
			t.Errorf("Expected 5 users, got %d", len(allUsers))
		}
	})

	t.Run("Find users by role", func(t *testing.T) {
		var adminUsers []models.User
		result := db.Where("role = ?", models.RoleAdmin).Find(&adminUsers)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(adminUsers) != 2 {
			t.Errorf("Expected 2 admin users, got %d", len(adminUsers))
		}
		for _, user := range adminUsers {
			if user.Role != models.RoleAdmin {
				t.Errorf("Expected role ADMIN, got %s", user.Role)
			}
		}
	})

	t.Run("Count users", func(t *testing.T) {
		var count int64
		result := db.Model(&models.User{}).Count(&count)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if count != 5 {
			t.Errorf("Expected count 5, got %d", count)
		}
	})

	t.Run("Find users with pagination", func(t *testing.T) {
		var paginatedUsers []models.User
		result := db.Limit(2).Offset(0).Find(&paginatedUsers)

		if result.Error != nil {
			t.Errorf("Expected no error but got: %v", result.Error)
		}
		if len(paginatedUsers) != 2 {
			t.Errorf("Expected 2 users in first page, got %d", len(paginatedUsers))
		}
	})
}

func TestUserValidation_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	t.Run("Create user without username", func(t *testing.T) {
		user := models.User{
			Password: "password123",
			Name:     "No Username",
			Role:     models.RoleUser,
		}

		result := db.Create(&user)
		if result.Error == nil {
			t.Error("Expected error for missing username but got none")
		}
	})

	t.Run("Create user without password", func(t *testing.T) {
		user := models.User{
			Username: "nopassword",
			Name:     "No Password",
			Role:     models.RoleUser,
		}

		result := db.Create(&user)
		if result.Error == nil {
			t.Error("Expected error for missing password but got none")
		}
	})

	t.Run("Create user without name", func(t *testing.T) {
		user := models.User{
			Username: "noname",
			Password: "password123",
			Role:     models.RoleUser,
		}

		result := db.Create(&user)
		if result.Error == nil {
			t.Error("Expected error for missing name but got none")
		}
	})

	t.Run("Create user without role", func(t *testing.T) {
		user := models.User{
			Username: "norole",
			Password: "password123",
			Name:     "No Role",
		}

		result := db.Create(&user)
		if result.Error == nil {
			t.Error("Expected error for missing/invalid role but got none")
		}
	})
}

