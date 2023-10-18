package helpers

import (
	"github.com/go-playground/validator/v10"
)

func ErrorMapValidation(err error) []map[string]string {
	var mapped = []map[string]string{} 
	for _, err := range err.(validator.ValidationErrors) {
		mapped = append(mapped, map[string]string {
			err.Field(): err.ActualTag(), 
		})
	}

	return mapped
}