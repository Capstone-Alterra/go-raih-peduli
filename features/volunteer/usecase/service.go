package usecase

import (
	"errors"
	"mime/multipart"
	"raihpeduli/features/volunteer"
	"raihpeduli/features/volunteer/dtos"
	"raihpeduli/helpers"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
)

type service struct {
	model      volunteer.Repository
	validation helpers.ValidationInterface
}

func New(model volunteer.Repository, validation helpers.ValidationInterface) volunteer.Usecase {
	return &service{
		model:      model,
		validation: validation,
	}
}

func (svc *service) FindAll(page, size int, title, skill, city string) ([]dtos.ResVolunteer, int64) {
	var volunteers []dtos.ResVolunteer

	var volunteersEnt []volunteer.VolunteerVacancies

	var totalData int64
	if title != "" {
		volunteersEnt = svc.model.SelectByTitle(page, size, title)
		totalData = svc.model.GetTotalDataByTitle(title)
	} else if skill != "" {
		volunteersEnt = svc.model.SelectBySkill(page, size, skill)
		totalData = svc.model.GetTotalDataBySkill(skill)
	} else if city != "" {
		volunteersEnt = svc.model.SelectByCity(page, size, city)
		totalData = svc.model.GetTotalDataByCity(city)
	} else {
		volunteersEnt = svc.model.Paginate(page, size)
		totalData = svc.model.GetTotalData()
	}

	for _, volunteer := range volunteersEnt {
		var data dtos.ResVolunteer

		if err := smapping.FillStruct(&data, smapping.MapFields(volunteer)); err != nil {
			log.Error(err.Error())
		}

		volunteers = append(volunteers, data)
	}

	return volunteers, totalData
}

func (svc *service) FindByID(volunteerID int) *dtos.ResVolunteer {
	res := dtos.ResVolunteer{}
	volunteer := svc.model.SelectByID(volunteerID)

	if volunteer == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(volunteer))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Modify(volunteerData dtos.InputVolunteer, volunteerID int) bool {
	newVolunteer := volunteer.VolunteerVacancies{}

	err := smapping.FillStruct(&newVolunteer, smapping.MapFields(volunteerData))
	if err != nil {
		log.Error(err)
		return false
	}

	newVolunteer.ID = volunteerID
	rowsAffected := svc.model.Update(newVolunteer)

	if rowsAffected <= 0 {
		log.Error("There is No Volunteer Updated!")
		return false
	}

	return true
}

func (svc *service) Remove(volunteerID int) bool {
	rowsAffected := svc.model.DeleteByID(volunteerID)

	if rowsAffected <= 0 {
		log.Error("There is No Volunteer Deleted!")
		return false
	}

	return true
}

func (svc *service) Create(newVolunteer dtos.InputVolunteer, UserID int, file multipart.File) (*dtos.ResVolunteer, []string, error) {
	if errMap := svc.validation.ValidateRequest(newVolunteer); errMap != nil {
		return nil, errMap, errors.New("validation error")
	}

	volun := volunteer.VolunteerVacancies{}

	var url string = ""

	if file != nil {
		imageURL, err := svc.model.UploadFile(file, "")

		if err != nil {
			logrus.Error(err)
			return nil, nil, err
		}

		url = imageURL
	}

	volun.UserID = UserID
	volun.Photo = url
	err := smapping.FillStruct(&volun, smapping.MapFields(newVolunteer))

	result, err := svc.model.Insert(&volun)

	if err != nil {
		log.Error(err)
		return nil, nil, errors.New("Use Case : failed to create volunteer")
	}

	resVolun := dtos.ResVolunteer{}
	resVolun.Photo = url
	errRes := smapping.FillStruct(&resVolun, smapping.MapFields(result))

	if errRes != nil {
		log.Error(errRes)
		return nil, nil, errors.New("Use Case : failed to mapping response")
	}

	return &resVolun, nil, nil
}

func (svc *service) Register(newApply dtos.ApplyVolunteer, userID int, file multipart.File) (bool, []string) {
	if errMap := svc.validation.ValidateRequest(newApply); errMap != nil {
		return false, errMap
	}

	registrar := volunteer.VolunteerRelations{}

	var url string = ""

	if file != nil {
		imageURL, err := svc.model.UploadFile(file, "")

		if err != nil {
			logrus.Error(err)
			return false, nil
		}

		url = imageURL
	}

	registrar.UserID = userID
	err := smapping.FillStruct(&registrar, smapping.MapFields(newApply))
	if err != nil {
		log.Error(err)
		return false, nil
	}

	registrar.Resume = url

	err = svc.model.Register(&registrar)
	if err != nil {
		return false, nil
	}

	return true, nil
}
