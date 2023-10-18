package _blueprint

import (
	"raihpeduli/features/_blueprint/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []Placeholder
	Insert(newPlaceholder Placeholder) int64
	SelectByID(placeholderID int) *Placeholder
	Update(placeholder Placeholder) int64
	DeleteByID(placeholderID int) int64
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResPlaceholder
	FindByID(placeholderID int) *dtos.ResPlaceholder
	Create(newPlaceholder dtos.InputPlaceholder) *dtos.ResPlaceholder
	Modify(placeholderData dtos.InputPlaceholder, placeholderID int) bool
	Remove(placeholderID int) bool
}

type Handler interface {
	GetPlaceholders() echo.HandlerFunc
	PlaceholderDetails() echo.HandlerFunc
	CreatePlaceholder() echo.HandlerFunc
	UpdatePlaceholder() echo.HandlerFunc
	DeletePlaceholder() echo.HandlerFunc
}
