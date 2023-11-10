package usecase

import (
	"raihpeduli/features/news"
	"raihpeduli/features/news/dtos"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
)

type service struct {
	model news.Repository
}

func New(model news.Repository) news.Usecase {
	return &service {
		model: model,
	}
}

func (svc *service) FindAll(page, size int) []dtos.ResNews {
	var newss []dtos.ResNews

	newssEnt := svc.model.Paginate(page, size)

	for _, news := range newssEnt {
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

func (svc *service) Create(newNews dtos.InputNews) *dtos.ResNews {
	news := news.News{}
	
	err := smapping.FillStruct(&news, smapping.MapFields(newNews))
	if err != nil {
		log.Error(err)
		return nil
	}

	newsID := svc.model.Insert(news)

	if newsID == -1 {
		return nil
	}

	resNews := dtos.ResNews{}
	errRes := smapping.FillStruct(&resNews, smapping.MapFields(newNews))
	if errRes != nil {
		log.Error(errRes)
		return nil
	}

	return &resNews
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