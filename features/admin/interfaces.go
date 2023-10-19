package admin

import (
	"raihpeduli/features/admin/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []Admin
	Insert(newAdmin *Admin) *Admin
	SelectByID(adminID int) *Admin
	Update(admin Admin) int64
	DeleteByID(adminID int) int64
	Login(email string, password string) (*Admin, error)
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResAdmin
	FindByID(adminID int) *dtos.ResAdmin
	Create(newAdmin dtos.InputAdmin) *dtos.ResAdmin
	Modify(adminData dtos.InputAdmin, adminID int) bool
	Remove(adminID int) bool
	Login(email, password string) (*dtos.ResLogin, error)
}

type Handler interface {
	GetAdmins() echo.HandlerFunc
	AdminDetails() echo.HandlerFunc
	CreateAdmin() echo.HandlerFunc
	UpdateAdmin() echo.HandlerFunc
	DeleteAdmin() echo.HandlerFunc
	LoginAdmin() echo.HandlerFunc
}
