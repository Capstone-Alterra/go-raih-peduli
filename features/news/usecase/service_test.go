package usecase

import (
	"errors"
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

		pagination.Page = 0
		pagination.PageSize = 0

		res, total := service.FindAll(pagination, searchAndFilter, ownerID)
		assert.Equal(t, res[0].ID, entities[0].ID)
		assert.Equal(t, total, int64(1))
		repository.AssertExpectations(t)
	})

	t.Run("Success V2", func(t *testing.T) {
		searchAndFilter.Title = ""
		pagination.Page = 1
		pagination.PageSize = 10
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

	var errMap = []string{"title must be atleast 20 characters", "description must be atleast 50 characters"}

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

	t.Run("Failed : Error When Validate Request", func(t *testing.T) {
		input.Title = "test"
		input.Description = "test"
		validation.On("ValidateRequest", input).Return(errMap).Once()

		result, errMap, err := service.Create(input, userID, mockFile)
		assert.Nil(t, result)
		assert.NotNil(t, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
	})
}

func TestModify(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var service = New(repository, validation)

	mockFile, err := os.Open("file-mock.jpg")
	if err != nil {
		logrus.Error(err)
	}
	defer mockFile.Close()


	var input = dtos.InputNews{
		Title:       "Ceritanya Berita asdas jdoasijkdpo askdpoaskd osapkd asdas dasd as dasda sijdaopsijkdpoaskdopaskdp oaskdpoaksdopkasod",
		Description: "Ceritanya Descriptionka smdpasokd poaskdpoask dpoaskdpokaspokdaspodkasokdpasokdpaoskdpasokdpasokd alsdjaoisjdaopsdkpoaskdpoaskdpoas kdpoaskdpoaskdpoaskdpoaskdpoaskd",
		Photo:       mockFile,
	}

	var oldFilename = "adawdawd"

	var oldData = dtos.ResNews{
		ID:          1,
		Title:       "Ceritanya Berita asdas jdoasijkdpo askdpoaskd osapkd asdas dasd as dasda sijdaopsijkdpoaskdopaskdp oaskdpoaksdopkasod",
		Description: "Ceritanya Descriptionka smdpasokd poaskdpoask dpoaskdpokaspokdaspodkasokdpasokdpaoskdpasokdpasokd alsdjaoisjdaopsdkpoaskdpoaskdpoas kdpoaskdpoaskdpoaskdpoaskdpoaskd",
		Photo:       "https://storage.googleapis.com//news/adawdawd",
		UserID:      1,
	}

	var entity = news.News{
		ID:          1,
		Title:       "Ceritanya Berita asdas jdoasijkdpo askdpoaskd osapkd asdas dasd as dasda sijdaopsijkdpoaskdopaskdp oaskdpoaksdopkasod",
		Description: "Ceritanya Descriptionka smdpasokd poaskdpoask dpoaskdpokaspokdaspodkasokdpasokdpaoskdpasokdpasokd alsdjaoisjdaopsdkpoaskdpoaskdpoas kdpoaskdpoaskdpoaskdpoaskdpoaskd",
		Photo:       "https://storage.googleapis.com//news/asdasdasd",
		UserID:      1,
	}

	var errMap = []string{"title is required"}

	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("DeleteFile", oldFilename).Return(nil).Once()
		repository.On("UploadFile", mockFile).Return(entity.Photo, nil).Once()
		repository.On("Update", entity).Return(nil).Once()
		
		errMap, err := service.Modify(input, mockFile, oldData)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})
	
	t.Run("Failed : Error When Update", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("DeleteFile", oldFilename).Return(errors.New("error when delete")).Once()
		repository.On("UploadFile", mockFile).Return(entity.Photo, nil).Once()
		repository.On("Update", entity).Return(errors.New("error when update")).Once()

		errMap, err := service.Modify(input, mockFile, oldData)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Upload File", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("DeleteFile", oldFilename).Return(nil).Once()
		repository.On("UploadFile", mockFile).Return("", errors.New("error when upload")).Once()

		errMap, err := service.Modify(input, mockFile, oldData)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})
	
	t.Run("Failed : Error When Validate Request", func(t *testing.T) {
		input.Title = "a"
		validation.On("ValidateRequest", input).Return(errMap).Once()

		errMap, err := service.Modify(input, mockFile, oldData)
		assert.NotNil(t, errMap)
		assert.Nil(t, err)
		validation.AssertExpectations(t)
		repository.AssertExpectations(t)
	})
}

func TestDeleteNews(t *testing.T) {
	repo := mocks.NewRepository(t)
	validation := helperMocks.NewValidationInterface(t)
	service := New(repo, validation)

	var oldFilename = "adawdawd"

	var oldData = dtos.ResNews{
		ID:          1,
		Title:       "Ceritanya Berita asdas jdoasijkdpo askdpoaskd osapkd asdas dasd as dasda sijdaopsijkdpoaskdopaskdp oaskdpoaksdopkasod",
		Description: "Ceritanya Descriptionka smdpasokd poaskdpoask dpoaskdpokaspokdaspodkasokdpasokdpaoskdpasokdpasokd alsdjaoisjdaopsdkpoaskdpoaskdpoas kdpoaskdpoaskdpoaskdpoaskdpoaskd",
		Photo:       "https://storage.googleapis.com//news/adawdawd",
		UserID:      1,
	}

	var newsID = 1

	t.Run("Success", func(t *testing.T) {
		repo.On("DeleteFile", oldFilename).Return(nil).Once()
		repo.On("DeleteByID", newsID).Return(nil).Once()

		err := service.Remove(newsID, oldData)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed : Error When Delete News", func(t *testing.T) {
		repo.On("DeleteFile", oldFilename).Return(nil).Once()
		repo.On("DeleteByID", newsID).Return(errors.New("error when delete")).Once()

		err := service.Remove(newsID, oldData)
		assert.NotNil(t, err)
		repo.AssertExpectations(t)
	})
}



