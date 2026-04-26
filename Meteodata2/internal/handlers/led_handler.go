package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"meteodata2/internal/models"
	"meteodata2/internal/usecase"
)

// LEDHandler handles HTTP requests for LED measurements
type LEDHandler struct {
	*BaseHandler
}

// NewLEDHandler creates a new instance of LEDHandler
func NewLEDHandler(usecase usecase.MeteoDataUsecase) *LEDHandler {
	return &LEDHandler{
		BaseHandler: NewBaseHandler(usecase),
	}
}

// Create handles POST requests to create new LED measurements
func (h *LEDHandler) Create(c *gin.Context) {
	var req models.LEDRequest
	if !h.validateRequest(c, &req) {
		return
	}

	// Convert request to model
	measurement := req.ToModel()

	// In InfluxDB, we'll store this as a point with a unique ID based on index and date
	id := "led_" + strconv.Itoa(measurement.Index) + "_" + strconv.FormatInt(measurement.Date.Unix(), 10)
	measurementWithId := &models.MeteoData{
		ID:          id,
		Temperature: 0,                               // Placeholder - not available in LED model
		Humidity:    0,                               // Placeholder - not available in LED model
		Pressure:    0,                               // Placeholder - not available in LED model
		WindSpeed:   0,                               // Placeholder - not available in LED model
		WindDir:     0,                               // Placeholder - not available in LED model
		Rainfall:    float64(measurement.Indication), // Using indication as a representative value
		CreatedAt:   measurement.Date,
	}

	if err := h.BaseHandler.Usecase.CreateMeteoData(measurementWithId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create LED measurement",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, measurement)
}

// Get handles GET requests to retrieve LED measurements with filtering
func (h *LEDHandler) Get(c *gin.Context) {
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
				// Extract index from ID (format: led_index_timestamp)
				if len(m.ID) > 4 && m.ID[:4] == "led_" {
					for _, idx := range indices {
						idWithoutPrefix := m.ID[4:] // Remove "led_" prefix
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

// GetByID handles GET requests to retrieve a specific LED measurement by index and date
func (h *LEDHandler) GetByID(c *gin.Context) {
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
	id := "led_" + strconv.Itoa(index) + "_" + dateStr

	// Try to get the measurement
	measurement, err := h.BaseHandler.Usecase.GetMeteoDataByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Measurement not found"})
		return
	}

	// Convert back to LED model (this is a simplified conversion)
	led := &models.LED{
		Index:      index,
		Indication: int16(measurement.Rainfall), // Using rainfall as indication placeholder
	}

	c.JSON(http.StatusOK, led)
}

// Update handles PUT requests to update existing LED measurements
func (h *LEDHandler) Update(c *gin.Context) {
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

	// Check if the measurement exists
	id := "led_" + strconv.Itoa(index) + "_" + dateStr
	_, err = h.BaseHandler.Usecase.GetMeteoDataByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Measurement not found"})
		return
	}

	var req models.LEDRequest
	if !h.validateRequest(c, &req) {
		return
	}

	// Update the measurement
	updated := req.ToModel()
	updated.Index = index

	// Convert to meteodata format for storage
	updatedWithId := &models.MeteoData{
		ID:          id,
		Temperature: 0,                           // Placeholder
		Humidity:    0,                           // Placeholder
		Pressure:    0,                           // Placeholder
		WindSpeed:   0,                           // Placeholder
		WindDir:     0,                           // Placeholder
		Rainfall:    float64(updated.Indication), // Using indication as rainfall placeholder
		CreatedAt:   time.Now(),                  // Updated time
	}

	if err := h.BaseHandler.Usecase.UpdateMeteoData(updatedWithId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// Delete handles DELETE requests to delete LED measurements by index and date
func (h *LEDHandler) Delete(c *gin.Context) {
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
	id := "led_" + strconv.Itoa(index) + "_" + dateStr

	if err := h.BaseHandler.Usecase.DeleteMeteoData(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
