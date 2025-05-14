package repositories

import (
	"fmt"
	"time"

	"github.com/Grubblex/weather-api/database"
	"github.com/Grubblex/weather-api/models"
)

func InsertWeather(weather models.WeatherData) (*models.WeatherData, error) {

	fmt.Print(weather)
	err := database.DB.Create(&weather).Error
	return &weather, err
}

func GetWeatherByDate(date time.Time) (*models.WeatherData, error) {
	var data models.WeatherData
	err := database.DB.Where("date = ?", date).First(&data).Error
	return &data, err
}

func GetWeatherByDateRange(start time.Time, end time.Time) ([]models.WeatherData, error) {
    var data []models.WeatherData  
    err := database.DB.Where("date BETWEEN ? AND ?", start, end).Find(&data).Error
    return data, err
}
