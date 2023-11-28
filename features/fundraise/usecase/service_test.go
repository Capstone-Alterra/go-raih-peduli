package usecase

import (
	"raihpeduli/features/fundraise/dtos"
	"raihpeduli/features/fundraise/mocks"
	helperMocks "raihpeduli/helpers/mocks"
	"testing"
)

func TestFindAll(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, validation)

	var fundraises = []dtos.ResFundraise{}

}
