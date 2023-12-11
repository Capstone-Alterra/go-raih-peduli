package usecase

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"raihpeduli/config"
	"raihpeduli/features/news"
	"raihpeduli/features/news/dtos"
	"raihpeduli/helpers"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
)

type service struct {
	model      news.Repository
	validation helpers.ValidationInterface
}

func New(model news.Repository, validation helpers.ValidationInterface) news.Usecase {
	return &service{
		model:      model,
		validation: validation,
	}
}

func (svc *service) FindAll(pagination dtos.Pagination, searchAndFilter dtos.SearchAndFilter, ownerID int) ([]dtos.ResNews, int64) {
	var newss []dtos.ResNews
	var bookmarkIDs map[int]string

	if pagination.Page == 0 {
		pagination.Page = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}

	entities, err := svc.model.Paginate(pagination, searchAndFilter)

	if err != nil {
		logrus.Error(err)
		return nil, 0
	}
	if ownerID != 0 {
		bookmarkIDs, err = svc.model.SelectBookmarkedNewsID(ownerID)
		if err != nil {
			logrus.Error(err)
			return nil, 0
		}
		fmt.Println(bookmarkIDs)
	}

	for _, news := range entities {
		var data dtos.ResNews

		if err := smapping.FillStruct(&data, smapping.MapFields(news)); err != nil {
			logrus.Error(err.Error())
		}

		if bookmarkIDs != nil {
			bookmarkID, ok := bookmarkIDs[data.ID]
			if ok {
				data.BookmarkID = &bookmarkID
			}
		}

		newss = append(newss, data)
	}

	var totalData int64 = 0

	if searchAndFilter.Title != "" {
		totalData = svc.model.GetTotalDataBySearchAndFilter(searchAndFilter)
	} else {
		totalData = svc.model.GetTotalData()
	}
	return newss, totalData
}

func (svc *service) FindByID(newsID, ownerID int) *dtos.ResNews {
	var res dtos.ResNews
	news, err := svc.model.SelectByID(newsID)

	if err != nil {
		logrus.Error(err)
		return nil
	}
	var bookmarkID string
	if ownerID != 0 {
		bookmarkID, err = svc.model.SelectBookmarkedByNewsAndOwnerID(newsID, ownerID)

		if bookmarkID != "" {
			res.BookmarkID = &bookmarkID
		}
	}
	if err := smapping.FillStruct(&res, smapping.MapFields(news)); err != nil {
		logrus.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newNews dtos.InputNews, userID int, file multipart.File) (*dtos.ResNews, []string, error) {
	if errorList, err := svc.validateInput(newNews, file); err != nil || len(errorList) > 0 {
		return nil, errorList, err
	}

	news := news.News{}
	var url string

	if file != nil {
		imgURL, err := svc.model.UploadFile(file)

		if err != nil {
			log.Error(err)
			return nil, nil, err
		}

		url = imgURL
	} else {
		config := config.LoadCloudStorageConfig()
		url = "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/news/default"
	}
	news.UserID = userID
	news.Photo = url
	news.Title = newNews.Title
	news.Description = newNews.Description

	inserted, err := svc.model.Insert(news)
	if err != nil {
		return nil, nil, err
	}

	resNews := dtos.ResNews{}
	// resNews.UserID = userID
	// resNews.Photo = url
	if err := smapping.FillStruct(&resNews, smapping.MapFields(inserted)); err != nil {
		return nil, nil, err
	}

	return &resNews, nil, nil
}

func (svc *service) Modify(newsData dtos.InputNews, file multipart.File, oldData dtos.ResNews) ([]string, error) {
	if errorList, err := svc.validateInput(newsData, file); err != nil || len(errorList) > 0 {
		return errorList, err
	}
	var newNews news.News
	var url string = ""
	var config = config.LoadCloudStorageConfig()
	var oldFilename string = oldData.Photo
	var urlLength int = len("https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/news/")

	if file != nil {
		if len(oldFilename) > urlLength {
			oldFilename = oldFilename[urlLength:]
		}

		if oldFilename != "default" {
			if err := svc.model.DeleteFile(oldFilename); err != nil {
				logrus.Error(err)
			}
		}

		imageURL, err := svc.model.UploadFile(file)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		url = imageURL
	}

	newNews.ID = oldData.ID
	newNews.Title = newsData.Title
	newNews.Description = newsData.Description
	newNews.Photo = url
	newNews.UserID = oldData.UserID

	if err := svc.model.Update(newNews); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return nil, nil

}

func (svc *service) Remove(newsID int, oldData dtos.ResNews) error {
	var config = config.LoadCloudStorageConfig()
	var oldFilename string = oldData.Photo
	var urlLength int = len("https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/news/")

	if len(oldFilename) > urlLength {
		oldFilename = oldFilename[urlLength:]
	}

	if oldFilename != "default" {
		svc.model.DeleteFile(oldFilename)
	}

	if err := svc.model.DeleteByID(newsID); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (svc *service) validateInput(input dtos.InputNews, file multipart.File) ([]string, error) {
	var errorList []string
	if errMap := svc.validation.ValidateRequest(input); errMap != nil {
		errorList = append(errorList, errMap...)
	}

	if len(input.Title) < 20 {
		errorList = append(errorList, "title must be at least 20 characters")
	}

	if len(input.Description) < 50 {
		errorList = append(errorList, "description must be at least 50 characters")
	}

	if file != nil {
		buffer := make([]byte, 512)

		if _, err := file.Read(buffer); err != nil {
			return nil, err
		}

		contentType := http.DetectContentType(buffer)
		isImage := contentType[:5] == "image"

		if !isImage {
			errorList = append(errorList, "photo file has to be an image (png, jpg, or jpeg)")
		}

		const maxFileSize = 5 * 1024 * 1024
		var fileSize int64

		buffer = make([]byte, 1024)
		for {
			n, err := file.Read(buffer)

			fileSize += int64(n)

			if err == io.EOF {
				break
			}

			if err != nil {
				errorList = append(errorList, "unknown file size")
			}
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}

		if fileSize > maxFileSize {
			errorList = append(errorList, "fize size exceeds the allowed limit (5MB)")
		}
	}

	return errorList, nil
}
