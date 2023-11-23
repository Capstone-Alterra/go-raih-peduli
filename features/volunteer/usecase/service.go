package usecase

import (
	"errors"
	"math"
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

func (svc *service) FindAll(page, size int, searchAndFilter dtos.SearchAndFilter) ([]dtos.ResVolunteer, int64) {
	var volunteers []dtos.ResVolunteer

	if searchAndFilter.MaxParticipant == 0 {
		searchAndFilter.MaxParticipant = math.MaxInt32
	}

	volunteersEnt := svc.model.Paginate(page, size, searchAndFilter)

	for _, volunteer := range volunteersEnt {
		var data dtos.ResVolunteer

		if err := smapping.FillStruct(&data, smapping.MapFields(volunteer)); err != nil {
			log.Error(err.Error())
		}

		volunteers = append(volunteers, data)
	}

	var totalData int64 = 0

	if searchAndFilter.Title != "" || searchAndFilter.Skill != "" || searchAndFilter.City != "" || searchAndFilter.MinParticipant != 0 || searchAndFilter.MaxParticipant != math.MaxInt32 {
		totalData = svc.model.GetTotalDataVacanciesBySearchAndFilter(searchAndFilter)
	} else {
		totalData = svc.model.GetTotalDataVacancies()
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

	res.TotalRegistrar = int(svc.model.GetTotalVolunteerByVacancyID(res.ID))

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

func (svc *service) UpdateStatusRegistrar(status string, registrarID int) bool {
	registrar := svc.model.SelectRegistrarByID(registrarID)

	if registrar == nil {
		return false
	}

	registrar.Status = status
	rowsAffected := svc.model.UpdateStatusRegistrar(*registrar)
	if rowsAffected <= 0 {
		log.Error("Update status registrar failed")
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

func (svc *service) FindAllVolunteerByVacancyID(page, size int, vacancyID int, name string) ([]dtos.ResRegistrantVacancy, int64) {
	var volunteers []dtos.ResRegistrantVacancy

	volunteerEnt := svc.model.SelectVolunteerByVacancyID(vacancyID, name, page, size)
	if volunteerEnt == nil {
		return nil, 0
	}

	for _, volunteer := range volunteerEnt {
		var data dtos.ResRegistrantVacancy

		if err := smapping.FillStruct(&data, smapping.MapFields(volunteer)); err != nil {
			log.Error(err.Error())
		}

		volunteers = append(volunteers, data)
	}

	totalData := svc.model.GetTotalVolunteer(vacancyID, name)

	return volunteers, totalData
}
