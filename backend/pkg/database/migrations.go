package database

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

// MigrationHistory tracks which migrations have been run
type MigrationHistory struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}

// TableName specifies the table name
func (MigrationHistory) TableName() string {
	return "migration_history"
}

// RunMigrations runs all migration files that haven't been run yet
func RunMigrations(db *gorm.DB, migrationsPath string) error {
	// Create migration history table if it doesn't exist
	if !db.Migrator().HasTable(&MigrationHistory{}) {
		if err := db.Migrator().CreateTable(&MigrationHistory{}); err != nil {
			return fmt.Errorf("failed to create migration history table: %w", err)
		}
	}

	// Get all .up.sql files
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Only process .up.sql files
		if !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}

		// Check if migration has already been run
		var history MigrationHistory
		result := db.Where("name = ?", file.Name()).First(&history)
		if result.RowsAffected > 0 {
			fmt.Printf("⏭️  Already executed: %s\n", file.Name())
			continue
		}

		// Read and execute migration
		filePath := filepath.Join(migrationsPath, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		// Execute the migration
		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}

		// Record migration execution
		if err := db.Create(&MigrationHistory{Name: file.Name()}).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file.Name(), err)
		}

		fmt.Printf("✅ Executed migration: %s\n", file.Name())
	}

	return nil
}