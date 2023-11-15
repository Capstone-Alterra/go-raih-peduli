package repository

import (
	"mime/multipart"
	"raihpeduli/config"
	"raihpeduli/features/fundraise"
	"raihpeduli/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type model struct {
	db        *gorm.DB
	clStorage helpers.CloudStorageInterface
}

func New(db *gorm.DB, clStorage helpers.CloudStorageInterface) fundraise.Repository {
	return &model{
		db:        db,
		clStorage: clStorage,
	}
}

func (mdl *model) Paginate(page int, size int, title string) ([]fundraise.Fundraise, error) {
	var fundraises []fundraise.Fundraise

	offset := (page - 1) * size
	titleName := "%" + title + "%"

	if err := mdl.db.Offset(offset).Limit(size).Where("title LIKE ?", titleName).Find(&fundraises).Error; err != nil {
		return nil, err
	}

	return fundraises, nil
}

func (mdl *model) Insert(newFundraise fundraise.Fundraise) (int, error) {
	if err := mdl.db.Create(&newFundraise).Error; err != nil {
		return 0, err
	}

	return newFundraise.ID, nil
}

func (mdl *model) SelectByID(fundraiseID int) (*fundraise.Fundraise, error) {
	var fundraise fundraise.Fundraise

	if err := mdl.db.First(&fundraise, fundraiseID).Error; err != nil {
		return nil, err
	}

	return &fundraise, nil
}

func (mdl *model) Update(fundraise fundraise.Fundraise) (int, error) {
	result := mdl.db.Save(&fundraise)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func (mdl *model) DeleteByID(fundraiseID int) (int, error) {
	result := mdl.db.Delete(&fundraise.Fundraise{}, fundraiseID)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
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

	return "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/fundraises/" + objectName, nil
}
