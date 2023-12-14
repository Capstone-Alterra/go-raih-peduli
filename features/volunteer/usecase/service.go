package usecase

import (
	"errors"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"raihpeduli/config"
	"raihpeduli/features/volunteer"
	"raihpeduli/features/volunteer/dtos"
	"raihpeduli/helpers"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type service struct {
	model      volunteer.Repository
	validation helpers.ValidationInterface
	nsRequest  helpers.NotificationInterface
}

func New(model volunteer.Repository, validation helpers.ValidationInterface, nsRequest helpers.NotificationInterface) volunteer.Usecase {
	return &service{
		model:      model,
		validation: validation,
		nsRequest:  nsRequest,
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
		data.SkillsRequired = strings.Split(volunteer.SkillsRequired, ", ")
		data.NumberOfVacancies = volunteer.NumberOfVacancies
		data.ApplicationDeadline = volunteer.ApplicationDeadline
		data.ContactEmail = volunteer.ContactEmail
		data.Province = volunteer.Province
		data.City = volunteer.City
		data.SubDistrict = volunteer.SubDistrict
		data.DetailLocation = volunteer.DetailLocation
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
	res.SkillsRequired = strings.Split(vacancy.SkillsRequired, ", ")
	res.NumberOfVacancies = vacancy.NumberOfVacancies
	res.ApplicationDeadline = vacancy.ApplicationDeadline
	res.ContactEmail = vacancy.ContactEmail
	res.Province = vacancy.Province
	res.City = vacancy.City
	res.SubDistrict = vacancy.SubDistrict
	res.DetailLocation = vacancy.DetailLocation
	res.Photo = vacancy.Photo
	res.Status = vacancy.Status
	res.RejectedReason = vacancy.RejectedReason
	res.CreatedAt = vacancy.CreatedAt
	res.UpdatedAt = vacancy.UpdatedAt
	res.DeletedAt = vacancy.DeletedAt

	res.TotalRegistrar = int(svc.model.GetTotalVolunteersByVacancyID(res.ID))

	return &res
}

func (svc *service) ModifyVacancy(vacancyData dtos.InputVacancy, file multipart.File, oldData dtos.ResVacancy) ([]string, error) {
	errMap := svc.validation.ValidateRequest(vacancyData)
	if errMap != nil {
		return errMap, nil
	}

	var newVacancy volunteer.VolunteerVacancies

	url, err := svc.model.UploadFile(file, oldData.Photo)
	if err != nil {
		return nil, errors.New("upload file failed")
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
		return nil, errors.New("there is no vacancy updated")
	}

	return nil, nil
}

func (svc *service) ModifyVacancyStatus(input dtos.StatusVacancies, oldData dtos.ResVacancy) (bool, []string) {
	errMap := svc.validation.ValidateRequest(input)
	if errMap != nil {
		return false, errMap
	}

	deviceToken := svc.model.GetDeviceToken(oldData.UserID)

	var newVacancy volunteer.VolunteerVacancies

	newVacancy.ID = oldData.ID
	newVacancy.Status = input.Status
	if input.Status == "rejected" {
		if input.RejectedReason == "" {
			return false, []string{"rejected_reason field is required when the status is rejected"}
		}
		newVacancy.RejectedReason = input.RejectedReason
		message := "Terima kasih sudah mengajukan lowongan relawan. Saat ini, kami belum bisa menyetujui permohonan ini karena.\n\nAlasan : " + input.RejectedReason + "\n\nTerima kasih atas partisipasinya"
		err := svc.nsRequest.SendNotifications(deviceToken, "Pengajuan Lowongan Relawan Ditolak", message)
		log.Error("Send Notifications Error: ", err)
	} else {
		message := "Kami ingin memberitahu bahwa pengajuan lowongan relawan Anda untuk " + oldData.Title + " telah diterima! Terima kasih atas langkah inisiatif Anda."
		err := svc.nsRequest.SendNotifications(deviceToken, "Pengajuan Lowongan Relawan Diterima", message)
		log.Error("Send Notifications Error: ", err)
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

	vacancy := svc.model.SelectVacancyByID(registrar.VolunteerID)

	deviceToken := svc.model.GetDeviceToken(registrar.UserID)

	registrar.Status = input.Status
	if input.Status == "rejected" {
		if input.RejectedReason == "" {
			return false, []string{"rejected_reason field is required when the status is rejected"}
		}
		registrar.RejectedReason = input.RejectedReason
		message := "Terima kasih sudah mengajukan diri sebagai relawan. Saat ini, kami belum bisa menyetujui permohonan Anda sebagai relawan di " + vacancy.Title + ".\n\nAlasan : " + input.RejectedReason + "\n\nTerima kasih atas partisipasinya"
		err := svc.nsRequest.SendNotifications(deviceToken, "Pengajuan Relawan Ditolak", message)
		log.Error("Send Notifications Error: ", err)
	} else {
		message := "Kami ingin memberitahu Anda bahwa pengajuan sebagai relawan di " + vacancy.Title + " telah diterima! Terima kasih atas minat dan dedikasi Anda."
		err := svc.nsRequest.SendNotifications(deviceToken, "Pengajuan Relawan Diterima", message)
		log.Error("Send Notifications Error: ", err)
	}

	rowsAffected := svc.model.UpdateStatusRegistrar(*registrar)
	if rowsAffected <= 0 {
		log.Error("Update status registrar failed")
		return false, nil
	}

	return true, nil
}

func (svc *service) RemoveVacancy(volunteerID int, oldData dtos.ResVacancy) error {
	var config = config.LoadCloudStorageConfig()
	var oldFilename string = oldData.Photo
	var urlLength int = len("https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/vacancies/")

	if len(oldFilename) > urlLength {
		oldFilename = oldFilename[urlLength:]
	}

	if oldFilename != "default" {
		svc.model.DeleteFile(oldFilename)
	}

	if err := svc.model.DeleteVacancyByID(volunteerID); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (svc *service) CreateVacancy(newVolunteer dtos.InputVacancy, UserID int, file multipart.File) (*dtos.ResVacancy, []string, error) {
	if err := svc.model.SelectByTittle(newVolunteer.Title); err == nil {
		return nil, nil, errors.New("title already used by another vacancy")
	}

	if errorList, err := svc.ValidateInput(newVolunteer, file); err != nil || len(errorList) > 0 {
		return nil, errorList, err
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

		data.ID = volunteer.ID
		data.Email = volunteer.Email
		data.Fullname = volunteer.Fullname
		data.Address = volunteer.Address
		data.PhoneNumber = volunteer.PhoneNumber
		data.Gender = volunteer.Gender
		data.Nik = volunteer.Nik
		data.Skills = strings.Split(volunteer.Skills, ", ")
		data.Resume = volunteer.Resume
		data.Reason = volunteer.Reason
		data.Photo = volunteer.Photo
		data.Status = volunteer.Status
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

	res.ID = volunteer.ID
	res.Email = volunteer.Email
	res.Fullname = volunteer.Fullname
	res.Address = volunteer.Address
	res.PhoneNumber = volunteer.PhoneNumber
	res.Gender = volunteer.Gender
	res.Nik = volunteer.Nik
	res.Skills = strings.Split(volunteer.Skills, ", ")
	res.Resume = volunteer.Resume
	res.Reason = volunteer.Reason
	res.Photo = volunteer.Photo
	res.Status = volunteer.Status

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

func (svc *service) ValidateInput(input dtos.InputVacancy, file multipart.File) ([]string, error) {
	var errorList []string
	if errMap := svc.validation.ValidateRequest(input); errMap != nil {
		errorList = append(errorList, errMap...)
	}
	if len(input.Title) <= 20 {
		errorList = append(errorList, "title must be at least 20 characters")
	}
	if len(input.Description) <= 50 {
		errorList = append(errorList, "description must be at least 50 characters")
	}

	if len(input.SkillsRequired) < 1 {
		errorList = append(errorList, "skillsRequired must be at least 1 word")
	}
	if input.NumberOfVacancies < 1 {
		errorList = append(errorList, "numberOfVacancies must be greater than 1")
	}
	if input.ApplicationDeadline.Before(time.Now()) {
		errorList = append(errorList, "applicationDeadline must be greater than today")
	}
	if file != nil {
		buffer := make([]byte, 512)

		if _, err := file.Read(buffer); err != nil {
			return nil, err
		}

		contentType := http.DetectContentType(buffer)
		isImage := contentType[:5] == "image"

		if !isImage {
			errorList = append(errorList, "photo file has to be an image (png, jpg, or jpeg)")
		}

		const maxFileSize = 5 * 1024 * 1024
		var fileSize int64

		buffer = make([]byte, 1024)
		for {
			n, err := file.Read(buffer)

			fileSize += int64(n)

			if err == io.EOF {
				break
			}

			if err != nil {
				errorList = append(errorList, "unknown file size")
			}
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}

		if fileSize > maxFileSize {
			errorList = append(errorList, "fize size exceeds the allowed limit (5MB)")
		}
	}
	return errorList, nil
}

func (svc *service) FindAllSkills() ([]dtos.Skill, error) {
	skills, err := svc.model.SelectAllSkills()

	if err != nil {
		return nil, err
	}

	return skills, nil
}
