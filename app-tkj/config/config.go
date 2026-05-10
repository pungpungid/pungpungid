package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var DB *sql.DB

type Config struct {
	Port        string
	DatabaseURL string
	SessionKey  string
}

// Load configuration from environment
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8090"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		SessionKey:  getEnv("SESSION_KEY", "default-session-key-change-in-production"),
	}
}

// InitDB initializes the database connection
func InitDB(databaseURL string) error {
	var err error
	DB, err = sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✅ Database connected successfully")
	return nil
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
