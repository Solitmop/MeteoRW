package config

import (
	"os"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig
	InfluxDB InfluxDBConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
}

// InfluxDBConfig holds InfluxDB-related configuration
type InfluxDBConfig struct {
	Host   string
	Port   string
	Token  string
	Org    string
	Bucket string
	URL    string
}

// NewConfig creates a new configuration instance
func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
		},
		InfluxDB: InfluxDBConfig{
			Host:   getEnvOrDefault("INFLUXDB_HOST", "localhost"),
			Port:   getEnvOrDefault("INFLUXDB_PORT", "8086"),
			Token:  getEnvOrDefault("INFLUXDB_TOKEN", ""),
			Org:    getEnvOrDefault("INFLUXDB_ORG", ""),
			Bucket: getEnvOrDefault("INFLUXDB_BUCKET", ""),
			URL:    "http://" + getEnvOrDefault("INFLUXDB_HOST", "localhost") + ":" + getEnvOrDefault("INFLUXDB_PORT", "8086"),
		},
	}
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
