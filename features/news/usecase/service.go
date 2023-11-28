package usecase

import (
	"errors"
	"fmt"
	"mime/multipart"
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

	if err != nil {
		logrus.Error(err)
		return nil, 0
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
		bookmarkID, err = svc.model.SelectBoockmarkByNewsAndOwnerID(newsID, ownerID)

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
	if errMap := svc.validation.ValidateRequest(newNews); errMap != nil {
		return nil, errMap, errors.New("error")
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
	err := smapping.FillStruct(&news, smapping.MapFields(newNews))
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	if _, err := svc.model.Insert(news); err != nil {
		return nil, nil, err
	}

	resNews := dtos.ResNews{}
	resNews.UserID = userID
	resNews.Photo = url
	if err := smapping.FillStruct(&resNews, smapping.MapFields(newNews)); err != nil {
		return nil, nil, err
	}

	return &resNews, nil, nil
}

func (svc *service) Modify(newsData dtos.InputNews, file multipart.File, oldData dtos.ResNews) ([]string, error) {
	if errMap := svc.validation.ValidateRequest(newsData); errMap != nil {
		return errMap, errors.New("error")
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

		if err := svc.model.DeleteFile(oldFilename); err != nil {
			return nil, err
		}

		imageURL, err := svc.model.UploadFile(file)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		url = imageURL
	}

	if err := smapping.FillStruct(&newNews, smapping.MapFields(newsData)); err != nil {
		logrus.Error(err)
		return nil, err
	}

	newNews.Photo = url
	newNews.ID = oldData.ID
	newNews.UserID = oldData.UserID
	err := svc.model.Update(newNews)

	if err != nil {
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
