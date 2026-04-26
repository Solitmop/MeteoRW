package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"meteodata2/internal/models"
	"meteodata2/internal/usecase"
)

// BaseHandler contains common methods for all handlers
type BaseHandler struct {
	Usecase usecase.MeteoDataUsecase
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(usecase usecase.MeteoDataUsecase) *BaseHandler {
	return &BaseHandler{Usecase: usecase}
}

// parseIndexList parses a comma-separated list of indices
func (h *BaseHandler) parseIndexList(indexStr string) []int {
	if indexStr == "" {
		return nil
	}

	parts := strings.Split(indexStr, ",")
	indices := make([]int, 0, len(parts))

	for _, part := range parts {
		if idx, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			indices = append(indices, idx)
		}
	}

	return indices
}

// validateRequest validates the request
func (h *BaseHandler) validateRequest(c *gin.Context, request interface{}) bool {
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	// For now, we'll skip validation in this implementation
	// In a real implementation, you would validate using a validation library

	return true
}

// validateQuery validates query parameters
func (h *BaseHandler) validateQuery(c *gin.Context, query interface{}) bool {
	if err := c.ShouldBindQuery(query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	// For now, we'll skip validation in this implementation

	return true
}

// parseTimeFromUnix parses a Unix timestamp from a string
func (h *BaseHandler) parseTimeFromUnix(timeStr string) (time.Time, error) {
	timeInt, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(timeInt, 0), nil
}

// buildBaseQuery builds a base query with filtering
func (h *BaseHandler) buildBaseQuery(filter models.BaseFilter, measurementType string) ([]interface{}, error) {
	// This is a simplified version - in the actual implementation,
	// we would call the usecase layer to query the data
	// For now, we'll return an empty slice
	return []interface{}{}, nil
}

// round rounds a float64 value to the specified precision
func round(value float64, precision int) float64 {
	if precision <= 0 {
		return value
	}

	multiplier := 1.0
	for i := 0; i < precision; i++ {
		multiplier *= 10
	}

	return float64(int(value*multiplier+0.5)) / multiplier
}