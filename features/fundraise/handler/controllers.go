package handler

import "raihpeduli/features/fundraise"

type controller struct {
	service fundraise.Usecase
}

func New(service fundraise.Usecase) fundraise.Handler {
	return &controller{
		service: service,
	}
}