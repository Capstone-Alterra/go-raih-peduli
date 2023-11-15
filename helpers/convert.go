package helpers

import "github.com/mashingan/smapping"

type converter struct{}

func NewConverter() ConverterInterface {
	return &converter{}
}

func (c converter) Convert(target any, value any) error {
	err := smapping.FillStruct(target, smapping.MapFields(value))
	if err != nil {
		return err
	}

	return nil
}
