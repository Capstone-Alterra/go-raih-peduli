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

	fundraisesEnt := svc.model.Paginate(page, size)

	for _, fundraise := range fundraisesEnt {
		var data dtos.ResFundraise

		if err := smapping.FillStruct(&data, smapping.MapFields(fundraise)); err != nil {
			log.Error(err.Error())
		} 
		
		fundraises = append(fundraises, data)
	}

	return fundraises
}

func (svc *service) FindByID(fundraiseID int) *dtos.ResFundraise {
	res := dtos.ResFundraise{}
	fundraise := svc.model.SelectByID(fundraiseID)

	if fundraise == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(fundraise))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newFundraise dtos.InputFundraise) *dtos.ResFundraise {
	fundraise := fundraise.Fundraise{}
	
	err := smapping.FillStruct(&fundraise, smapping.MapFields(newFundraise))
	if err != nil {
		log.Error(err)
		return nil
	}

	fundraiseID := svc.model.Insert(fundraise)

	if fundraiseID == -1 {
		return nil
	}

	resFundraise := dtos.ResFundraise{}
	errRes := smapping.FillStruct(&resFundraise, smapping.MapFields(newFundraise))
	if errRes != nil {
		log.Error(errRes)
		return nil
	}

	return &resFundraise
}

func (svc *service) Modify(fundraiseData dtos.InputFundraise, fundraiseID int) bool {
	newFundraise := fundraise.Fundraise{}

	err := smapping.FillStruct(&newFundraise, smapping.MapFields(fundraiseData))
	if err != nil {
		log.Error(err)
		return false
	}

	newFundraise.ID = fundraiseID
	rowsAffected := svc.model.Update(newFundraise)

	if rowsAffected <= 0 {
		log.Error("There is No Fundraise Updated!")
		return false
	}
	
	return true
}

func (svc *service) Remove(fundraiseID int) bool {
	rowsAffected := svc.model.DeleteByID(fundraiseID)

	if rowsAffected <= 0 {
		log.Error("There is No Fundraise Deleted!")
		return false
	}

	return true
}