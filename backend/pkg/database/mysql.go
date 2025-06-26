package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"concert-booking/pkg/models" // adjust the import path as needed

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connected successfully")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	modelsList := []interface{}{
		&models.User{},
		&models.Concert{},
		&models.Booking{},
		&models.Payment{},
	}

	for _, model := range modelsList {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	// Seed initial concerts
	db.Create(&models.Concert{
		Title:    "Rock Festival 2025",
		Date:     time.Now().AddDate(0, 1, 0),
		Venue:    "National Stadium",
		Capacity: 10000,
	})

	return nil
}
