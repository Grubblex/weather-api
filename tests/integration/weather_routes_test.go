package routes_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Grubblex/weather-api/handlers"
	"github.com/Grubblex/weather-api/models"
	"github.com/Grubblex/weather-api/repositories"
	"github.com/Grubblex/weather-api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)


func setupTestDB() *gorm.DB {
	
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	date, err := time.Parse("2006-01-02", "2023-10-01")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.WeatherData{})

	result := db.Create(&models.WeatherData{
		Date:        date,
		Temperature: 63.14634353949799,
		Humidity:    16.522612484526494,
	})

	if result.Error != nil {
		panic("failed to insert test data: " + result.Error.Error())
	}

	return db
}

func TestWeatherRoutesWithDB(t *testing.T) {

	db := setupTestDB()


	repo := repositories.NewWeatherRepository(db)
	handler := &handlers.WeatherHandler{Repo: repo}


	app := fiber.New()
	routes.SetupRoutes(app, handler)


	t.Run("GET /api/v1/weather/date/:date (Date exists)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/weather/date/2023-10-01", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode) // 200 OK
	})


	t.Run("GET /api/v1/weather/date/:date ( Date doesnt exist )", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/weather/date/1890-01-01", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode) // 404 Not Found
	})

	t.Run("GET /api/v1/weather/date/:date ( Wrong date format )", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/weather/date/1890-01-0", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode) // 400 Not Found
	})

	t.Run("POST /api/v1/addWeatherEntry ( Date already exists in database )", func(t *testing.T) {
    invalidData := "2023-10-01\t22.526308288\t75.5390732852052"
    req := httptest.NewRequest("POST", "/api/v1/addWeatherEntry", strings.NewReader(invalidData))
    req.Header.Set("Content-Type", "text/plain")
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 500, resp.StatusCode) // 500 Internal Server Error
	})

	t.Run("POST /api/v1/addWeatherEntry ( No data )", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/addWeatherEntry", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode) // 400 Bad Request
	})


	t.Run("POST /api/v1/addWeatherEntry ( Wrong data fromat (date) )", func(t *testing.T) {
    invalidData := "2023-02-2\t22.526308288\t75.5390732852052"
    req := httptest.NewRequest("POST", "/api/v1/addWeatherEntry", strings.NewReader(invalidData))
    req.Header.Set("Content-Type", "text/plain")
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 400, resp.StatusCode)
	})
}