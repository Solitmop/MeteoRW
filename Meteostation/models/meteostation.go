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

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Meteostation{})
}