package models

import "gorm.io/gorm"

type Meteostation struct {
    Index     uint    `json:"index" gorm:"primaryKey"`
    Name      string  `json:"name"`
    Latitude  float64 `json:"latitude"`  // широта
	Longitude float64 `json:"longitude"` // долгота
	Altitude  int     `json:"altitude"`  // высота над уровнем моря
	Geohash   string  `json:"geohash"`
}

// MeteostationCreateRequest для создания
type MeteostationCreateRequest struct {
	Index     uint    `json:"index" validate:"required,gt=0"`
	Name      string  `json:"name" validate:"required,min=2,max=100"`
	Latitude  float64 `json:"latitude" validate:"required,latitude"`  // широта
	Longitude float64 `json:"longitude" validate:"required,longitude"` // долгота
	Altitude  int     `json:"altitude" validate:"min=-417,max=8848"`  // высота над уровнем моря
	Geohash   string  `json:"geohash" validate:"omitempty,max=12,alphanum"`  // geohash опционален, если не генерируется вручную
}

// MeteostationUpdateRequest для обновления
type MeteostationUpdateRequest struct {
	Name      *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Latitude  *float64 `json:"latitude,omitempty" validate:"omitempty,latitude"`  // широта
	Longitude *float64 `json:"longitude,omitempty" validate:"omitempty,longitude"` // долгота
	Altitude  *int     `json:"altitude,omitempty" validate:"omitempty,min=-417,max=8848"`  // высота над уровнем моря
	Geohash   *string  `json:"geohash,omitempty" validate:"omitempty,max=12,alphanum"`  // geohash опционален
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Meteostation{})
}