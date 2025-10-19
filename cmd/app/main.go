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
}