package usecase

import (
	"errors"
	"mime/multipart"
	"os"
	"raihpeduli/features/news"
	"raihpeduli/features/news/dtos"
	"raihpeduli/features/news/mocks"
	helperMocks "raihpeduli/helpers/mocks"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

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

	var bookmarkedNews = map[int]string{
		1: "awdwadwad",
	}

	var pagination = dtos.Pagination{
		Page: 1,
		PageSize: 10,
	}
	var searchAndFilter = dtos.SearchAndFilter{
		Title: "Ceritanya",
	}
	var ownerID = 1

	t.Run("Success V1", func(t *testing.T) {
		repository.On("Paginate", pagination, searchAndFilter).Return(entities, nil).Once()
		repository.On("SelectBookmarkedNewsID", ownerID).Return(bookmarkedNews, nil).Once()
		repository.On("GetTotalDataBySearchAndFilter", searchAndFilter).Return(int64(1)).Once()

		res, total := service.FindAll(pagination, searchAndFilter, ownerID)
		assert.Equal(t, res[0].ID, entities[0].ID)
		assert.Equal(t, total, int64(1))
		repository.AssertExpectations(t)
	})

	t.Run("Success V2", func(t *testing.T) {
		searchAndFilter.Title = ""
		repository.On("Paginate", pagination, searchAndFilter).Return(entities, nil).Once()
		repository.On("SelectBookmarkedNewsID", ownerID).Return(bookmarkedNews, nil).Once()
		repository.On("GetTotalData").Return(int64(1)).Once()

		res, total := service.FindAll(pagination, searchAndFilter, ownerID)
		assert.Equal(t, res[0].ID, entities[0].ID)
		assert.Equal(t, total, int64(1))
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Select Bookmarked Post IDs", func(t *testing.T) {
		repository.On("Paginate", pagination, searchAndFilter).Return(entities, nil).Once()
		repository.On("SelectBookmarkedNewsID", ownerID).Return(nil, errors.New("error when select")).Once()

		res, total := service.FindAll(pagination, searchAndFilter, ownerID)
		assert.Nil(t, res)
		assert.Zero(t, total)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Select All", func(t *testing.T) {
		repository.On("Paginate", pagination, searchAndFilter).Return(nil, errors.New("error when select")).Once()

		res, total := service.FindAll(pagination, searchAndFilter, ownerID)
		assert.Nil(t, res)
		assert.Zero(t, total)
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

	var ownerID = 1

	var newsID = 1

	var bookmarkID = "asdaopwkdowa" 

	t.Run("Success", func(t *testing.T) {
		repository.On("SelectByID", newsID).Return(&entity, nil).Once()
		repository.On("SelectBookmarkedByNewsAndOwnerID", newsID, ownerID).Return(bookmarkID, nil).Once()

		res := service.FindByID(newsID, ownerID)
		assert.Equal(t, res.ID, entity.ID)
		repository.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repository.On("SelectByID", newsID).Return(nil, errors.New("error when select")).Once()

		res := service.FindByID(newsID, ownerID)
		assert.Nil(t, res)
		repository.AssertExpectations(t)
	})
}

func TestCreateNews(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	service := New(repo, validation)

	mockFile, err := os.Open("file-mock.jpg")
	if err != nil {
		logrus.Error(err)
	}
	defer mockFile.Close()

	var input = dtos.InputNews{
		Title:       "Ceritanya Berita oaisjdopaskdopaksopkdaops aoskdpaoskdpoa sk",
		Description: "Ceritanya Description aosidjaops jdopask dpoaskd poaskdop aksdpoaskdpoaskdpaoksdpoaskdpoaskdpoaskd poakdpoaskd poaksdp okaspdkaspodkapsok dpaoskd poaskdpaoskdp aoskdpoaskd poaskdpoaskd poaskdpaoskdpasokd",
		Photo:       mockFile,
	}

	var entity = news.News{
		UserID: 1,
		Title:       "Ceritanya Berita oaisjdopaskdopaksopkdaops aoskdpaoskdpoa sk",
		Description: "Ceritanya Description aosidjaops jdopask dpoaskd poaskdop aksdpoaskdpoaskdpaoksdpoaskdpoaskdpoaskd poakdpoaskd poaksdp okaspdkaspodkapsok dpaoskd poaskdpaoskdp aoskdpoaskd poaskdpoaskd poaskdpaoskdpasokd",
		Photo: "https://storage.googleapis.com/news/alskdmaslmd",
	}
	
	var photo = "https://storage.googleapis.com/news/alskdmaslmd"

	var userID = 1

	var errMap = []string{"title is required"}

	t.Run("Success V1", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repo.On("UploadFile", mockFile).Return(photo, nil).Once()
		repo.On("Insert", entity).Return(&entity, nil).Once()

		result, errMap, err := service.Create(input, userID, mockFile)
		assert.Equal(t, result.Title, input.Title)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Failed : Error When Insert", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repo.On("UploadFile", mockFile).Return(photo, nil).Once()
		repo.On("Insert", entity).Return(nil, errors.New("error when insert")).Once()

		result, errMap, err := service.Create(input, userID, mockFile)
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Failed : Error When Upload File", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repo.On("UploadFile", mockFile).Return("", errors.New("error when upload")).Once()

		result, errMap, err := service.Create(input, userID, mockFile)
		assert.Nil(t, result)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Failed : Error When Upload File", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(errMap).Once()

		result, errMap, err := service.Create(input, userID, mockFile)
		assert.Nil(t, result)
		assert.NotNil(t, errMap)
		assert.NotNil(t, err)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Success V2", func(t *testing.T) {
		input.Photo = nil
		validation.On("ValidateRequest", input).Return(nil).Once()
		entity.Photo = "https://storage.googleapis.com//news/default"
		repo.On("Insert", entity).Return(&entity, nil).Once()

		result, errMap, err := service.Create(input, userID, nil)
		assert.Equal(t, result.Title, input.Title)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

func TestModify(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, validation)

	var file multipart.File

	var input = dtos.InputNews{
		Title:       "Ceritanya Berita",
		Description: "Ceritanya Description",
		Photo:       file,
	}

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
	// var errMap = []string{"title is required"}

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



