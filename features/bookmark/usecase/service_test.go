package usecase

import (
	"errors"
	"raihpeduli/features/bookmark"
	"raihpeduli/features/bookmark/dtos"
	"raihpeduli/features/bookmark/mocks"
	helperMocks "raihpeduli/helpers/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFindAll(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, validation)

	var bookmarked_news = []dtos.ResNews{}
	var bookmarked_vacancies = []dtos.ResVolunteerVacancy{}
	var bookmarked_fundraises = []dtos.ResFundraise{
		{
			BookmarkID: "09123890awdaw0192",
			ID: 1,
			Title: "Pembangunan Masjid",
			Description: "Lorem ipsum",
			Photo: "https://googleapis.com/awdadwd",
			Target: 500000000,
			StartDate: time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC), 
			EndDate: time.Date(2024, time.April, 19, 15, 30, 0, 0, time.UTC),
		},
	}

	var bookmarked_posts = dtos.ResBookmark{
		Fundraise: bookmarked_fundraises,
		News: bookmarked_news,
		Vacancy: bookmarked_vacancies,
	}

	var size = 10
	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("Paginate", size, userID).Return(&bookmarked_posts, nil).Once()

		result := service.FindAll(size, userID)
		assert.NotNil(t, result)
		assert.Equal(t, result.Fundraise[0].BookmarkID, bookmarked_fundraises[0].BookmarkID)
		repository.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repository.On("Paginate", size, userID).Return(nil, errors.New("error decoding")).Once()

		result := service.FindAll(size, userID)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestFindByID(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, validation)

	var bookmarked_fundraise = bson.M{
		"_id": "paowdpa1l2knj1i2kj3",
		"post_id": 1,
		"post_type": "fundraise",
		"title": "Pembangunan Masjid",
		"description": "Lorem Ipsum",
		"photo": "https://googleapis.com/awdadwd",
		"target": 500000000,
		"start_date": time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC), 
		"end_date": time.Date(2024, time.April, 19, 15, 30, 0, 0, time.UTC),
		"status": "live",
	}

	var bookmarkID = "paowdpa1l2knj1i2kj3"

	t.Run("Success", func(t *testing.T) {
		repository.On("SelectByID", bookmarkID).Return(&bookmarked_fundraise, nil).Once()

		result := service.FindByID(bookmarkID)
		assert.NotNil(t, result)
		assert.Equal(t, (*result)["_id"], bookmarked_fundraise["_id"])
		repository.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repository.On("SelectByID", bookmarkID).Return(nil, errors.New("not found")).Once()

		result := service.FindByID(bookmarkID)
		assert.Nil(t, result)
		repository.AssertExpectations(t)
	})
}

func TestSetBookmark(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, validation)

	var input = dtos.InputBookmarkPost{
		PostID: 1,
		PostType: "fundraise",
	}

	var fundraise = bookmark.Fundraise{
		ID: 1,
		Title: "Pembangunan Masjid",
		Description: "Lorem ipsum",
		Photo: "https://googleapis.com/awdadwd",
		Target: 500000000,
		StartDate: time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC), 
		EndDate: time.Date(2024, time.April, 19, 15, 30, 0, 0, time.UTC),
		Status: "live",
		UserID: 1,
	}

	var document = &bookmark.FundraiseBookmark{
		PostID: 1,
		Title: "Pembangunan Masjid",
		Description: "Lorem ipsum",
		Photo: "https://googleapis.com/awdadwd",
		Target: 500000000,
		StartDate: time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC), 
		EndDate: time.Date(2024, time.April, 19, 15, 30, 0, 0, time.UTC),
		Status: "live",
		PostType: "fundraise",
		OwnerID: 1,
	}

	var bookmarked_fundraise = bson.M{
		"_id": "paowdpa1l2knj1i2kj3",
		"post_id": 1,
		"post_type": "fundraise",
		"title": "Pembangunan Masjid",
		"description": "Lorem Ipsum",
		"photo": "https://googleapis.com/awdadwd",
		"target": 500000000,
		"start_date": time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC), 
		"end_date": time.Date(2024, time.April, 19, 15, 30, 0, 0, time.UTC),
		"status": "live",
	}

	var postID = 1
	var ownerID = 1
	var postType = "fundraise"

	var errorValidation = []string{"PostID is required"}

	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByPostAndOwnerID", postID, ownerID, postType).Return(nil, errors.New("not found")).Once()
		repository.On("SelectFundraiseByID", postID).Return(&fundraise, nil).Once()
		repository.On("Insert", document).Return(true, nil).Once()

		result, errMap, err := service.SetBookmark(input, ownerID)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		assert.True(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Insert", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByPostAndOwnerID", postID, ownerID, postType).Return(nil, errors.New("not found")).Once()
		repository.On("SelectFundraiseByID", postID).Return(&fundraise, nil).Once()
		repository.On("Insert", document).Return(false, errors.New("error when insert")).Once()

		result, errMap, err := service.SetBookmark(input, ownerID)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		assert.False(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error Fundraise Not Found", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByPostAndOwnerID", postID, ownerID, postType).Return(nil, errors.New("not found")).Once()
		repository.On("SelectFundraiseByID", postID).Return(nil, errors.New("not found")).Once()

		result, errMap, err := service.SetBookmark(input, ownerID)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		assert.False(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Unknown Post Type", func(t *testing.T) {
		input.PostType = "unknown_post_type"
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByPostAndOwnerID", postID, ownerID, input.PostType).Return(nil, errors.New("not found")).Once()

		result, errMap, err := service.SetBookmark(input, ownerID)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		assert.False(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Post Already Bookmarked By This User", func(t *testing.T) {
		input.PostType = "fundraise"
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("SelectByPostAndOwnerID", postID, ownerID, postType).Return(&bookmarked_fundraise, nil).Once()

		result, errMap, err := service.SetBookmark(input, ownerID)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		assert.False(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Validate Request", func(t *testing.T) {
		input.PostID = 0
		validation.On("ValidateRequest", input).Return(errorValidation).Once()

		result, errMap, err := service.SetBookmark(input, ownerID)
		assert.NotNil(t, errMap)
		assert.NotNil(t, err)
		assert.False(t, result)
		repository.AssertExpectations(t)
	})
}

func TestUnsetBookmark(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, validation)

	var bookmarked_fundraise = bson.M{
		"_id": "awdopjakwpdopok123opk",
		"post_id": 1,
		"post_type": "fundraise",
		"title": "Pembangunan Masjid",
		"description": "Lorem Ipsum",
		"photo": "https://googleapis.com/awdadwd",
		"target": 500000000,
		"start_date": time.Date(2023, time.April, 19, 15, 30, 0, 0, time.UTC), 
		"end_date": time.Date(2024, time.April, 19, 15, 30, 0, 0, time.UTC),
		"status": "live",
		"owner_id": 1,
	}

	var bookmarkID = "awdopjakwpdopok123opk"
	var ownerID = 1


	t.Run("Success", func(t *testing.T) {
		repository.On("DeleteByID", bookmarkID).Return(1, nil).Once()

		result, err := service.UnsetBookmark(bookmarkID, &bookmarked_fundraise, ownerID)
		assert.Nil(t, err)
		assert.True(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Delete", func(t *testing.T) {
		repository.On("DeleteByID", bookmarkID).Return(0, errors.New("error when delete")).Once()

		result, err := service.UnsetBookmark(bookmarkID, &bookmarked_fundraise, ownerID)
		assert.NotNil(t, err)
		assert.False(t, result)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Owner / UserID Mismatch", func(t *testing.T) {
		result, err := service.UnsetBookmark(bookmarkID, &bookmarked_fundraise, 2)
		assert.NotNil(t, err)
		assert.False(t, result)
	})
}