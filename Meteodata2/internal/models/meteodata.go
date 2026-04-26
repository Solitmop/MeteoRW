package models

import (
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// MeteoData represents meteorological data
type MeteoData struct {
	ID          string    `json:"id,omitempty"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Pressure    float64   `json:"pressure"`
	WindSpeed   float64   `json:"wind_speed"`
	WindDir     float64   `json:"wind_dir"`
	Rainfall    float64   `json:"rainfall"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToPoint converts MeteoData to InfluxDB Point
func (m *MeteoData) ToPoint() *write.Point {
	return write.NewPoint(
		"meteodata",
		map[string]string{
			"id": m.ID,
		},
		map[string]interface{}{
			"temperature": m.Temperature,
			"humidity":    m.Humidity,
			"pressure":    m.Pressure,
			"wind_speed":  m.WindSpeed,
			"wind_dir":    m.WindDir,
			"rainfall":    m.Rainfall,
		},
		m.CreatedAt,
	)
}

// FromValues creates a MeteoData from InfluxDB values
func FromValues(values map[string]interface{}, timestamp time.Time) *MeteoData {
	m := &MeteoData{
		CreatedAt: timestamp,
	}

	if id, ok := values["id"].(string); ok {
		m.ID = id
	}
	if temperature, ok := values["temperature"].(float64); ok {
		m.Temperature = temperature
	} else if temperature, ok := values["temperature"].(int64); ok {
		m.Temperature = float64(temperature)
	}
	if humidity, ok := values["humidity"].(float64); ok {
		m.Humidity = humidity
	} else if humidity, ok := values["humidity"].(int64); ok {
		m.Humidity = float64(humidity)
	}
	if pressure, ok := values["pressure"].(float64); ok {
		m.Pressure = pressure
	} else if pressure, ok := values["pressure"].(int64); ok {
		m.Pressure = float64(pressure)
	}
	if windSpeed, ok := values["wind_speed"].(float64); ok {
		m.WindSpeed = windSpeed
	} else if windSpeed, ok := values["wind_speed"].(int64); ok {
		m.WindSpeed = float64(windSpeed)
	}
	if windDir, ok := values["wind_dir"].(float64); ok {
		m.WindDir = windDir
	} else if windDir, ok := values["wind_dir"].(int64); ok {
		m.WindDir = float64(windDir)
	}
	if rainfall, ok := values["rainfall"].(float64); ok {
		m.Rainfall = rainfall
	} else if rainfall, ok := values["rainfall"].(int64); ok {
		m.Rainfall = float64(rainfall)
	}

	return m
}
