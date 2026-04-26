package repository

import (
	"context"
	"fmt"
	"time"

	"meteodata2/internal/models"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// MeteoDataRepository defines the interface for meteo data operations
type MeteoDataRepository interface {
	Create(meteoData *models.MeteoData) error
	GetByID(id string) (*models.MeteoData, error)
	GetAll(limit int, offset int) ([]*models.MeteoData, error)
	Update(meteoData *models.MeteoData) error
	Delete(id string) error
	GetByTimeRange(start, end time.Time) ([]*models.MeteoData, error)
}

// meteoDataRepository implements MeteoDataRepository
type meteoDataRepository struct {
	client   influxdb2.Client
	bucket   string
	org      string
	writeAPI api.WriteAPI
	queryAPI api.QueryAPI
}

// NewMeteoDataRepository creates a new instance of MeteoDataRepository
func NewMeteoDataRepository(client influxdb2.Client, bucket, org string) MeteoDataRepository {
	return &meteoDataRepository{
		client:   client,
		bucket:   bucket,
		org:      org,
		writeAPI: client.WriteAPI(org, bucket),
		queryAPI: client.QueryAPI(org),
	}
}

// Create adds a new meteo data record to the database
func (r *meteoDataRepository) Create(meteoData *models.MeteoData) error {
	if meteoData.ID == "" {
		// Generate a unique ID if not provided
		meteoData.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	point := meteoData.ToPoint()
	r.writeAPI.WritePoint(point)

	// Flush to ensure data is written
	r.writeAPI.Flush()

	return nil
}

// GetByID retrieves a meteo data record by its ID
func (r *meteoDataRepository) GetByID(id string) (*models.MeteoData, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: -30d)
		|> filter(fn: (r) => r["_measurement"] == "meteodata")
		|> filter(fn: (r) => r["id"] == "%s")
		|> sort(columns: ["_time"], desc: true)
		|> limit(n: 1)
	`, r.bucket, id)

	result, err := r.queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var meteoData *models.MeteoData

	for result.Next() {
		values := result.Record().Values()
		timestamp := result.Record().Time()

		meteoData = models.FromValues(values, timestamp)
		break // Get only the first record
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	if meteoData == nil {
		return nil, fmt.Errorf("meteo data with ID %s not found", id)
	}

	return meteoData, nil
}

// GetAll retrieves all meteo data records with pagination
func (r *meteoDataRepository) GetAll(limit int, offset int) ([]*models.MeteoData, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: -30d)
		|> filter(fn: (r) => r["_measurement"] == "meteodata")
		|> sort(columns: ["_time"], desc: true)
		|> limit(n: %d, offset: %d)
	`, r.bucket, limit, offset)

	result, err := r.queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var meteoDataList []*models.MeteoData

	for result.Next() {
		values := result.Record().Values()
		timestamp := result.Record().Time()

		meteoData := models.FromValues(values, timestamp)
		meteoDataList = append(meteoDataList, meteoData)
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	return meteoDataList, nil
}

// Update updates an existing meteo data record
func (r *meteoDataRepository) Update(meteoData *models.MeteoData) error {
	// In InfluxDB, updates are essentially writes with the same measurement and tag
	// Since InfluxDB is a time-series DB, we'll write a new point with the same ID
	// but a newer timestamp

	// Check if the record exists
	_, err := r.GetByID(meteoData.ID)
	if err != nil {
		return fmt.Errorf("cannot update: meteo data with ID %s not found", meteoData.ID)
	}

	// Set the CreatedAt to now to reflect the update time
	meteoData.CreatedAt = time.Now()

	point := meteoData.ToPoint()
	r.writeAPI.WritePoint(point)

	// Flush to ensure data is written
	r.writeAPI.Flush()

	return nil
}

// Delete removes a meteo data record by its ID
func (r *meteoDataRepository) Delete(id string) error {
	// InfluxDB doesn't have a direct delete by ID operation in the Go client
	// We'll need to use the delete API with a predicate

	startTime := time.Now().AddDate(0, 0, -30) // Last 30 days
	endTime := time.Now().Add(24 * time.Hour)  // Future date to ensure current time is covered

	err := r.client.DeleteAPI().DeleteWithName(
		context.Background(),
		r.org,
		r.bucket,
		startTime,
		endTime,
		fmt.Sprintf(`id="%s"`, id),
	)

	if err != nil {
		return fmt.Errorf("failed to delete meteo data with ID %s: %w", id, err)
	}

	return nil
}

// GetByTimeRange retrieves meteo data records within a specific time range
func (r *meteoDataRepository) GetByTimeRange(start, end time.Time) ([]*models.MeteoData, error) {
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: time(v: %d), stop: time(v: %d))
		|> filter(fn: (r) => r["_measurement"] == "meteodata")
		|> sort(columns: ["_time"], desc: false)
	`, r.bucket, start.UnixNano(), end.UnixNano())

	result, err := r.queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var meteoDataList []*models.MeteoData

	for result.Next() {
		values := result.Record().Values()
		timestamp := result.Record().Time()

		meteoData := models.FromValues(values, timestamp)
		meteoDataList = append(meteoDataList, meteoData)
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	return meteoDataList, nil
}
