package config

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	Validate      *validator.Validate
	onceValidator sync.Once
)

func InitValidator() {
	onceValidator.Do(func() {
		o := validator.New()
		Validate = o
	})

}
