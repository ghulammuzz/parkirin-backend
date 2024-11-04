package form

import "github.com/go-playground/validator/v10"

func ValidationErrorResponse(err error) map[string]string {
	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := err.StructField()
		errors[field] = err.Tag()
	}
	return errors
}
