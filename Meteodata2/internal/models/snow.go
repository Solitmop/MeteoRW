package models

import (
	"strconv"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// SNOW represents snow measurements
type SNOW struct {
	Index  int       `json:"index"`
	Date   time.Time `json:"date"`
	Height float64   `json:"height"`
	Grade  int16     `json:"grade"`
}

// SNOWRequest represents the request structure for validation
type SNOWRequest struct {
	Index  int     `json:"index" validate:"required,min=1"`
	Date   int64   `json:"date" validate:"required,min=0"` // Unix timestamp
	Height float64 `json:"height" validate:"min=0,max=1000"`
	Grade  int16   `json:"grade" validate:"min=0,max=10"`
}

// ToModel converts SNOWRequest to SNOW model
func (s *SNOWRequest) ToModel() SNOW {
	return SNOW{
		Index:  s.Index,
		Date:   time.Unix(s.Date, 0),
		Height: s.Height,
		Grade:  s.Grade,
	}
}

// ToPoint converts SNOW to InfluxDB Point
func (s *SNOW) ToPoint() *write.Point {
	return write.NewPoint(
		"snow",
		map[string]string{
			"index": strconv.Itoa(s.Index),
		},
		map[string]interface{}{
			"height": s.Height,
			"grade":  s.Grade,
		},
		s.Date,
	)
}

// FromValues creates a SNOW from InfluxDB values
func SNOWFromValues(tags map[string]string, fields map[string]interface{}, timestamp time.Time) *SNOW {
	s := &SNOW{
		Date: timestamp,
	}

	if indexStr, ok := tags["index"]; ok {
		if index, err := strconv.Atoi(indexStr); err == nil {
			s.Index = index
		}
	}

	if height, ok := fields["height"].(float64); ok {
		s.Height = height
	}

	if grade, ok := fields["grade"].(float64); ok {
		s.Grade = int16(grade)
	} else if grade, ok := fields["grade"].(int64); ok {
		s.Grade = int16(grade)
	}

	return s
}

// BaseFilter represents a filter for searching measurements
type BaseFilter struct {
	Index    string `form:"index"`
	DateFrom int64  `form:"date_from"` // Unix timestamp
	DateTo   int64  `form:"date_to"`   // Unix timestamp
	Limit    int    `form:"limit" validate:"min=1,max=1000"`
	Offset   int    `form:"offset" validate:"min=0"`
}
