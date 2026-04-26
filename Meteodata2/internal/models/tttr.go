package models

import (
	"strconv"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// TTTR represents temperature measurements
type TTTR struct {
	Index    int       `json:"index"`
	Date     time.Time `json:"date"`
	Quality  int16     `json:"quality"`
	TMin     float64   `json:"t_min"`
	TAvg     float64   `json:"t_avg"`
	TMax     float64   `json:"t_max"`
	Rainfall float64   `json:"rainfall"`
}

// TTTRRequest represents the request structure for validation
type TTTRRequest struct {
	Index    int     `json:"index" validate:"required,min=1"`
	Date     int64   `json:"date" validate:"required,min=0"` // Unix timestamp
	Quality  int16   `json:"quality" validate:"min=0,max=10"`
	TMin     float64 `json:"t_min" validate:"min=-100,max=100"`
	TAvg     float64 `json:"t_avg" validate:"min=-100,max=100"`
	TMax     float64 `json:"t_max" validate:"min=-100,max=100"`
	Rainfall float64 `json:"rainfall" validate:"min=0,max=1000"`
}

// ToModel converts TTTRRequest to TTTR model
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

// ToPoint converts TTTR to InfluxDB Point
func (t *TTTR) ToPoint() *write.Point {
	return write.NewPoint(
		"tttr",
		map[string]string{
			"index": strconv.Itoa(t.Index),
		},
		map[string]interface{}{
			"quality":  t.Quality,
			"t_min":    t.TMin,
			"t_avg":    t.TAvg,
			"t_max":    t.TMax,
			"rainfall": t.Rainfall,
		},
		t.Date,
	)
}

// FromValues creates a TTTR from InfluxDB values
func TTTRFromValues(tags map[string]string, fields map[string]interface{}, timestamp time.Time) *TTTR {
	t := &TTTR{
		Date: timestamp,
	}

	if indexStr, ok := tags["index"]; ok {
		if index, err := strconv.Atoi(indexStr); err == nil {
			t.Index = index
		}
	}

	if quality, ok := fields["quality"].(float64); ok {
		t.Quality = int16(quality)
	} else if quality, ok := fields["quality"].(int64); ok {
		t.Quality = int16(quality)
	}

	if tMin, ok := fields["t_min"].(float64); ok {
		t.TMin = tMin
	}

	if tAvg, ok := fields["t_avg"].(float64); ok {
		t.TAvg = tAvg
	}

	if tMax, ok := fields["t_max"].(float64); ok {
		t.TMax = tMax
	}

	if rainfall, ok := fields["rainfall"].(float64); ok {
		t.Rainfall = rainfall
	}

	return t
}
