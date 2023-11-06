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
	hash  helpers.HashInterface
}

func New(model customer.Repository, jwt helpers.JWTInterface, hash helpers.HashInterface) customer.Usecase {
	return &service{
		model: model,
		jwt:   jwt,
		hash:  hash,
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

func (svc *service) Create(newData dtos.InputCustomer) (*dtos.ResCustomer, error) {
	newUser := customer.User{}
	newCustomer := customer.Customer{}

	err := smapping.FillStruct(&newUser, smapping.MapFields(newData))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = smapping.FillStruct(&newCustomer, smapping.MapFields(newData))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	newUser.Password = svc.hash.HashPassword(newUser.Password)
	user, err := svc.model.InsertUser(&newUser)
	if user == nil {
		return nil, err
	}

	newCustomer.UserID = user.ID
	customer, err := svc.model.InsertCustomer(&newCustomer)
	if customer == nil {
		return nil, err
	}

	resCustomer := dtos.ResCustomer{}
	errRes := smapping.FillStruct(&resCustomer, smapping.MapFields(customer))
	if errRes != nil {
		log.Error(errRes)
		return nil, err
	}
	resCustomer.RoleID = user.RoleID
	resCustomer.Email = user.Email

	userID := strconv.Itoa(resCustomer.UserID)
	roleID := strconv.Itoa(resCustomer.RoleID)
	tokenData := svc.jwt.GenerateJWT(userID, roleID)

	if tokenData == nil {
		log.Error("Token process failed")
	}

	resCustomer.AccessToken = tokenData["access_token"].(string)
	resCustomer.RefreshToken = tokenData["refresh_token"].(string)

	return &resCustomer, nil
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
