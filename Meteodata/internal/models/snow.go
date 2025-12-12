package models

import (
    "time"
    _ "gorm.io/gorm"
)

// SNOW - снеговые измерения
type SNOW struct {
    Index  int       `gorm:"primaryKey;column:index" json:"index"`
    Date   time.Time `gorm:"primaryKey;column:date" json:"date"`
    Height float64   `gorm:"column:height;type:decimal(5,2)" json:"height"`
    Grade  int16     `gorm:"column:grade" json:"grade"`
}

func (SNOW) TableName() string {
    return "snow"
}

// SNOWRequest - структура для валидации запросов SNOW (Unix timestamp для даты)
type SNOWRequest struct {
    Index  int     `json:"index" validate:"required,min=1"`
    Date   int64   `json:"date" validate:"required,min=0"` // Unix timestamp
    Height float64 `json:"height" validate:"min=0,max=1000"`
    Grade  int16   `json:"grade" validate:"min=0,max=10"`
}

// ToModel преобразует SNOWRequest в SNOW модель
func (s *SNOWRequest) ToModel() SNOW {
    return SNOW{
        Index:  s.Index,
        Date:   time.Unix(s.Date, 0),
        Height: s.Height,
        Grade:  s.Grade,
    }
}

// MeasurementFilter - фильтр для поиска измерений
type MeasurementFilter struct {
    Index    string `form:"index"`
    DateFrom int64  `form:"date_from"` // Unix timestamp
    DateTo   int64  `form:"date_to"`   // Unix timestamp
    Limit    int    `form:"limit" validate:"min=1,max=1000"`
    Offset   int    `form:"offset" validate:"min=0"`
}