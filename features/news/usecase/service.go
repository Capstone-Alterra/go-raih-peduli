package usecase

import (
	"mime/multipart"
	"raihpeduli/features/news"
	"raihpeduli/features/news/dtos"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
)

type service struct {
	model news.Repository
}

func New(model news.Repository) news.Usecase {
	return &service{
		model: model,
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
	res := dtos.ResNews{}
	news := svc.model.SelectByID(newsID)

	if news == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(news))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newNews dtos.InputNews, userID int, file multipart.File) (*dtos.ResNews, error) {
	news := news.News{}
	var url string

	if file != nil {
		imgURL, err := svc.model.UploadFile(file, "")

		if err != nil {
			return nil, err
		}

		url = imgURL
	}
	news.UserID = userID
	news.Photo = url
	err := smapping.FillStruct(&news, smapping.MapFields(newNews))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if _, err := svc.model.Insert(news); err != nil {
		return nil, err
	}

	resNews := dtos.ResNews{}
	resNews.UserID = userID
	resNews.Photo = url
	if err := smapping.FillStruct(&resNews, smapping.MapFields(newNews)); err != nil {
		return nil, err
	}

	return &resNews, nil
}

func (svc *service) Modify(newsData dtos.InputNews, newsID int) bool {
	newNews := news.News{}

	err := smapping.FillStruct(&newNews, smapping.MapFields(newsData))
	if err != nil {
		log.Error(err)
		return false
	}

	newNews.ID = newsID
	rowsAffected := svc.model.Update(newNews)

	if rowsAffected <= 0 {
		log.Error("There is No News Updated!")
		return false
	}

	return true
}

func (svc *service) Remove(newsID int) bool {
	rowsAffected := svc.model.DeleteByID(newsID)

	if rowsAffected <= 0 {
		log.Error("There is No News Deleted!")
		return false
	}

	return true
}
