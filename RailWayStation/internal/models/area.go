// internal/models/area.go
package models

import _ "gorm.io/gorm"

type Area struct {
    ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
    Name     string `json:"name" gorm:"size:100;not null"`
    RegionID uint   `json:"region_id" gorm:"not null;index"`
    
    Region    *Region     `json:"region,omitempty" gorm:"foreignKey:RegionID"`
    Districts []District  `json:"districts,omitempty" gorm:"foreignKey:AreaID"`
}

// AreaCreateRequest для создания
type AreaCreateRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100,notblank"`
    RegionID uint   `json:"region_id" validate:"required,positive"`
}

// AreaUpdateRequest для обновления
type AreaUpdateRequest struct {
    Name     *string `json:"name,omitempty" validate:"omitempty,min=2,max=100,notblank"`
    RegionID *uint   `json:"region_id,omitempty" validate:"omitempty,positive"`
}