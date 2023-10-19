package repository

import (
	"raihpeduli/features/customer"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) customer.Repository {
	return &model{
		db: db,
	}
}

func (mdl *model) Paginate(page, size int) []customer.Customer {
	var customers []customer.Customer

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&customers)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return customers
}

func (mdl *model) Insert(newCustomer *customer.Customer) *customer.Customer {
	result := mdl.db.Create(newCustomer)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return newCustomer
}

func (mdl *model) SelectByID(customerID int) *customer.Customer {
	var customer customer.Customer
	result := mdl.db.First(&customer, customerID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &customer
}

func (mdl *model) Update(customer customer.Customer) int64 {
	result := mdl.db.Save(&customer)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteByID(customerID int) int64 {
	result := mdl.db.Delete(&customer.Customer{}, customerID)

	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}

func (mdl *model) Login(email string, password string) (*customer.Customer, error) {
	var customer customer.Customer
	result := mdl.db.Where("email = ?", email).First(&customer)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error(result.Error)
		return nil, result.Error
	}
	return &customer, nil
}
