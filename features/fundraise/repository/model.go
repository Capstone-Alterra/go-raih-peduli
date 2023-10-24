package repository

import (
	"fmt"
	"raihpeduli/features/fundraise"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) fundraise.Repository {
	return &model{
		db: db,
	}
}

func (mdl *model) Paginate(page, size int) []fundraise.Fundraise {
	var fundraises []fundraise.Fundraise

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&fundraises)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return fundraises
}

func (mdl *model) Insert(newFundraise fundraise.Fundraise) int64 {
	result := mdl.db.Create(&newFundraise)
	fmt.Println(newFundraise.ID)
	if result.Error != nil {
		log.Error(result.Error)
		return -1
	}

	return int64(newFundraise.ID)
}

func (mdl *model) SelectByID(fundraiseID int) *fundraise.Fundraise {
	var fundraise fundraise.Fundraise
	result := mdl.db.First(&fundraise, fundraiseID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &fundraise
}

func (mdl *model) Update(fundraise fundraise.Fundraise) int64 {
	result := mdl.db.Save(&fundraise)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteByID(fundraiseID int) int64 {
	result := mdl.db.Delete(&fundraise.Fundraise{}, fundraiseID)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}
