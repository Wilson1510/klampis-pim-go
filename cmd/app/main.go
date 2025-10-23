package main

import (
	"fmt"
	"github.com/Wilson1510/klampis-pim-go/internal/config"
	"github.com/Wilson1510/klampis-pim-go/internal/database"
	"github.com/Wilson1510/klampis-pim-go/internal/models"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := database.NewConnection(&config.Database)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	defer sqlDB.Close()

	// Run database migrations
	err = database.AutoMigrate(db)
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}

	fmt.Println("Database connected successfully")
	fmt.Println("Database migration completed successfully")

	// Create or get system user for audit fields
	var systemUser models.User
	err = db.Where("username = ?", "system").First(&systemUser).Error
	if err != nil {
		// Create system user if not exists
		systemUser = models.User{
			Username: "system",
			Password: "system123", // In production, use proper hashing
			Name:     "System User",
			Role:     models.RoleSystem,
		}
		err = db.Create(&systemUser).Error
		if err != nil {
			panic(fmt.Sprintf("Failed to create system user: %v", err))
		}
		fmt.Println("System user created successfully")
	} else {
		fmt.Println("System user already exists")
	}
}
