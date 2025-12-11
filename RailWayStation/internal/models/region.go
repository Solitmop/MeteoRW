// internal/models/region.go
package models

import _ "gorm.io/gorm"

type Region struct {
    ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
    Name string `json:"name" gorm:"size:100;not null;uniqueIndex"`
    
    Areas []Area `json:"areas,omitempty" gorm:"foreignKey:RegionID"`
}

// RegionCreateRequest для создания
type RegionCreateRequest struct {
    Name string `json:"name" validate:"required,min=2,max=100,notblank"`
}

// RegionUpdateRequest для обновления
type RegionUpdateRequest struct {
    Name *string `json:"name,omitempty" validate:"omitempty,min=2,max=100,notblank"`
}