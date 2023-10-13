package repository

import (
	"raihpeduli/features/fundraise"

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