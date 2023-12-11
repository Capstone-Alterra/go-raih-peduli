package usecase

import (
	"errors"
	"raihpeduli/features/history"
	"raihpeduli/features/history/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindAllHistoryFundraiseCreatedByUser(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var service = New(repository)

	var entities = []history.Fundraise{
		{
			ID:             1,
			Title:          "Pembangunan Masjid",
			Description:    "Pembangunan Masjid Sukamaju",
			Photo:          "https://storage.googleapis.com/raih-peduli/fundraises/c",
			Target:         50000000,
			StartDate:      time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC),
			EndDate:        time.Date(2023, time.May, 19, 15, 30, 0, 0, time.UTC),
			Status:         "accepted",
			RejectedReason: "",
			UserID:         1,
			CreatedAt:      time.Date(2023, time.April, 15, 15, 30, 0, 0, time.UTC),
			UpdatedAt:      time.Date(2023, time.April, 15, 15, 30, 0, 0, time.UTC),
		},
	}

	var bookmarkIDs = map[int]string{
		1: "Aasodjpoawkpdow12km",
	}

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("HistoryFundraiseCreatedByUser", userID).Return(entities, nil).Once()
		repository.On("SelectBookmarkedFundraiseID", userID).Return(bookmarkIDs, nil).Once()
		repository.On("TotalFundAcquired", entities[0].ID).Return(int32(100), nil).Once()

		res, err := service.FindAllHistoryFundraiseCreatedByUser(userID)
		assert.Equal(t, res[0].ID, entities[0].ID)
		assert.Equal(t, res[0].FundAcquired, int32(100))
		assert.Equal(t, *res[0].BookmarkID, bookmarkIDs[1])
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Select Fund", func(t *testing.T) {
		repository.On("HistoryFundraiseCreatedByUser", userID).Return(entities, nil).Once()
		repository.On("SelectBookmarkedFundraiseID", userID).Return(bookmarkIDs, nil).Once()
		repository.On("TotalFundAcquired", entities[0].ID).Return(int32(0), errors.New("error when select")).Once()

		res, err := service.FindAllHistoryFundraiseCreatedByUser(userID)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Select Bookmarked Fundraise", func(t *testing.T) {
		repository.On("HistoryFundraiseCreatedByUser", userID).Return(entities, nil).Once()
		repository.On("SelectBookmarkedFundraiseID", userID).Return(nil, errors.New("error when select")).Once()

		res, err := service.FindAllHistoryFundraiseCreatedByUser(userID)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
	t.Run("Failed : Error When Select History", func(t *testing.T) {
		repository.On("HistoryFundraiseCreatedByUser", userID).Return(nil, errors.New("error when select")).Once()

		res, err := service.FindAllHistoryFundraiseCreatedByUser(userID)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
}

func TestFindAllHistoryVolunteerVacanciesCreatedByUser(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var service = New(repository)

	var entities = []history.VolunteerVacancies{
		{
			ID:                  1,
			Title:               "Pembangunan Masjid",
			Description:         "Pembangunan Masjid Sukamaju",
			Photo:               "https://storage.googleapis.com/raih-peduli/fundraises/c",
			ApplicationDeadline: time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC),
			Province:            "Jawa Tengah",
			City:                "Jakarta",
			SubDistrict:         "Cipedak",
			DetailLocation:      "jalan raisun no 53",
			NumberOfVacancies:   100,
			ContactEmail:        "sorosan@example.com",
			SkillsRequired:      "pendidikan,dapur,",
			Status:              "accepted",
			RejectedReason:      "",
			UserID:              1,
			CreatedAt:           time.Date(2023, time.April, 15, 15, 30, 0, 0, time.UTC),
			UpdatedAt:           time.Date(2023, time.April, 15, 15, 30, 0, 0, time.UTC),
		},
	}

	var bookmarkIDs = map[int]string{
		1: "Aasodjpoawkpdow12km",
	}

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("HistoryVolunteerVacanciesCreatedByUser", userID).Return(entities, nil).Once()
		repository.On("SelectBookmarkedVacancyID", userID).Return(bookmarkIDs, nil).Once()
		repository.On("GetTotalVolunteersByVacancyID", entities[0].ID).Return(int64(40), nil).Once()

		res, err := service.FindAllHistoryVolunteerVacanciesCreatedByUser(userID)
		assert.Equal(t, res[0].ID, entities[0].ID)
		assert.Equal(t, *res[0].BookmarkID, bookmarkIDs[1])
		assert.Equal(t, res[0].TotalRegistrar, 40)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Select Total Volunteer", func(t *testing.T) {
		repository.On("HistoryVolunteerVacanciesCreatedByUser", userID).Return(entities, nil).Once()
		repository.On("SelectBookmarkedVacancyID", userID).Return(bookmarkIDs, nil).Once()
		repository.On("GetTotalVolunteersByVacancyID", entities[0].ID).Return(int64(0), errors.New("error when select")).Once()

		res, err := service.FindAllHistoryVolunteerVacanciesCreatedByUser(userID)
		assert.Equal(t, res[0].ID, entities[0].ID)
		assert.Equal(t, *res[0].BookmarkID, bookmarkIDs[1])
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Select Bookmarked Vacancy", func(t *testing.T) {
		repository.On("HistoryVolunteerVacanciesCreatedByUser", userID).Return(entities, nil).Once()
		repository.On("SelectBookmarkedVacancyID", userID).Return(nil, errors.New("error when select")).Once()

		res, err := service.FindAllHistoryVolunteerVacanciesCreatedByUser(userID)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Select History", func(t *testing.T) {
		repository.On("HistoryVolunteerVacanciesCreatedByUser", userID).Return(nil, errors.New("error when select")).Once()

		res, err := service.FindAllHistoryVolunteerVacanciesCreatedByUser(userID)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
}

func TestFindAllHistoryVolunteerVacanciewsRegisterByUser(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var service = New(repository)

	var entities = []history.Volunteer{
		{
			ID:          1,
			Email:       "johndoe@gmail.com",
			Fullname:    "John Doe",
			Address:     "Jepang, Jepang Utara",
			PhoneNumber: "xxxxxxxxxxxxxxx",
			Gender:      "Male",
			Skills:      "bisa terbang",
			Nik:         "xxxxxxxxxxxxxxxx",
			Resume:      "Resume",
			Reason:      "Mau Flexing",
			Photo:       "https://storage.googleapis.com/raih-peduli/fundraises/c",
			Status:      "accepted",
		},
	}

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("HistoryVolunteerVacanciesRegisterByUser", userID).Return(entities, nil).Once()

		res, err := service.FindAllHistoryVolunteerVacanciesRegisterByUser(userID)
		assert.Equal(t, res[0].ID, entities[0].ID)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Select History", func(t *testing.T) {
		repository.On("HistoryVolunteerVacanciesRegisterByUser", userID).Return(nil, errors.New("error when select")).Once()

		res, err := service.FindAllHistoryVolunteerVacanciesRegisterByUser(userID)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
}

func TestFindAllHistoryUserTransaction(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var service = New(repository)

	var entities = []history.Transaction{
		{
			ID:             1,
			FundraiseID:    1,
			UserID:         1,
			PaymentType:    "Credit Card",
			Amount:         100000,
			Status:         "success",
			PaidAt:         "xxxxxxxxxxxxxxxxxxx",
			VirtualAccount: "xxxxxxxxxxxxxxxxx",
			UrlCallback:    "google.com",
			ValidUntil:     "2 Hari",
		},
	}

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("HistoryUserTransaction", userID).Return(entities, nil).Once()

		res, err := service.FindAllHistoryUserTransaction(userID)
		assert.Equal(t, res[0].ID, entities[0].ID)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repository.On("HistoryUserTransaction", userID).Return(nil, errors.New("error when select")).Once()

		res, err := service.FindAllHistoryUserTransaction(userID)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
}
