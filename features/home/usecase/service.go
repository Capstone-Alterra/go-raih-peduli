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
	validation helpers.ValidationInterface
}

func New(model home.Repository, validation helpers.ValidationInterface) home.Usecase {
	return &service{
		model:      model,
		validation: validation,
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

/*
func (svc *service) AddPersonalization(userID int, data dtos.InputPersonalization) error {
	user := svc.model.SelectByID(userID)
	if user == nil {
		return errors.New("user not found")
	}

	user.Personalization = strings.Join(data.Personalization, ", ")
	rowsAffected := svc.model.UpdateUser(*user)
	if rowsAffected == 0 {
		return errors.New("add personalization failed")
	}

	return nil
}
*/

func (svc *service) GetPersonalization(userID int) []string {
	userById := svc.model.SelectUserByID(userID)

	if userById == nil {
		return nil
	}
	var personalization []string
	personalization = strings.Split(*userById.Personalization, ",")
	return personalization
}

func (svc *service) FindByID(homeID int) *dtos.ResHome {
	res := dtos.ResHome{}
	homeService := svc.model.SelectByID(homeID)

	if homeService == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(homeService))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newHome dtos.InputHome) *dtos.ResHome {
	homeStruct := home.Home{}

	err := smapping.FillStruct(&homeStruct, smapping.MapFields(newHome))
	if err != nil {
		log.Error(err)
		return nil
	}

	homeID := svc.model.Insert(homeStruct)

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
