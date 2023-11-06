package repository

import (
	"raihpeduli/features/admin"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) admin.Repository {
	return &model{
		db: db,
	}
}

func (mdl *model) Paginate(page, size int) []admin.Admin {
	var admins []admin.Admin

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&admins)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return admins
}

func (mdl *model) InsertAdmin(newAdmin *admin.Admin) (*admin.Admin, error) {
	result := mdl.db.Table("admins").Create(newAdmin)

	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}

	return newAdmin, nil
}

func (mdl *model) InsertUser(newUser *admin.User) (*admin.User, error) {
	result := mdl.db.Table("users").Create(newUser)

	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}

	return newUser, nil
}

func (mdl *model) SelectByID(adminID int) *admin.Admin {
	var admin admin.Admin
	result := mdl.db.First(&admin, adminID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &admin
}

func (mdl *model) Update(admin admin.Admin) int64 {
	result := mdl.db.Save(&admin)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteByID(adminID int) int64 {
	result := mdl.db.Delete(&admin.Admin{}, adminID)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}
