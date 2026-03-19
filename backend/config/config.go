package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Log      LogConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Env  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string
}

// New loads and returns configuration from environment variables
func New() *Config {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "lms_user"),
			Password: getEnv("DB_PASSWORD", "lms_password_dev"),
			DBName:   getEnv("DB_NAME", "leave_management_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("SERVER_ENV", "development"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your_secret_key_change_in_production"),
			Expiration: getEnv("JWT_EXPIRATION", "24h"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}

// GetDatabaseURL returns the PostgreSQL connection string
func (dc *DatabaseConfig) GetDatabaseURL() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dc.User,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.DBName,
		dc.SSLMode,
	)
}

// GetLibPQConnectionString returns the libpq connection string
func (dc *DatabaseConfig) GetLibPQConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dc.Host,
		dc.Port,
		dc.User,
		dc.Password,
		dc.DBName,
		dc.SSLMode,
	)
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}