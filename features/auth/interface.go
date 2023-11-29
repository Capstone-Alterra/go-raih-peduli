package auth

import (
	"raihpeduli/features/auth/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Login(email string) (*User, error)
	Register(newUser *User) (*User, error)
	SelectByEmail(email string) (*User, error)
	SendOTPByEmail(email string, otp string) error
	InsertVerification(email string, verificationKey string) error
}

type Usecase interface {
	Login(dtos.RequestLogin) (*dtos.LoginResponse, []string, error)
	Register(newUser dtos.InputUser) (*dtos.ResUser, []string, error)
	ResendOTP(email string) bool
	RefreshJWT(jwt dtos.RefreshJWT) (*dtos.ResJWT, error)
}

type Handler interface {
	Login() echo.HandlerFunc
	RegisterUser() echo.HandlerFunc
	ResendOTP() echo.HandlerFunc
	RefreshJWT() echo.HandlerFunc
}
