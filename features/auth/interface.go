package auth

import (
	"raihpeduli/features/auth/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Login(email string) (*User, error)
	Register(newUser *User) (*User, error)
	GetNameAdmin(id int) (string, error)
	GetNameCustomer(id int) (string, error)
	SelectByEmail(email string) (*User, error)
	SendOTPByEmail(email string, otp string) error
	InsertVerification(email string, verificationKey string) error
}

type Usecase interface {
	Login(dtos.RequestLogin) (*dtos.LoginResponse, error)
	Register(newUser dtos.InputUser) (*dtos.ResUser, error)
}

type Handler interface {
	LoginCustomer() echo.HandlerFunc
	RegisterUser() echo.HandlerFunc
}
