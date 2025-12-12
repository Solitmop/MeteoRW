package models

import (
    "time"
    _ "gorm.io/gorm"
)

// TTTR - температурные измерения
type TTTR struct {
    Index     int       `gorm:"primaryKey;column:index" json:"index"`
    Date      time.Time `gorm:"primaryKey;column:date" json:"date"`
    Quality   int16     `gorm:"column:quality" json:"quality"`
    TMin      float64   `gorm:"column:t_min;type:decimal(5,2)" json:"t_min"`
    TAvg      float64   `gorm:"column:t_avg;type:decimal(5,2)" json:"t_avg"`
    TMax      float64   `gorm:"column:t_max;type:decimal(5,2)" json:"t_max"`
    Rainfall  float64   `gorm:"column:rainfall;type:decimal(5,2)" json:"rainfall"`
}

func (TTTR) TableName() string {
    return "tttr"
}

// TTTRRequest - структура для валидации запросов TTTR (Unix timestamp для даты)
type TTTRRequest struct {
    Index    int     `json:"index" validate:"required,min=1"`
    Date     int64   `json:"date" validate:"required,min=0"` // Unix timestamp
    Quality  int16   `json:"quality" validate:"min=0,max=10"`
    TMin     float64 `json:"t_min" validate:"min=-100,max=100"`
    TAvg     float64 `json:"t_avg" validate:"min=-100,max=100"`
    TMax     float64 `json:"t_max" validate:"min=-100,max=100"`
    Rainfall float64 `json:"rainfall" validate:"min=0,max=1000"`
}

// ToModel преобразует TTTRRequest в TTTR модель
func (t *TTTRRequest) ToModel() TTTR {
    return TTTR{
        Index:    t.Index,
        Date:     time.Unix(t.Date, 0),
        Quality:  t.Quality,
        TMin:     t.TMin,
        TAvg:     t.TAvg,
        TMax:     t.TMax,
        Rainfall: t.Rainfall,
    }
}
