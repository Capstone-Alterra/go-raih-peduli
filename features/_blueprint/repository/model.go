package repository

import (
	"raihpeduli/features/_blueprint"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) _blueprint.Repository {
	return &model {
		db: db,
	}
}

func (mdl *model) Paginate(page, size int) []_blueprint.Placeholder {
	var placeholders []_blueprint.Placeholder

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&placeholders)
	
	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return placeholders
}

func (mdl *model) Insert(newPlaceholder _blueprint.Placeholder) int64 {
	result := mdl.db.Create(&newPlaceholder)

	if result.Error != nil {
		log.Error(result.Error)
		return -1
	}

	return int64(newPlaceholder.ID)
}

func (mdl *model) SelectByID(placeholderID int) *_blueprint.Placeholder {
	var placeholder _blueprint.Placeholder
	result := mdl.db.First(&placeholder, placeholderID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &placeholder
}

func (mdl *model) Update(placeholder _blueprint.Placeholder) int64 {
	result := mdl.db.Save(&placeholder)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteByID(placeholderID int) int64 {
	result := mdl.db.Delete(&_blueprint.Placeholder{}, placeholderID)
	
	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}