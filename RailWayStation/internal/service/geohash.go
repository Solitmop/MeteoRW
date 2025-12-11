// internal/service/geohash.go
package service

import (
    "github.com/aoliveti/geohash"
)

type GeoHashService struct {
    precision geohash.Precision
}

func NewGeoHashService(precision geohash.Precision) *GeoHashService {
    return &GeoHashService{precision: precision}
}

func (s *GeoHashService) Generate(lat, lon float64) (string, error) {
    return geohash.Encode(lat, lon, s.precision)
}

func (s *GeoHashService) GenerateNeighbors(hash string) ([]string, error) {
    return geohash.Neighbors(hash)
}

func (s *GeoHashService) Decode(hash string) (float64, float64, error) {
    return geohash.Decode(hash)
}