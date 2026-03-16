package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBCredentials struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Schema   string
	SSLMode  string
}

func LoadConfig() (*DBCredentials, error) {
	// Try to load .env file from multiple locations
	envPaths := []string{
		".env",                    
		"../.env",                  // Parent directory (if running from cmd/)
	}

	loaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Loaded .env from: %s", path)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Println("No .env file found, using environment variables")
	}
	// For testing purpose, the values are hardcoded here.
	creds := &DBCredentials{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "postgres"),
		Schema:   getEnv("DB_SCHEMA", "public"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	if creds.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	return creds, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func BuildDBUrl(creds DBCredentials) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=%s",
		creds.Host,
		creds.Port,
		creds.User,
		creds.Password,
		creds.DBName,
		creds.Schema,
		creds.SSLMode,
	)
}
