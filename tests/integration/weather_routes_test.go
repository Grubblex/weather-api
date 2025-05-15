package routes_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Grubblex/weather-api/handlers"
	"github.com/Grubblex/weather-api/middleware"
	"github.com/Grubblex/weather-api/models"
	"github.com/Grubblex/weather-api/repositories"
	"github.com/Grubblex/weather-api/routes"
	"github.com/gofiber/contrib/websocket"
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

    db.AutoMigrate(&models.WeatherData{})

    testData := []models.WeatherData{
        {
            Date:        time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
            Temperature: 63.14634353949799,
            Humidity:    16.522612484526494,
        },
        {
            Date:        time.Date(2023, time.October, 2, 0, 0, 0, 0, time.UTC),
            Temperature: 65.23456789012345,
            Humidity:    18.76543210987654,
        },
        {
            Date:        time.Date(2023, time.October, 3, 0, 0, 0, 0, time.UTC),
            Temperature: 60.98765432123456,
            Humidity:    20.12345678901234,
        },
    }

    if result := db.Create(&testData); result.Error != nil {
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

	// POST /api/v1/weather

	t.Run("POST /api/v1/weather ( Date exists and format is correct )", func(t *testing.T) {
    invalidData := "2023-12-01\t22.526308288\t75.5390732852052"
    req := httptest.NewRequest("POST", "/api/v1/weather", strings.NewReader(invalidData))
    req.Header.Set("Content-Type", "text/plain")
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode) // 200 OK
	})

	
	t.Run("POST /api/v1/weather ( No data )", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/weather", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode) // 400 Bad Request
	})


	t.Run("POST /api/v1/weather ( Wrong data format (date) )", func(t *testing.T) {
    invalidData := "2023-02-2\t22.526308288\t75.5390732852052"
    req := httptest.NewRequest("POST", "/api/v1/weather", strings.NewReader(invalidData))
    req.Header.Set("Content-Type", "text/plain")
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("POST /api/v1/weather ( Date already exists in database )", func(t *testing.T) {
    invalidData := "2023-10-01\t22.526308288\t75.5390732852052"
    req := httptest.NewRequest("POST", "/api/v1/weather", strings.NewReader(invalidData))
    req.Header.Set("Content-Type", "text/plain")
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 500, resp.StatusCode) // 500 Internal Server Error
	})


	// GET /api/v1/weather/date/:date

	t.Run("GET /api/v1/weather/date/:date (Date exists)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/weather/date/2023-10-01", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode) // 200 OK
	})

	t.Run("GET /api/v1/weather/date/:date ( Wrong date format )", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/weather/date/1890-01-0", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode) // 400 Bad Request
	})

	t.Run("GET /api/v1/weather/date/:date ( Date doesnt exist )", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/weather/date/1890-01-01", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode) // 404 Not Found
	})


	// GET /api/v1/weather/range?start=<date>&end=<date>

	t.Run("GET /api/v1/weather/range?start=<date>&end=<date> (Dates exist and date format is correct)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/weather/range?start=2023-10-01&end=2023-10-03", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode) // 200 OK
	})

	t.Run("GET /api/v1/weather/range?start=<date>&end=<date> ( Wrong date format )", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/weather/range?start=2023-1-01&end=2023-10-03", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode) // 400 Bad Request
	})

	t.Run("GET /api/v1/weather/range?start=<date>&end=<date> ( Start date must be before end date )", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/weather/range?start=2023-10-03&end=2023-10-01", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode) // 400 Bad Request
	})

	t.Run("GET /api/v1/weather/range?start=<date>&end=<date> ( Date range doesnt exist in database )", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/weather/range?start=2021-01-05&end=2021-01-30", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode) // 404 Not found
	})


	// ws /api/v1/ws
	t.Run("ws /api/v1/ws ( Successfull websocket connection )", func(t *testing.T) {

		app := fiber.New()
		app.Use("/ws", middleware.WebSocketUpgrade)
		
	
		app.Get("/ws", websocket.New(func(c *websocket.Conn) {
	
			c.Close()
		}))

		req := httptest.NewRequest("GET", "/ws", nil)
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Upgrade", "websocket") 
		req.Header.Set("Sec-WebSocket-Version", "13")
		req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

		resp, _ := app.Test(req)
		defer resp.Body.Close()
		
		assert.Equal(t, fiber.StatusSwitchingProtocols, resp.StatusCode) // 101 Switching Protocols
		assert.Equal(t, "websocket", resp.Header.Get("Upgrade"))
	})
	
}