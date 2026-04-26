package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"meteodata2/internal/models"
	"meteodata2/internal/usecase"
)

// TTTRHandler handles HTTP requests for TTTR measurements
type TTTRHandler struct {
	*BaseHandler
}

// NewTTTRHandler creates a new instance of TTTRHandler
func NewTTTRHandler(usecase usecase.MeteoDataUsecase) *TTTRHandler {
	return &TTTRHandler{
		BaseHandler: NewBaseHandler(usecase),
	}
}

// Create handles POST requests to create new TTTR measurements
func (h *TTTRHandler) Create(c *gin.Context) {
	var req models.TTTRRequest
	if !h.validateRequest(c, &req) {
		return
	}

	// Convert request to model
	measurement := req.ToModel()

	// In InfluxDB, we'll store this as a point with a unique ID based on index and date
	id := "tttr_" + strconv.Itoa(measurement.Index) + "_" + strconv.FormatInt(measurement.Date.Unix(), 10)
	measurementWithId := &models.MeteoData{
		ID:          id,
		Temperature: measurement.TAvg, // Using average temperature as the main value
		Humidity:    0,                // Placeholder - not available in TTTR model
		Pressure:    0,                // Placeholder - not available in TTTR model
		WindSpeed:   0,                // Placeholder - not available in TTTR model
		WindDir:     0,                // Placeholder - not available in TTTR model
		Rainfall:    measurement.Rainfall,
		CreatedAt:   measurement.Date,
	}

	if err := h.BaseHandler.Usecase.CreateMeteoData(measurementWithId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create TTTR measurement",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, measurement)
}

// Get handles GET requests to retrieve TTTR measurements with filtering
func (h *TTTRHandler) Get(c *gin.Context) {
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
				// Extract index from ID (format: tttr_index_timestamp)
				if len(m.ID) > 5 && m.ID[:5] == "tttr_" {
					for _, idx := range indices {
						idWithoutPrefix := m.ID[5:] // Remove "tttr_" prefix
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

// GetByID handles GET requests to retrieve a specific TTTR measurement by index and date
func (h *TTTRHandler) GetByID(c *gin.Context) {
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
	id := "tttr_" + strconv.Itoa(index) + "_" + dateStr

	// Try to get the measurement
	measurement, err := h.BaseHandler.Usecase.GetMeteoDataByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Measurement not found"})
		return
	}

	// Convert back to TTTR model (this is a simplified conversion)
	tttr := &models.TTTR{
		Index:    index,
		TAvg:     measurement.Temperature, // Using temperature as average temp
		Rainfall: measurement.Rainfall,
	}

	c.JSON(http.StatusOK, tttr)
}

// Delete handles DELETE requests to delete TTTR measurements by index and date
func (h *TTTRHandler) Delete(c *gin.Context) {
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
	id := "tttr_" + strconv.Itoa(index) + "_" + dateStr

	if err := h.BaseHandler.Usecase.DeleteMeteoData(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// MonthlyTemperatureRequest represents the request structure for monthly temperature calculation
type MonthlyTemperatureRequest struct {
	Month     int    `form:"month" validate:"required,min=1,max=12" json:"month"`
	StartYear int    `form:"start_year" validate:"required,min=1900,max=2100" json:"start_year"`
	EndYear   int    `form:"end_year" validate:"required,min=1900,max=2100" json:"end_year"`
	Index     string `form:"index"` // Optional: station index (can be multiple separated by commas)
}

// MonthlyTemperatureResponse represents the response structure for monthly temperature
type MonthlyTemperatureResponse struct {
	Month              int     `json:"month"`
	MonthName          string  `json:"month_name"`
	StartYear          int     `json:"start_year"`
	EndYear            int     `json:"end_year"`
	Indexes            []int   `json:"indexes,omitempty"`
	AverageTemperature float64 `json:"average_temperature"`
	MinTemperature     float64 `json:"min_temperature"`
	MaxTemperature     float64 `json:"max_temperature"`
	DaysCount          int64   `json:"days_count"`
	YearsCount         int     `json:"years_count"`
	YearsInRange       []int   `json:"years_in_range"`
}

// GetMonthlyTemperature handles GET requests to calculate average temperature for a month
func (h *TTTRHandler) GetMonthlyTemperature(c *gin.Context) {
	var req MonthlyTemperatureRequest

	// Get parameters from request
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate parameters
	if req.Month < 1 || req.Month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Month must be between 1 and 12"})
		return
	}

	if req.StartYear < 1900 || req.StartYear > 2100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start year must be between 1900 and 2100"})
		return
	}

	if req.EndYear < 1900 || req.EndYear > 2100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End year must be between 1900 and 2100"})
		return
	}

	if req.StartYear > req.EndYear {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start year cannot be greater than end year"})
		return
	}

	// Calculate monthly temperature
	result, err := h.calculateMonthlyTemperature(req.Month, req.StartYear, req.EndYear, req.Index)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to calculate monthly temperature",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// calculateMonthlyTemperature calculates average temperature for a month
func (h *TTTRHandler) calculateMonthlyTemperature(month, startYear, endYear int, indexStr string) (*MonthlyTemperatureResponse, error) {
	// This is a simplified implementation
	// In a real implementation, we would query the database for the specific data

	// For now, we'll return a mock response
	response := &MonthlyTemperatureResponse{
		Month:              month,
		MonthName:          time.Month(month).String(),
		StartYear:          startYear,
		EndYear:            endYear,
		AverageTemperature: 15.0, // Mock value
		MinTemperature:     -5.0, // Mock value
		MaxTemperature:     30.0, // Mock value
		DaysCount:          100,  // Mock value
		YearsCount:         2,    // Mock value
		YearsInRange:       []int{startYear, endYear},
	}

	// Parse indexes if provided
	if indexStr != "" {
		indices := h.parseIndexList(indexStr)
		if len(indices) > 0 {
			response.Indexes = indices
		}
	}

	return response, nil
}
