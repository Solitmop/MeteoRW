// pkg/validator/custom_validators.go
package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidations регистрирует кастомные валидаторы
func RegisterCustomValidations(v *validator.Validate) {
	// Валидатор для широты
	v.RegisterValidation("latitude", func(fl validator.FieldLevel) bool {
		if lat, ok := fl.Field().Interface().(float64); ok {
			return lat >= -90 && lat <= 90
		}
		return false
	})

	// Валидатор для долготы
	v.RegisterValidation("longitude", func(fl validator.FieldLevel) bool {
		if lon, ok := fl.Field().Interface().(float64); ok {
			return lon >= -180 && lon <= 180
		}
		return false
	})

	// Валидатор для геохеша
	v.RegisterValidation("geohash", func(fl validator.FieldLevel) bool {
		if hash, ok := fl.Field().Interface().(string); ok {
			// Geohash содержит только base32 символы (0-9, b-z кроме a, i, l, o)
			if len(hash) == 0 || len(hash) > 12 {
				return false
			}
			matched, _ := regexp.MatchString(`^[0-9bcdefghjkmnpqrstuvwxyz]+$`, hash)
			return matched
		}
		return false
	})

	// Валидатор для ID линии (только буквы, цифры, дефисы)
	v.RegisterValidation("lineid", func(fl validator.FieldLevel) bool {
		if id, ok := fl.Field().Interface().(string); ok {
			matched, _ := regexp.MatchString(`^[A-Za-z0-9\-_]+$`, id)
			return matched && len(id) >= 1 && len(id) <= 20
		}
		return false
	})

	// Валидатор для проверки, что строка не пустая после trim
	v.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		if str, ok := fl.Field().Interface().(string); ok {
			return strings.TrimSpace(str) != ""
		}
		return false
	})

	// Валидатор для положительного числа
	v.RegisterValidation("positive", func(fl validator.FieldLevel) bool {
        field := fl.Field()
        
        switch field.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            return field.Int() > 0
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            return field.Uint() > 0
        case reflect.Float32, reflect.Float64:
            return field.Float() > 0
        default:
            return false
        }
    })
}

// ValidateID проверяет валидность строкового ID
func ValidateID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: must be a number")
	}

	if id == 0 {
		return 0, fmt.Errorf("ID cannot be zero")
	}

	return uint(id), nil
}

// IsDuplicateKeyError проверяет ошибку дублирования ключа
func IsDuplicateKeyError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "23505") ||
		strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "unique constraint")
}

// FormatError форматирует ошибку для ответа
func FormatError(err error) gin.H {
	return gin.H{
		"error":   "Validation failed",
		"details": FormatValidationError(err),
	}
}
