package repository

import (
	"raihpeduli/features/volunteer"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) volunteer.Repository {
	return &model{
		db: db,
	}
}

func (mdl *model) Paginate(page, size int, skill string) []volunteer.VolunteerVacancies {
	var volunteers []volunteer.VolunteerVacancies

	offset := (page - 1) * size

	result := mdl.db.Where("skills_required LIKE ?", "%"+skill+"%").Offset(offset).Limit(size).Find(&volunteers)

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
	result := mdl.db.Save(&volunteer)

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
