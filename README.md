## Weather-Api

A simple RESTful Api written in GO.

## Tech Stack

- [Go 1.23.3](https://go.dev/dl/)
- [Fiber v2](https://docs.gofiber.io/)
- [PostgreSQL 17](https://www.postgresql.org/docs/)
- [GORM](https://gorm.io/) (ORM)


## Project Structure


```
weather-api/
├── cmd/
│   ├── parser/
│   └── server/
├── database/
├── handlers/
├── middleware/
├── models/
├── repositories/
├── routes/
├── services/
└── tests/
    └── integration/

```

## Setup

Make sure you have a existing `.env` file in your root dir. It should look like the `.env.example` file.

```
# Database Connection
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=secret
DB_NAME=weather_db

# Connection Pool Settings
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=1h
DB_CONN_MAX_IDLE_TIME=30m

# Server Configuration
PORT=3000
```

1. Start the server
```
go run ./cmd/server/main.go
```
2. Start the parser

```
go run ./cmd/parser/main.go
```

## Endpoints

### POST `/api/v1/weather`

Creates a new weather entry. The endpoint accepts raw data in plain text format.

\<date> \<humidity> \<temperature> 

```
2023-10-01 22.526308288 75.5390732852052
```

### GET `/api/v1/weather/date/<date>` 

Returns the data for a specific day.

Example response:

```json
{
    "id": 3,
    "date": "2023-10-01T00:00:00Z",
    "humidity": 22.526308288,
    "temperature": 75.5390732852052
}
```

### GET `/api/v1/weather/range?start=<date>&end=<date>` 

Returns the data for a range of days.

Example response:

```json
[
    {
        "id": 1,
        "date": "2023-10-02T00:00:00Z",
        "humidity": 22.526308288,
        "temperature": 75.5390732852052
    },
    {
        "id": 2,
        "date": "2023-10-03T00:00:00Z",
        "humidity": 22.526308288,
        "temperature": 75.5390732852052
    },
    {
        "id": 3,
        "date": "2023-10-01T00:00:00Z",
        "humidity": 22.526308288,
        "temperature": 75.5390732852052
    }
]
```

### Websocket `/api/v1/ws` 

Establishes a websocket connection and broadcasts the latest entry.

Example:

```json
{"id":3,"date":"2023-10-01T00:00:00Z","humidity":22.526308288,"temperature":75.5390732852052}
```

## Testing

To run a simple test run the following command

```
go test -v tests/integration/weather_routes_test.go
```

If the test was Successful it should return this

```
--- PASS: TestWeatherRoutesWithDB (0.04s)
    --- PASS: TestWeatherRoutesWithDB/POST_/api/v1/weather_(_Date_exists_and_format_is_correct_) (0.00s)
    --- PASS: TestWeatherRoutesWithDB/POST_/api/v1/weather_(_No_data_) (0.00s)
    --- PASS: TestWeatherRoutesWithDB/POST_/api/v1/weather_(_Wrong_data_format_(date)_) (0.00s)
    --- PASS: TestWeatherRoutesWithDB/POST_/api/v1/weather_(_Date_already_exists_in_database_) (0.00s)
    --- PASS: TestWeatherRoutesWithDB/GET_/api/v1/weather/date/:date_(Date_exists) (0.00s)
    --- PASS: TestWeatherRoutesWithDB/GET_/api/v1/weather/date/:date_(_Wrong_date_format_) (0.00s)
    --- PASS: TestWeatherRoutesWithDB/GET_/api/v1/weather/date/:date_(_Date_doesnt_exist_) (0.00s)
    --- PASS: TestWeatherRoutesWithDB/GET_/api/v1/weather/range?start=<date>&end=<date>_(Dates_exist_and_date_format_is_correct) (0.00s)
    --- PASS: TestWeatherRoutesWithDB/GET_/api/v1/weather/range?start=<date>&end=<date>_(_Wrong_date_format_) (0.00s)
    --- PASS: TestWeatherRoutesWithDB/GET_/api/v1/weather/range?start=<date>&end=<date>_(_Start_date_must_be_before_end_date_) (0.00s)
    --- PASS: TestWeatherRoutesWithDB/GET_/api/v1/weather/range?start=<date>&end=<date>_(_Date_range_doesnt_exist_in_database_) (0.00s)
    --- PASS: TestWeatherRoutesWithDB/ws_/api/v1/ws_(_Successfull_websocket_connection_) (0.00s)
PASS
ok      command-line-arguments  0.401s
```

> Make sure that cgo is enabled as go-sqlite3 requires cgo to work.
Also check if the gcc compiler is properly installed and added to path.



