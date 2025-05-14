package main

import (
	"log"

	"github.com/Grubblex/weather-api/database"
	"github.com/Grubblex/weather-api/handlers"
	"github.com/Grubblex/weather-api/repositories"
	"github.com/Grubblex/weather-api/routes"
	"github.com/gofiber/fiber/v2"
)


func main() {
	
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal("Fehler bei der DB-Verbindung: ", err)
	}
	
	repo := repositories.NewWeatherRepository(db)
	handler := &handlers.WeatherHandler{Repo: repo}

	app := fiber.New()

	routes.SetupRoutes(app, handler)
	
	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}

}