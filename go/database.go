/*
 * Grocery List Management API
 *
 * Database configuration and setup for SQLAlchemy ORM
 */

package openapi

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global database instance
var DB *gorm.DB

// GroceryItem represents a grocery item in the database
type GroceryItem struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Quantity  int32     `json:"quantity" gorm:"not null"`
	Unit      string    `json:"unit"`
	Category  string    `json:"category"`
	Notes     string    `json:"notes"`
	Purchased bool      `json:"purchased" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// SetupDB initializes the database connection and migrations
func SetupDB() error {
	// Get database URL from environment variable or use default for local development
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres:postgres@localhost:5432/grocerydb"
	}

	// Open database connection
	var err error
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Migrate the schema
	err = DB.AutoMigrate(&GroceryItem{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return nil
}
