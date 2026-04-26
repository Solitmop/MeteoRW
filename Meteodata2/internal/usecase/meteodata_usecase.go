package usecase

import (
	"errors"
	"time"

	"meteodata2/internal/models"
	"meteodata2/internal/repository"
)

// MeteoDataUsecase defines the interface for business logic operations
type MeteoDataUsecase interface {
	CreateMeteoData(meteoData *models.MeteoData) error
	GetMeteoDataByID(id string) (*models.MeteoData, error)
	GetAllMeteoData(limit int, offset int) ([]*models.MeteoData, error)
	UpdateMeteoData(meteoData *models.MeteoData) error
	DeleteMeteoData(id string) error
	GetMeteoDataByTimeRange(start, end time.Time) ([]*models.MeteoData, error)
}

// meteoDataUsecase implements MeteoDataUsecase
type meteoDataUsecase struct {
	repo repository.MeteoDataRepository
}

// NewMeteoDataUsecase creates a new instance of MeteoDataUsecase
func NewMeteoDataUsecase(repo repository.MeteoDataRepository) MeteoDataUsecase {
	return &meteoDataUsecase{
		repo: repo,
	}
}

// CreateMeteoData creates a new meteo data record
func (uc *meteoDataUsecase) CreateMeteoData(meteoData *models.MeteoData) error {
	// Validate required fields
	if meteoData.Temperature < -273.15 || meteoData.Temperature > 100 {
		return errors.New("invalid temperature value: must be between -273.15 and 100 degrees Celsius")
	}

	if meteoData.Humidity < 0 || meteoData.Humidity > 100 {
		return errors.New("invalid humidity value: must be between 0 and 100 percent")
	}

	if meteoData.Pressure < 0 {
		return errors.New("invalid pressure value: must be positive")
	}

	if meteoData.WindSpeed < 0 {
		return errors.New("invalid wind speed value: must be non-negative")
	}

	if meteoData.WindDir < 0 || meteoData.WindDir >= 360 {
		return errors.New("invalid wind direction value: must be between 0 and 359 degrees")
	}

	if meteoData.Rainfall < 0 {
		return errors.New("invalid rainfall value: must be non-negative")
	}

	// Set created at time if not already set
	if meteoData.CreatedAt.IsZero() {
		meteoData.CreatedAt = time.Now()
	}

	return uc.repo.Create(meteoData)
}

// GetMeteoDataByID retrieves a meteo data record by its ID
func (uc *meteoDataUsecase) GetMeteoDataByID(id string) (*models.MeteoData, error) {
	if id == "" {
		return nil, errors.New("ID cannot be empty")
	}

	return uc.repo.GetByID(id)
}

// GetAllMeteoData retrieves all meteo data records with pagination
func (uc *meteoDataUsecase) GetAllMeteoData(limit int, offset int) ([]*models.MeteoData, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	if offset < 0 {
		offset = 0 // Default offset
	}

	return uc.repo.GetAll(limit, offset)
}

// UpdateMeteoData updates an existing meteo data record
func (uc *meteoDataUsecase) UpdateMeteoData(meteoData *models.MeteoData) error {
	// Validate required fields
	if meteoData.ID == "" {
		return errors.New("ID cannot be empty for update operation")
	}

	if meteoData.Temperature < -273.15 || meteoData.Temperature > 100 {
		return errors.New("invalid temperature value: must be between -273.15 and 100 degrees Celsius")
	}

	if meteoData.Humidity < 0 || meteoData.Humidity > 100 {
		return errors.New("invalid humidity value: must be between 0 and 100 percent")
	}

	if meteoData.Pressure < 0 {
		return errors.New("invalid pressure value: must be positive")
	}

	if meteoData.WindSpeed < 0 {
		return errors.New("invalid wind speed value: must be non-negative")
	}

	if meteoData.WindDir < 0 || meteoData.WindDir >= 360 {
		return errors.New("invalid wind direction value: must be between 0 and 359 degrees")
	}

	if meteoData.Rainfall < 0 {
		return errors.New("invalid rainfall value: must be non-negative")
	}

	return uc.repo.Update(meteoData)
}

// DeleteMeteoData deletes a meteo data record by its ID
func (uc *meteoDataUsecase) DeleteMeteoData(id string) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	return uc.repo.Delete(id)
}

// GetMeteoDataByTimeRange retrieves meteo data records within a specific time range
func (uc *meteoDataUsecase) GetMeteoDataByTimeRange(start, end time.Time) ([]*models.MeteoData, error) {
	if start.After(end) {
		return nil, errors.New("start time cannot be after end time")
	}

	return uc.repo.GetByTimeRange(start, end)
}
