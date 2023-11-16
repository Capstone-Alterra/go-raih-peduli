package repository

import (
	"raihpeduli/features/bookmark"

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

func (mdl *model) Insert(newBookmark bookmark.Bookmark) int64 {
	result := mdl.db.Create(&newBookmark)

	if result.Error != nil {
		log.Error(result.Error)
		return -1
	}

	return int64(newBookmark.ID)
}

func (mdl *model) SelectByID(bookmarkID int) *bookmark.Bookmark {
	var bookmark bookmark.Bookmark
	result := mdl.db.First(&bookmark, bookmarkID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &bookmark
}

func (mdl *model) DeleteByID(bookmarkID int) int64 {
	result := mdl.db.Delete(&bookmark.Bookmark{}, bookmarkID)
	
	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}