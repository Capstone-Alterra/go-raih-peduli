package usecase

import (
	"errors"
	"mime/multipart"
	"os"
	"raihpeduli/features/volunteer"
	"raihpeduli/features/volunteer/dtos"
	"raihpeduli/features/volunteer/mocks"
	helperMocks "raihpeduli/helpers/mocks"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFindAllVacancies(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	var bookMarkIDs = map[int]string{
		1: "123456",
	}

	vacancies := []volunteer.VolunteerVacancies{
		{
			ID:                1,
			UserID:            1,
			Title:             "bencana alam gunung meletus",
			Description:       "terjadi tsunami di suatu pantai",
			SkillsRequired:    "berenang",
			NumberOfVacancies: 15,
			ContactEmail:      "081221278393",
			Province:          "Jawa Barat",
			City:              "Banten",
			SubDistrict:       "Rangkasbitung",
			DetailLocation:    "suatu daerah di banten",
			Photo:             "https://exampleurl.com/example",
			Status:            "pending",
		},
		{
			ID:                2,
			UserID:            1,
			Title:             "tsunami di pantai utara",
			Description:       "terjadi tsunami di suatu pantai",
			SkillsRequired:    "berenang",
			NumberOfVacancies: 15,
			ContactEmail:      "081221278393",
			Province:          "Jawa Barat",
			City:              "Banten",
			SubDistrict:       "Rangkasbitung",
			DetailLocation:    "suatu daerah di banten",
			Photo:             "https://exampleurl.com/example",
			Status:            "pending",
		},
	}

	var pagination = dtos.Pagination{
		Page:     1,
		PageSize: 10,
	}

	var searchAndFilter = dtos.SearchAndFilter{
		Title:          "example title",
		City:           "example city",
		Skill:          "example skill",
		MinParticipant: 1,
		MaxParticipant: 15,
	}

	var emptySearAndFilter = dtos.SearchAndFilter{
		Title:          "",
		City:           "",
		Skill:          "",
		MinParticipant: 0,
		MaxParticipant: 0,
	}

	t.Run("Succes With Searching (mobile)", func(t *testing.T) {
		repo.On("PaginateMobile", pagination.Page, pagination.PageSize, searchAndFilter).Return(vacancies).Once()
		repo.On("SelectBookmarkedVacancyID", 1).Return(bookMarkIDs, nil).Once()
		repo.On("GetTotalVolunteersByVacancyID", 1).Return(int64(15)).Once()
		repo.On("GetTotalVolunteersByVacancyID", 2).Return(int64(15)).Once()
		repo.On("GetTotalDataVacanciesBySearchAndFilterMobile", searchAndFilter).Return(int64(2)).Once()

		result, totalData := service.FindAllVacancies(pagination.Page, pagination.PageSize, searchAndFilter, 1, "mobile")
		assert.Equal(t, len(vacancies), len(result))
		assert.Equal(t, int64(2), totalData)
	})

	t.Run("Success (mobile)", func(t *testing.T) {
		emptySearAndFilter.MaxParticipant = 2147483647
		repo.On("PaginateMobile", pagination.Page, pagination.PageSize, emptySearAndFilter).Return(vacancies).Once()
		repo.On("SelectBookmarkedVacancyID", 1).Return(bookMarkIDs, nil).Once()
		repo.On("GetTotalVolunteersByVacancyID", 1).Return(int64(15)).Once()
		repo.On("GetTotalVolunteersByVacancyID", 2).Return(int64(15)).Once()
		repo.On("GetTotalDataVacanciesMobile").Return(int64(2)).Once()
		emptySearAndFilter.MaxParticipant = 0

		result, totalData := service.FindAllVacancies(pagination.Page, pagination.PageSize, emptySearAndFilter, 1, "mobile")
		assert.Equal(t, len(vacancies), len(result))
		assert.Equal(t, int64(2), totalData)
	})

	t.Run("Succes With Searching", func(t *testing.T) {
		repo.On("Paginate", pagination.Page, pagination.PageSize, searchAndFilter).Return(vacancies).Once()
		repo.On("SelectBookmarkedVacancyID", 1).Return(bookMarkIDs, nil).Once()
		repo.On("GetTotalVolunteersByVacancyID", 1).Return(int64(15)).Once()
		repo.On("GetTotalVolunteersByVacancyID", 2).Return(int64(15)).Once()
		repo.On("GetTotalDataVacanciesBySearchAndFilter", searchAndFilter).Return(int64(2)).Once()

		result, totalData := service.FindAllVacancies(pagination.Page, pagination.PageSize, searchAndFilter, 1, "")
		assert.Equal(t, len(vacancies), len(result))
		assert.Equal(t, int64(2), totalData)
	})

	t.Run("Success", func(t *testing.T) {
		emptySearAndFilter.MaxParticipant = 2147483647
		repo.On("Paginate", pagination.Page, pagination.PageSize, emptySearAndFilter).Return(vacancies).Once()
		repo.On("SelectBookmarkedVacancyID", 1).Return(bookMarkIDs, nil).Once()
		repo.On("GetTotalVolunteersByVacancyID", 1).Return(int64(15)).Once()
		repo.On("GetTotalVolunteersByVacancyID", 2).Return(int64(15)).Once()
		repo.On("GetTotalDataVacancies").Return(int64(2)).Once()

		result, totalData := service.FindAllVacancies(pagination.Page, pagination.PageSize, emptySearAndFilter, 1, "")
		assert.Equal(t, len(vacancies), len(result))
		assert.Equal(t, int64(2), totalData)
	})

	t.Run("Select Bookmarked Error", func(t *testing.T) {
		repo.On("PaginateMobile", pagination.Page, pagination.PageSize, searchAndFilter).Return(vacancies).Once()
		repo.On("SelectBookmarkedVacancyID", 1).Return(nil, errors.New("select bookmarked error")).Once()

		result, totalData := service.FindAllVacancies(pagination.Page, pagination.PageSize, searchAndFilter, 1, "mobile")
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalData)
		repo.AssertExpectations(t)
	})
}

func TestFindVacancyByID(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	vacancy := volunteer.VolunteerVacancies{
		ID:                1,
		UserID:            1,
		Title:             "bencana alam gunung meletus",
		Description:       "terjadi tsunami di suatu pantai",
		SkillsRequired:    "berenang",
		NumberOfVacancies: 15,
		ContactEmail:      "081221278393",
		Province:          "Jawa Barat",
		City:              "Banten",
		SubDistrict:       "Rangkasbitung",
		DetailLocation:    "suatu daerah di banten",
		Photo:             "https://exampleurl.com/example",
		Status:            "pending",
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("SelectVacancyByID", vacancy.ID).Return(&vacancy).Once()
		repo.On("SelectBookmarkByVacancyAndOwnerID", vacancy.ID, 1).Return("123456").Once()
		repo.On("FindUserInVacancy", vacancy.ID, 1).Return(true).Once()
		repo.On("GetTotalVolunteersByVacancyID", vacancy.ID).Return(int64(15)).Once()

		result := service.FindVacancyByID(1, 1)
		assert.Equal(t, vacancy.Title, result.Title)
		repo.AssertExpectations(t)
	})

	t.Run("Bookmark empty", func(t *testing.T) {
		repo.On("SelectVacancyByID", vacancy.ID).Return(&vacancy).Once()
		repo.On("SelectBookmarkByVacancyAndOwnerID", vacancy.ID, 1).Return("").Once()
		repo.On("FindUserInVacancy", vacancy.ID, 1).Return(true).Once()
		repo.On("GetTotalVolunteersByVacancyID", vacancy.ID).Return(int64(15)).Once()

		result := service.FindVacancyByID(1, 1)
		assert.Equal(t, vacancy.Title, result.Title)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("SelectVacancyByID", vacancy.ID).Return(nil).Once()

		result := service.FindVacancyByID(1, 1)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestModifyVacancy(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	var file multipart.File

	var vacancy = dtos.InputVacancy{
		Title:               "ini bencana alam kebakaran hutan gambut",
		Description:         "telah terjadi bencana alam disuatu lokasi yang mengakibatkan kerusakan berat dan menimbulkan korban jiwa",
		SkillsRequired:      []string{"pengendali air"},
		NumberOfVacancies:   15,
		ApplicationDeadline: time.Date(2023, time.December, 30, 15, 30, 0, 0, time.UTC),
		ContactEmail:        "example@example.com",
		Province:            "knowwhere",
		City:                "knowwhere",
		SubDistrict:         "knowwhere",
		DetailLocation:      "knowwhere",
	}

	var oldData = dtos.ResVacancy{
		ID:                  1,
		UserID:              1,
		Title:               "ini bencana alam kebakaran hutan gambut",
		Description:         "telah terjadi bencana alam disuatu lokasi yang mengakibatkan kerusakan berat dan menimbulkan korban jiwa",
		SkillsRequired:      []string{"pengendali air"},
		NumberOfVacancies:   15,
		ApplicationDeadline: time.Date(2023, time.December, 30, 15, 30, 0, 0, time.UTC),
		ContactEmail:        "example@example.com",
		Province:            "knowwhere",
		City:                "knowwhere",
		SubDistrict:         "knowwhere",
		DetailLocation:      "knowwhere",
		Photo:               "https://storage.googleapis.com//vacancies/volunteer-vacancy.jpg",
	}

	var newVacancy = volunteer.VolunteerVacancies{
		ID:                  1,
		UserID:              1,
		Title:               "ini bencana alam kebakaran hutan gambut",
		Description:         "telah terjadi bencana alam disuatu lokasi yang mengakibatkan kerusakan berat dan menimbulkan korban jiwa",
		SkillsRequired:      "pengendali air",
		NumberOfVacancies:   15,
		ApplicationDeadline: time.Date(2023, time.December, 30, 15, 30, 0, 0, time.UTC),
		ContactEmail:        "example@example.com",
		Province:            "knowwhere",
		City:                "knowwhere",
		SubDistrict:         "knowwhere",
		DetailLocation:      "knowwhere",
		Photo:               "https://storage.googleapis.com//vacancies/volunteer-vacancy.jpg",
	}

	var errValidation = []string{
		"title required",
		"description required",
	}

	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", vacancy).Return(nil).Once()
		repo.On("UploadFile", file, oldData.Photo).Return(oldData.Photo, nil).Once()
		repo.On("UpdateVacancy", newVacancy).Return(int64(1)).Once()

		errMap, err := service.ModifyVacancy(vacancy, file, oldData)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
		validation.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		validation.On("ValidateRequest", vacancy).Return(errValidation).Once()

		errMap, err := service.ModifyVacancy(vacancy, file, oldData)
		assert.Nil(t, err)
		assert.NotNil(t, errMap)
		assert.Equal(t, errValidation, errMap)
		repo.AssertExpectations(t)
		validation.AssertExpectations(t)
	})

	t.Run("Upload File Error", func(t *testing.T) {
		validation.On("ValidateRequest", vacancy).Return(nil).Once()
		repo.On("UploadFile", file, oldData.Photo).Return("", errors.New("upload file failed")).Once()

		errMap, err := service.ModifyVacancy(vacancy, file, oldData)
		assert.Nil(t, errMap)
		assert.Error(t, err)
		assert.EqualError(t, err, "upload file failed")
		repo.AssertExpectations(t)
		validation.AssertExpectations(t)
	})

	t.Run("Update Data Error", func(t *testing.T) {
		validation.On("ValidateRequest", vacancy).Return(nil).Once()
		repo.On("UploadFile", file, oldData.Photo).Return(oldData.Photo, nil).Once()
		repo.On("UpdateVacancy", newVacancy).Return(int64(0)).Once()

		errMap, err := service.ModifyVacancy(vacancy, file, oldData)
		assert.Nil(t, errMap)
		assert.Error(t, err)
		assert.EqualError(t, err, "there is no vacancy updated")
		repo.AssertExpectations(t)
		validation.AssertExpectations(t)
	})
}

func TestModifyVacancyStatus(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	var statusAccepted = dtos.StatusVacancies{
		Status: "accepted",
	}

	var statusRejected = dtos.StatusVacancies{
		Status:         "rejected",
		RejectedReason: "too kind",
	}

	var statusRejectedWithoutReason = dtos.StatusVacancies{
		Status: "rejected",
	}

	var errValidation = []string{
		"status required",
	}

	var oldData = dtos.ResVacancy{
		ID:                  1,
		UserID:              1,
		Title:               "ini bencana alam kebakaran hutan gambut",
		Description:         "telah terjadi bencana alam disuatu lokasi yang mengakibatkan kerusakan berat dan menimbulkan korban jiwa",
		SkillsRequired:      []string{"pengendali air"},
		NumberOfVacancies:   15,
		ApplicationDeadline: time.Date(2023, time.December, 30, 15, 30, 0, 0, time.UTC),
		ContactEmail:        "example@example.com",
		Province:            "knowwhere",
		City:                "knowwhere",
		SubDistrict:         "knowwhere",
		DetailLocation:      "knowwhere",
		Photo:               "https://storage.googleapis.com//vacancies/volunteer-vacancy.jpg",
	}

	var newVacancyAccepted = volunteer.VolunteerVacancies{
		ID:     1,
		Status: "accepted",
	}

	var newVacancyRejected = volunteer.VolunteerVacancies{
		ID:             1,
		Status:         "rejected",
		RejectedReason: "too kind",
	}

	t.Run("Success Update Status Accepted", func(t *testing.T) {
		repo.On("SelectByTittle", oldData.Title).Return(errors.New("data not found")).Once()
		validation.On("ValidateRequest", statusAccepted).Return(nil).Once()
		repo.On("GetDeviceToken", oldData.ID).Return("deviceToken").Once()
		message := "Kami ingin memberitahu bahwa pengajuan lowongan relawan Anda untuk " + oldData.Title + " telah diterima! Terima kasih atas langkah inisiatif Anda."
		notification.On("SendNotifications", "deviceToken", "Pengajuan Lowongan Relawan Diterima", message).Return(nil).Once()
		repo.On("UpdateVacancy", newVacancyAccepted).Return(int64(1)).Once()

		err, errMap := service.ModifyVacancyStatus(statusAccepted, oldData)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Success Update Status Rejected", func(t *testing.T) {
		repo.On("SelectByTittle", oldData.Title).Return(errors.New("data not found")).Once()
		validation.On("ValidateRequest", statusRejected).Return(nil).Once()
		repo.On("GetDeviceToken", oldData.ID).Return("deviceToken").Once()
		message := "Terima kasih sudah mengajukan lowongan relawan. Saat ini, kami belum bisa menyetujui permohonan ini karena.\n\nAlasan : " + statusRejected.RejectedReason + "\n\nTerima kasih atas partisipasinya"
		notification.On("SendNotifications", "deviceToken", "Pengajuan Lowongan Relawan Ditolak", message).Return(nil).Once()
		repo.On("UpdateVacancy", newVacancyRejected).Return(int64(1)).Once()

		err, errMap := service.ModifyVacancyStatus(statusRejected, oldData)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Title Already Used", func(t *testing.T) {
		repo.On("SelectByTittle", oldData.Title).Return(nil).Once()

		err, errMap := service.ModifyVacancyStatus(statusRejected, oldData)
		assert.Nil(t, errMap)
		assert.Error(t, err)
		assert.EqualError(t, err, "title already used by another vacancy")
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		repo.On("SelectByTittle", oldData.Title).Return(errors.New("data not found")).Once()
		validation.On("ValidateRequest", statusAccepted).Return(errValidation).Once()

		err, errMap := service.ModifyVacancyStatus(statusAccepted, oldData)
		assert.NotNil(t, errMap)
		assert.Equal(t, errValidation, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
	})

	t.Run("Status Rejected Not Found", func(t *testing.T) {
		repo.On("SelectByTittle", oldData.Title).Return(errors.New("data not found")).Once()
		validation.On("ValidateRequest", statusRejectedWithoutReason).Return(nil).Once()
		repo.On("GetDeviceToken", oldData.ID).Return("deviceToken").Once()

		err, errMap := service.ModifyVacancyStatus(statusRejectedWithoutReason, oldData)
		assert.NotNil(t, errMap)
		assert.Equal(t, []string{"rejected_reason field is required when the status is rejected"}, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
	})

	t.Run("Update Data Failed", func(t *testing.T) {
		repo.On("SelectByTittle", oldData.Title).Return(errors.New("data not found")).Once()
		validation.On("ValidateRequest", statusAccepted).Return(nil).Once()
		repo.On("GetDeviceToken", oldData.ID).Return("deviceToken").Once()
		message := "Kami ingin memberitahu bahwa pengajuan lowongan relawan Anda untuk " + oldData.Title + " telah diterima! Terima kasih atas langkah inisiatif Anda."
		notification.On("SendNotifications", "deviceToken", "Pengajuan Lowongan Relawan Diterima", message).Return(nil).Once()
		repo.On("UpdateVacancy", newVacancyAccepted).Return(int64(0)).Once()

		err, errMap := service.ModifyVacancyStatus(statusAccepted, oldData)
		assert.Nil(t, errMap)
		assert.Error(t, err)
		assert.EqualError(t, err, "there is no vacancy updated")
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

func TestUpdateStatusRegistrar(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	var statusAccepted = dtos.StatusRegistrar{
		Status: "accepted",
	}

	var statusRejected = dtos.StatusRegistrar{
		Status:         "rejected",
		RejectedReason: "too kind",
	}

	var statusRejectedWithoutReason = dtos.StatusRegistrar{
		Status: "rejected",
	}

	var registrarAccepted = volunteer.VolunteerRelations{
		ID:          1,
		UserID:      1,
		VolunteerID: 1,
		Status:      "accepted",
	}

	var registrarRejected = volunteer.VolunteerRelations{
		ID:             1,
		UserID:         1,
		VolunteerID:    1,
		Status:         "rejected",
		RejectedReason: "too kind",
	}

	var errValidation = []string{
		"rejected_reason required",
	}

	var vacancy = volunteer.VolunteerVacancies{
		ID:                1,
		UserID:            1,
		Title:             "bencana alam gunung meletus",
		Description:       "terjadi tsunami di suatu pantai",
		SkillsRequired:    "berenang",
		NumberOfVacancies: 15,
		ContactEmail:      "081221278393",
		Province:          "Jawa Barat",
		City:              "Banten",
		SubDistrict:       "Rangkasbitung",
		DetailLocation:    "suatu daerah di banten",
		Photo:             "https://exampleurl.com/example",
		Status:            "pending",
	}

	t.Run("Success Update Status Accepted", func(t *testing.T) {
		validation.On("ValidateRequest", statusAccepted).Return(nil).Once()
		repo.On("SelectRegistrarByID", registrarAccepted.ID).Return(&registrarAccepted).Once()
		repo.On("SelectVacancyByID", registrarAccepted.VolunteerID).Return(&vacancy).Once()
		repo.On("GetDeviceToken", registrarAccepted.UserID).Return("deviceToken").Once()
		message := "Kami ingin memberitahu Anda bahwa pengajuan sebagai relawan di " + vacancy.Title + " telah diterima! Terima kasih atas minat dan dedikasi Anda."
		notification.On("SendNotifications", "deviceToken", "Pengajuan Relawan Diterima", message).Return(nil).Once()
		repo.On("UpdateStatusRegistrar", registrarAccepted).Return(int64(1)).Once()

		result, errMap := service.UpdateStatusRegistrar(statusAccepted, 1)
		assert.Nil(t, errMap)
		assert.NotNil(t, result)
		assert.Equal(t, true, result)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Success Update Status Rejected", func(t *testing.T) {
		validation.On("ValidateRequest", statusRejected).Return(nil).Once()
		repo.On("SelectRegistrarByID", registrarRejected.ID).Return(&registrarRejected).Once()
		repo.On("SelectVacancyByID", registrarRejected.VolunteerID).Return(&vacancy).Once()
		repo.On("GetDeviceToken", registrarRejected.UserID).Return("deviceToken").Once()
		message := "Terima kasih sudah mengajukan diri sebagai relawan. Saat ini, kami belum bisa menyetujui permohonan Anda sebagai relawan di " + vacancy.Title + ".\n\nAlasan : " + registrarRejected.RejectedReason + "\n\nTerima kasih atas partisipasinya"
		notification.On("SendNotifications", "deviceToken", "Pengajuan Relawan Ditolak", message).Return(nil).Once()
		repo.On("UpdateStatusRegistrar", registrarRejected).Return(int64(1)).Once()

		result, errMap := service.UpdateStatusRegistrar(statusRejected, 1)
		assert.Nil(t, errMap)
		assert.NotNil(t, result)
		assert.Equal(t, true, result)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		validation.On("ValidateRequest", statusAccepted).Return(errValidation).Once()

		result, errMap := service.UpdateStatusRegistrar(statusAccepted, 1)
		assert.NotNil(t, errMap)
		assert.Equal(t, errValidation, errMap)
		assert.Equal(t, false, result)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Registrar Not Found", func(t *testing.T) {
		validation.On("ValidateRequest", statusAccepted).Return(nil).Once()
		repo.On("SelectRegistrarByID", 1).Return(nil).Once()

		result, errMap := service.UpdateStatusRegistrar(statusAccepted, 1)
		assert.Nil(t, errMap)
		assert.Equal(t, false, result)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Rejected Reason Not Found", func(t *testing.T) {
		validation.On("ValidateRequest", statusRejectedWithoutReason).Return(nil).Once()
		repo.On("SelectRegistrarByID", registrarRejected.ID).Return(&registrarRejected).Once()
		repo.On("SelectVacancyByID", registrarRejected.VolunteerID).Return(&vacancy).Once()
		repo.On("GetDeviceToken", registrarRejected.UserID).Return("deviceToken").Once()

		result, errMap := service.UpdateStatusRegistrar(statusRejectedWithoutReason, 1)
		assert.NotNil(t, errMap)
		assert.Equal(t, []string{"rejected_reason field is required when the status is rejected"}, errMap)
		assert.Equal(t, false, result)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Update Data Failed", func(t *testing.T) {
		validation.On("ValidateRequest", statusAccepted).Return(nil).Once()
		repo.On("SelectRegistrarByID", registrarAccepted.ID).Return(&registrarAccepted).Once()
		repo.On("SelectVacancyByID", registrarAccepted.VolunteerID).Return(&vacancy).Once()
		repo.On("GetDeviceToken", registrarAccepted.UserID).Return("deviceToken").Once()
		message := "Kami ingin memberitahu Anda bahwa pengajuan sebagai relawan di " + vacancy.Title + " telah diterima! Terima kasih atas minat dan dedikasi Anda."
		notification.On("SendNotifications", "deviceToken", "Pengajuan Relawan Diterima", message).Return(nil).Once()
		repo.On("UpdateStatusRegistrar", registrarAccepted).Return(int64(0)).Once()

		result, errMap := service.UpdateStatusRegistrar(statusAccepted, 1)
		assert.Nil(t, errMap)
		assert.NotNil(t, result)
		assert.Equal(t, false, result)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

func TestCreateVacancy(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	file, err := os.Open("file-mock.jpg")
	if err != nil {
		logrus.Error(err)
	}
	defer file.Close()

	var newVacancy = dtos.InputVacancy{
		Title:               "ini bencana alam kebakaran hutan gambut",
		Description:         "telah terjadi bencana alam disuatu lokasi yang mengakibatkan kerusakan berat dan menimbulkan korban jiwa",
		SkillsRequired:      []string{"skil1", "skill2"},
		NumberOfVacancies:   20,
		ApplicationDeadline: time.Date(2023, time.December, 30, 15, 30, 0, 0, time.UTC),
		ContactEmail:        "didadejan45@gmail.com",
		Province:            "jawa barat",
		City:                "banten",
		SubDistrict:         "kec banten",
		DetailLocation:      "di suatu pantai daerah banten",
		Photo:               file,
	}

	var wrongData = dtos.InputVacancy{
		Title:               "ini bencana",
		Description:         "telah terjadi bencana",
		SkillsRequired:      []string{},
		NumberOfVacancies:   0,
		ApplicationDeadline: time.Date(2023, time.December, 1, 15, 30, 0, 0, time.UTC),
		ContactEmail:        "didadejan45@gmail.com",
		Province:            "jawa barat",
		City:                "banten",
		SubDistrict:         "kec banten",
		DetailLocation:      "di suatu pantai daerah banten",
		Photo:               file,
	}

	var vacancyData = volunteer.VolunteerVacancies{
		ID:                  1,
		UserID:              1,
		Title:               "ini bencana alam kebakaran hutan gambut",
		Description:         "telah terjadi bencana alam disuatu lokasi yang mengakibatkan kerusakan berat dan menimbulkan korban jiwa",
		SkillsRequired:      "skil1, skill2",
		NumberOfVacancies:   20,
		ApplicationDeadline: time.Date(2023, time.December, 30, 15, 30, 0, 0, time.UTC),
		ContactEmail:        "didadejan45@gmail.com",
		Province:            "jawa barat",
		City:                "banten",
		SubDistrict:         "kec banten",
		DetailLocation:      "di suatu pantai daerah banten",
		Status:              "pending",
		Photo:               "https://storage.googleapis.com//vacancies/volunteer-vacancy.jpg",
	}

	var errValidation = []string{
		"title required",
		"title must be at least 20 characters",
		"description must be at least 50 characters",
		"skillsRequired must be at least 1 word",
		"numberOfVacancies must be greater than 1",
		"applicationDeadline must be greater than today",
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
		assert.Equal(t, "ini bencana alam kebakaran hutan gambut", result.Title)
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

	t.Run("Duplicate title", func(t *testing.T) {
		repo.On("SelectByTittle", newVacancy.Title).Return(nil).Once()

		result, errMap, err := service.CreateVacancy(newVacancy, 1, file)
		assert.Error(t, err)
		assert.EqualError(t, err, "title already used by another vacancy")
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		repo.AssertExpectations(t)
	})

	t.Run("Validation error", func(t *testing.T) {
		repo.On("SelectByTittle", wrongData.Title).Return(errors.New("data not found")).Once()
		validation.On("ValidateRequest", wrongData).Return(errValidation).Once()

		result, errMap, err := service.CreateVacancy(wrongData, 1, file)
		assert.NotNil(t, errMap)
		assert.Nil(t, result)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
		validation.AssertExpectations(t)
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

	t.Run("Insert Data Error", func(t *testing.T) {
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
		repo.On("InsertVacancy", &vacancy).Return(nil, errors.New("insert data failed")).Once()

		result, errMap, err := service.CreateVacancy(newVacancy, 1, file)
		assert.Error(t, err)
		assert.EqualError(t, err, "Use Case : failed to create volunteer")
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		repo.AssertExpectations(t)
		validation.AssertExpectations(t)
	})
}

func TestDeleteVacancy(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

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

func TestRegisterVacancy(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	var file multipart.File

	newApply := dtos.ApplyVacancy{
		VacancyID: 1,
		Skills:    []string{"skill1", "skill2"},
		Resume:    "ini adalah suatu cv saya yang sangat lengkap karena saya ingin daftar menjadi volunteer pada vacancy ini",
		Reason:    "alasan saya mendaftar karena saya jago di bidang tersebut dan kebetulan saya pengangguran",
		Photo:     file,
	}

	registrar := volunteer.VolunteerRelations{
		UserID:      1,
		VolunteerID: 1,
		Skills:      "skill1, skill2",
		Reason:      "alasan saya mendaftar karena saya jago di bidang tersebut dan kebetulan saya pengangguran",
		Resume:      "ini adalah suatu cv saya yang sangat lengkap karena saya ingin daftar menjadi volunteer pada vacancy ini",
		Photo:       "https://storage.googleapis.com//vacancies/volunteer-vacancy.jpg",
	}

	var errValidation = []string{
		"vacancy_id required",
		"skills required",
		"resume required",
		"reason required",
		"photo required",
	}

	t.Run("Succes Register Vacancy", func(t *testing.T) {
		validation.On("ValidateRequest", newApply).Return(nil).Once()
		repo.On("UploadFile", file, "").Return("https://storage.googleapis.com//vacancies/volunteer-vacancy.jpg", nil).Once()
		repo.On("RegisterVacancy", &registrar).Return(nil).Once()

		result, errMap := service.RegisterVacancy(newApply, 1)
		assert.Nil(t, errMap)
		assert.NotNil(t, result)
		assert.Equal(t, true, result)
		repo.AssertExpectations(t)
		validation.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		validation.On("ValidateRequest", newApply).Return(errValidation).Once()

		result, errMap := service.RegisterVacancy(newApply, 1)
		assert.NotNil(t, errMap)
		assert.Equal(t, errValidation, errMap)
		assert.NotNil(t, result)
		assert.Equal(t, false, result)
		validation.AssertExpectations(t)
	})

	t.Run("Upload File Error", func(t *testing.T) {
		validation.On("ValidateRequest", newApply).Return(nil).Once()
		repo.On("UploadFile", file, "").Return("", errors.New("upload file failed")).Once()

		result, errMap := service.RegisterVacancy(newApply, 1)
		assert.Nil(t, errMap)
		assert.NotNil(t, result)
		assert.Equal(t, false, result)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Register Data Error", func(t *testing.T) {
		validation.On("ValidateRequest", newApply).Return(nil).Once()
		repo.On("UploadFile", file, "").Return("https://storage.googleapis.com//vacancies/volunteer-vacancy.jpg", nil).Once()
		repo.On("RegisterVacancy", &registrar).Return(errors.New("register data failed")).Once()

		result, errMap := service.RegisterVacancy(newApply, 1)
		assert.Nil(t, errMap)
		assert.NotNil(t, result)
		assert.Equal(t, false, result)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

func TestFindAllVolunteersByVacancyID(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	var volunteers = []volunteer.Volunteer{
		{
			ID:          1,
			Email:       "test1@gmail.com",
			Fullname:    "Test1",
			Address:     "Bogor",
			PhoneNumber: "081234567890",
			Gender:      "Male",
			Skills:      "Berenang",
			Nik:         "3212345678901234",
			Resume:      "something",
			Reason:      "something",
			Photo:       "www.example.com",
			Status:      "accepted",
		},
		{
			ID:          2,
			Email:       "test2@gmail.com",
			Fullname:    "Test2",
			Address:     "Bogor",
			PhoneNumber: "081234567890",
			Gender:      "Male",
			Skills:      "Berenang",
			Nik:         "3212345678901234",
			Resume:      "something",
			Reason:      "something",
			Photo:       "www.example.com",
			Status:      "accepted",
		},
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("SelectVolunteersByVacancyID", 1, "", 1, 10).Return(volunteers).Once()
		repo.On("GetTotalVolunteers", 1, "").Return(int64(2)).Once()

		result, totalData := service.FindAllVolunteersByVacancyID(1, 10, 1, "")
		assert.NotNil(t, result)
		assert.Equal(t, len(volunteers), len(result))
		assert.Equal(t, int64(2), totalData)
		repo.AssertExpectations(t)
	})

	t.Run("Data Not Found", func(t *testing.T) {
		repo.On("SelectVolunteersByVacancyID", 1, "", 1, 10).Return(nil, 0).Once()

		result, totalData := service.FindAllVolunteersByVacancyID(1, 10, 1, "")
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalData)
		repo.AssertExpectations(t)
	})

	t.Run("Get Total Data Error", func(t *testing.T) {
		repo.On("SelectVolunteersByVacancyID", 1, "", 1, 10).Return(volunteers).Once()
		repo.On("GetTotalVolunteers", 1, "").Return(int64(0)).Once()

		result, totalData := service.FindAllVolunteersByVacancyID(1, 10, 1, "")
		assert.NotNil(t, result)
		assert.Equal(t, len(volunteers), len(result))
		assert.Equal(t, int64(0), totalData)
		repo.AssertExpectations(t)
	})
}

func TestFindDetailVolunteers(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	var volunteer = volunteer.Volunteer{
		ID:          1,
		Email:       "test1@gmail.com",
		Fullname:    "Test1",
		Address:     "Bogor",
		PhoneNumber: "081234567890",
		Gender:      "Male",
		Skills:      "Berenang",
		Nik:         "3212345678901234",
		Resume:      "something",
		Reason:      "something",
		Photo:       "www.example.com",
		Status:      "accepted",
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("SelectVolunteerDetails", 1, 1).Return(&volunteer).Once()

		result := service.FindDetailVolunteers(1, 1)
		assert.NotNil(t, result)
		assert.Equal(t, volunteer.Email, result.Email)
		repo.AssertExpectations(t)
	})

	t.Run("Data Not Found", func(t *testing.T) {
		repo.On("SelectVolunteerDetails", 1, 1).Return(nil).Once()

		result := service.FindDetailVolunteers(1, 1)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestCheckUser(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	t.Run("Success", func(t *testing.T) {
		repo.On("CheckUser", 1).Return(true).Once()

		result := service.CheckUser(1)
		assert.Equal(t, true, result)
		repo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		repo.On("CheckUser", 1).Return(false).Once()

		result := service.CheckUser(1)
		assert.Equal(t, false, result)
		repo.AssertExpectations(t)
	})
}

func TestFindUserInVacancy(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	t.Run("Success", func(t *testing.T) {
		repo.On("FindUserInVacancy", 1, 1).Return(true).Once()

		result := service.FindUserInVacancy(1, 1)
		assert.Equal(t, true, result)
		repo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		repo.On("FindUserInVacancy", 1, 1).Return(false).Once()

		result := service.FindUserInVacancy(1, 1)
		assert.Equal(t, false, result)
		repo.AssertExpectations(t)
	})
}

func TestFindAllSkills(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	notification := helperMocks.NewNotificationInterface(t)
	service := New(repo, validation, notification)

	var skills = []dtos.Skill{
		{
			ID:   1,
			Name: "Berenang",
		},
		{
			ID:   2,
			Name: "Mengemudi",
		},
		{
			ID:   3,
			Name: "Memanjat",
		},
	}

	t.Run("Success", func(t *testing.T) {
		repo.On("SelectAllSkills").Return(skills, nil).Once()

		result, err := service.FindAllSkills()
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, len(skills), len(result))
		repo.AssertExpectations(t)
	})

	t.Run("Data Not Found", func(t *testing.T) {
		repo.On("SelectAllSkills").Return(nil, errors.New("data not found")).Once()

		result, err := service.FindAllSkills()
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
