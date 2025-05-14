package routes

import (
	"github.com/Grubblex/weather-api/handlers"
	"github.com/Grubblex/weather-api/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api", middleware.Logger())
  	v1 := api.Group("/v1") 

  	v1.Get("/hello", handlers.HelloWorld)

  	v1.Post("/addWeatherEntry",
	middleware.ValidateRawData(),
	handlers.HandleInsertWeather)

  	v1.Get("/weather/date/:date", 
		middleware.ValidateDate("date"), 
		handlers.HandleWeatherByDate,    
	)
  	v1.Get("/weather/range", 
		middleware.ValidateDateRange(),  
		handlers.HandleWeatherByDateRange,   
	)

	v1.Use("/ws", middleware.WebSocketUpgrade)
	v1.Get("/ws", websocket.New(handlers.HandleWebSocket))
	 
}