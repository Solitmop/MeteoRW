package models

import (
	"time"

	"gorm.io/gorm"
)

type Meteostation struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Index     uint      `json:"index" gorm:"uniqueIndex"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Latitude  float64   `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	Altitude  int       `json:"altitude" gorm:"type:integer"`
	Geohash   string    `json:"geohash" gorm:"type:varchar(12)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Meteostation) TableName() string {
	return "meteostations"
}
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Meteostation{})
}
