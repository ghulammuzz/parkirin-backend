package config

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func InitValidator() {
	o := validator.New()
	Validate = o
}
