package usecase

import (
	"errors"
	"mime/multipart"
	"raihpeduli/features/volunteer"
	"raihpeduli/features/volunteer/dtos"
	"raihpeduli/features/volunteer/mocks"
	helperMocks "raihpeduli/helpers/mocks"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindAllVacancies(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	service := New(repo, validation)

	var bookMarkIDs = map[int]string{
		1 : "123456",
	}
	var Bookmark = bookMarkIDs[1]

	vacancies := []dtos.ResVacancy{
		{
			ID: 1,
			UserID: 1,
			Title: "bencana alam",
			Description: "terjadi tsunami di suatu pantai",
			SkillsRequired: []string{"berenang"},
			NumberOfVacancies: 0,
			ContactEmail: "081221278393",
			Province: "Jawa Barat",
			City: "Banten",
			SubDistrict: "Rangkasbitung",
			DetailLocation: "suatu daerah di banten",
			Photo: "https://exampleurl.com/example",
			Status: "pending",
			TotalRegistrar: 21,
			BookmarkID: &Bookmark,
		},
		{
			ID: 2,
			UserID: 1,
			Title: "bencana berbahaya",
			Description: "terjadi tsunami di suatu pantai",
			SkillsRequired: []string{"berenang"},
			NumberOfVacancies: 0,
			ContactEmail: "081221278393",
			Province: "Jawa Barat",
			City: "Banten",
			SubDistrict: "Rangkasbitung",
			DetailLocation: "suatu daerah di banten",
			Photo: "https://exampleurl.com/example",
			Status: "pending",
			TotalRegistrar: 21,
			BookmarkID: nil,
		},
	}

	var pagination = dtos.Pagination{
		Page: 1,
		PageSize: 10,
	}

	var searchAndFilter = dtos.SearchAndFilter{
		Title: "example title",
		City: "example city",
		Skill: "example skill",
		MinParticipant: 1,
		MaxParticipant: 10,
	}

	t.Run("Succes (mobile)", func(t *testing.T){
		repo.On("PaginateMobile", pagination.Page, pagination.PageSize, searchAndFilter).Once()
		for i := 1; i <= len(vacancies); i++ {
			repo.On("SelectBookmarkedVacancyID", i).Return(bookMarkIDs[i], nil).Once()
			repo.On("GetTotalVolunteersByVacancyID", i).Return(int64(5)).Once()
		}
		repo.On("GetTotalDataVacanciesMobile").Return(int64(len(vacancies))).Once()

		result, totalData := service.FindAllVacancies(pagination.Page, pagination.PageSize, searchAndFilter, 1, "mobile")
		assert.Equal(t, len(vacancies), len(result))
		assert.Equal(t, 2, totalData)
	})
}

func TestCreateVacancy(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	service := New(repo, validation)

	var file multipart.File

	var newVacancy = dtos.InputVacancy{
		Title:               "ini bencana alam yang sedang terjadi ada di suatu wilayah",
		Description:         "telah terjadi bencana alam disuatu lokasi telah terjadi bencana alam disuatu lokasi telah terjadi bencana alam disuatu lokasi telah terjadi bencana alam disuatu lokasi telah terjadi bencana alam disuatu lokasi",
		SkillsRequired:      []string{"skil1", "skill2"},
		NumberOfVacancies:   20,
		ApplicationDeadline: time.Date(2024, time.April, 19, 15, 30, 0, 0, time.UTC),
		ContactEmail:        "didadejan45@gmail.com",
		Province:            "jawa barat",
		City:                "banten",
		SubDistrict:         "kec banten",
		DetailLocation:      "di suatu pantai daerah banten",
		Status:              "pending",
		Photo:               file,
	}
	var invalidData = dtos.InputVacancy{
		Title: "hahahahah",
	}

	var vacancyData = volunteer.VolunteerVacancies{
		ID:                  1,
		UserID:              1,
		Title:               "ini bencana alam yang sedang terjadi ada di suatu wilayah",
		Description:         "telah terjadi bencana alam disuatu lokasi telah terjadi bencana alam disuatu lokasi telah terjadi bencana alam disuatu lokasi telah terjadi bencana alam disuatu lokasi telah terjadi bencana alam disuatu lokasi",
		SkillsRequired:      "skil1, skill2",
		NumberOfVacancies:   20,
		ApplicationDeadline: time.Date(2024, time.April, 19, 15, 30, 0, 0, time.UTC),
		ContactEmail:        "didadejan45@gmail.com",
		Province:            "jawa barat",
		City:                "banten",
		SubDistrict:         "kec banten",
		DetailLocation:      "di suatu pantai daerah banten",
		Status:              "pending",
		Photo:               "https://storage.googleapis.com//vacancies/volunteer-vacancy.jpg",
	}

	var errValidation = []string{
		"title required", "description required", "title must be at least 20 characters", "description must be at least 50 characters", "skillsRequired must be at least 1 word", "numberOfVacancies must be greater than 1", "applicationDeadline must be greater than today", "title must be at least 20 characters", "description must be at least 50 characters", "skillsRequired must be at least 1 word", "numberOfVacancies must be greater than 1", "applicationDeadline must be greater than today",
	}

	t.Run("Succes Create", func(t *testing.T) {
		repo.On("SelectByTittle", newVacancy.Title).Return(errors.New("data not found")).Once()
		validation.On("ValidateRequest", newVacancy).Return(nil).Once()
		repo.On("UploadFile", file, "").Return("https://storage.googleapis.com//vacancies/volunteer-vacancy.jpg", nil).Once()

		var vacancy volunteer.VolunteerVacancies
		vacancy.UserID = 1
		vacancy.Title = newVacancy.Title
		vacancy.Description = newVacancy.Description
		vacancy.SkillsRequired = strings.Join(newVacancy.SkillsRequired, ", ")
		vacancy.NumberOfVacancies = newVacancy.NumberOfVacancies
		vacancy.ApplicationDeadline = newVacancy.ApplicationDeadline
		vacancy.ContactEmail = newVacancy.ContactEmail
		vacancy.Province = newVacancy.Province
		vacancy.City = newVacancy.City
		vacancy.SubDistrict = newVacancy.SubDistrict
		vacancy.DetailLocation = newVacancy.DetailLocation
		vacancy.Photo = "https://storage.googleapis.com//vacancies/volunteer-vacancy.jpg"
		repo.On("InsertVacancy", &vacancy).Return(&vacancyData, nil).Once()

		result, errMap, err := service.CreateVacancy(newVacancy, 1, file)
		assert.Nil(t, err)
		assert.Nil(t, errMap)
		assert.NotNil(t, result)
		assert.Equal(t,"ini bencana alam yang sedang terjadi ada di suatu wilayah", result.Title)
		repo.AssertExpectations(t)
		validation.AssertExpectations(t)
	})

	t.Run("Title Already Exists", func(t *testing.T) {
		repo.On("SelectByTittle", newVacancy.Title).Return(nil).Once()

		result, errMap, err := service.CreateVacancy(newVacancy, 1, file)
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		assert.Error(t, err)
		assert.EqualError(t, err, "title already used by another vacancy")
		repo.AssertExpectations(t)
	})

	t.Run("Failed upload file", func(t *testing.T) {
		repo.On("SelectByTittle", newVacancy.Title).Return(errors.New("data not found")).Once()
		validation.On("ValidateRequest", newVacancy).Return(nil).Once()
		repo.On("UploadFile", file, "").Return("", errors.New("failed to upload file")).Once()

		result, errMap, err := service.CreateVacancy(newVacancy, 1, file)

		assert.Error(t, err)
		assert.Nil(t, errMap)
		assert.Nil(t, result)
		assert.Equal(t, "failed to upload file", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Validation error", func(t *testing.T) {
		repo.On("SelectByTittle", invalidData.Title).Return(errors.New("data not found")).Once()
		validation.On("ValidateRequest", invalidData).Return(errValidation).Once()

		result, errMap, err := service.CreateVacancy(invalidData, 1, file)
		assert.Nil(t, err)
		assert.NotNil(t, errMap)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestDeleteVacancy(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	service := New(repo, validation)

	t.Run("Succes deleted vacancy", func(t *testing.T) {
		volunteerID := 1
		oldData := dtos.ResVacancy{
			Photo: "https://storage.googleapis.com/bucket-name/vacancies/photo.jpg",
		}
		repo.On("DeleteFile", "/vacancies/photo.jpg").Return(nil).Once()
		repo.On("DeleteVacancyByID", volunteerID).Return(nil).Once()

		err := service.RemoveVacancy(volunteerID, oldData)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Success removing vacancy without photo", func(t *testing.T) {
		volunteerID := 1
		oldData := dtos.ResVacancy{
			Photo: "default",
		}

		repo.On("DeleteVacancyByID", volunteerID).Return(nil).Once()

		err := service.RemoveVacancy(volunteerID, oldData)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Error deleting vacancy by ID", func(t *testing.T) {
		volunteerID := 1
		oldData := dtos.ResVacancy{
			Photo: "https://storage.googleapis.com/bucket-name/vacancies/photo.jpg",
		}

		repo.On("DeleteFile", "/vacancies/photo.jpg").Return(nil).Once()
		repo.On("DeleteVacancyByID", volunteerID).Return(errors.New("some error")).Once()

		err := service.RemoveVacancy(volunteerID, oldData)

		assert.Error(t, err)
		assert.Equal(t, "some error", err.Error())
		repo.AssertExpectations(t)
	})
}

// func TestFindVacancyByID(t *testing.T) {
// 	repo := mocks.NewRepository(t)
// 	validation := helperMocks.NewValidationInterface(t)
// 	service := New(repo, validation)
// }
