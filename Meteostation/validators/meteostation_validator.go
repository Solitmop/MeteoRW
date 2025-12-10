package validators

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode/utf8"
)

// Магические константы
const (
	MaxNameLength    = 100  // максимальная длина имени
	MinNameLength    = 2    // минимальная длина имени
	MaxLatitude      = 90   // максимальная широта (северная)
	MinLatitude      = -90  // минимальная широта (южная)
	MaxLongitude     = 180  // максимальная долгота (восточная)
	MinLongitude     = -180 // минимальная долгота (западная)
	MaxAltitude      = 8848 // максимальная высота над уровнем моря
	MinAltitude      = -417 // минимальная высота ниже уровня моря
	MaxGeohashLength = 12
)

// MeteostationCreateRequest представляет структуру для создания продукта
type MeteostationCreateRequest struct {
	Index     uint    `json:"index" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`  // широта
	Longitude float64 `json:"longitude" binding:"required"` // долгота
	Altitude  int     `json:"altitude" binding:"required"`  // высота над уровнем моря
	Geohash   string  `json:"geohash"`
}

// MeteostationUpdateRequest представляет структуру для обновления продукта
type MeteostationUpdateRequest struct {
	Index     uint    `json:"index"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`  // широта
	Longitude float64 `json:"longitude"` // долгота
	Altitude  int     `json:"altitude"`  // высота над уровнем моря
	Geohash   string  `json:"geohash"`
}

// ValidateMeteostationIndex проверяет валидность индекса продукта
func ValidateMeteostationIndex(indexStr string) (uint, error) {
	index, err := strconv.ParseUint(indexStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid index")
	}

	return uint(index), nil
}

// ValidateMeteostationName проверяет валидность имени продукта
func ValidateMeteostationName(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("empty name")
	}

	if utf8.RuneCountInString(name) > MaxNameLength {
		return "", fmt.Errorf("name cannot be longer than %d characters", MaxNameLength)
	}

	if utf8.RuneCountInString(name) < MinNameLength {
		return "", fmt.Errorf("name must be at least %d characters long", MinNameLength)
	}

	// Проверяем, что имя состоит из допустимых символов
	matched, _ := regexp.MatchString(`^[a-zA-Zа-яА-Я0-9\s\.,!?()-]+$`, name)
	if !matched {
		return "", fmt.Errorf("invalid name format")
	}

	return name, nil
}

func ValidateLatitude(latitudeStr string) (float64, error) {
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid latitude")
	}

	if latitude < MinLatitude || latitude > MaxLatitude {
		return 0, fmt.Errorf("latitude must be between %f and %f", MinLatitude, MaxLatitude)
	}

	return latitude, nil
}

func ValidateLongitude(longitudeStr string) (float64, error) {
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid longitude")
	}

	if longitude < MinLongitude || longitude > MaxLongitude {
		return 0, fmt.Errorf("longitude must be between %f and %f", MinLongitude, MaxLongitude)
	}

	return longitude, nil
}

func ValidateAltitude(altitudeStr string) (int, error) {
	altitude, err := strconv.Atoi(altitudeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid altitude")
	}

	if altitude < MinAltitude || altitude > MaxAltitude {
		return 0, fmt.Errorf("altitude must be between %d and %d", MinAltitude, MaxAltitude)
	}

	return altitude, nil
}

func ValidateGeohash(geohash string) (string, error) {
	if geohash == "" {
		return "", fmt.Errorf("empty geohash")
	}

	if len(geohash) > MaxGeohashLength {
		return "", fmt.Errorf("geohash cannot be longer than %d characters", MaxGeohashLength)
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, geohash)
	if !matched {
		return "", fmt.Errorf("invalid geohash format")
	}

	return geohash, nil
}
