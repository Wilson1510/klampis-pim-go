package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name  string
	Port  int
	Debug bool
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret            string
	AccessTokenExpiry time.Duration
}

func getEnvBool(key string) bool {
	return strings.ToLower(os.Getenv(key)) == "true"
}

func getEnvInt(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		panic(err)
	}
	return value
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		App: AppConfig{
			Name:  os.Getenv("APP_NAME"),
			Port:  getEnvInt("APP_PORT"),
			Debug: getEnvBool("APP_DEBUG"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     getEnvInt("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret:            os.Getenv("JWT_SECRET"),
			AccessTokenExpiry: time.Duration(getEnvInt("JWT_ACCESS_TOKEN_EXPIRY")) * time.Second,
		},
	}

	return config, nil
}
