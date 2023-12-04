package usecase

import (
	"errors"
	"mime/multipart"
	// "os"
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
	repo := new(mocks.Repository)
	service := New(repo, nil)

	t.Run("TestFindAll_Success", func(t *testing.T) {
		expectedEntities := []news.News{}
		pagination := dtos.Pagination{Page: 1, PageSize: 10}
		searchAndFilter := dtos.SearchAndFilter{}
		ownerID := 1

		repo.On("Paginate", pagination, searchAndFilter).Return(expectedEntities, nil)
		repo.On("SelectBookmarkedNewsID", ownerID).Return(map[int]string{}, nil)
		repo.On("GetTotalData").Return(int64(10))

		result, total := service.FindAll(pagination, searchAndFilter, ownerID)

		assert.NotNil(t, result)
		assert.Equal(t, int64(10), total)
	})

	t.Run("TestFindAll_Error", func(t *testing.T) {
		pagination := dtos.Pagination{}
		searchAndFilter := dtos.SearchAndFilter{}
		ownerID := 1

		repo.On("Paginate", pagination, searchAndFilter).Return(nil, errors.New("error"))
		repo.On("SelectBookmarkedNewsID", ownerID).Return(nil, errors.New("error"))

		result, total := service.FindAll(pagination, searchAndFilter, ownerID)

		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
	})
}

func TestFindByID(t *testing.T) {
	repo := new(mocks.Repository)
	service := New(repo, nil)

	t.Run("TestFindByID_Success", func(t *testing.T) {
		expectedNews := &news.News{}
		newsID := 1
		ownerID := 1

		repo.On("SelectByID", newsID).Return(expectedNews, nil)
		repo.On("SelectBoockmarkByNewsAndOwnerID", newsID, ownerID).Return("bookmarkID", nil)

		result := service.FindByID(newsID, ownerID)

		assert.NotNil(t, result)
	})

	t.Run("TestFindByID_Error", func(t *testing.T) {
		newsID := 1
		ownerID := 1

		repo.On("SelectByID", newsID).Return(nil, errors.New("error"))
		repo.On("SelectBoockmarkByNewsAndOwnerID", newsID, ownerID).Return("", errors.New("error"))

		result := service.FindByID(newsID, ownerID)

		assert.Nil(t, result)
	})
}
// func TestModify(t *testing.T) {
//     // Create mocks
//     repoMock := &mocks.Repository{}
//     validationMock := &helperMocks.ValidationInterface{}

//     // Create service instance
//     svc := New(repoMock, validationMock)

//     // Set up input data
//     newsData := dtos.InputNews{
//         Title: "Test Title",
//         Description: "Test description",
//     }
//     file, _ := os.Open("test.jpg")
//     oldData := dtos.ResNews{
//         ID:    1,
//         Photo: "https://storage.googleapis.com/test-bucket/news/test.jpg",
//     }

//     // Set up expected output
//     expectedErr := errors.New("error")

//     // Set up mock expectations
//     validationMock.On("ValidateInput", newsData, file).Return([]string{}, nil)
//     repoMock.On("DeleteFile", "test.jpg").Return(nil)
//     repoMock.On("UploadFile", file).Return("https://storage.googleapis.com/test-bucket/news/new-test.jpg", nil)
//     repoMock.On("Update", mocks.NewRepository("*news.News")).Return(expectedErr)

//     // Call function and check output
//     errorList, err := svc.Modify(newsData, file, oldData)
//     assert.Equal(t, []string{}, errorList)
//     assert.Equal(t, expectedErr, err)

//     // Assert mock expectations
//     validationMock.AssertExpectations(t)
//     repoMock.AssertExpectations(t)
// }


