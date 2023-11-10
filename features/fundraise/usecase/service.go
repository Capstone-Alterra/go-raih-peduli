package usecase

import (
	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/dtos"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
)

type service struct {
	model fundraise.Repository
}

func New(model fundraise.Repository) fundraise.Usecase {
	return &service {
		model: model,
	}
}

func (svc *service) FindAll(page, size int) []dtos.ResFundraise {
	var fundraises []dtos.ResFundraise

	entites, err := svc.model.Paginate(page, size)

	if err != nil {
		log.Error(err)
		return nil
	}

	for _, fundraise := range entites {
		var data dtos.ResFundraise

		if err := smapping.FillStruct(&data, smapping.MapFields(fundraise)); err != nil {
			log.Error(err.Error())
		} 
		
		fundraises = append(fundraises, data)
	}

	return fundraises
}

func (svc *service) FindByID(fundraiseID int) *dtos.ResFundraise {
	var res dtos.ResFundraise
	fundraise, err := svc.model.SelectByID(fundraiseID)

	if err != nil {
		log.Error(err)
		return nil
	}
	
	if err := smapping.FillStruct(&res, smapping.MapFields(fundraise)); err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newFundraise dtos.InputFundraise) (*dtos.ResFundraise, error) {
	var fundraise fundraise.Fundraise
	
	if err := smapping.FillStruct(&fundraise, smapping.MapFields(newFundraise)); err != nil {
		log.Error(err)
		return nil, err
	}

	_, err := svc.model.Insert(fundraise)

	if err != nil {
		return nil, err
	}

	var res dtos.ResFundraise
	
	if err := smapping.FillStruct(&res, smapping.MapFields(newFundraise)); err != nil {
		return nil, err
	}

	return &res, nil
}

func (svc *service) Modify(fundraiseData dtos.InputFundraise, fundraiseID int) bool {
	var newFundraise fundraise.Fundraise
	
	if err := smapping.FillStruct(&newFundraise, smapping.MapFields(fundraiseData)); err != nil {
		log.Error(err)
		return false
	}

	newFundraise.ID = fundraiseID
	_, err := svc.model.Update(newFundraise)

	if err != nil {
		log.Error(err)
		return false
	}
	
	return true
}

func (svc *service) Remove(fundraiseID int) bool {
	_, err := svc.model.DeleteByID(fundraiseID)

	if err != nil {
		log.Error(err)
		return false
	}

	return true
}