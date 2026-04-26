package service

import (
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"railwaystation/internal/models"
)

// RailwayStationService handles railway station operations
type RailwayStationService struct {
	db *gorm.DB
}

// NewRailwayStationService creates a new railway station service
func NewRailwayStationService(db *gorm.DB) *RailwayStationService {
	return &RailwayStationService{db: db}
}

// CreateRailwayStation creates a new railway station
func (s *RailwayStationService) CreateRailwayStation(station *models.RailwayStation) error {
	result := s.db.Create(station)
	if result.Error != nil {
		return fmt.Errorf("failed to create railway station: %w", result.Error)
	}
	return nil
}

// GetRailwayStationByID gets a railway station by ID
func (s *RailwayStationService) GetRailwayStationByID(id uint) (*models.RailwayStation, error) {
	var station models.RailwayStation
	result := s.db.First(&station, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get railway station: %w", result.Error)
	}
	return &station, nil
}

// GetAllRailwayStations gets all railway stations
func (s *RailwayStationService) GetAllRailwayStations() ([]*models.RailwayStation, error) {
	var stations []*models.RailwayStation
	result := s.db.Find(&stations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get railway stations: %w", result.Error)
	}
	return stations, nil
}

// UpdateRailwayStation updates a railway station
func (s *RailwayStationService) UpdateRailwayStation(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.RailwayStation{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update railway station: %w", result.Error)
	}
	return nil
}

// DeleteRailwayStation deletes a railway station
func (s *RailwayStationService) DeleteRailwayStation(id uint) error {
	result := s.db.Delete(&models.RailwayStation{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete railway station: %w", result.Error)
	}
	return nil
}

// SearchStationsByLocation searches stations by location (using PostGIS)
func (s *RailwayStationService) SearchStationsByLocation(lat, lon float64, radius float64) ([]*models.RailwayStation, error) {
	var stations []*models.RailwayStation
	query := `
		SELECT * FROM railway_stations 
		WHERE ST_Distance(location, ST_SetSRID(ST_MakePoint(?, ?), 4326)) <= ?
	`
	result := s.db.Raw(query, lon, lat, radius).Scan(&stations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to search stations by location: %w", result.Error)
	}
	return stations, nil
}

// GetStationsByRegion gets stations by region
func (s *RailwayStationService) GetStationsByRegion(regionID uint) ([]*models.RailwayStation, error) {
	var stations []*models.RailwayStation
	result := s.db.Where("region_id = ?", regionID).Find(&stations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get stations by region: %w", result.Error)
	}
	return stations, nil
}

// GetStationsByArea gets stations by area
func (s *RailwayStationService) GetStationsByArea(areaID uint) ([]*models.RailwayStation, error) {
	var stations []*models.RailwayStation
	result := s.db.Where("area_id = ?", areaID).Find(&stations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get stations by area: %w", result.Error)
	}
	return stations, nil
}

// GetStationsByDistrict gets stations by district
func (s *RailwayStationService) GetStationsByDistrict(districtID uint) ([]*models.RailwayStation, error) {
	var stations []*models.RailwayStation
	result := s.db.Where("district_id = ?", districtID).Find(&stations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get stations by district: %w", result.Error)
	}
	return stations, nil
}

// GetStationsByLine gets stations by line
func (s *RailwayStationService) GetStationsByLine(lineID uint) ([]*models.RailwayStation, error) {
	var stations []*models.RailwayStation
	result := s.db.Where("line_id = ?", lineID).Find(&stations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get stations by line: %w", result.Error)
	}
	return stations, nil
}

// GetStationDetails gets station with related details
func (s *RailwayStationService) GetStationDetails(stationID uint) (*models.RailwayStationWithDetails, error) {
	var stationDetails models.RailwayStationWithDetails
	result := s.db.Table("railway_stations rs").
		Select("rs.*, r.name as region_name, a.name as area_name, d.name as district_name, l.name as line_name").
		Joins("LEFT JOIN regions r ON rs.region_id = r.id").
		Joins("LEFT JOIN areas a ON rs.area_id = a.id").
		Joins("LEFT JOIN districts d ON rs.district_id = d.id").
		Joins("LEFT JOIN lines l ON rs.line_id = l.id").
		Where("rs.id = ?", stationID).
		Scan(&stationDetails)
	
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get station details: %w", result.Error)
	}
	
	return &stationDetails, nil
}

// CreateRegion creates a new region
func (s *RailwayStationService) CreateRegion(region *models.Region) error {
	result := s.db.Create(region)
	if result.Error != nil {
		return fmt.Errorf("failed to create region: %w", result.Error)
	}
	return nil
}

// CreateArea creates a new area
func (s *RailwayStationService) CreateArea(area *models.Area) error {
	result := s.db.Create(area)
	if result.Error != nil {
		return fmt.Errorf("failed to create area: %w", result.Error)
	}
	return nil
}

// CreateDistrict creates a new district
func (s *RailwayStationService) CreateDistrict(district *models.District) error {
	result := s.db.Create(district)
	if result.Error != nil {
		return fmt.Errorf("failed to create district: %w", result.Error)
	}
	return nil
}

// CreateLine creates a new line
func (s *RailwayStationService) CreateLine(line *models.Line) error {
	result := s.db.Create(line)
	if result.Error != nil {
		return fmt.Errorf("failed to create line: %w", result.Error)
	}
	return nil
}

// GetRegionByID gets a region by ID
func (s *RailwayStationService) GetRegionByID(id uint) (*models.Region, error) {
	var region models.Region
	result := s.db.First(&region, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get region: %w", result.Error)
	}
	return &region, nil
}

// GetAreaByID gets an area by ID
func (s *RailwayStationService) GetAreaByID(id uint) (*models.Area, error) {
	var area models.Area
	result := s.db.First(&area, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get area: %w", result.Error)
	}
	return &area, nil
}

// GetDistrictByID gets a district by ID
func (s *RailwayStationService) GetDistrictByID(id uint) (*models.District, error) {
	var district models.District
	result := s.db.First(&district, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get district: %w", result.Error)
	}
	return &district, nil
}

// GetLineByID gets a line by ID
func (s *RailwayStationService) GetLineByID(id uint) (*models.Line, error) {
	var line models.Line
	result := s.db.First(&line, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get line: %w", result.Error)
	}
	return &line, nil
}

// UpdateRegion updates a region
func (s *RailwayStationService) UpdateRegion(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.Region{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update region: %w", result.Error)
	}
	return nil
}

// UpdateArea updates an area
func (s *RailwayStationService) UpdateArea(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.Area{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update area: %w", result.Error)
	}
	return nil
}

// UpdateDistrict updates a district
func (s *RailwayStationService) UpdateDistrict(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.District{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update district: %w", result.Error)
	}
	return nil
}

// UpdateLine updates a line
func (s *RailwayStationService) UpdateLine(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.Line{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update line: %w", result.Error)
	}
	return nil
}

// DeleteRegion deletes a region
func (s *RailwayStationService) DeleteRegion(id uint) error {
	result := s.db.Delete(&models.Region{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete region: %w", result.Error)
	}
	return nil
}

// DeleteArea deletes an area
func (s *RailwayStationService) DeleteArea(id uint) error {
	result := s.db.Delete(&models.Area{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete area: %w", result.Error)
	}
	return nil
}

// DeleteDistrict deletes a district
func (s *RailwayStationService) DeleteDistrict(id uint) error {
	result := s.db.Delete(&models.District{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete district: %w", result.Error)
	}
	return nil
}

// DeleteLine deletes a line
func (s *RailwayStationService) DeleteLine(id uint) error {
	result := s.db.Delete(&models.Line{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete line: %w", result.Error)
	}
	return nil
}

// GetStationCount returns the total count of railway stations
func (s *RailwayStationService) GetStationCount() (int64, error) {
	var count int64
	result := s.db.Model(&models.RailwayStation{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to get station count: %w", result.Error)
	}
	return count, nil
}

// GetRegionCount returns the total count of regions
func (s *RailwayStationService) GetRegionCount() (int64, error) {
	var count int64
	result := s.db.Model(&models.Region{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to get region count: %w", result.Error)
	}
	return count, nil
}

// GetAreaCount returns the total count of areas
func (s *RailwayStationService) GetAreaCount() (int64, error) {
	var count int64
	result := s.db.Model(&models.Area{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to get area count: %w", result.Error)
	}
	return count, nil
}

// GetDistrictCount returns the total count of districts
func (s *RailwayStationService) GetDistrictCount() (int64, error) {
	var count int64
	result := s.db.Model(&models.District{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to get district count: %w", result.Error)
	}
	return count, nil
}

// GetLineCount returns the total count of lines
func (s *RailwayStationService) GetLineCount() (int64, error) {
	var count int64
	result := s.db.Model(&models.Line{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to get line count: %w", result.Error)
	}
	return count, nil
}