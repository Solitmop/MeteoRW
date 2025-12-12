package models

import (
    "time"
    _ "gorm.io/gorm"
)

// LED - измерения осадков типа LED
type LED struct {
    Index     int       `gorm:"primaryKey;column:index" json:"index"`
    Date      time.Time `gorm:"primaryKey;column:date" json:"date"`
    Indication int16    `gorm:"column:indication" json:"indication"`
    Duration   int16    `gorm:"column:duration" json:"duration"`
    Time       int16    `gorm:"column:time" json:"time"`
    Diameter   int16    `gorm:"column:diameter" json:"diameter"`
    Thickness  int16    `gorm:"column:thickness" json:"thickness"`
}

func (LED) TableName() string {
    return "led"
}

// LEDRequest - структура для валидации запросов LED (Unix timestamp для даты)
type LEDRequest struct {
    Index      int   `json:"index" validate:"required,min=1"`
    Date       int64 `json:"date" validate:"required,min=0"` // Unix timestamp
    Indication int16 `json:"indication" validate:"min=0,max=1000"`
    Duration   int16 `json:"duration" validate:"min=0,max=1440"`
    Time       int16 `json:"time" validate:"min=0,max=1440"`
    Diameter   int16 `json:"diameter" validate:"min=0,max=100"`
    Thickness  int16 `json:"thickness" validate:"min=0,max=100"`
}

// ToModel преобразует LEDRequest в LED модель
func (l *LEDRequest) ToModel() LED {
    return LED{
        Index:      l.Index,
        Date:       time.Unix(l.Date, 0),
        Indication: l.Indication,
        Duration:   l.Duration,
        Time:       l.Time,
        Diameter:   l.Diameter,
        Thickness:  l.Thickness,
    }
}
