package models

import (
	"strconv"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// Regular represents regular meteorological measurements
type Regular struct {
	Index      int       `json:"index"`
	Date       time.Time `json:"date"`
	Visibility int       `json:"visibility"`
	BeforeCode int16     `json:"before_code"`
	DuringCode int16     `json:"during_code"`
	WindAvg    int16     `json:"wind_avg"`
	WindMax    int16     `json:"wind_max"`
	Rainfall   float64   `json:"rainfall"`
	TDry       float64   `json:"t_dry"`
	TWet       float64   `json:"t_wet"`
	TMin       float64   `json:"t_min"`
	TMax       float64   `json:"t_max"`
}

// RegularRequest represents the request structure for validation
type RegularRequest struct {
	Index      int     `json:"index" validate:"required,min=1"`
	Date       int64   `json:"date" validate:"required,min=0"` // Unix timestamp
	Visibility int     `json:"visibility" validate:"min=0,max=100"`
	BeforeCode int16   `json:"before_code" validate:"min=0,max=99"`
	DuringCode int16   `json:"during_code" validate:"min=0,max=99"`
	WindAvg    int16   `json:"wind_avg" validate:"min=0,max=100"`
	WindMax    int16   `json:"wind_max" validate:"min=0,max=150"`
	Rainfall   float64 `json:"rainfall" validate:"min=0,max=1000"`
	TDry       float64 `json:"t_dry" validate:"min=-100,max=100"`
	TWet       float64 `json:"t_wet" validate:"min=-100,max=100"`
	TMin       float64 `json:"t_min" validate:"min=-100,max=100"`
	TMax       float64 `json:"t_max" validate:"min=-100,max=100"`
}

// ToModel converts RegularRequest to Regular model
func (r *RegularRequest) ToModel() Regular {
	return Regular{
		Index:      r.Index,
		Date:       time.Unix(r.Date, 0),
		Visibility: r.Visibility,
		BeforeCode: r.BeforeCode,
		DuringCode: r.DuringCode,
		WindAvg:    r.WindAvg,
		WindMax:    r.WindMax,
		Rainfall:   r.Rainfall,
		TDry:       r.TDry,
		TWet:       r.TWet,
		TMin:       r.TMin,
		TMax:       r.TMax,
	}
}

// ToPoint converts Regular to InfluxDB Point
func (r *Regular) ToPoint() *write.Point {
	return write.NewPoint(
		"regular",
		map[string]string{
			"index": strconv.Itoa(r.Index),
		},
		map[string]interface{}{
			"visibility":  r.Visibility,
			"before_code": r.BeforeCode,
			"during_code": r.DuringCode,
			"wind_avg":    r.WindAvg,
			"wind_max":    r.WindMax,
			"rainfall":    r.Rainfall,
			"t_dry":       r.TDry,
			"t_wet":       r.TWet,
			"t_min":       r.TMin,
			"t_max":       r.TMax,
		},
		r.Date,
	)
}

// FromValues creates a Regular from InfluxDB values
func RegularFromValues(tags map[string]string, fields map[string]interface{}, timestamp time.Time) *Regular {
	r := &Regular{
		Date: timestamp,
	}

	if indexStr, ok := tags["index"]; ok {
		if index, err := strconv.Atoi(indexStr); err == nil {
			r.Index = index
		}
	}

	if visibility, ok := fields["visibility"].(float64); ok {
		r.Visibility = int(visibility)
	} else if visibility, ok := fields["visibility"].(int64); ok {
		r.Visibility = int(visibility)
	}

	if beforeCode, ok := fields["before_code"].(float64); ok {
		r.BeforeCode = int16(beforeCode)
	} else if beforeCode, ok := fields["before_code"].(int64); ok {
		r.BeforeCode = int16(beforeCode)
	}

	if duringCode, ok := fields["during_code"].(float64); ok {
		r.DuringCode = int16(duringCode)
	} else if duringCode, ok := fields["during_code"].(int64); ok {
		r.DuringCode = int16(duringCode)
	}

	if windAvg, ok := fields["wind_avg"].(float64); ok {
		r.WindAvg = int16(windAvg)
	} else if windAvg, ok := fields["wind_avg"].(int64); ok {
		r.WindAvg = int16(windAvg)
	}

	if windMax, ok := fields["wind_max"].(float64); ok {
		r.WindMax = int16(windMax)
	} else if windMax, ok := fields["wind_max"].(int64); ok {
		r.WindMax = int16(windMax)
	}

	if rainfall, ok := fields["rainfall"].(float64); ok {
		r.Rainfall = rainfall
	}

	if tDry, ok := fields["t_dry"].(float64); ok {
		r.TDry = tDry
	}

	if tWet, ok := fields["t_wet"].(float64); ok {
		r.TWet = tWet
	}

	if tMin, ok := fields["t_min"].(float64); ok {
		r.TMin = tMin
	}

	if tMax, ok := fields["t_max"].(float64); ok {
		r.TMax = tMax
	}

	return r
}
