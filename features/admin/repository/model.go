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

func (mdl *model) Insert(newAdmin *admin.Admin) *admin.Admin {
	result := mdl.db.Create(newAdmin)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return newAdmin
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

func (mdl *model) Login(email string, password string) (*admin.Admin, error) {
	var admin admin.Admin
	result := mdl.db.Where("email = ?", email).First(&admin)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error(result.Error)
		return nil, result.Error
	}
	return &admin, nil
}
