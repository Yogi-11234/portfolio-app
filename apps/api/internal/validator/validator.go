package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// FormatValidationErrors converts validator.ValidationErrors into a clean map of field -> message.
func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	fieldErrors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			fieldName := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				fieldErrors[fieldName] = fmt.Sprintf("Kolom %s wajib diisi", fieldName)
			case "email":
				fieldErrors[fieldName] = "Format email tidak valid"
			case "min":
				fieldErrors[fieldName] = fmt.Sprintf("Panjang minimal %s karakter", e.Param())
			case "max":
				fieldErrors[fieldName] = fmt.Sprintf("Panjang maksimal %s karakter", e.Param())
			case "uuid":
				fieldErrors[fieldName] = "Format UUID tidak valid"
			default:
				fieldErrors[fieldName] = fmt.Sprintf("Nilai %s tidak valid", fieldName)
			}
		}
	} else {
		fieldErrors["general"] = err.Error()
	}

	return fieldErrors
}
