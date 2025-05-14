package routes

import (
	"github.com/Grubblex/weather-api/handlers"
	"github.com/Grubblex/weather-api/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, handler *handlers.WeatherHandler) {

	api := app.Group("/api", middleware.Logger())
  	v1 := api.Group("/v1") 

  	v1.Post("/addWeatherEntry",
	middleware.ValidateRawData(),
	handler.HandleInsertWeather)

  	v1.Get("/weather/date/:date", 
		middleware.ValidateDate("date"), 
		handler.HandleWeatherByDate,    
	)
  	v1.Get("/weather/range", 
		middleware.ValidateDateRange(),  
		handler.HandleWeatherByDateRange,   
	)

	v1.Use("/ws", middleware.WebSocketUpgrade)
	v1.Get("/ws", websocket.New(handler.HandleWebSocket))
	 
}