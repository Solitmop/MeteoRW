// internal/models/district.go
package models

import _ "gorm.io/gorm"

type District struct {
    ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
    Name   string `json:"name" gorm:"size:100;not null"`
    AreaID uint   `json:"area_id" gorm:"not null;index"`
    
    Area         *Area         `json:"area,omitempty" gorm:"foreignKey:AreaID"`
    RailStations []RailStation `json:"rail_stations,omitempty" gorm:"foreignKey:DistrictID"`
}

// DistrictCreateRequest для создания
type DistrictCreateRequest struct {
    Name   string `json:"name" validate:"required,min=2,max=100,notblank"`
    AreaID uint   `json:"area_id" validate:"required,positive"`
}

// DistrictUpdateRequest для обновления
type DistrictUpdateRequest struct {
    Name   *string `json:"name,omitempty" validate:"omitempty,min=2,max=100,notblank"`
    AreaID *uint   `json:"area_id,omitempty" validate:"omitempty,positive"`
}