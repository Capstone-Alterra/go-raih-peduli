package customer

import (
	"raihpeduli/features/customer/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []Customer
	InsertCustomer(newCustomer *Customer) (*Customer, error)
	InsertUser(newUser *User) (*User, error)
	SelectByID(customerID int) *Customer
	Update(customer Customer) int64
	DeleteByID(customerID int) int64
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResCustomer
	FindByID(customerID int) *dtos.ResCustomer
	Create(newCustomer dtos.InputCustomer) (*dtos.ResCustomer, error)
	Modify(customerData dtos.InputCustomer, customerID int) bool
	Remove(customerID int) bool
}

type Handler interface {
	GetCustomers() echo.HandlerFunc
	CustomerDetails() echo.HandlerFunc
	CreateCustomer() echo.HandlerFunc
	UpdateCustomer() echo.HandlerFunc
	DeleteCustomer() echo.HandlerFunc
}
