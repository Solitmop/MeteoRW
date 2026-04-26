# Meteo Data Service

A microservice for managing meteorological data using Gin framework and InfluxDB as the time-series database.

## Features

- RESTful API for CRUD operations on meteorological data
- Integration with InfluxDB for efficient time-series data storage
- Clean Architecture implementation
- Environment-based configuration
- Input validation

## Prerequisites

- Go 1.21+
- InfluxDB server running

## Installation

1. Clone the repository
2. Navigate to the project directory
3. Install dependencies:

```bash
go mod tidy
```

4. Copy `.env.example` to `.env` and configure your settings:

```bash
cp .env.example .env
```

## Configuration

Edit the `.env` file to configure your InfluxDB connection:

```env
INFLUXDB_HOST=localhost
INFLUXDB_PORT=8086
INFLUXDB_TOKEN=your_token_here
INFLUXDB_ORG=your_org_here
INFLUXDB_BUCKET=your_bucket_here
SERVER_PORT=8080
GIN_MODE=release
```

## Usage

Run the service:

```bash
go run main.go
```

## API Endpoints

### Health Check

- `GET /health` - Check if the service is running

### Meteo Data Operations

- `POST /meteodata` - Create new meteo data
- `GET /meteodata/:id` - Get meteo data by ID
- `GET /meteodata` - Get all meteo data with pagination
- `PUT /meteodata/:id` - Update meteo data
- `DELETE /meteodata/:id` - Delete meteo data
- `GET /meteodata/range` - Get meteo data within a time range

### Examples

#### Create Meteo Data

```bash
curl -X POST http://localhost:8080/meteodata \
  -H "Content-Type: application/json" \
  -d '{
    "temperature": 25.5,
    "humidity": 60.0,
    "pressure": 1013.25,
    "wind_speed": 5.2,
    "wind_dir": 180.0,
    "rainfall": 0.0
  }'
```

#### Get Meteo Data by ID

```bash
curl http://localhost:8080/meteodata/{id}
```

#### Get All Meteo Data with Pagination

```bash
curl "http://localhost:8080/meteodata?limit=10&offset=0"
```

#### Get Meteo Data by Time Range

```bash
curl "http://localhost:8080/meteodata/range?start=2023-01-01T00:00:00Z&end=2023-01-02T00:00:00Z"
```

## Data Model

The meteo data includes the following fields:

- `id`: Unique identifier (auto-generated if not provided)
- `temperature`: Temperature in Celsius
- `humidity`: Humidity percentage
- `pressure`: Atmospheric pressure in hPa
- `wind_speed`: Wind speed in m/s
- `wind_dir`: Wind direction in degrees (0-359)
- `rainfall`: Rainfall amount in mm
- `created_at`: Timestamp of data creation

## Architecture

This service follows Clean Architecture principles:

- **Handlers**: HTTP request/response handling
- **Use Cases**: Business logic layer
- **Repository**: Data access layer
- **Models**: Data structures
- **Database**: InfluxDB client wrapper

## Error Handling

The service returns appropriate HTTP status codes:

- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid input
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Testing

To run tests:

```bash
go test ./...