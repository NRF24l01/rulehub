package schemas

import (
	"unicode"

	"github.com/go-playground/validator"
)

func RegisterCustomValidations(v *validator.Validate) error {
	err := v.RegisterValidation("strongpwd", validateStrongPassword)
	if err != nil {
		return err
	}
	err = v.RegisterValidation("validusername", validUsername)
	if err != nil {
		return err
	}
	return nil
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var hasLower, hasNumber bool
	for _, ch := range password {
		switch {
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasNumber = true
		}
	}

	return hasLower && hasNumber
}

func validUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	for _, ch := range username {
		if unicode.IsSpace(ch) || (!unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' && ch != '-') {
			return false
		}
	}
	return true
}
