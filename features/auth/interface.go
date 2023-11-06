package auth

import (
	"raihpeduli/features/auth/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Login(email string) (*User, error)
	GetNameAdmin(id int) (string, error)
	GetNameCustomer(id int) (string, error)
}

type Usecase interface {
	Login(dtos.RequestLogin) (*dtos.LoginResponse, error)
}

type Handler interface {
	LoginCustomer() echo.HandlerFunc
}
