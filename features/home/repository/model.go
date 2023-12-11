package repository

import (
	"raihpeduli/features/auth"
	"raihpeduli/features/fundraise"
	"raihpeduli/features/home"
	"raihpeduli/features/news"
	"raihpeduli/features/user"
	"raihpeduli/features/volunteer"
	"raihpeduli/helpers"

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

func (mdl *model) PaginateFundraise(page, size int, personalization []string) []fundraise.Fundraise {
	var likeQuery string = ""
	var notLikeQuery string = ""
	
	if personalization != nil || len(personalization) != 0 {
		likeQuery = helpers.BuildLikeQuery("title", personalization)
		notLikeQuery = helpers.BuildNotLikeQuery("title", personalization)
	}

	var fundraises []fundraise.Fundraise
	var additionalFundraises []fundraise.Fundraise
	offset := (page - 1) * size
	var result *gorm.DB
	if likeQuery == "webVersion" || likeQuery == "" {
		result = mdl.db.Offset(offset).Limit(size).Find(&fundraises)
	} else {
		result = mdl.db.Offset(offset).Where(likeQuery).Limit(size).Find(&fundraises)
		if len(fundraises) < size {
			additionalSize := size - len(fundraises)
			additionalResult := mdl.db.Offset(offset).Where(notLikeQuery).Limit(additionalSize).Find(&additionalFundraises)

			if additionalResult.Error != nil {
				log.Error(additionalResult.Error)
				return nil
			}
			// Gabungkan hasil tambahan dengan hasil sebelumnya
			fundraises = append(fundraises, additionalFundraises...)
		}
	}

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return fundraises
}
func (mdl *model) PaginateVolunteer(page, size int, personalization []string) []volunteer.VolunteerVacancies {
	var likeQuery string = ""
	var notLikeQuery string = ""
	
	if personalization != nil || len(personalization) != 0 {
		likeQuery = helpers.BuildLikeQuery("title", personalization)
		notLikeQuery = helpers.BuildNotLikeQuery("title", personalization)
	}

	var volunteers []volunteer.VolunteerVacancies
	var additionalVolunteers []volunteer.VolunteerVacancies
	offset := (page - 1) * size
	var result *gorm.DB
	if likeQuery == "webVersion" || likeQuery == "" {
		result = mdl.db.Offset(offset).Limit(size).Find(&volunteers)
	} else {
		result = mdl.db.Offset(offset).Where(likeQuery).Limit(size).Find(&volunteers)
		if len(volunteers) < size {
			additionalSize := size - len(volunteers)
			additionalResult := mdl.db.Offset(offset).Where(notLikeQuery).Limit(additionalSize).Find(&additionalVolunteers)

			if additionalResult.Error != nil {
				log.Error(additionalResult.Error)
				return nil
			}
			// Gabungkan hasil tambahan dengan hasil sebelumnya
			volunteers = append(volunteers, additionalVolunteers...)
		}
	}

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}
	return volunteers
}

func (mdl *model) PaginateNews(page, size int, personalization []string) []news.News {
	var likeQuery string = ""
	var notLikeQuery string = ""
	
	if personalization != nil || len(personalization) != 0 {
		likeQuery = helpers.BuildLikeQuery("title", personalization)
		notLikeQuery = helpers.BuildNotLikeQuery("title", personalization)
	}

	var newses []news.News

	var additionalNews []news.News
	offset := (page - 1) * size
	var result *gorm.DB
	if likeQuery == "webVersion" || likeQuery == "" {
		result = mdl.db.Offset(offset).Limit(size).Find(&newses)
	} else {
		result = mdl.db.Offset(offset).Where(likeQuery).Limit(size).Find(&newses)
		if len(newses) < size {
			additionalSize := size - len(newses)
			additionalResult := mdl.db.Offset(offset).Where(notLikeQuery).Limit(additionalSize).Find(&additionalNews)

			if additionalResult.Error != nil {
				log.Error(additionalResult.Error)
				return nil
			}
			// Gabungkan hasil tambahan dengan hasil sebelumnya
			newses = append(newses, additionalNews...)
		}
	}

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

func (mdl *model) SelectUserByID(userID int) *user.User {
	var userById user.User
	result := mdl.db.First(&userById, userID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &userById
}