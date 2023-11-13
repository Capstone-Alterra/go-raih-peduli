package helpers

import (
	"github.com/go-playground/validator/v10"
)

func ErrorMapValidation(err error) []string {
	var mapped []string 
	for _, err := range err.(validator.ValidationErrors) {
		mapped = append(mapped, err.Error())
	}

	return mapped
}