package models

import (
	"time"
)

type WeatherData struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Date        time.Time `gorm:"type:date;uniqueIndex" json:"date"`
	Humidity    float64     `json:"humidity"`
	Temperature float64  `json:"temperature"`
}