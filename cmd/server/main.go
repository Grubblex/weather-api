package main

import (
	"github.com/Grubblex/weather-api/database"
	"github.com/Grubblex/weather-api/routes"
	"github.com/gofiber/fiber/v2"
)


func main() {
	
	database.ConnectDb()
	
	app := fiber.New()
	routes.SetupRoutes(app)
	

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}

}