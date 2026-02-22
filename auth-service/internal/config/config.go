package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DatabaseURL      string
	Environment      string
	JWTSecret        string
	JWTAccessExpiry  string
	JWTRefreshExpiry string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config := &Config{
		Port:             getEnv("PORT", ":8001"),
		Environment:      getEnv("ENVIRONMENT", "development"),
		JWTSecret:        getEnv("JWT_SECRET", "default-secret-change-in-production"),
		JWTAccessExpiry:  getEnv("JWT_ACCESS_EXPIRY", "15m"),
		JWTRefreshExpiry: getEnv("JWT_REFRESH_EXPIRY", "168h"),
	}

	config.DatabaseURL = buildDatabaseURL()

	return config
}

func buildDatabaseURL() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "auth_user")
	password := getEnv("DB_PASSWORD", "auth_dev_password")
	dbname := getEnv("DB_NAME", "auth_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
