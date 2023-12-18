package usecase

import (
	"errors"
	"math"
	"os"
	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/dtos"
	"raihpeduli/features/fundraise/mocks"
	helperMocks "raihpeduli/helpers/mocks"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(repository, validation, nsRequest)
	var bookmarkId = "09123890awdaw0192"
	var fundraises = []fundraise.Fundraise{
		{
			ID:          1,
			Title:       "Pembangunan Masjid",
			Description: "Lorem ipsum",
			Photo:       "https://googleapis.com/awdadwd",
			Target:      500000000,
			StartDate:   time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC),
			EndDate:     time.Date(2024, time.April, 19, 15, 30, 0, 0, time.UTC),
			Status:      "accepted",
			UserID:      1,
			CreatedAt:   time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC),
		},
	}
	pagination := dtos.Pagination{
		Page:     1,
		PageSize: 10,
	}
	searchAndFilter := dtos.SearchAndFilter{
		Title:     "elgjrut",
		MinTarget: 1,
		MaxTarget: 10,
	}
	t.Run("Success Web", func(t *testing.T) {
		repository.On("Paginate", pagination, searchAndFilter).Return(fundraises, nil).Once()
		repository.On("SelectBookmarkedFundraiseID", fundraises[0].UserID).Return(map[int]string{
			1: bookmarkId,
		}, nil).Once()
		repository.On("TotalFundAcquired", 1).Return(int32(0), nil).Once()
		repository.On("GetTotalDataBySearchAndFilter", searchAndFilter).Return(int64(1), nil).Once()
		result, totalData := service.FindAll(pagination, searchAndFilter, 1, "web")
		assert.NotNil(t, result)
		assert.Equal(t, result[0].ID, fundraises[0].ID)
		assert.Equal(t, totalData, int64(len(fundraises)))
		assert.Equal(t, result[0].BookmarkID, bookmarkId)
		repository.AssertExpectations(t)
	})

	t.Run("Success Mobile", func(t *testing.T) {
		repository.On("PaginateMobile", dtos.Pagination{
			Page:     1,
			PageSize: 10,
		}, dtos.SearchAndFilter{
			Title:     "elgjrut",
			MinTarget: 1,
			MaxTarget: math.MaxInt32,
		}).Return(fundraises, nil).Once()
		repository.On("SelectBookmarkedFundraiseID", fundraises[0].UserID).Return(map[int]string{
			1: bookmarkId,
		}, nil).Once()
		repository.On("TotalFundAcquired", 1).Return(int32(0), nil).Once()
		repository.On("GetTotalDataBySearchAndFilterMobile", dtos.SearchAndFilter{
			Title:     "elgjrut",
			MinTarget: 1,
			MaxTarget: math.MaxInt32,
		}).Return(int64(1), nil).Once()
		result, totalData := service.FindAll(dtos.Pagination{
			Page:     0,
			PageSize: 0,
		}, dtos.SearchAndFilter{
			Title:     "elgjrut",
			MinTarget: 1,
			MaxTarget: 0,
		}, 1, "mobile")
		assert.NotNil(t, result)
		assert.Equal(t, result[0].ID, fundraises[0].ID)
		assert.Equal(t, totalData, int64(len(fundraises)))
		assert.Equal(t, result[0].BookmarkID, bookmarkId)
		repository.AssertExpectations(t)
	})
	t.Run("Error Select Bookmark", func(t *testing.T) {
		repository.On("PaginateMobile", dtos.Pagination{
			Page:     1,
			PageSize: 10,
		}, dtos.SearchAndFilter{
			Title:     "elgjrut",
			MinTarget: 1,
			MaxTarget: math.MaxInt32,
		}).Return(fundraises, nil).Once()
		repository.On("SelectBookmarkedFundraiseID", fundraises[0].UserID).Return(nil, errors.New("error")).Once()
		result, totalData := service.FindAll(dtos.Pagination{
			Page:     0,
			PageSize: 0,
		}, dtos.SearchAndFilter{
			Title:     "elgjrut",
			MinTarget: 1,
			MaxTarget: 0,
		}, 1, "mobile")
		assert.Nil(t, result)
		assert.Equal(t, totalData, int64(0))
		repository.AssertExpectations(t)
	})
	t.Run("Error Paginate", func(t *testing.T) {
		repository.On("PaginateMobile", dtos.Pagination{
			Page:     1,
			PageSize: 10,
		}, dtos.SearchAndFilter{
			Title:     "elgjrut",
			MinTarget: 1,
			MaxTarget: math.MaxInt32,
		}).Return(nil, errors.New("error")).Once()
		result, totalData := service.FindAll(dtos.Pagination{
			Page:     0,
			PageSize: 0,
		}, dtos.SearchAndFilter{
			Title:     "elgjrut",
			MinTarget: 1,
			MaxTarget: 0,
		}, 1, "mobile")
		assert.Nil(t, result)
		assert.NotNil(t, totalData)
		repository.AssertExpectations(t)
	})

}

func TestFindById(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(repository, validation, nsRequest)
	var bookmarkId = "09123890awdaw0192"
	var fundraiseItem = &dtos.FundraiseDetails{
		ID:          1,
		Title:       "Pembangunan Masjid",
		Description: "Lorem ipsum",
		Photo:       "https://googleapis.com/awdadwd",
		Target:      500000000,
		StartDate:   time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.April, 19, 15, 30, 0, 0, time.UTC),
		Status:      "accepted",
		UserID:      1,
		CreatedAt:   time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC),
	}
	t.Run("Success", func(t *testing.T) {
		repository.On("SelectByID", fundraiseItem.ID).Return(fundraiseItem, nil).Once()
		repository.On("SelectBookmarkByFundraiseAndOwnerID", fundraiseItem.ID, fundraiseItem.UserID).Return(bookmarkId, nil).Once()
		repository.On("TotalFundAcquired", 1).Return(int32(0), nil).Once()
		result := service.FindByID(1, 1)
		assert.NotNil(t, result)
		assert.Equal(t, result.ID, fundraiseItem.ID)
		assert.Equal(t, fundraiseItem.BookmarkID, bookmarkId)
		repository.AssertExpectations(t)
	})
	t.Run("Error Total Fund Acquired", func(t *testing.T) {
		repository.On("SelectByID", fundraiseItem.ID).Return(fundraiseItem, nil).Once()
		repository.On("SelectBookmarkByFundraiseAndOwnerID", fundraiseItem.ID, fundraiseItem.UserID).Return(bookmarkId, nil).Once()
		repository.On("TotalFundAcquired", 1).Return(int32(0), errors.New("error")).Once()
		result := service.FindByID(1, 1)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
	t.Run("Error SelectById", func(t *testing.T) {
		repository.On("SelectByID", fundraiseItem.ID).Return(nil, errors.New("error")).Once()
		result := service.FindByID(1, 1)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

}

func TestCreate(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(repository, validation, nsRequest)

	mockFile, err := os.Open("file-mock.jpg")
	if err != nil {
		logrus.Error(err)
	}
	defer mockFile.Close()
	var fundraiseItem = fundraise.Fundraise{
		ID:          0,
		Title:       "Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam gravida, turpis consequat malesuada luctus, nisl dolor dignissim justo, a molestie sem massa et nulla. Duis diam ligula, iaculis lacinia iaculis sed, finibus eu urna. Mauris et auctor est. Etiam elementum tortor ac velit porttitor semper. Pellentesque habitant morbi tristique senectus.",
		Photo:       "https://googleapis.com/awdadwd",
		Target:      500000000,
		StartDate:   time.Date(2023, time.December, 19, 15, 30, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.December, 23, 15, 30, 0, 0, time.UTC),
		Status:      "pending",
		UserID:      1,
		CreatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	var fundraiseItemFileNil = fundraise.Fundraise{
		ID:          0,
		Title:       "Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam gravida, turpis consequat malesuada luctus, nisl dolor dignissim justo, a molestie sem massa et nulla. Duis diam ligula, iaculis lacinia iaculis sed, finibus eu urna. Mauris et auctor est. Etiam elementum tortor ac velit porttitor semper. Pellentesque habitant morbi tristique senectus.",
		Photo:       "https://storage.googleapis.com//fundraises/default",
		Target:      500000000,
		StartDate:   time.Date(2023, time.December, 19, 15, 30, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.December, 23, 15, 30, 0, 0, time.UTC),
		Status:      "pending",
		UserID:      1,
		CreatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	var fundraiseItemInput = dtos.InputFundraise{
		Title:       "Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam gravida, turpis consequat malesuada luctus, nisl dolor dignissim justo, a molestie sem massa et nulla. Duis diam ligula, iaculis lacinia iaculis sed, finibus eu urna. Mauris et auctor est. Etiam elementum tortor ac velit porttitor semper. Pellentesque habitant morbi tristique senectus.",
		Target:      500000000,
		StartDate:   time.Date(2023, time.December, 19, 15, 30, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.December, 23, 15, 30, 0, 0, time.UTC),
	}
	t.Run("Success", func(t *testing.T) {
		repository.On("SelectByTitle", fundraiseItemInput.Title).Return(nil, errors.New("error")).Once()
		validation.On("ValidateRequest", fundraiseItemInput).Return(nil).Once()
		repository.On("UploadFile", mockFile).Return("https://googleapis.com/awdadwd", nil).Once()
		repository.On("Insert", fundraiseItem).Return(&fundraiseItem, nil).Once()
		result, _, _ := service.Create(fundraiseItemInput, 1, mockFile)
		assert.NotNil(t, result)
		assert.Equal(t, result.ID, fundraiseItem.ID)
		assert.Equal(t, result.UserID, fundraiseItem.UserID)
		repository.AssertExpectations(t)
	})
	t.Run("Success with file nil", func(t *testing.T) {
		repository.On("SelectByTitle", fundraiseItemInput.Title).Return(nil, errors.New("error")).Once()
		validation.On("ValidateRequest", fundraiseItemInput).Return(nil).Once()
		repository.On("Insert", fundraiseItemFileNil).Return(&fundraiseItemFileNil, nil).Once()
		result, _, _ := service.Create(fundraiseItemInput, 1, nil)
		assert.NotNil(t, result)
		assert.Equal(t, result.ID, fundraiseItem.ID)
		assert.Equal(t, result.UserID, fundraiseItem.UserID)
		repository.AssertExpectations(t)
	})
	t.Run("Error Upload", func(t *testing.T) {
		repository.On("SelectByTitle", fundraiseItemInput.Title).Return(nil, errors.New("error")).Once()
		validation.On("ValidateRequest", fundraiseItemInput).Return(nil).Once()
		repository.On("UploadFile", mockFile).Return("error", errors.New("error")).Once()
		result, _, _ := service.Create(fundraiseItemInput, 1, mockFile)
		assert.Nil(t, result)

		repository.AssertExpectations(t)
	})
	t.Run("Error Validate", func(t *testing.T) {
		repository.On("SelectByTitle", fundraiseItemInput.Title).Return(nil, errors.New("error")).Once()
		validation.On("ValidateRequest", fundraiseItemInput).Return([]string{"error"}, errors.New("error")).Once()
		result, _, _ := service.Create(fundraiseItemInput, 1, mockFile)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
	t.Run("Error title used", func(t *testing.T) {
		repository.On("SelectByTitle", fundraiseItemInput.Title).Return(nil, nil).Once()
		result, _, _ := service.Create(fundraiseItemInput, 1, mockFile)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})

}

func TestModify(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(repository, validation, nsRequest)

	mockFile, err := os.Open("file-mock.jpg")
	if err != nil {
		logrus.Error(err)
	}
	defer mockFile.Close()
	var fundraiseItem = fundraise.Fundraise{
		ID:          0,
		Title:       "Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam gravida, turpis consequat malesuada luctus, nisl dolor dignissim justo, a molestie sem massa et nulla. Duis diam ligula, iaculis lacinia iaculis sed, finibus eu urna. Mauris et auctor est. Etiam elementum tortor ac velit porttitor semper. Pellentesque habitant morbi tristique senectus.",
		Photo:       "https://googleapis.com/awdadwd",
		Target:      500000000,
		StartDate:   time.Date(2023, time.December, 19, 15, 30, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.December, 23, 15, 30, 0, 0, time.UTC),
		Status:      "pending",
		UserID:      1,
		CreatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	var oldFundraiseItem = dtos.FundraiseDetails{
		ID:          0,
		Title:       "Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam gravida, turpis consequat malesuada luctus, nisl dolor dignissim justo, a molestie sem massa et nulla. Duis diam ligula, iaculis lacinia iaculis sed, finibus eu urna. Mauris et auctor est. Etiam elementum tortor ac velit porttitor semper. Pellentesque habitant morbi tristique senectus.",
		Photo:       "https://googleapis.com/awdadwd",
		Target:      500000000,
		StartDate:   time.Date(2023, time.December, 19, 15, 30, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.December, 23, 15, 30, 0, 0, time.UTC),
		Status:      "pending",
		UserID:      1,
		CreatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	var fundraiseItemInput = dtos.InputFundraise{
		Title:       "Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam gravida, turpis consequat malesuada luctus, nisl dolor dignissim justo, a molestie sem massa et nulla. Duis diam ligula, iaculis lacinia iaculis sed, finibus eu urna. Mauris et auctor est. Etiam elementum tortor ac velit porttitor semper. Pellentesque habitant morbi tristique senectus.",
		Target:      500000000,
		StartDate:   time.Date(2023, time.December, 19, 15, 30, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.December, 23, 15, 30, 0, 0, time.UTC),
	}
	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", fundraiseItemInput).Return(nil).Once()
		repository.On("DeleteFile", "https://googleapis.com/awdadwd").Return(nil).Once()
		repository.On("UploadFile", mockFile).Return("https://googleapis.com/awdadwd", nil).Once()
		repository.On("Update", fundraiseItem).Return(nil).Once()
		result, _ := service.Modify(fundraiseItemInput, mockFile, oldFundraiseItem)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
	t.Run("Error Update", func(t *testing.T) {
		validation.On("ValidateRequest", fundraiseItemInput).Return(nil).Once()
		repository.On("DeleteFile", "https://googleapis.com/awdadwd").Return(nil).Once()
		repository.On("UploadFile", mockFile).Return("https://googleapis.com/awdadwd", nil).Once()
		repository.On("Update", fundraiseItem).Return(errors.New("error")).Once()
		result, _ := service.Modify(fundraiseItemInput, mockFile, oldFundraiseItem)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
	t.Run("Error Upload", func(t *testing.T) {
		validation.On("ValidateRequest", fundraiseItemInput).Return(nil).Once()
		repository.On("DeleteFile", "https://googleapis.com/awdadwd").Return(nil).Once()
		repository.On("UploadFile", mockFile).Return("", errors.New("error")).Once()
		result, _ := service.Modify(fundraiseItemInput, mockFile, oldFundraiseItem)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
	t.Run("Error Delete", func(t *testing.T) {
		validation.On("ValidateRequest", fundraiseItemInput).Return(nil).Once()
		repository.On("DeleteFile", "https://googleapis.com/awdadwd").Return(errors.New("error")).Once()
		repository.On("UploadFile", mockFile).Return("https://googleapis.com/awdadwd", nil).Once()
		repository.On("Update", fundraiseItem).Return(nil).Once()
		result, _ := service.Modify(fundraiseItemInput, mockFile, oldFundraiseItem)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
	t.Run("Error Validation", func(t *testing.T) {
		validation.On("ValidateRequest", fundraiseItemInput).Return([]string{"error"}).Once()
		result, err := service.Modify(fundraiseItemInput, mockFile, oldFundraiseItem)
		assert.NotNil(t, result)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})

}

func TestModifyStatus(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(repository, validation, nsRequest)
	var fundraiseItemOld = fundraise.Fundraise{
		ID:          1,
		Title:       "Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam gravida, turpis consequat malesuada luctus, nisl dolor dignissim justo, a molestie sem massa et nulla. Duis diam ligula, iaculis lacinia iaculis sed, finibus eu urna. Mauris et auctor est. Etiam elementum tortor ac velit porttitor semper. Pellentesque habitant morbi tristique senectus.",
		Photo:       "https://googleapis.com/awdadwd",
		Target:      500000000,
		StartDate:   time.Date(2023, time.December, 19, 15, 30, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.December, 23, 15, 30, 0, 0, time.UTC),
		Status:      "rejected",
		UserID:      1,
		CreatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	var fundraiseItemNew = fundraise.Fundraise{
		ID:          1,
		Title:       "Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam gravida, turpis consequat malesuada luctus, nisl dolor dignissim justo, a molestie sem massa et nulla. Duis diam ligula, iaculis lacinia iaculis sed, finibus eu urna. Mauris et auctor est. Etiam elementum tortor ac velit porttitor semper. Pellentesque habitant morbi tristique senectus.",
		Photo:       "https://googleapis.com/awdadwd",
		Target:      500000000,
		StartDate:   time.Date(2023, time.December, 19, 15, 30, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.December, 23, 15, 30, 0, 0, time.UTC),
		Status:      "accepted",
		UserID:      1,
		CreatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	var fundraiseItemNewRejected = fundraise.Fundraise{
		ID:             1,
		Title:          "Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid",
		Description:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam gravida, turpis consequat malesuada luctus, nisl dolor dignissim justo, a molestie sem massa et nulla. Duis diam ligula, iaculis lacinia iaculis sed, finibus eu urna. Mauris et auctor est. Etiam elementum tortor ac velit porttitor semper. Pellentesque habitant morbi tristique senectus.",
		Photo:          "https://googleapis.com/awdadwd",
		Target:         500000000,
		StartDate:      time.Date(2023, time.December, 19, 15, 30, 0, 0, time.UTC),
		EndDate:        time.Date(2024, time.December, 23, 15, 30, 0, 0, time.UTC),
		Status:         "rejected",
		RejectedReason: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		UserID:         1,
		CreatedAt:      time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:      time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	var oldFundraiseItem = dtos.FundraiseDetails{
		ID:          1,
		Title:       "Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam gravida, turpis consequat malesuada luctus, nisl dolor dignissim justo, a molestie sem massa et nulla. Duis diam ligula, iaculis lacinia iaculis sed, finibus eu urna. Mauris et auctor est. Etiam elementum tortor ac velit porttitor semper. Pellentesque habitant morbi tristique senectus.",
		Photo:       "https://googleapis.com/awdadwd",
		Target:      500000000,
		StartDate:   time.Date(2023, time.December, 19, 15, 30, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.December, 23, 15, 30, 0, 0, time.UTC),
		Status:      "pending",
		UserID:      1,
		CreatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	var fundraiseItemInputStatus = dtos.InputFundraiseStatus{
		Status:         "accepted",
		RejectedReason: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	}
	var fundraiseItemInputStatusRejected = dtos.InputFundraiseStatus{
		Status:         "rejected",
		RejectedReason: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	}
	var fundraiseItemInputStatusRejectedRequired = dtos.InputFundraiseStatus{
		Status:         "rejected",
		RejectedReason: "",
	}
	notificationMessage := "Kami ingin memberitahu bahwa pengajuan penggalangan dana Anda untuk Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid telah diterima! Terima kasih atas langkah inisiatif Anda."
	notificationTitle := "Penerimaan Penggalangan Dana"
	notificationMessageRejected := "Terima kasih sudah mengajukan penggalangan dana. Saat ini, kami belum bisa menyetujui permohonan ini.\n\nAlasan : " + fundraiseItemInputStatusRejected.RejectedReason + "\n\nTerima kasih atas partisipasinya"
	//notificationTitleRejected := fundraiseItemInputStatusRejected.RejectedReason
	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", fundraiseItemInputStatus).Return(nil).Once()
		repository.On("SelectByTitle", fundraiseItemOld.Title).Return(&fundraiseItemOld, errors.New("error")).Once()
		repository.On("Update", fundraiseItemNew).Return(nil).Once()
		repository.On("GetDeviceToken", 1).Return("vxgeeiuwtbl").Once()
		nsRequest.On("SendNotifications", "vxgeeiuwtbl", notificationTitle, notificationMessage).Return(nil).Once()
		result, _ := service.ModifyStatus(fundraiseItemInputStatus, oldFundraiseItem)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
	t.Run("Success Rejected", func(t *testing.T) {
		validation.On("ValidateRequest", fundraiseItemInputStatusRejected).Return(nil).Once()
		//repository.On("SelectByTitle", fundraiseItemOldAccepted.Title).Return(&fundraiseItemOldAccepted, errors.New("error")).Once()
		repository.On("Update", fundraiseItemNewRejected).Return(nil).Once()
		repository.On("GetDeviceToken", 1).Return("vxgeeiuwtbl").Once()
		nsRequest.On("SendNotifications", "vxgeeiuwtbl", notificationTitle, notificationMessageRejected).Return(nil).Once()
		result, _ := service.ModifyStatus(fundraiseItemInputStatusRejected, oldFundraiseItem)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
	t.Run("Success Rejected Field Required", func(t *testing.T) {
		validation.On("ValidateRequest", fundraiseItemInputStatusRejectedRequired).Return(nil).Once()
		result, _ := service.ModifyStatus(fundraiseItemInputStatusRejectedRequired, oldFundraiseItem)
		assert.NotNil(t, result)
		assert.Equal(t, "rejected_reason field is required when the status is rejected", result[0])
		repository.AssertExpectations(t)
	})
	t.Run("Error Validation", func(t *testing.T) {
		validation.On("ValidateRequest", fundraiseItemInputStatusRejected).Return([]string{"error"}).Once()
		result, _ := service.ModifyStatus(fundraiseItemInputStatusRejected, oldFundraiseItem)
		assert.NotNil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestRemove(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var nsRequest = helperMocks.NewNotificationInterface(t)
	var service = New(repository, validation, nsRequest)
	var oldFundraiseItem = dtos.FundraiseDetails{
		ID:          1,
		Title:       "Pembangunan Masjid Pembangunan Masjid Pembangunan Masjid",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam gravida, turpis consequat malesuada luctus, nisl dolor dignissim justo, a molestie sem massa et nulla. Duis diam ligula, iaculis lacinia iaculis sed, finibus eu urna. Mauris et auctor est. Etiam elementum tortor ac velit porttitor semper. Pellentesque habitant morbi tristique senectus.",
		Photo:       "https://googleapis.com/awdadwd",
		Target:      500000000,
		StartDate:   time.Date(2023, time.December, 19, 15, 30, 0, 0, time.UTC),
		EndDate:     time.Date(2024, time.December, 23, 15, 30, 0, 0, time.UTC),
		Status:      "pending",
		UserID:      1,
		CreatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	t.Run("Success", func(t *testing.T) {
		repository.On("DeleteFile", oldFundraiseItem.Photo).Return(nil).Once()
		repository.On("DeleteByID", oldFundraiseItem.ID).Return(nil).Once()
		result := service.Remove(1, oldFundraiseItem)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}
