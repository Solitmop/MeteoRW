package database

import (
	"context"

	"meteodata2/config"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// InfluxDBClient wraps the InfluxDB client
type InfluxDBClient struct {
	client influxdb2.Client
	config config.InfluxDBConfig
}

// NewInfluxDBClient creates a new InfluxDB client
func NewInfluxDBClient(cfg config.InfluxDBConfig) (*InfluxDBClient, error) {
	client := influxdb2.NewClient(cfg.URL, cfg.Token)

	// Test the connection
	_, err := client.Health(context.Background())
	if err != nil {
		return nil, err
	}

	return &InfluxDBClient{
		client: client,
		config: cfg,
	}, nil
}

// GetClient returns the underlying InfluxDB client
func (c *InfluxDBClient) GetClient() influxdb2.Client {
	return c.client
}

// GetConfig returns the InfluxDB configuration
func (c *InfluxDBClient) GetConfig() config.InfluxDBConfig {
	return c.config
}

// Close closes the InfluxDB client connection
func (c *InfluxDBClient) Close() {
	c.client.Close()
}
