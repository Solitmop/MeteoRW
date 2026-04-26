package models

import (
	"strconv"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// LED represents LED precipitation measurements
type LED struct {
	Index      int    `json:"index"`
	Date       time.Time `json:"date"`
	Indication int16  `json:"indication"`
	Duration   int16  `json:"duration"`
	Time       int16  `json:"time"`
	Diameter   int16  `json:"diameter"`
	Thickness  int16  `json:"thickness"`
}

// LEDRequest represents the request structure for validation
type LEDRequest struct {
	Index      int   `json:"index" validate:"required,min=1"`
	Date       int64 `json:"date" validate:"required,min=0"` // Unix timestamp
	Indication int16 `json:"indication" validate:"min=0,max=1000"`
	Duration   int16 `json:"duration" validate:"min=0,max=1440"`
	Time       int16 `json:"time" validate:"min=0,max=1440"`
	Diameter   int16 `json:"diameter" validate:"min=0,max=100"`
	Thickness  int16 `json:"thickness" validate:"min=0,max=100"`
}

// ToModel converts LEDRequest to LED model
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

// ToPoint converts LED to InfluxDB Point
func (l *LED) ToPoint() *write.Point {
	return write.NewPoint(
		"led",
		map[string]string{
			"index": strconv.Itoa(l.Index),
		},
		map[string]interface{}{
			"indication": l.Indication,
			"duration":   l.Duration,
			"time":       l.Time,
			"diameter":   l.Diameter,
			"thickness":  l.Thickness,
		},
		l.Date,
	)
}

// FromValues creates a LED from InfluxDB values
func LEDFromValues(tags map[string]string, fields map[string]interface{}, timestamp time.Time) *LED {
	l := &LED{
		Date: timestamp,
	}

	if indexStr, ok := tags["index"]; ok {
		if index, err := strconv.Atoi(indexStr); err == nil {
			l.Index = index
		}
	}

	if indication, ok := fields["indication"].(float64); ok {
		l.Indication = int16(indication)
	} else if indication, ok := fields["indication"].(int64); ok {
		l.Indication = int16(indication)
	}

	if duration, ok := fields["duration"].(float64); ok {
		l.Duration = int16(duration)
	} else if duration, ok := fields["duration"].(int64); ok {
		l.Duration = int16(duration)
	}

	if timeVal, ok := fields["time"].(float64); ok {
		l.Time = int16(timeVal)
	} else if timeVal, ok := fields["time"].(int64); ok {
		l.Time = int16(timeVal)
	}

	if diameter, ok := fields["diameter"].(float64); ok {
		l.Diameter = int16(diameter)
	} else if diameter, ok := fields["diameter"].(int64); ok {
		l.Diameter = int16(diameter)
	}

	if thickness, ok := fields["thickness"].(float64); ok {
		l.Thickness = int16(thickness)
	} else if thickness, ok := fields["thickness"].(int64); ok {
		l.Thickness = int16(thickness)
	}

	return l
}
