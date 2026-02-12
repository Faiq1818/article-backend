package pkg

import "github.com/go-playground/validator/v10"

func FormatValidationError(err error) map[string]string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{
			"error": "invalid request",
		}
	}

	errMap := make(map[string]string)

	for _, e := range validationErrors {
		switch e.Tag() {
		case "required":
			errMap[e.Field()] = "field is required"

		case "email":
			errMap[e.Field()] = "invalid email format"

		case "min":
			errMap[e.Field()] = "minimum length is " + e.Param()

		case "max":
			errMap[e.Field()] = "maximum length is " + e.Param()

		default:
			errMap[e.Field()] = "invalid value"
		}
	}

	return errMap
}
