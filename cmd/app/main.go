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

	// Test Category creation with auto slug generation
	category := models.Category{
		Base: models.Base{
			CreatedBy: systemUser.ID,
			UpdatedBy: systemUser.ID,
			IsActive:  true,
		},
		Name:        "Electronics & Gadgets",
		// Slug will be auto-generated from Name
		Description: "Electronic products and gadgets",
	}

	err = db.Create(&category).Error
	if err != nil {
		panic(err)
	}
	fmt.Printf("Category created successfully with slug: %s\n", category.Slug)

	// Test child category creation
	childCategory := models.Category{
		Base: models.Base{
			CreatedBy: systemUser.ID,
			UpdatedBy: systemUser.ID,
			IsActive:  true,
		},
		Name:        "Smart Phones & Accessories",
		// Slug will be auto-generated from Name
		Description: "Mobile phones and accessories",
		ParentID:    &category.ID,
	}

	err = db.Create(&childCategory).Error
	if err != nil {
		panic(err)
	}
	fmt.Printf("Child category created successfully with slug: %s\n", childCategory.Slug)

	// Test Product creation with auto slug generation
	product := models.Product{
		Base: models.Base{
			CreatedBy: systemUser.ID,
			UpdatedBy: systemUser.ID,
			IsActive:  true,
		},
		Name:        "iPhone 15 Pro Max",
		// Slug will be auto-generated from Name
		Description: "Latest iPhone with advanced features",
		CategoryID:  childCategory.ID,
	}

	err = db.Create(&product).Error
	if err != nil {
		panic(err)
	}
	fmt.Printf("Product created successfully with slug: %s\n", product.Slug)

	// Test SKU creation with auto slug generation
	sku := models.Sku{
		Base: models.Base{
			CreatedBy: systemUser.ID,
			UpdatedBy: systemUser.ID,
			IsActive:  true,
		},
		Name:        "iPhone 15 Pro Max 256GB Space Black",
		// Slug will be auto-generated from Name
		Description: "iPhone 15 Pro Max with 256GB storage in Space Black",
		SkuNumber:   "IPH15PM-256-SB",
		Price:       15999000.00,
		ProductID:   product.ID,
	}

	err = db.Create(&sku).Error
	if err != nil {
		panic(err)
	}
	fmt.Printf("SKU created successfully with slug: %s\n", sku.Slug)

	// Test duplicate name handling
	duplicateCategory := models.Category{
		Base: models.Base{
			CreatedBy: systemUser.ID,
			UpdatedBy: systemUser.ID,
			IsActive:  true,
		},
		Name:        "Electronics & Gadgets", // Same name as first category
		Description: "Another electronics category",
	}

	err = db.Create(&duplicateCategory).Error
	if err != nil {
		panic(err)
	}
	fmt.Printf("Duplicate category created with unique slug: %s\n", duplicateCategory.Slug)
}