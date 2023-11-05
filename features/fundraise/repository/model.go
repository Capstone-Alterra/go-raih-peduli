package repository

import (
	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/dtos"

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

func (mdl *model) Paginate(page int, size int, title string) ([]dtos.ResFundraise, error) {
	var res []dtos.ResFundraise

	offset := (page - 1) * size
	titleName := "%" + title + "%"	

	if err := mdl.db.Table("fundraises").
	Joins("LEFT JOIN contents ON fundraises.id = contents.fundraise_id").
	Where("contents.title = ?", titleName).
	Offset(offset).Limit(size).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (mdl *model) SelectByID(fundraiseID int) (*dtos.ResFundraise, error) {
	var fundraise dtos.ResFundraise

	if err := mdl.db.First(&fundraise, fundraiseID).Error; err != nil {
		return nil, err
	}

	return &fundraise, nil
}

func (mdl *model) DeleteByID(fundraiseID int) (int, error) {
	result := mdl.db.Delete(&fundraise.Fundraise{}, fundraiseID)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}
