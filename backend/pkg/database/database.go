package database

import (
	"fmt"

	"github.com/KalinduBihan/leave-management-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Database.GetLibPQConnectionString()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Note: We're NOT using AutoMigrate because migrations are already
	// applied via docker-compose entrypoint scripts (SQL files in migrations/)
	// AutoMigrate would cause conflicts with existing constraints

	return db, nil
}