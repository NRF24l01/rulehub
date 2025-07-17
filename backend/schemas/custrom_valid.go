package schemas

import (
	"unicode"

	"github.com/go-playground/validator"
)

func RegisterCustomValidations(v *validator.Validate) error {
	return v.RegisterValidation("strongpwd", validateStrongPassword)
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var hasUpper, hasLower, hasNumber bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasNumber = true
		}
	}

	return hasUpper && hasLower && hasNumber
}
