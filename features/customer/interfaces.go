package customer

import (
	"raihpeduli/features/customer/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []Customer
	Insert(newCustomer *Customer) *Customer
	SelectByID(customerID int) *Customer
	Update(customer Customer) int64
	DeleteByID(customerID int) int64
	Login(email string, password string) (*Customer, error)
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResCustomer
	FindByID(customerID int) *dtos.ResCustomer
	Create(newCustomer dtos.InputCustomer) *dtos.ResCustomer
	Modify(customerData dtos.InputCustomer, customerID int) bool
	Remove(customerID int) bool
	Login(email, password string) (*dtos.ResLogin, error)
}

type Handler interface {
	GetCustomers() echo.HandlerFunc
	CustomerDetails() echo.HandlerFunc
	CreateCustomer() echo.HandlerFunc
	UpdateCustomer() echo.HandlerFunc
	DeleteCustomer() echo.HandlerFunc
	LoginCustomer() echo.HandlerFunc
}
