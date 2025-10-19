package main

import (
	"fmt"
	"github.com/Wilson1510/klampis-pim-go/internal/config"
	"github.com/Wilson1510/klampis-pim-go/internal/database"
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

	fmt.Println("Database connected successfully")
}