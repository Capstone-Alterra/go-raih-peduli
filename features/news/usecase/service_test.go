package usecase

import (
	"errors"
	"mime/multipart"
	"raihpeduli/features/news"
	"raihpeduli/features/news/dtos"
	"raihpeduli/features/news/mocks"
	helperMocks "raihpeduli/helpers/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNews(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	service := New(repo, validation)

	var file multipart.File

	var newNews = dtos.InputNews{
		Title:       "Ceritanya Berita",
		Description: "Ceritanya Description",
		Photo:       file,
	}

	var invalidData = dtos.InputNews{
		Title: "Fake News",
	}

	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", newNews).Return(nil).Once()
		repo.On("UploadFile", file, "").Return("https://storage.googleapis.com//news/News.jpg", nil).Once()

		var news news.News
		news.ID = 1
		news.Title = newNews.Title
		news.Description = newNews.Description
		news.Photo = "https://storage.googleapis.com//news/News.jpg"
		repo.On("InsertNews", &news).Return(&newNews, nil).Once()

		result, errMap, err := service.Create(newNews, 1, file)
		assert.Nil(t, err)
		assert.Nil(t, errMap)
		assert.NotNil(t, result)
		assert.Equal(t, newNews.Title, result.Title)
		repo.AssertExpectations(t)
	})

	t.Run("Failed upload file", func(t *testing.T) {
		validation.On("ValidateRequest", newNews).Return(nil).Once()
		repo.On("UploadFile", file, "").Return("", errors.New("failed")).Once()

		result, errMap, err := service.Create(newNews, 1, file)

		assert.Error(t, err)
		assert.Nil(t, errMap)
		assert.Nil(t, result)
		assert.Equal(t, "failed", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Validation error", func(t *testing.T) {
		validation.On("ValidateRequest", invalidData).Return(nil).Once()
		repo.On("UploadFile", file, "").Return("", errors.New("validation error")).Once()

		result, errMap, err := service.Create(invalidData, 1, file)
		assert.Error(t, err)
		assert.Nil(t, errMap)
		assert.Nil(t, result)
		assert.Equal(t, "validation error", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestDeleteNews(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	service := New(repo, validation)

	t.Run("Success delete news", func(t *testing.T) {
		newsID := 1
		oldData := dtos.ResNews{
			ID:          1,
			Title:       "Ceritanya Berita",
			Description: "Ceritanya Description",
			Photo:       "https://storage.googleapis.com/bucket-name/news/photo.jpg",
		}
		repo.On("DeleteFile", "/news/photo.jpg").Return(nil).Once()
		repo.On("DeleteNewsByID", newsID).Return(nil).Once()

		err := service.Remove(newsID, oldData)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Error deleteing news by ID", func(t *testing.T) {
		newsID := 1
		oldData := dtos.ResNews{
			Photo: "https://storage.googleapis.com/bucket-name/news/photo.jpg",
		}
		repo.On("DeleteFile", "/news/photo.jpg").Return(nil).Once()
		repo.On("DeleteNewsByID", newsID).Return(errors.New("some error")).Once()

		err := service.Remove(newsID, oldData)

		assert.Error(t, err)
		assert.Equal(t, "some error", err.Error())
		repo.AssertExpectations(t)
	})
}

func TestFindAll(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var service = New(repository, nil)

	var entities = []news.News{
		{
			ID:          1,
			Title:       "Ceritanya Berita",
			Description: "Ceritanya Description",
			Photo:       "https://storage.googleapis.com/bucket-name/news",
			UserID:      1,
		},
	}

	var page = 1
	var pageSize = 10

	t.Run("Success", func(t *testing.T) {
		repository.On("Paginate", page, pageSize).Return(entities).Once()
		repository.On("GetTotalData").Return(int64(1)).Once()

		res, total := service.FindAll(page, pageSize)
		assert.Equal(t, res[0].ID, entities[0].ID)
		assert.Equal(t, total, int64(1))
		repository.AssertExpectations(t)
	})
}

func TestFindByID(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var service = New(repository, nil)

	var entity = news.News{
		ID:          1,
		Title:       "Ceritanya Berita",
		Description: "Ceritanya Description",
		Photo:       "https://storage.googleapis.com/bucket-name/news",
		UserID:      1,
	}

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("FindByID", userID).Return(&entity).Once()

		res := service.FindByID(userID)
		assert.Equal(t, res.ID, entity.ID)
		repository.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repository.On("FindByID", userID).Return(nil).Once()

		res := service.FindByID(userID)
		assert.Nil(t, res)
		repository.AssertExpectations(t)
	})
}

func TestModify(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, validation)

	var input = dtos.InputNews{
		Title:       "Ceritanya Berita",
		Description: "Ceritanya Description",
		Photo:       "https://storage.googleapis.com/bucket-name/news",
	}

	var file multipart.File

	var oldData = dtos.ResNews{
		ID:          1,
		Title:       "Ceritanya Berita",
		Description: "Ceritanya Description",
		Photo:       "https://storage.googleapis.com/bucket-name/news",
		UserID:      1,
	}

	var entity = news.News{
		ID:          1,
		Title:       "Ceritanya Berita",
		Description: "Ceritanya Description",
		Photo:       "https://storage.googleapis.com/bucket-name/news",
		UserID:      1,
	}
	var errMap = []string{"title is required"}

	t.Run("Success", func(t *testing.T) {
		repository.On("UploadFile", file, oldData.Photo).Return(oldData.Photo, nil).Once()
		repository.On("UpdateNews", entity).Return(int64(1)).Once()
		
		err, errMap := service.Modify(input, file, oldData)
		assert.Nil(t, err)
		assert.Nil(t, errMap)
		repository.AssertExpectations(t)
	})
	
	t.Run("Failed : error when Update", func(t *testing.T) {
		repository.On("UploadFile", file, oldData.Photo).Return(oldData.Photo, nil).Once()
		repository.On("UpdateNews", entity).Return(int64(0)).Once()

		err, errmap := service.Modify(input, file, oldData)
		assert.NotNil(t, err)
		assert.Nil(t, errmap)
		repository.AssertExpectations(t)
	})

	t.Run("failed : Error When Upload File", func(t *testing.T) {
		repository.On("UploadFile", file, oldData.Photo).Return("",errors.New("error when upload")).Once()

		err, errMap := service.Modify(input, file, oldData)
		assert.NotNil(t, err)
		assert.Nil(t, errMap)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error title Required", func(t *testing.T) {
		input.Title = "Test News"
		err, errMap := service.Modify(input, file, oldData)
		assert.NotNil(t, err)
		assert.NotNil(t, errMap)
	})
}