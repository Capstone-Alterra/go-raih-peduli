package repository

import (
	"raihpeduli/features/auth"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) auth.Repository {
	return &model{
		db: db,
	}
}

func (mdl *model) Login(email string) (*auth.User, error) {
	var customer auth.User
	result := mdl.db.Table("users").Where("email = ?", email).First(&customer)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error(result.Error)
		return nil, result.Error
	}
	return &customer, nil
}

func (mdl *model) GetNameAdmin(id int) (string, error) {
	var fullname string

	result := mdl.db.Table("admins").Where("user_id = ?", id).Select("fullname").Scan(&fullname)
	if result.Error != nil {
		log.Error(result.Error)
		return "", result.Error
	}

	return fullname, nil
}

func (mdl *model) GetNameCustomer(id int) (string, error) {
	var fullname string

	result := mdl.db.Table("customers").Where("user_id = ?", id).Select("fullname").Scan(&fullname)
	if result.Error != nil {
		log.Error(result.Error)
		return "", result.Error
	}

	return fullname, nil
}
