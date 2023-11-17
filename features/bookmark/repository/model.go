package repository

import (
	"raihpeduli/features/bookmark"
	"raihpeduli/features/news"
	"raihpeduli/features/volunteer"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) bookmark.Repository {
	return &model {
		db: db,
	}
}

func (mdl *model) Paginate(page, size int) []bookmark.Bookmark {
	var bookmarks []bookmark.Bookmark

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&bookmarks)
	
	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return bookmarks
}

func (mdl *model) Insert(inputBookmark any) int {
	result := mdl.db.Create(&inputBookmark)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return 1
}

// func (mdl *model) SelectByID(bookmarkID int) *bookmark.Bookmark {
// 	var bookmark bookmark.Bookmark
// 	result := mdl.db.First(&bookmark, bookmarkID)

// 	if result.Error != nil {
// 		log.Error(result.Error)
// 		return nil
// 	}

// 	return &bookmark
// }

func (mdl *model) SelectNewsByID(newsID int) (*news.News, error) {
	var news news.News

	if err := mdl.db.First(&news, newsID).Error; err != nil {
		return nil, err
	}

	return &news, nil
}

func (mdl *model) SelectVolunteerByID(volunteerID int) (*volunteer.VolunteerVacancies, error) {
	var volunteer volunteer.VolunteerVacancies

	if err := mdl.db.First(&volunteer, volunteerID).Error; err != nil {
		return nil, err
	}

	return &volunteer, nil
}

func (mdl *model) DeleteByID(bookmarkID int) int64 {
	result := mdl.db.Delete(&bookmark.Bookmark{}, bookmarkID)
	
	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}