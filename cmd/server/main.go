package main

import (
	"log"
	"os"

	"github.com/Grubblex/weather-api/database"
	"github.com/Grubblex/weather-api/handlers"
	"github.com/Grubblex/weather-api/repositories"
	"github.com/Grubblex/weather-api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)


func main() {

	err := godotenv.Load()
	if err != nil {
		log.Print("Loading .env file failed!")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default
	}
	
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal("Couldnt establish database connection: ", err)
	}
	
	repo := repositories.NewWeatherRepository(db)
	handler := &handlers.WeatherHandler{Repo: repo}

	app := fiber.New()

	routes.SetupRoutes(app, handler)
	
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Couldnt start server: ", err)
	}

}