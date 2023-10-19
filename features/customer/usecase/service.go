package usecase

import (
	"raihpeduli/features/customer"
	"raihpeduli/features/customer/dtos"
	"raihpeduli/helpers"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
)

type service struct {
	model customer.Repository
	jwt   helpers.JWTInterface
}

func New(model customer.Repository, jwt helpers.JWTInterface) customer.Usecase {
	return &service{
		model: model,
		jwt:   jwt,
	}
}

func (svc *service) FindAll(page, size int) []dtos.ResCustomer {
	var customers []dtos.ResCustomer

	customersEnt := svc.model.Paginate(page, size)

	for _, customer := range customersEnt {
		var data dtos.ResCustomer

		if err := smapping.FillStruct(&data, smapping.MapFields(customer)); err != nil {
			log.Error(err.Error())
		}

		customers = append(customers, data)
	}

	return customers
}

func (svc *service) FindByID(customerID int) *dtos.ResCustomer {
	res := dtos.ResCustomer{}
	customer := svc.model.SelectByID(customerID)

	if customer == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(customer))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newCustomer dtos.InputCustomer) *dtos.ResCustomer {
	customer := customer.Customer{}

	err := smapping.FillStruct(&customer, smapping.MapFields(newCustomer))
	if err != nil {
		log.Error(err)
		return nil
	}

	customer.Password = helpers.HashPassword(customer.Password)

	customerID := svc.model.Insert(&customer)

	if customerID == nil {
		return nil
	}

	resCustomer := dtos.ResCustomer{}
	errRes := smapping.FillStruct(&resCustomer, smapping.MapFields(customer))
	if errRes != nil {
		log.Error(errRes)
		return nil
	}

	ID := strconv.Itoa(resCustomer.ID)
	tokenData := svc.jwt.GenerateJWT(ID)

	if tokenData == nil {
		log.Error("Token process failed")
	}

	resCustomer.Token = tokenData

	return &resCustomer
}

func (svc *service) Modify(customerData dtos.InputCustomer, customerID int) bool {
	newCustomer := customer.Customer{}

	err := smapping.FillStruct(&newCustomer, smapping.MapFields(customerData))
	if err != nil {
		log.Error(err)
		return false
	}

	newCustomer.ID = customerID
	rowsAffected := svc.model.Update(newCustomer)

	if rowsAffected <= 0 {
		log.Error("There is No Customer Updated!")
		return false
	}

	return true
}

func (svc *service) Remove(customerID int) bool {
	rowsAffected := svc.model.DeleteByID(customerID)

	if rowsAffected <= 0 {
		log.Error("There is No Customer Deleted!")
		return false
	}

	return true
}
