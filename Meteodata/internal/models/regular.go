package models

import (
    "time"
    _ "gorm.io/gorm"
)

// Regular - регулярные метеорологические измерения
type Regular struct {
    Index       int       `gorm:"primaryKey;column:index" json:"index"`
    Date        time.Time `gorm:"primaryKey;column:date" json:"date"`
    Visibility  int       `gorm:"column:visibility" json:"visibility"`
    BeforeCode  int16     `gorm:"column:before_code" json:"before_code"`
    DuringCode  int16     `gorm:"column:during_code" json:"during_code"`
    WindAvg     int16     `gorm:"column:wind_avg" json:"wind_avg"`
    WindMax     int16     `gorm:"column:wind_max" json:"wind_max"`
    Rainfall    float64   `gorm:"column:rainfall;type:decimal(5,2)" json:"rainfall"`
    TDry        float64   `gorm:"column:t_dry;type:decimal(5,2)" json:"t_dry"`
    TWet        float64   `gorm:"column:t_wet;type:decimal(5,2)" json:"t_wet"`
    TMin        float64   `gorm:"column:t_min;type:decimal(5,2)" json:"t_min"`
    TMax        float64   `gorm:"column:t_max;type:decimal(5,2)" json:"t_max"`
}

// RegularRequest - структура для валидации запросов Regular (Unix timestamp для даты)
type RegularRequest struct {
    Index       int     `json:"index" validate:"required,min=1"`
    Date        int64   `json:"date" validate:"required,min=0"` // Unix timestamp
    Visibility  int     `json:"visibility" validate:"min=0,max=100"`
    BeforeCode  int16   `json:"before_code" validate:"min=0,max=99"`
    DuringCode  int16   `json:"during_code" validate:"min=0,max=99"`
    WindAvg     int16   `json:"wind_avg" validate:"min=0,max=100"`
    WindMax     int16   `json:"wind_max" validate:"min=0,max=150"`
    Rainfall    float64 `json:"rainfall" validate:"min=0,max=1000"`
    TDry        float64 `json:"t_dry" validate:"min=-100,max=100"`
    TWet        float64 `json:"t_wet" validate:"min=-100,max=100"`
    TMin        float64 `json:"t_min" validate:"min=-100,max=100"`
    TMax        float64 `json:"t_max" validate:"min=-100,max=100"`
}

// ToModel преобразует RegularRequest в Regular модель
func (r *RegularRequest) ToModel() Regular {
    return Regular{
        Index:       r.Index,
        Date:        time.Unix(r.Date, 0),
        Visibility:  r.Visibility,
        BeforeCode:  r.BeforeCode,
        DuringCode:  r.DuringCode,
        WindAvg:     r.WindAvg,
        WindMax:     r.WindMax,
        Rainfall:    r.Rainfall,
        TDry:        r.TDry,
        TWet:        r.TWet,
        TMin:        r.TMin,
        TMax:        r.TMax,
    }
}

