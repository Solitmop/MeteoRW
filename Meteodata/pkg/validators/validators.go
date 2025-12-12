package validators

import (
    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Init инициализирует валидатор
func Init() {
    validate = validator.New()
}

// Validate проверяет структуру
func Validate(s interface{}) error {
    if validate == nil {
        Init()
    }
    return validate.Struct(s)
}