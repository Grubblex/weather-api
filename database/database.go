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

var DB *gorm.DB

func ConnectDb () {
	
	err := godotenv.Load()
	if err != nil {
		log.Print("Loading .env file failed !")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Fehler bei der DB-Verbindung: ", err)
	}

	db.Logger = logger.Default.LogMode(logger.Info)
	db.AutoMigrate(&models.WeatherData{})

	DB = db
	log.Print("Successfully connected to db")

}