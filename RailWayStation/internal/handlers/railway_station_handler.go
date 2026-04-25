package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"railwaystation/internal/service"
)

// RailwayStationHandler handles HTTP requests for railway stations
type RailwayStationHandler struct {
	service *service.RailwayStationService
}

// NewRailwayStationHandler creates a new railway station handler
func NewRailwayStationHandler(stationService *service.RailwayStationService) *RailwayStationHandler {
	return &RailwayStationHandler{service: stationService}
}

// GetAllStations gets all railway stations
func (h *RailwayStationHandler) GetAllStations(c *gin.Context) {
	stations, err := h.service.GetAllRailwayStations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stations)
}

// GetStationByID gets a railway station by ID
func (h *RailwayStationHandler) GetStationByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid station ID"})
		return
	}

	station, err := h.service.GetRailwayStationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, station)
}

// CreateStation creates a new railway station
func (h *RailwayStationHandler) CreateStation(c *gin.Context) {
	var station service.RailwayStation
	if err := c.ShouldBindJSON(&station); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateRailwayStation(&station); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, station)
}

// UpdateStation updates a railway station
func (h *RailwayStationHandler) UpdateStation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid station ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateRailwayStation(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Station updated successfully"})
}

// DeleteStation deletes a railway station
func (h *RailwayStationHandler) DeleteStation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid station ID"})
		return
	}

	if err := h.service.DeleteRailwayStation(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Station deleted successfully"})
}

// SearchStations searches stations by location
func (h *RailwayStationHandler) SearchStations(c *gin.Context) {
	lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
	lon, err2 := strconv.ParseFloat(c.Query("lon"), 64)
	radius, err3 := strconv.ParseFloat(c.Query("radius"), 64)

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coordinates or radius"})
		return
	}

	stations, err := h.service.SearchStationsByLocation(lat, lon, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stations)
}

// GetAllRegions gets all regions
func (h *RailwayStationHandler) GetAllRegions(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Get all regions"})
}

// GetRegionByID gets a region by ID
func (h *RailwayStationHandler) GetRegionByID(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Get region by ID"})
}

// CreateRegion creates a new region
func (h *RailwayStationHandler) CreateRegion(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Create region"})
}

// UpdateRegion updates a region
func (h *RailwayStationHandler) UpdateRegion(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Update region"})
}

// DeleteRegion deletes a region
func (h *RailwayStationHandler) DeleteRegion(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Delete region"})
}

// GetAllAreas gets all areas
func (h *RailwayStationHandler) GetAllAreas(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Get all areas"})
}

// GetAreaByID gets an area by ID
func (h *RailwayStationHandler) GetAreaByID(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Get area by ID"})
}

// CreateArea creates a new area
func (h *RailwayStationHandler) CreateArea(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Create area"})
}

// UpdateArea updates an area
func (h *RailwayStationHandler) UpdateArea(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Update area"})
}

// DeleteArea deletes an area
func (h *RailwayStationHandler) DeleteArea(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Delete area"})
}

// GetAllDistricts gets all districts
func (h *RailwayStationHandler) GetAllDistricts(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Get all districts"})
}

// GetDistrictByID gets a district by ID
func (h *RailwayStationHandler) GetDistrictByID(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Get district by ID"})
}

// CreateDistrict creates a new district
func (h *RailwayStationHandler) CreateDistrict(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Create district"})
}

// UpdateDistrict updates a district
func (h *RailwayStationHandler) UpdateDistrict(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Update district"})
}

// DeleteDistrict deletes a district
func (h *RailwayStationHandler) DeleteDistrict(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Delete district"})
}

// GetAllLines gets all lines
func (h *RailwayStationHandler) GetAllLines(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Get all lines"})
}

// GetLineByID gets a line by ID
func (h *RailwayStationHandler) GetLineByID(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Get line by ID"})
}

// CreateLine creates a new line
func (h *RailwayStationHandler) CreateLine(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Create line"})
}

// UpdateLine updates a line
func (h *RailwayStationHandler) UpdateLine(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Update line"})
}

// DeleteLine deletes a line
func (h *RailwayStationHandler) DeleteLine(c *gin.Context) {
	// Implementation would go here
	c.JSON(http.StatusOK, gin.H{"message": "Delete line"})
}