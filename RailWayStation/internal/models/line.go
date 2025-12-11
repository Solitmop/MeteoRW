// internal/models/line.go
package models

import _ "gorm.io/gorm"

type Line struct {
    ID   string `json:"id" gorm:"primaryKey;type:varchar(20)"`
    Name string `json:"name" gorm:"size:100;not null"`
    
    RailStations []RailStation `json:"rail_stations,omitempty" gorm:"foreignKey:LineID"`
}

// LineCreateRequest для создания
type LineCreateRequest struct {
    ID   string `json:"id" validate:"required,min=1,max=20,lineid"`
    Name string `json:"name" validate:"required,min=2,max=100,notblank"`
}

// LineUpdateRequest для обновления
type LineUpdateRequest struct {
    Name *string `json:"name,omitempty" validate:"omitempty,min=2,max=100,notblank"`
}