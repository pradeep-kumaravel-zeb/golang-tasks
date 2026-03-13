package database

import (
	"fmt"
	"student-enrollment/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func EstablishPostgresConnection(dburl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dburl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run AutoMigrate
	if err := models.AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	fmt.Println("Database connected and tables migrated successfully")
	return db, nil
}
