package admin

import (
	"raihpeduli/features/admin/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []Admin
	InsertAdmin(newAdmin *Admin) (*Admin, error)
	InsertUser(newUser *User) (*User, error)
	SelectByID(adminID int) *Admin
	Update(admin Admin) int64
	DeleteByID(adminID int) int64
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResAdmin
	FindByID(adminID int) *dtos.ResAdmin
	Create(newAdmin dtos.InputAdmin) (*dtos.ResAdmin, error)
	Modify(adminData dtos.InputAdmin, adminID int) bool
	Remove(adminID int) bool
}

type Handler interface {
	GetAdmins() echo.HandlerFunc
	AdminDetails() echo.HandlerFunc
	CreateAdmin() echo.HandlerFunc
	UpdateAdmin() echo.HandlerFunc
	DeleteAdmin() echo.HandlerFunc
}
