package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"meteodata2/internal/models"
)

// Mock usecase for testing
type mockMeteoUsecase struct {
	createFunc         func(meteoData *models.MeteoData) error
	getFunc            func(id string) (*models.MeteoData, error)
	getAllFunc         func(limit int, offset int) ([]*models.MeteoData, error)
	updateFunc         func(meteoData *models.MeteoData) error
	deleteFunc         func(id string) error
	getByTimeRangeFunc func(start, end time.Time) ([]*models.MeteoData, error)
}

func (m *mockMeteoUsecase) CreateMeteoData(meteoData *models.MeteoData) error {
	if m.createFunc != nil {
		return m.createFunc(meteoData)
	}
	return nil
}

func (m *mockMeteoUsecase) GetMeteoDataByID(id string) (*models.MeteoData, error) {
	if m.getFunc != nil {
		return m.getFunc(id)
	}
	return nil, nil
}

func (m *mockMeteoUsecase) GetAllMeteoData(limit int, offset int) ([]*models.MeteoData, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(limit, offset)
	}
	return nil, nil
}

func (m *mockMeteoUsecase) UpdateMeteoData(meteoData *models.MeteoData) error {
	if m.updateFunc != nil {
		return m.updateFunc(meteoData)
	}
	return nil
}

func (m *mockMeteoUsecase) DeleteMeteoData(id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func (m *mockMeteoUsecase) GetMeteoDataByTimeRange(start, end time.Time) ([]*models.MeteoData, error) {
	if m.getByTimeRangeFunc != nil {
		return m.getByTimeRangeFunc(start, end)
	}
	return nil, nil
}

func TestCreateMeteoData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := &mockMeteoUsecase{
		createFunc: func(meteoData *models.MeteoData) error {
			return nil
		},
	}

	handler := NewMeteoHandler(mockUsecase)

	router := gin.Default()
	router.POST("/meteodata", handler.CreateMeteoData)

	meteoData := models.MeteoData{
		Temperature: 25.5,
		Humidity:    60.0,
		Pressure:    1013.25,
		WindSpeed:   5.2,
		WindDir:     180.0,
		Rainfall:    0.0,
		CreatedAt:   time.Now(),
	}

	jsonBytes, _ := json.Marshal(meteoData)

	req, _ := http.NewRequest(http.MethodPost, "/meteodata", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetMeteoDataByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expectedData := &models.MeteoData{
		ID:          "test-id",
		Temperature: 25.5,
		Humidity:    60.0,
		Pressure:    1013.25,
		WindSpeed:   5.2,
		WindDir:     180.0,
		Rainfall:    0.0,
		CreatedAt:   time.Now(),
	}

	mockUsecase := &mockMeteoUsecase{
		getFunc: func(id string) (*models.MeteoData, error) {
			return expectedData, nil
		},
	}

	handler := NewMeteoHandler(mockUsecase)

	router := gin.Default()
	router.GET("/meteodata/:id", handler.GetMeteoDataByID)

	req, _ := http.NewRequest(http.MethodGet, "/meteodata/test-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateMeteoData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := &mockMeteoUsecase{
		updateFunc: func(meteoData *models.MeteoData) error {
			return nil
		},
	}

	handler := NewMeteoHandler(mockUsecase)

	router := gin.Default()
	router.PUT("/meteodata/:id", handler.UpdateMeteoData)

	meteoData := models.MeteoData{
		ID:          "test-id",
		Temperature: 26.0,
		Humidity:    65.0,
		Pressure:    1012.50,
		WindSpeed:   6.0,
		WindDir:     190.0,
		Rainfall:    1.0,
		CreatedAt:   time.Now(),
	}

	jsonBytes, _ := json.Marshal(meteoData)

	req, _ := http.NewRequest(http.MethodPut, "/meteodata/test-id", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteMeteoData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := &mockMeteoUsecase{
		deleteFunc: func(id string) error {
			return nil
		},
	}

	handler := NewMeteoHandler(mockUsecase)

	router := gin.Default()
	router.DELETE("/meteodata/:id", handler.DeleteMeteoData)

	req, _ := http.NewRequest(http.MethodDelete, "/meteodata/test-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
