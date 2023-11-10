package repository

import (
	"raihpeduli/features/news"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) news.Repository {
	return &model {
		db: db,
	}
}

func (mdl *model) Paginate(page, size int) []news.News {
	var newss []news.News

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&newss)
	
	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return newss
}

func (mdl *model) Insert(newNews news.News) int64 {
	result := mdl.db.Create(&newNews)

	if result.Error != nil {
		log.Error(result.Error)
		return -1
	}

	return int64(newNews.ID)
}

func (mdl *model) SelectByID(newsID int) *news.News {
	var news news.News
	result := mdl.db.First(&news, newsID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &news
}

func (mdl *model) Update(news news.News) int64 {
	result := mdl.db.Save(&news)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteByID(newsID int) int64 {
	result := mdl.db.Delete(&news.News{}, newsID)
	
	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}