package repositories

import (
	"time"

	"github.com/Grubblex/weather-api/models"
	"gorm.io/gorm"
)

type WeatherRepository struct {
    DB *gorm.DB
}

func NewWeatherRepository(db *gorm.DB) *WeatherRepository {
    return &WeatherRepository{DB: db}
}

func (r *WeatherRepository) InsertWeather(weather models.WeatherData) (*models.WeatherData, error) {
	err := r.DB.Create(&weather).Error
	return &weather, err
}

func (r *WeatherRepository) GetWeatherByDate(date time.Time) (*models.WeatherData, error) {
	var data models.WeatherData
	err := r.DB.Where("date = ?", date).First(&data).Error
	return &data, err
}

func (r *WeatherRepository) GetWeatherByDateRange(start time.Time, end time.Time) ([]models.WeatherData, error) {
    var data []models.WeatherData
    if err := r.DB.Where("date BETWEEN ? AND ?", start, end).Find(&data).Error; err != nil {
        return nil, err
    }
    if len(data) == 0 {
        return nil, gorm.ErrRecordNotFound 
    }
    return data, nil
}