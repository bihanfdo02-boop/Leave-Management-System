package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/KalinduBihan/leave-management-api/config"
	"github.com/KalinduBihan/leave-management-api/internal/repository"
	"github.com/KalinduBihan/leave-management-api/internal/routes"
	"github.com/KalinduBihan/leave-management-api/internal/service"
	"github.com/KalinduBihan/leave-management-api/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Run migrations
	migrationsPath := getMigrationsPath()
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	services := service.NewServices(repos, cfg, db)

	// Create Gin router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, services, cfg)

	// Start server
	port := cfg.Server.Port
	fmt.Printf("✅ Application started successfully\n")
	fmt.Printf("📊 Database: %s\n", cfg.Database.DBName)
	fmt.Printf("🚀 Server listening on port %s\n", port)
	fmt.Printf("🏥 Health check: http://localhost:%s/health\n", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getMigrationsPath returns the path to migrations directory
func getMigrationsPath() string {
	// Try to find migrations directory relative to current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// Check if migrations exist in current directory
	migPath := filepath.Join(cwd, "migrations")
	if _, err := os.Stat(migPath); err == nil {
		return migPath
	}

	// Check parent directory (in case running from backend folder)
	migPath = filepath.Join(cwd, "..", "migrations")
	if _, err := os.Stat(migPath); err == nil {
		return migPath
	}

	// Default to ./migrations
	return "./migrations"
}