package repository

import (
	"mime/multipart"
	"raihpeduli/config"
	"raihpeduli/features/news"
	"raihpeduli/helpers"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db        *gorm.DB
	clStorage helpers.CloudStorageInterface
}

func New(db *gorm.DB, clStorage helpers.CloudStorageInterface) news.Repository {
	return &model{
		db:        db,
		clStorage: clStorage,
	}
}

func (mdl *model) Paginate(page, size int, keyword string) ([]news.News, error) {
	var news []news.News

	offset := (page - 1) * size
	searching := "%" + keyword + "%"

	if err := mdl.db.Offset(offset).Limit(size).Where("title LIKE ?", searching).Find(&news).Error; err != nil {
		return nil, err
	}

	return news, nil
}

func (mdl *model) Insert(newNews news.News) (int, error) {
	if err := mdl.db.Create(&newNews).Error; err != nil {
		return 0, err
	}

	return newNews.ID, nil
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

func (mdl *model) UploadFile(file multipart.File, objectName string) (string, error) {
	config := config.LoadCloudStorageConfig()
	randomChar := uuid.New().String()
	if objectName == "" {
		objectName = randomChar
	}

	if err := mdl.clStorage.UploadFile(file, objectName); err != nil {
		return "", err
	}

	return "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/news/" + objectName, nil
}
