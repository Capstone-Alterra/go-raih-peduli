package repository

import (
	"raihpeduli/features/fundraise"

	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) fundraise.Repository {
	return &model {
		db: db,
	}
}

func (mdl *model) Paginate(page, size int) ([]fundraise.Fundraise, error) {
	var fundraises []fundraise.Fundraise

	offset := (page - 1) * size

	if err := mdl.db.Offset(offset).Limit(size).Find(&fundraises).Error; err != nil {
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