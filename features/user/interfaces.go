package user

import (
	"raihpeduli/features/user/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []User
	InsertUser(newUser *User) (*User, error)
	SelectByID(customerID int) *User
	SelectByEmail(email string) (*User, error)
	UpdateUser(user User) int64
	UpdateUserstatus(user User) error
	DeleteByID(customerID int) int64
	SendOTPByEmail(email string, otp string) error
	InsertVerification(email string, verificationKey string) error
	ValidateVerification(verificationKey string) bool
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResUser
	FindByID(customerID int) *dtos.ResUser
	Create(newUser dtos.InputUser) (*dtos.ResUser, error)
	Modify(customerData dtos.InputUser, customerID int) bool
	Remove(customerID int) bool
	InsertVerification(email string, verificationKey string) error
	ValidateVerification(verificationKey string) bool
}

type Handler interface {
	GetUsers() echo.HandlerFunc
	UserDetails() echo.HandlerFunc
	CreateUser() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
	DeleteUser() echo.HandlerFunc
	VerifyEmail() echo.HandlerFunc
}
