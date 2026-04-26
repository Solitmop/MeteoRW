package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"meteodata2/internal/models"
	"meteodata2/internal/usecase"
)

// MeteoHandler handles HTTP requests for meteo data
type MeteoHandler struct {
	usecase usecase.MeteoDataUsecase
}

// NewMeteoHandler creates a new instance of MeteoHandler
func NewMeteoHandler(usecase usecase.MeteoDataUsecase) *MeteoHandler {
	return &MeteoHandler{
		usecase: usecase,
	}
}

// CreateMeteoData handles POST requests to create new meteo data
func (h *MeteoHandler) CreateMeteoData(c *gin.Context) {
	var meteoData models.MeteoData

	if err := c.ShouldBindJSON(&meteoData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	if err := h.usecase.CreateMeteoData(&meteoData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meteo data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Meteo data created successfully", "id": meteoData.ID})
}

// GetMeteoDataByID handles GET requests to retrieve meteo data by ID
func (h *MeteoHandler) GetMeteoDataByID(c *gin.Context) {
	id := c.Param("id")

	meteoData, err := h.usecase.GetMeteoDataByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meteo data not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meteoData)
}

// GetAllMeteoData handles GET requests to retrieve all meteo data with pagination
func (h *MeteoHandler) GetAllMeteoData(c *gin.Context) {
	// Parse query parameters for pagination
	limitParam := c.DefaultQuery("limit", "10")
	offsetParam := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	meteoDataList, err := h.usecase.GetAllMeteoData(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meteo data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": meteoDataList,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
		},
	})
}

// UpdateMeteoData handles PUT requests to update existing meteo data
func (h *MeteoHandler) UpdateMeteoData(c *gin.Context) {
	id := c.Param("id")

	var meteoData models.MeteoData
	if err := c.ShouldBindJSON(&meteoData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	// Set the ID from the URL parameter
	meteoData.ID = id

	if err := h.usecase.UpdateMeteoData(&meteoData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meteo data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meteo data updated successfully"})
}

// DeleteMeteoData handles DELETE requests to delete meteo data by ID
func (h *MeteoHandler) DeleteMeteoData(c *gin.Context) {
	id := c.Param("id")

	if err := h.usecase.DeleteMeteoData(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meteo data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meteo data deleted successfully"})
}

// GetMeteoDataByTimeRange handles GET requests to retrieve meteo data within a time range
func (h *MeteoHandler) GetMeteoDataByTimeRange(c *gin.Context) {
	startStr := c.Query("start")
	endStr := c.Query("end")

	var start, end time.Time
	var err error

	// Parse start time
	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format, use RFC3339 format", "details": err.Error()})
			return
		}
	} else {
		// Default to 24 hours ago if not provided
		start = time.Now().Add(-24 * time.Hour)
	}

	// Parse end time
	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format, use RFC3339 format", "details": err.Error()})
			return
		}
	} else {
		// Default to now if not provided
		end = time.Now()
	}

	meteoDataList, err := h.usecase.GetMeteoDataByTimeRange(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meteo data by time range", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": meteoDataList,
		"time_range": gin.H{
			"start": start.Format(time.RFC3339),
			"end":   end.Format(time.RFC3339),
		},
	})
}
