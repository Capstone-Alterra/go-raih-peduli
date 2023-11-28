package usecase

import (
	"errors"
	"math"
	"mime/multipart"
	"raihpeduli/features/volunteer"
	"raihpeduli/features/volunteer/dtos"
	"raihpeduli/helpers"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
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

func (svc *service) FindAllVacancies(page, size int, searchAndFilter dtos.SearchAndFilter, ownerID int, suffix string) ([]dtos.ResVacancy, int64) {
	var volunteers []dtos.ResVacancy
	var bookmarkIDs map[int]string
	var err error

	if searchAndFilter.MaxParticipant == 0 {
		searchAndFilter.MaxParticipant = math.MaxInt32
	}

	var volunteersEnt []volunteer.VolunteerVacancies

	if suffix == "mobile" {
		volunteersEnt = svc.model.PaginateMobile(page, size, searchAndFilter)
	} else {
		volunteersEnt = svc.model.Paginate(page, size, searchAndFilter)
	}

	if ownerID != 0 {
		bookmarkIDs, err = svc.model.SelectBookmarkedVacancyID(ownerID)
		if err != nil {
			return nil, 0
		}
	}

	for _, volunteer := range volunteersEnt {
		var data dtos.ResVacancy

		data.ID = volunteer.ID
		data.UserID = volunteer.UserID
		data.Title = volunteer.Title
		data.Description = volunteer.Description
		data.SkillsRequired = strings.Split(volunteer.SkillsRequired, ",")
		data.NumberOfVacancies = volunteer.NumberOfVacancies
		data.ApplicationDeadline = volunteer.ApplicationDeadline
		data.ContactEmail = volunteer.ContactEmail
		data.Province = volunteer.Province
		data.City = volunteer.City
		data.SubDistrict = volunteer.SubDistrict
		data.Photo = volunteer.Photo
		data.Status = volunteer.Status
		data.RejectedReason = volunteer.RejectedReason
		data.CreatedAt = volunteer.CreatedAt
		data.UpdatedAt = volunteer.UpdatedAt
		data.DeletedAt = volunteer.DeletedAt

		if bookmarkIDs != nil {
			bookmardID, ok := bookmarkIDs[data.ID]

			if ok {
				data.BookmarkID = &bookmardID
			}
		}

		data.TotalRegistrar = int(svc.model.GetTotalVolunteersByVacancyID(data.ID))

		volunteers = append(volunteers, data)
	}

	var totalData int64 = 0

	if searchAndFilter.Title != "" || searchAndFilter.Skill != "" || searchAndFilter.City != "" || searchAndFilter.MinParticipant != 0 || searchAndFilter.MaxParticipant != math.MaxInt32 {
		if suffix == "mobile" {
			totalData = svc.model.GetTotalDataVacanciesBySearchAndFilterMobile(searchAndFilter)
		} else {
			totalData = svc.model.GetTotalDataVacanciesBySearchAndFilter(searchAndFilter)
		}
	} else {
		if suffix == "mobile" {
			totalData = svc.model.GetTotalDataVacanciesMobile()
		} else {
			totalData = svc.model.GetTotalDataVacancies()
		}
	}

	return volunteers, totalData
}

func (svc *service) FindVacancyByID(vacancyID, ownerID int) *dtos.ResVacancy {
	res := dtos.ResVacancy{}
	vacancy := svc.model.SelectVacancyByID(vacancyID)

	if vacancy == nil {
		return nil
	}

	var bookmarkID string

	if ownerID != 0 {
		bookmarkID = svc.model.SelectBookmarkByVacancyAndOwnerID(vacancyID, ownerID)

		if bookmarkID != "" {
			res.BookmarkID = &bookmarkID
		}
	}

	res.ID = vacancy.ID
	res.UserID = vacancy.UserID
	res.Title = vacancy.Title
	res.Description = vacancy.Description
	res.SkillsRequired = strings.Split(vacancy.SkillsRequired, ",")
	res.NumberOfVacancies = vacancy.NumberOfVacancies
	res.ApplicationDeadline = vacancy.ApplicationDeadline
	res.ContactEmail = vacancy.ContactEmail
	res.Province = vacancy.Province
	res.City = vacancy.City
	res.SubDistrict = vacancy.SubDistrict
	res.Photo = vacancy.Photo
	res.Status = vacancy.Status
	res.RejectedReason = vacancy.RejectedReason
	res.CreatedAt = vacancy.CreatedAt
	res.UpdatedAt = vacancy.UpdatedAt
	res.DeletedAt = vacancy.DeletedAt

	res.TotalRegistrar = int(svc.model.GetTotalVolunteersByVacancyID(res.ID))

	return &res
}

func (svc *service) ModifyVacancy(vacancyData dtos.InputVacancy, file multipart.File, oldData dtos.ResVacancy) (bool, []string) {
	errMap := svc.validation.ValidateRequest(vacancyData)
	if errMap != nil {
		return false, errMap
	}

	var newVacancy volunteer.VolunteerVacancies

	url, err := svc.model.UploadFile(file, oldData.Photo)
	if err != nil {
		return false, nil
	}

	newVacancy.ID = oldData.ID
	newVacancy.UserID = oldData.UserID
	newVacancy.Title = vacancyData.Title
	newVacancy.Description = vacancyData.Description
	newVacancy.SkillsRequired = strings.Join(vacancyData.SkillsRequired, ", ")
	newVacancy.NumberOfVacancies = vacancyData.NumberOfVacancies
	newVacancy.ApplicationDeadline = vacancyData.ApplicationDeadline
	newVacancy.ContactEmail = vacancyData.ContactEmail
	newVacancy.Province = vacancyData.Province
	newVacancy.City = vacancyData.City
	newVacancy.SubDistrict = vacancyData.SubDistrict
	newVacancy.DetailLocation = vacancyData.DetailLocation
	newVacancy.Photo = url

	rowsAffected := svc.model.UpdateVacancy(newVacancy)

	if rowsAffected <= 0 {
		log.Error("There is No Volunteer Updated!")
		return false, nil
	}

	return true, nil
}

func (svc *service) ModifyVacancyStatus(input dtos.StatusVacancies, oldData dtos.ResVacancy) (bool, []string) {
	errMap := svc.validation.ValidateRequest(input)
	if errMap != nil {
		return false, errMap
	}

	var newVacancy volunteer.VolunteerVacancies

	newVacancy.ID = oldData.ID
	newVacancy.Status = input.Status
	if input.Status == "rejected" {
		if input.RejectedReason == "" {
			return false, []string{"rejected_reason field is required when the status is rejected"}
		}
		newVacancy.RejectedReason = input.RejectedReason
	}

	rowsAffected := svc.model.UpdateVacancy(newVacancy)

	if rowsAffected <= 0 {
		log.Error("There is No Volunteer Updated!")
		return false, nil
	}

	return true, nil
}

func (svc *service) UpdateStatusRegistrar(input dtos.StatusRegistrar, registrarID int) (bool, []string) {
	errMap := svc.validation.ValidateRequest(input)
	if errMap != nil {
		return false, errMap
	}

	registrar := svc.model.SelectRegistrarByID(registrarID)

	if registrar == nil {
		return false, nil
	}

	registrar.Status = input.Status
	if input.Status == "rejected" {
		if input.RejectedReason == "" {
			return false, []string{"rejected_reason field is required when the status is rejected"}
		}
		registrar.RejectedReason = input.RejectedReason
	}

	rowsAffected := svc.model.UpdateStatusRegistrar(*registrar)
	if rowsAffected <= 0 {
		log.Error("Update status registrar failed")
		return false, nil
	}

	return true, nil
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

	url, err := svc.model.UploadFile(file, "")
	if err != nil {
		return nil, nil, err
	}

	vacancy.UserID = UserID
	vacancy.Title = newVolunteer.Title
	vacancy.Description = newVolunteer.Description
	vacancy.SkillsRequired = strings.Join(newVolunteer.SkillsRequired, ", ")
	vacancy.NumberOfVacancies = newVolunteer.NumberOfVacancies
	vacancy.ApplicationDeadline = newVolunteer.ApplicationDeadline
	vacancy.ContactEmail = newVolunteer.ContactEmail
	vacancy.Province = newVolunteer.Province
	vacancy.City = newVolunteer.City
	vacancy.SubDistrict = newVolunteer.SubDistrict
	vacancy.DetailLocation = newVolunteer.DetailLocation
	vacancy.Photo = url

	result, err := svc.model.InsertVacancy(&vacancy)
	if err != nil {
		log.Error(err)
		return nil, nil, errors.New("Use Case : failed to create volunteer")
	}

	resVolun := dtos.ResVacancy{}
	resVolun.ID = result.ID
	resVolun.UserID = result.UserID
	resVolun.Title = result.Title
	resVolun.Description = result.Description
	resVolun.SkillsRequired = strings.Split(result.SkillsRequired, ",")
	resVolun.NumberOfVacancies = result.NumberOfVacancies
	resVolun.ApplicationDeadline = result.ApplicationDeadline
	resVolun.ContactEmail = result.ContactEmail
	resVolun.Province = result.Province
	resVolun.City = result.City
	resVolun.SubDistrict = result.SubDistrict
	resVolun.Photo = result.Photo
	resVolun.Status = result.Status
	resVolun.CreatedAt = result.CreatedAt
	resVolun.UpdatedAt = result.UpdatedAt
	resVolun.DeletedAt = result.DeletedAt

	return &resVolun, nil, nil
}

func (svc *service) RegisterVacancy(newApply dtos.ApplyVacancy, userID int) (bool, []string) {
	if errMap := svc.validation.ValidateRequest(newApply); errMap != nil {
		return false, errMap
	}

	registrar := volunteer.VolunteerRelations{}

	url, err := svc.model.UploadFile(newApply.Photo, "")
	if err != nil {
		return false, nil
	}

	registrar.UserID = userID
	registrar.VolunteerID = newApply.VacancyID
	registrar.Skills = strings.Join(newApply.Skills, ", ")
	registrar.Reason = newApply.Reason
	registrar.Resume = newApply.Resume
	registrar.Photo = url

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

func (svc *service) FindDetailVolunteers(vacancyID, volunteerID int) *dtos.ResRegistrantVacancy {
	res := dtos.ResRegistrantVacancy{}
	volunteer := svc.model.SelectVolunteerDetails(vacancyID, volunteerID)

	if volunteer == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(volunteer))
	if err != nil {
		log.Error("Failed mapping into dtos")
		return nil
	}
	return &res

}

func (svc *service) CheckUser(userID int) bool {
	result := svc.model.CheckUser(userID)
	if !result {
		return false
	}

	return true
}

func (svc *service) FindUserInVacancy(vacancyID, userID int) bool {
	result := svc.model.FindUserInVacancy(vacancyID, userID)
	if !result {
		return false
	}

	return true
}
