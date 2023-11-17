package repository

import (
	"mime/multipart"
	"raihpeduli/config"
	"raihpeduli/features/volunteer"
	"raihpeduli/helpers"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
	clStorage helpers.CloudStorageInterface
}

func New(db *gorm.DB, clStorage helpers.CloudStorageInterface) volunteer.Repository {
	return &model{
		db: db,
		clStorage: clStorage,
	}
}

func (mdl *model) Paginate(page, size int) []volunteer.VolunteerVacancies {
	var volunteers []volunteer.VolunteerVacancies

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&volunteers)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return volunteers
}

func (mdl *model) SelectByTitle(page, size int, title string) []volunteer.VolunteerVacancies {
	var volunteers []volunteer.VolunteerVacancies

	offset := (page - 1) * size

	result := mdl.db.Where("title LIKE ?", "%"+title+"%").Offset(offset).Limit(size).Find(&volunteers)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return volunteers
}

func (mdl *model) SelectBySkill(page, size int, title string) []volunteer.VolunteerVacancies {
	var volunteers []volunteer.VolunteerVacancies

	offset := (page - 1) * size

	result := mdl.db.Where("skills_required LIKE ?", "%"+title+"%").Offset(offset).Limit(size).Find(&volunteers)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return volunteers
}

func (mdl *model) SelectByID(volunteerID int) *volunteer.VolunteerVacancies {
	var volunteer volunteer.VolunteerVacancies
	result := mdl.db.First(&volunteer, volunteerID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &volunteer
}

func (mdl *model) Update(volunteer volunteer.VolunteerVacancies) int64 {
	result := mdl.db.Updates(&volunteer)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteByID(volunteerID int) int64 {
	result := mdl.db.Delete(&volunteer.VolunteerVacancies{}, volunteerID)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}

func (mdl *model) Insert(newVolunteer *volunteer.VolunteerVacancies) (*volunteer.VolunteerVacancies, error){
	result := mdl.db.Create(newVolunteer)
	
	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}
	return newVolunteer, nil
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
