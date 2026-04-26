package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"meteodata2/internal/models"
	"meteodata2/internal/usecase"
)

// RegularHandler handles HTTP requests for regular measurements
type RegularHandler struct {
	*BaseHandler
}

// NewRegularHandler creates a new instance of RegularHandler
func NewRegularHandler(usecase usecase.MeteoDataUsecase) *RegularHandler {
	return &RegularHandler{
		BaseHandler: NewBaseHandler(usecase),
	}
}

// Create handles POST requests to create new regular measurements
func (h *RegularHandler) Create(c *gin.Context) {
	var req models.RegularRequest
	if !h.validateRequest(c, &req) {
		return
	}

	// Convert request to model
	measurement := req.ToModel()

	// In InfluxDB, we'll store this as a point with a unique ID based on index and date
	id := strconv.Itoa(measurement.Index) + "_" + strconv.FormatInt(measurement.Date.Unix(), 10)
	measurementWithId := &models.MeteoData{
		ID:          id,
		Temperature: measurement.TDry,             // Using dry temperature as a representative value
		Humidity:    0,                            // Placeholder - not available in Regular model
		Pressure:    0,                            // Placeholder - not available in Regular model
		WindSpeed:   float64(measurement.WindAvg), // Using average wind as representative
		WindDir:     0,                            // Placeholder - not available in Regular model
		Rainfall:    measurement.Rainfall,
		CreatedAt:   measurement.Date,
	}

	if err := h.BaseHandler.Usecase.CreateMeteoData(measurementWithId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create regular measurement",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, measurement)
}

// Get handles GET requests to retrieve regular measurements with filtering
func (h *RegularHandler) Get(c *gin.Context) {
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
				// Extract index from ID (format: index_timestamp)
				if len(m.ID) > 0 {
					for _, idx := range indices {
						if len(m.ID) >= len(strconv.Itoa(idx)) && m.ID[:len(strconv.Itoa(idx))] == strconv.Itoa(idx) {
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

// GetByID handles GET requests to retrieve a specific regular measurement by index and date
func (h *RegularHandler) GetByID(c *gin.Context) {
	indexStr := c.Param("index")
	dateStr := c.Param("date")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid index"})
		return
	}

	dateUnix, err := strconv.ParseInt(dateStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected Unix timestamp"})
		return
	}

	// Create ID based on index and date
	id := strconv.Itoa(index) + "_" + dateStr

	// Try to get the measurement
	measurement, err := h.BaseHandler.Usecase.GetMeteoDataByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Measurement not found"})
		return
	}

	// Convert back to Regular model (this is a simplified conversion)
	regular := &models.Regular{
		Index:    index,
		Date:     time.Unix(dateUnix, 0),
		TDry:     measurement.Temperature,
		WindAvg:  int16(measurement.WindSpeed),
		Rainfall: measurement.Rainfall,
	}

	c.JSON(http.StatusOK, regular)
}

// Update handles PUT requests to update existing regular measurements
func (h *RegularHandler) Update(c *gin.Context) {
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
	id := strconv.Itoa(index) + "_" + dateStr
	_, err = h.BaseHandler.Usecase.GetMeteoDataByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Measurement not found"})
		return
	}

	var req models.RegularRequest
	if !h.validateRequest(c, &req) {
		return
	}

	// Update the measurement
	updated := req.ToModel()
	updated.Index = index

	// Convert to meteodata format for storage
	updatedWithId := &models.MeteoData{
		ID:          id,
		Temperature: updated.TDry,
		Humidity:    0, // Placeholder
		Pressure:    0, // Placeholder
		WindSpeed:   float64(updated.WindAvg),
		WindDir:     0, // Placeholder
		Rainfall:    updated.Rainfall,
		CreatedAt:   time.Now(), // Updated time
	}

	if err := h.BaseHandler.Usecase.UpdateMeteoData(updatedWithId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// Delete handles DELETE requests to delete regular measurements by index and date
func (h *RegularHandler) Delete(c *gin.Context) {
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
	id := strconv.Itoa(index) + "_" + dateStr

	if err := h.BaseHandler.Usecase.DeleteMeteoData(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
