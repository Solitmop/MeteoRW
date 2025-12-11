// internal/models/railstation.go
package models

import _ "gorm.io/gorm"

type RailStation struct {
    ID         uint    `json:"id" gorm:"primaryKey;;not null"`
    Name       string  `json:"name" gorm:"size:100;not null"`
    Lat        float64 `json:"lat" gorm:"type:decimal(10,8);not null"`
    Lon        float64 `json:"lon" gorm:"type:decimal(11,8);not null"`
    DistrictID uint    `json:"district_id" gorm:"not null;index"`
    Hash       string  `json:"hash" gorm:"size:12;index"`
    LineID     string  `json:"line_id" gorm:"type:varchar(20);not null;index"`
    
    District *District `json:"district,omitempty" gorm:"foreignKey:DistrictID"`
    Line     *Line     `json:"line,omitempty" gorm:"foreignKey:LineID"`
}

// RailStationCreateRequest для создания
type RailStationCreateRequest struct {
    ID         uint    `json:"id" validate:"required,gt=0"`
    Name       string  `json:"name" validate:"required,min=2,max=100"`
    Lat        float64 `json:"lat" validate:"required,latitude"`
    Lon        float64 `json:"lon" validate:"required,longitude"`
    DistrictID uint    `json:"district_id" validate:"required,min=1"`
    LineID     string  `json:"line_id" validate:"required,min=1,max=20"`
}

// RailStationUpdateRequest для обновления
type RailStationUpdateRequest struct {
    Name       *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
    Lat        *float64 `json:"lat,omitempty" validate:"omitempty,latitude"`
    Lon        *float64 `json:"lon,omitempty" validate:"omitempty,longitude"`
    DistrictID *uint    `json:"district_id,omitempty" validate:"omitempty,min=1"`
    LineID     *string  `json:"line_id,omitempty" validate:"omitempty,min=1,max=20"`
}