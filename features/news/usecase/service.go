package usecase

import (
	"errors"
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

func (svc *service) FindAll(page, size int, keyword string) []dtos.ResNews {
	var newss []dtos.ResNews

	newsEnt, err := svc.model.Paginate(page, size, keyword)

	if err != nil {
		log.Error(err)
		return nil
	}

	for _, news := range newsEnt {
		var data dtos.ResNews

		if err := smapping.FillStruct(&data, smapping.MapFields(news)); err != nil {
			log.Error(err.Error())
		}

		newss = append(newss, data)
	}

	return newss
}

func (svc *service) FindByID(newsID int) *dtos.ResNews {
	var res dtos.ResNews
	news, err := svc.model.SelectByID(newsID)

	if err != nil {
		logrus.Error(err)
		return nil
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
		imgURL, err := svc.model.UploadFile(file, "")

		if err != nil {
			return nil, nil, err
		}

		url = imgURL
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

func (svc *service) Modify(newsData dtos.InputNews, file multipart.File, oldData dtos.ResNews) bool {
	var newNews news.News
	var url string = ""
	var config = config.LoadCloudStorageConfig()
	var oldFilename string = oldData.Photo
	var urlLength int = len("https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/news/")

	if file != nil {
		if len(oldFilename) > urlLength {
			oldFilename = oldFilename[urlLength:]
		}
		imageURL, err := svc.model.UploadFile(file, oldFilename)

		if err != nil {
			logrus.Error(err)
			return false
		}

		url = imageURL
	}

	if err := smapping.FillStruct(&newNews, smapping.MapFields(newsData)); err != nil {
		logrus.Error(err)
		return false
	}

	newNews.Photo = url
	newNews.ID = oldData.ID
	newNews.UserID = oldData.UserID
	_, err := svc.model.Update(newNews)

	if err != nil {
		logrus.Error(err)
		return false
	}

	return true

}

func (svc *service) Remove(newsID int) bool {
	_, err := svc.model.DeleteByID(newsID)

	if err != nil {
		logrus.Error(err)
		return false
	}

	return true
}
