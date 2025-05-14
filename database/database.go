package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Grubblex/weather-api/models"
)


func NewDatabase() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Print("Loading .env file failed!")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	if err := db.AutoMigrate(&models.WeatherData{}); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	log.Print("Successfully connected to DB")
	return db, nil
}