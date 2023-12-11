package usecase

import (
	"raihpeduli/features/fundraise"
	"raihpeduli/features/home"
	"raihpeduli/features/home/dtos"
	"raihpeduli/features/news"
	"raihpeduli/features/volunteer"
	"raihpeduli/helpers"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
)

type service struct {
	model      home.Repository
}

func New(model home.Repository) home.Usecase {
	return &service{
		model:      model,
	}
}

func (svc *service) FindAll(page, size int, personalization []string) dtos.ResGetHome {
	var resGetHome dtos.ResGetHome
	var fundraises []dtos.ResFundraise
	var volunteers []dtos.ResVolunteer
	var newses []dtos.ResNews
	var fundraiseEnt []fundraise.Fundraise
	var volunteerEnt []volunteer.VolunteerVacancies
	var newsEnt []news.News
	if len(personalization) == 0 {
		fundraiseEnt = svc.model.PaginateFundraise(page, size, "", "")
		volunteerEnt = svc.model.PaginateVolunteer(page, size, "", "")
		newsEnt = svc.model.PaginateNews(page, size, "", "")
	} else {
		fundraiseEnt = svc.model.PaginateFundraise(page, size, helpers.BuildLikeQuery("title", personalization), helpers.BuildNotLikeQuery("title", personalization))
		volunteerEnt = svc.model.PaginateVolunteer(page, size, helpers.BuildLikeQuery("title", personalization), helpers.BuildNotLikeQuery("title", personalization))
		newsEnt = svc.model.PaginateNews(page, size, helpers.BuildLikeQuery("title", personalization), helpers.BuildNotLikeQuery("title", personalization))
	}

	for _, fundraiseItem := range fundraiseEnt {
		var data dtos.ResFundraise

		if err := smapping.FillStruct(&data, smapping.MapFields(fundraiseItem)); err != nil {
			log.Error(err.Error())
		}

		fundraises = append(fundraises, data)
	}

	for _, volunteerItem := range volunteerEnt {
		var data dtos.ResVolunteer

		if err := smapping.FillStruct(&data, smapping.MapFields(volunteerItem)); err != nil {
			log.Error(err.Error())
		}

		volunteers = append(volunteers, data)
	}

	for _, newsItem := range newsEnt {
		var data dtos.ResNews

		if err := smapping.FillStruct(&data, smapping.MapFields(newsItem)); err != nil {
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

	fundraiseEnt := svc.model.PaginateFundraise(page, size, "", "")

	for _, fundraiseItem := range fundraiseEnt {
		var data dtos.ResFundraise

		if err := smapping.FillStruct(&data, smapping.MapFields(fundraiseItem)); err != nil {
			log.Error(err.Error())
		}

		fundraises = append(fundraises, data)
	}

	volunteerEnt := svc.model.PaginateVolunteer(page, size, "", "")

	for _, volunteerItem := range volunteerEnt {
		var data dtos.ResVolunteer

		if err := smapping.FillStruct(&data, smapping.MapFields(volunteerItem)); err != nil {
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

func (svc *service) GetPersonalization(userID int) []string {
	userById := svc.model.SelectUserByID(userID)

	if userById == nil {
		return nil
	}
	var personalization []string
	personalization = strings.Split(*userById.Personalization, ",")
	return personalization
}

