package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// RailwayStation represents a railway station with geospatial data
type RailwayStation struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Location    string    `json:"location" gorm:"type:geometry(Point, 4326);not null"`
	RegionID    *uint     `json:"region_id" gorm:"type:integer"`
	AreaID      *uint     `json:"area_id" gorm:"type:integer"`
	DistrictID  *uint     `json:"district_id" gorm:"type:integer"`
	LineID      *uint     `json:"line_id" gorm:"type:integer"`
	StationType string    `json:"station_type" gorm:"type:varchar(100)"`
	Active      bool      `json:"active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Region represents a region
type Region struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null;unique"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Area represents an area within a region
type Area struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	RegionID  uint      `json:"region_id" gorm:"type:integer;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// District represents a district within an area
type District struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	AreaID    uint      `json:"area_id" gorm:"type:integer;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Line represents a railway line
type Line struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// RailwayStationWithDetails includes additional details for a railway station
type RailwayStationWithDetails struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Region      *Region   `json:"region,omitempty"`
	Area        *Area     `json:"area,omitempty"`
	District    *District `json:"district,omitempty"`
	Line        *Line     `json:"line,omitempty"`
	StationType string    `json:"station_type"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}