//go:build integration
// +build integration

package testutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/Wilson1510/klampis-pim-go/internal/config"
	"github.com/Wilson1510/klampis-pim-go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB creates a test database for integration tests
// The database is automatically created before the test and dropped after
//
// By default, it uses the same PostgreSQL connection as your main app (from .env)
// but creates a separate database for testing (default: klampis_pim_test)
//
// Environment variables:
//   DB_HOST, DB_PORT, DB_USER, DB_PASSWORD - same as main app (from .env)
//   TEST_DB_NAME - optional, defaults to klampis_pim_test
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Get base configuration
	baseConfig := getBaseDBConfig()

	dbTestName := getEnvOrDefault("TEST_DB_NAME", "klampis_pim_test")

	// Step 1: Connect to postgres database to create the test database
	if err := createDatabase(baseConfig, dbTestName); err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Step 2: Connect to the newly created test database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		baseConfig.Host,
		baseConfig.Port,
		baseConfig.User,
		baseConfig.Password,
		dbTestName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Try to cleanup the database if connection fails
		_ = dropDatabase(baseConfig, dbTestName)
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Verify connection
	sqlDB, err := db.DB()
	if err != nil {
		closeDB(db)
		_ = dropDatabase(baseConfig, dbTestName)
		t.Fatalf("Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		closeDB(db)
		_ = dropDatabase(baseConfig, dbTestName)
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Step 3: Run migrations
	if err := runMigrations(db); err != nil {
		// Cleanup on migration failure
		closeDB(db)
		_ = dropDatabase(baseConfig, dbTestName)
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

// CleanupTestDB drops the test database and closes the connection
// This is called automatically via defer after each test
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Close the connection
	closeDB(db)
	dbTestName := getEnvOrDefault("TEST_DB_NAME", "klampis_pim_test")

	// Get base config to connect to postgres database
	baseConfig := getBaseDBConfig()

	// Drop the test database
	if err := dropDatabase(baseConfig, dbTestName); err != nil {
		t.Logf("Warning: Failed to drop test database %s: %v", dbTestName, err)
	}
}

// createDatabase creates a new database
func createDatabase(config *config.DatabaseConfig, dbName string) error {
	// Connect to the postgres database to create the test database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer closeDB(db)

	// Create the database
	sql := fmt.Sprintf("CREATE DATABASE %s", dbName)
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

// dropDatabase drops a database
func dropDatabase(config *config.DatabaseConfig, dbName string) error {
	// Connect to the postgres database to drop the test database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer closeDB(db)

	// Terminate all connections to the database
	terminateSQL := fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '%s'
		AND pid <> pg_backend_pid()
	`, dbName)
	db.Exec(terminateSQL)

	// Drop the database
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	return nil
}

// runMigrations runs database migrations for all models
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Sku{},
		&models.Attribute{},
		&models.SkuAttributeValue{},
		&models.Image{},
	)
}

// getBaseDBConfig returns the base database configuration for creating/dropping databases
// It uses the same connection settings as the main app (from .env or defaults)
// but connects to the 'postgres' database to create/drop test databases
func getBaseDBConfig() *config.DatabaseConfig {
	return &config.DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvIntOrDefault("DB_PORT", 5432),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "sofia"),
		Name:     "postgres", // Always connect to postgres database for admin operations
	}
}

// closeDB closes the database connection
func closeDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault returns environment variable as int or default
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

