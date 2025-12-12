package models

import (
    "time"
    _ "gorm.io/gorm"
)

// BaseMeasurement базовый интерфейс для всех измерений
type BaseMeasurement interface {
    GetIndex() int
    GetDate() time.Time
    GetTableName() string
}

// BaseRequest интерфейс для запросов
type BaseRequest interface {
    GetIndex() int
    GetDate() int64 // Unix timestamp
    Validate() error
    ToModel() interface{}
}

// BaseFilter фильтр для поиска
type BaseFilter struct {
    Index    string `form:"index"`
    DateFrom int64  `form:"date_from"`
    DateTo   int64  `form:"date_to"`
    Limit    int    `form:"limit" validate:"min=1,max=1000"`
    Offset   int    `form:"offset" validate:"min=0"`
}