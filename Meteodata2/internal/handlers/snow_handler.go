package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"meteodata2/internal/models"
	"meteodata2/internal/usecase"
)

// SNOWHandler handles HTTP requests for SNOW measurements
type SNOWHandler struct {
	*BaseHandler
}

// NewSNOWHandler creates a new instance of SNOWHandler
func NewSNOWHandler(usecase usecase.MeteoDataUsecase) *SNOWHandler {
	return &SNOWHandler{
		BaseHandler: NewBaseHandler(usecase),
	}
}

// Create handles POST requests to create new SNOW measurements
func (h *SNOWHandler) Create(c *gin.Context) {
	var req models.SNOWRequest
	if !h.validateRequest(c, &req) {
		return
	}

	// Convert request to model
	measurement := req.ToModel()

	// In InfluxDB, we'll store this as a point with a unique ID based on index and date
	id := "snow_" + strconv.Itoa(measurement.Index) + "_" + strconv.FormatInt(measurement.Date.Unix(), 10)
	measurementWithId := &models.MeteoData{
		ID:          id,
		Temperature: 0,                  // Placeholder - not available in SNOW model
		Humidity:    0,                  // Placeholder - not available in SNOW model
		Pressure:    0,                  // Placeholder - not available in SNOW model
		WindSpeed:   0,                  // Placeholder - not available in SNOW model
		WindDir:     0,                  // Placeholder - not available in SNOW model
		Rainfall:    measurement.Height, // Using height as a representative value
		CreatedAt:   measurement.Date,
	}

	if err := h.BaseHandler.Usecase.CreateMeteoData(measurementWithId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create SNOW measurement",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, measurement)
}

// Get handles GET requests to retrieve SNOW measurements with filtering
func (h *SNOWHandler) Get(c *gin.Context) {
	var filter models.BaseFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate start and end times from the filter
	startTime := time.Now().AddDate(-1, 0, 0) // Default to last year
	endTime := time.Now()

	if filter.DateFrom > 0 {
		startTime = time.Unix(filter.DateFrom, 0)
	}
	if filter.DateTo > 0 {
		endTime = time.Unix(filter.DateTo, 0)
	}

	// Get measurements from usecase
	measurements, err := h.BaseHandler.Usecase.GetMeteoDataByTimeRange(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter measurements based on index if specified
	if filter.Index != "" {
		indices := h.parseIndexList(filter.Index)
		if len(indices) > 0 {
			// Filter measurements by index
			filtered := make([]*models.MeteoData, 0)
			for _, m := range measurements {
				// Extract index from ID (format: snow_index_timestamp)
				if len(m.ID) > 5 && m.ID[:5] == "snow_" {
					for _, idx := range indices {
						idWithoutPrefix := m.ID[5:] // Remove "snow_" prefix
						if len(idWithoutPrefix) >= len(strconv.Itoa(idx)) && idWithoutPrefix[:len(strconv.Itoa(idx))] == strconv.Itoa(idx) {
							filtered = append(filtered, m)
							break
						}
					}
				}
			}
			measurements = filtered
		}
	}

	// Apply pagination
	startIdx := filter.Offset
	if startIdx < 0 {
		startIdx = 0
	}
	endIdx := startIdx + filter.Limit
	if endIdx > len(measurements) {
		endIdx = len(measurements)
	}

	if startIdx >= len(measurements) {
		measurements = []*models.MeteoData{}
	} else {
		measurements = measurements[startIdx:endIdx]
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   measurements,
		"count":  len(measurements),
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetByID handles GET requests to retrieve a specific SNOW measurement by index and date
func (h *SNOWHandler) GetByID(c *gin.Context) {
	indexStr := c.Param("index")
	dateStr := c.Param("date")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid index"})
		return
	}

	_, err = strconv.ParseInt(dateStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected Unix timestamp"})
		return
	}

	// Create ID based on index and date
	id := "snow_" + strconv.Itoa(index) + "_" + dateStr

	// Try to get the measurement
	measurement, err := h.BaseHandler.Usecase.GetMeteoDataByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Measurement not found"})
		return
	}

	// Convert back to SNOW model (this is a simplified conversion)
	snow := &models.SNOW{
		Index:  index,
		Height: measurement.Rainfall, // Using rainfall as height placeholder
	}

	c.JSON(http.StatusOK, snow)
}

// Delete handles DELETE requests to delete SNOW measurements by index and date
func (h *SNOWHandler) Delete(c *gin.Context) {
	indexStr := c.Param("index")
	dateStr := c.Param("date")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid index"})
		return
	}

	_, err = strconv.ParseInt(dateStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected Unix timestamp"})
		return
	}

	// Create ID based on index and date
	id := "snow_" + strconv.Itoa(index) + "_" + dateStr

	if err := h.BaseHandler.Usecase.DeleteMeteoData(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
