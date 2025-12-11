package geoservice

import (
	"github.com/aoliveti/geohash"
	"errors"
)

// GeoHashClient использует локальную библиотеку для расчета geohash
type GeoHashClient struct {
	// Precision определяет детализацию геохеша (например, geohash.City)
	Precision geohash.Precision
}

// NewGeoHashClient создает новый локальный клиент
func NewGeoHashClient(precisionStr string) (*GeoHashClient, error) {
	var precision geohash.Precision
	if p, ok := PrecisionMapNames[precisionStr]; ok {
		precision = p
	} else if p, ok := PrecisionMapNumbers[precisionStr]; ok {
		precision = p
	} else {
		return nil, errors.New("invalid precision")
	}
	return &GeoHashClient{
		Precision: precision,
	}, nil
}

// GetGeohash вычисляет geohash по координатам
func (c *GeoHashClient) GetGeohash(lat, lon float64) (string, error) {
	// Прямой вызов функции из библиотеки
	return geohash.Encode(lat, lon, c.Precision)
}

// GenerateNeighbors вычисляет соседние геохеши
func (s *GeoHashClient) GenerateNeighbors(hash string) ([]string, error) {
    return geohash.Neighbors(hash)
}

var PrecisionMapNames = map[string]geohash.Precision{
	"house":    geohash.House,
	"block":    geohash.Block,
	"building": geohash.Building,
	"street":   geohash.Street,
	"city":     geohash.City,
	"region":   geohash.Region,
	"state":    geohash.State,
	"country":  geohash.Country,
	"global":   geohash.Global,
}

var PrecisionMapNumbers = map[string]geohash.Precision{
	"9": geohash.House,
	"8": geohash.Block,
	"7": geohash.Building,
	"6": geohash.Street,
	"5": geohash.City,
	"4": geohash.Region,
	"3": geohash.State,
	"2": geohash.Country,
	"1": geohash.Global,
}
