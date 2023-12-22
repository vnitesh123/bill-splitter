package utils

import (
	"github.com/go-playground/validator/v10"
)

func GetErrorCode(validationErrors validator.ValidationErrors) string {
	errorMessage := ""

	for _, err := range validationErrors {
		code := ""
		f := err.Field()
		t := err.Tag()
		p := err.Param()
		switch t {
		case "required":
			code = f + " is missing in the request, "
		case "email":
			code = "Email provided is invalid, "
		case "max":
			code = f + " should contain atmost " + p + " characters, "
		}

		errorMessage += code

	}

	n := len(errorMessage)

	return errorMessage[0 : n-2]
}
