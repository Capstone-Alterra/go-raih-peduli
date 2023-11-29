package helpers

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type validation struct{}

func NewValidationRequest() ValidationInterface {
	return &validation{}
}

func (v validation) ValidateRequest(request any) []string {
	var validate = validator.New()

	if err := validate.Struct(request); err != nil {
		var errMap = []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errMap = append(errMap, strings.ToLower(err.Error()))
		}

		return errMap
	}

	return nil
}
