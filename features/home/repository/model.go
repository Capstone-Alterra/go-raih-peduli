package repository

import (
	"raihpeduli/features/auth"
	"raihpeduli/features/fundraise"
	"raihpeduli/features/home"
	"raihpeduli/features/news"
	"raihpeduli/features/volunteer"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) home.Repository {
	return &model{
		db: db,
	}
}

func (mdl *model) PaginateFundraise(page, size int) []fundraise.Fundraise {
	var fundraises []fundraise.Fundraise

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&fundraises)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return fundraises
}
func (mdl *model) PaginateVolunteer(page, size int) []volunteer.VolunteerVacancies {
	var volunteers []volunteer.VolunteerVacancies

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&volunteers)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return volunteers
}
func (mdl *model) PaginateNews(page, size int) []news.News {
	var newses []news.News

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&newses)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return newses
}
func (mdl *model) CountUser() int {
	var totalUser int64
	result := mdl.db.Model(&auth.User{}).Where("role_id = ?", "1").Count(&totalUser)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return int(totalUser)
}
func (mdl *model) CountFundraise() int {
	var totalFundraise int64
	result := mdl.db.Model(&fundraise.Fundraise{}).Count(&totalFundraise)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return int(totalFundraise)
}
func (mdl *model) CountVolunteer() int {
	var totalVolunteer int64
	result := mdl.db.Model(&volunteer.VolunteerVacancies{}).Count(&totalVolunteer)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return int(totalVolunteer)
}
func (mdl *model) CountNews() int {
	var totalNews int64
	result := mdl.db.Model(&news.News{}).Count(&totalNews)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return int(totalNews)
}

func (mdl *model) Insert(newHome home.Home) int64 {
	result := mdl.db.Create(&newHome)

	if result.Error != nil {
		log.Error(result.Error)
		return -1
	}

	return int64(newHome.ID)
}

func (mdl *model) SelectByID(homeID int) *home.Home {
	var home home.Home
	result := mdl.db.First(&home, homeID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &home
}

func (mdl *model) Update(home home.Home) int64 {
	result := mdl.db.Save(&home)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteByID(homeID int) int64 {
	result := mdl.db.Delete(&home.Home{}, homeID)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}
