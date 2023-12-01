package usecase

import (
	"raihpeduli/features/home"
	"raihpeduli/features/home/dtos"
	"raihpeduli/helpers"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
)

type service struct {
	model      home.Repository
	validation helpers.ValidationInterface
}

func New(model home.Repository, validation helpers.ValidationInterface) home.Usecase {
	return &service{
		model:      model,
		validation: validation,
	}
}

func (svc *service) FindAll(page, size int) dtos.ResGetHome {
	var resGetHome dtos.ResGetHome
	var fundraises []dtos.ResFundraise
	var volunteers []dtos.ResVolunteer
	var newses []dtos.ResNews

	fundraiseEnt := svc.model.PaginateFundraise(page, size)

	for _, fundraise := range fundraiseEnt {
		var data dtos.ResFundraise

		if err := smapping.FillStruct(&data, smapping.MapFields(fundraise)); err != nil {
			log.Error(err.Error())
		}

		fundraises = append(fundraises, data)
	}

	volunteerEnt := svc.model.PaginateVolunteer(page, size)

	for _, volunteer := range volunteerEnt {
		var data dtos.ResVolunteer

		if err := smapping.FillStruct(&data, smapping.MapFields(volunteer)); err != nil {
			log.Error(err.Error())
		}

		volunteers = append(volunteers, data)
	}

	newsEnt := svc.model.PaginateVolunteer(page, size)

	for _, news := range newsEnt {
		var data dtos.ResNews

		if err := smapping.FillStruct(&data, smapping.MapFields(news)); err != nil {
			log.Error(err.Error())
		}

		newses = append(newses, data)
	}
	resGetHome.Fundraise = fundraises
	resGetHome.Volunteer = volunteers
	resGetHome.News = newses
	return resGetHome
}

func (svc *service) FindAllWeb(page, size int) dtos.ResWebGetHome {
	var resWebGetHome dtos.ResWebGetHome
	var fundraises []dtos.ResFundraise
	var volunteers []dtos.ResVolunteer

	fundraiseEnt := svc.model.PaginateFundraise(page, size)

	for _, fundraise := range fundraiseEnt {
		var data dtos.ResFundraise

		if err := smapping.FillStruct(&data, smapping.MapFields(fundraise)); err != nil {
			log.Error(err.Error())
		}

		fundraises = append(fundraises, data)
	}

	volunteerEnt := svc.model.PaginateVolunteer(page, size)

	for _, volunteer := range volunteerEnt {
		var data dtos.ResVolunteer

		if err := smapping.FillStruct(&data, smapping.MapFields(volunteer)); err != nil {
			log.Error(err.Error())
		}

		volunteers = append(volunteers, data)
	}

	totalUser := svc.model.CountUser()
	totalFundraise := svc.model.CountFundraise()
	totalVolunteer := svc.model.CountVolunteer()
	totalNews := svc.model.CountNews()

	resWebGetHome.UserAmount = totalUser
	resWebGetHome.FundraiseAmount = totalFundraise
	resWebGetHome.VolunteerAmount = totalVolunteer
	resWebGetHome.NewsAmount = totalNews
	resWebGetHome.Fundraise = fundraises
	resWebGetHome.Volunteer = volunteers
	return resWebGetHome
}

func (svc *service) FindByID(homeID int) *dtos.ResHome {
	res := dtos.ResHome{}
	home := svc.model.SelectByID(homeID)

	if home == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(home))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newHome dtos.InputHome) *dtos.ResHome {
	home := home.Home{}

	err := smapping.FillStruct(&home, smapping.MapFields(newHome))
	if err != nil {
		log.Error(err)
		return nil
	}

	homeID := svc.model.Insert(home)

	if homeID == -1 {
		return nil
	}

	resHome := dtos.ResHome{}
	errRes := smapping.FillStruct(&resHome, smapping.MapFields(newHome))
	if errRes != nil {
		log.Error(errRes)
		return nil
	}

	return &resHome
}

func (svc *service) Modify(homeData dtos.InputHome, homeID int) bool {
	newHome := home.Home{}

	err := smapping.FillStruct(&newHome, smapping.MapFields(homeData))
	if err != nil {
		log.Error(err)
		return false
	}

	newHome.ID = homeID
	rowsAffected := svc.model.Update(newHome)

	if rowsAffected <= 0 {
		log.Error("There is No Home Updated!")
		return false
	}

	return true
}

func (svc *service) Remove(homeID int) bool {
	rowsAffected := svc.model.DeleteByID(homeID)

	if rowsAffected <= 0 {
		log.Error("There is No Home Deleted!")
		return false
	}

	return true
}
