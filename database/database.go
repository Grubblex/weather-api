package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting sql.DB from gorm.DB: %w", err)
	}

	
	maxOpen, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS")) 
	if maxOpen <= 0 {
		maxOpen = 100 
	}


	maxIdle, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS")) 
	if maxIdle <= 0 {
		maxIdle = 10 
	}

	
	maxLifetime, _ := time.ParseDuration(os.Getenv("DB_CONN_MAX_LIFETIME")) 
	if maxLifetime <= 0 {
		maxLifetime = time.Hour 
	}

	
	maxIdleTime, _ := time.ParseDuration(os.Getenv("DB_CONN_MAX_IDLE_TIME")) 
	if maxIdleTime <= 0 {
		maxIdleTime = 30 * time.Minute 
	}

	log.Printf("Pool config: maxOpen=%d, maxIdle=%d, maxLifetime=%v, maxIdleTime=%v", 
    maxOpen, maxIdle, maxLifetime, maxIdleTime)
	
	//Pool config
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetConnMaxIdleTime(maxIdleTime) 

	log.Print("Successfully connected to DB")
	return db, nil
}