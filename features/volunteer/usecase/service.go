package usecase

import (
	"errors"
	"math"
	"mime/multipart"
	"raihpeduli/config"
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

func (svc *service) FindAllVacancies(page, size int, searchAndFilter dtos.SearchAndFilter) ([]dtos.ResVacancy, int64) {
	var volunteers []dtos.ResVacancy

	if searchAndFilter.MaxParticipant == 0 {
		searchAndFilter.MaxParticipant = math.MaxInt32
	}

	volunteersEnt := svc.model.Paginate(page, size, searchAndFilter)

	for _, volunteer := range volunteersEnt {
		var data dtos.ResVacancy

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

func (svc *service) FindVacancyByID(volunteerID int) *dtos.ResVacancy {
	res := dtos.ResVacancy{}
	volunteer := svc.model.SelectVacancyByID(volunteerID)

	if volunteer == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(volunteer))
	if err != nil {
		log.Error(err)
		return nil
	}

	res.TotalRegistrar = int(svc.model.GetTotalVolunteersByVacancyID(res.ID))

	return &res
}

func (svc *service) ModifyVacancy(vacancyData dtos.InputVacancy, file multipart.File, oldData dtos.ResVacancy) (bool, []string) {
	errMap := svc.validation.ValidateRequest(vacancyData)
	if errMap != nil {
		return false, errMap
	}

	var newVacancy volunteer.VolunteerVacancies
	var url string = ""
	var config = config.LoadCloudStorageConfig()
	var oldFilename string = oldData.Photo
	var urlLength int = len("https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/vacancies/")

	if file != nil {
		if oldFilename == "https://storage.googleapis.com/raih_peduli/vacancies/volunteer-vacancy.jpg" {
			oldFilename = ""
		} else if len(oldFilename) > urlLength {
			oldFilename = oldFilename[urlLength:]
		}

		imageURL, err := svc.model.UploadFile(file, oldFilename)

		if err != nil {
			logrus.Error(err)
			return false, nil
		}

		url = imageURL
	}

	err := smapping.FillStruct(&newVacancy, smapping.MapFields(vacancyData))
	if err != nil {
		log.Error(err)
		return false, nil
	}

	newVacancy.ID = oldData.ID
	newVacancy.Photo = url
	rowsAffected := svc.model.UpdateVacancy(newVacancy)

	if rowsAffected <= 0 {
		log.Error("There is No Volunteer Updated!")
		return false, nil
	}

	return true, nil
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

func (svc *service) RemoveVacancy(volunteerID int) bool {
	rowsAffected := svc.model.DeleteVacancyByID(volunteerID)

	if rowsAffected <= 0 {
		log.Error("There is No Volunteer Deleted!")
		return false
	}

	return true
}

func (svc *service) CreateVacancy(newVolunteer dtos.InputVacancy, UserID int, file multipart.File) (*dtos.ResVacancy, []string, error) {
	if errMap := svc.validation.ValidateRequest(newVolunteer); errMap != nil {
		return nil, errMap, errors.New("validation error")
	}

	vacancy := volunteer.VolunteerVacancies{}
	var url string = ""

	if file != nil {
		imageURL, err := svc.model.UploadFile(file, "")

		if err != nil {
			logrus.Error(err)
			return nil, nil, err
		}

		url = imageURL
	} else {
		config := config.LoadCloudStorageConfig()
		url = "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/vacancies/volunteer-vacancy.jpg"
	}

	vacancy.UserID = UserID
	vacancy.Photo = url
	err := smapping.FillStruct(&vacancy, smapping.MapFields(newVolunteer))
	if err != nil {
		return nil, nil, err
	}

	result, err := svc.model.InsertVacancy(&vacancy)
	if err != nil {
		log.Error(err)
		return nil, nil, errors.New("Use Case : failed to create volunteer")
	}

	resVolun := dtos.ResVacancy{}
	resVolun.Photo = url
	errRes := smapping.FillStruct(&resVolun, smapping.MapFields(result))

	if errRes != nil {
		log.Error(errRes)
		return nil, nil, errors.New("Use Case : failed to mapping response")
	}

	return &resVolun, nil, nil
}

func (svc *service) RegisterVacancy(newApply dtos.ApplyVacancy, userID int, file multipart.File) (bool, []string) {
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

	err = svc.model.RegisterVacancy(&registrar)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (svc *service) FindAllVolunteersByVacancyID(page, size int, vacancyID int, name string) ([]dtos.ResRegistrantVacancy, int64) {
	var volunteers []dtos.ResRegistrantVacancy

	volunteerEnt := svc.model.SelectVolunteersByVacancyID(vacancyID, name, page, size)
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

	totalData := svc.model.GetTotalVolunteers(vacancyID, name)

	return volunteers, totalData
}
