package usecase

import (
	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/dtos"
)

type service struct {
	model fundraise.Repository
}

func New(model fundraise.Repository) fundraise.Usecase {
	return &service{
		model: model,
	}
}

func (svc *service) FindAll(page int, size int, title string) []dtos.ResFundraise {
	fundraises, err := svc.model.Paginate(page, size, title)

	if err != nil {
		return nil
	}

	return fundraises
}

func (svc *service) FindByID(fundraiseID int) *dtos.ResFundraise {
	res, err := svc.model.SelectByID(fundraiseID)
	
	if err != nil {
		return nil
	}

	return res
}

func (svc *service) Remove(fundraiseID int) bool {
	_, err := svc.model.DeleteByID(fundraiseID)

	if err != nil {
		return false
	}

	return true
}
