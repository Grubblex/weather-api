package handlers

import (
	"errors"
	"sync"
	"time"

	"github.com/Grubblex/weather-api/models"
	"github.com/Grubblex/weather-api/repositories"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
    clients = make(map[*websocket.Conn]struct{})
    mu      sync.RWMutex
)


func HandleInsertWeather(c *fiber.Ctx) error {

	date:= c.Locals("date").(time.Time) 
    humidity := c.Locals("humidity").(float64)
    temperature := c.Locals("temperature").(float64)

	weather := models.WeatherData{
		Date:        date,
		Humidity:    humidity,
		Temperature: temperature,
	}

	weatherData, err := repositories.InsertWeather(weather)

	if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Server Error"})
    }

	go BroadcastWeather(*weatherData)

	return c.JSON(weatherData)
}


func HandleWeatherByDate(c *fiber.Ctx) error {
    
	date := c.Locals("date").(time.Time) 

    weatherData, err := repositories.GetWeatherByDate(date)

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return c.Status(404).JSON(fiber.Map{"error": "Date not found"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "Server Error"})
    }

    return c.JSON(weatherData)
}


func HandleWeatherByDateRange(c *fiber.Ctx) error {
    
	start := c.Locals("startDate").(time.Time) 
    end := c.Locals("endDate").(time.Time)
    
	weatherData, err := repositories.GetWeatherByDateRange(start, end)

	if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return c.Status(404).JSON(fiber.Map{"error": "Date not found"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "Server Error"})
    }

    return c.JSON(weatherData)
}


func HandleWebSocket(c *websocket.Conn) {

    mu.Lock()
    clients[c] = struct{}{}
    mu.Unlock()

    defer func() {
        mu.Lock()
        delete(clients, c)
        mu.Unlock()
        c.Close()
    }()

    for {
        if _, _, err := c.ReadMessage(); err != nil {
            break
        }
    }
}


func BroadcastWeather(data models.WeatherData) {
    mu.Lock() 
    defer mu.Unlock()
    
    for client := range clients {
        if err := client.WriteJSON(data); err != nil {
            client.Close()
            delete(clients, client)
        }
    }
}

