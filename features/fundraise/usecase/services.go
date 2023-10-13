package usecase

import "raihpeduli/features/fundraise"

type service struct {
	model fundraise.Repository
}

func New(model fundraise.Repository) fundraise.Usecase {
	return &service{
		model: model,
	}
}