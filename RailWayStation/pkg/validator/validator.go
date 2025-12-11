// pkg/validator/validator.go
package validator

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
		RegisterCustomValidations(v)
	} else {
		validate = validator.New()
		RegisterCustomValidations(validate)
	}
}

// CustomValidator структура для валидации
type CustomValidator struct {
	Validator *validator.Validate
}

// ValidateStruct валидирует структуру
func (cv *CustomValidator) ValidateStruct(s interface{}) error {
	return cv.Validator.Struct(s)
}

// New создает новый валидатор
func New() *CustomValidator {
	return &CustomValidator{Validator: validate}
}

// BindAndValidate биндит и валидирует запрос
func BindAndValidate(c *gin.Context, obj interface{}) error {
	// Биндим JSON
	if err := c.ShouldBindJSON(obj); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Валидируем структуру
	return New().ValidateStruct(obj)
}

// ValidationError форматированная ошибка валидации
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
}

// FormatValidationError форматирует ошибки валидации
func FormatValidationError(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   fieldError.Field(),
				Message: getValidationMessage(fieldError),
				Tag:     fieldError.Tag(),
				Value:   fieldError.Param(),
			})
		}
	} else {
		// Если ошибка не валидации, возвращаем общую
		errors = append(errors, ValidationError{
			Message: err.Error(),
		})
	}

	return errors
}

// getValidationMessage возвращает понятное сообщение
func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return fmt.Sprintf("Must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("Cannot exceed %s characters", fe.Param())
	case "latitude":
		return "Latitude must be between -90 and 90"
	case "longitude":
		return "Longitude must be between -180 and 180"
	case "unique":
		return "This value already exists"
	case "alphanum":
		return "Must contain only letters and numbers"
	case "numeric":
		return "Must be a valid number"
	case "gt":
		return fmt.Sprintf("Must be greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be %s or greater", fe.Param())
	case "lt":
		return fmt.Sprintf("Must be less than %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be %s or less", fe.Param())
	default:
		return fmt.Sprintf("Field validation failed: %s", fe.Tag())
	}
}
