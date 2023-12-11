package usecase

import (
	"raihpeduli/features/fundraise"
	"raihpeduli/features/home/mocks"
	"raihpeduli/features/news"
	"raihpeduli/features/user"
	"raihpeduli/features/volunteer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var service = New(repository)

	var fundraises = []fundraise.Fundraise{
		{
			ID: 1,
			Title: "Pembangunan Mushola",
			Description: "deskripsi asdnka sijdopas dkpoaskd aopskd",
			Photo: "google.com",
			Target: 1000000,
		},
	}

	var volunteers = []volunteer.VolunteerVacancies{
		{
			ID: 1,
			Title: "Pembangunan Mushola",
			Description: "deskripsi asdnka sijdopas dkpoaskd aopskd",
			Photo: "google.com",
			NumberOfVacancies: 50,
		},
	}

	var news = []news.News{
		{
			ID: 1,
			Title: "Pembangunan Mushola",
			Description: "deskripsi asdnka sijdopas dkpoaskd aopskd",
			Photo: "google.com",
		},
	}

	var page = 1
	var pageSize = 5
	
	var personalization = []string{
		"pendidikan",
	}

	t.Run("Success", func(t *testing.T) {
		repository.On("PaginateFundraise", page, pageSize, personalization).Return(fundraises).Once()
		repository.On("PaginateVolunteer", page, pageSize, personalization).Return(volunteers).Once()
		repository.On("PaginateNews", page, pageSize, personalization).Return(news).Once()

		res := service.FindAll(page, pageSize, personalization)
		assert.Equal(t, res.Fundraise[0].ID, fundraises[0].ID)
		assert.Equal(t, res.Volunteer[0].ID, volunteers[0].ID)
		assert.Equal(t, res.News[0].ID, news[0].ID)
		repository.AssertExpectations(t)
	})
}

func TestFindAllWeb(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var service = New(repository)

	var fundraises = []fundraise.Fundraise{
		{
			ID: 1,
			Title: "Pembangunan Mushola",
			Description: "deskripsi asdnka sijdopas dkpoaskd aopskd",
			Photo: "google.com",
			Target: 1000000,
		},
	}

	var volunteers = []volunteer.VolunteerVacancies{
		{
			ID: 1,
			Title: "Pembangunan Mushola",
			Description: "deskripsi asdnka sijdopas dkpoaskd aopskd",
			Photo: "google.com",
			NumberOfVacancies: 50,
		},
	}
	
	var page = 1
	var pageSize = 5

	t.Run("Success", func(t *testing.T) {
		repository.On("PaginateFundraise", page, pageSize, []string(nil)).Return(fundraises).Once()
		repository.On("PaginateVolunteer", page, pageSize, []string(nil)).Return(volunteers).Once()
		repository.On("CountUser").Return(10).Once()
		repository.On("CountFundraise").Return(10).Once()
		repository.On("CountVolunteer").Return(10).Once()
		repository.On("CountNews").Return(10).Once()

		res := service.FindAllWeb(page, pageSize)
		assert.Equal(t, res.Fundraise[0].ID, fundraises[0].ID)
		assert.Equal(t, res.Volunteer[0].ID, volunteers[0].ID)
		assert.Equal(t, res.FundraiseAmount, 10)
		repository.AssertExpectations(t)
	})
}

func TestGetPersonalization(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var service = New(repository)

	var userID = 1

	var personalization = "koki"

	var user = user.User{
		ID: 1,
		Fullname: "sarbin",
		Email: "sarbin@example.com",
		Personalization: &personalization,
	}

	t.Run("Success", func(t *testing.T) {
		repository.On("SelectUserByID", userID).Return(&user).Once()

		res := service.GetPersonalization(userID)
		assert.Equal(t, res[0], personalization)
		repository.AssertExpectations(t)
	})

	
	t.Run("Failed", func(t *testing.T) {
		repository.On("SelectUserByID", userID).Return(nil).Once()

		res := service.GetPersonalization(userID)
		assert.Nil(t, res)
		repository.AssertExpectations(t)
	})

}