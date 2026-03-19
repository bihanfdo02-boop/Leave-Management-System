package database

import (
	"fmt"

	"github.com/KalinduBihan/leave-management-api/config"
	"github.com/KalinduBihan/leave-management-api/internal/domain"
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

	// Auto migrate all models
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// autoMigrate migrates all domain models
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Department{},
		&domain.Employee{},
		&domain.LeaveType{},
		&domain.LeaveBalance{},
		&domain.LeaveRequest{},
		&domain.AuditLog{},
	)
}