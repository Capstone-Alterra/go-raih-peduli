package volunteer

import (
	"raihpeduli/features/volunteer/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int, skill string) []VolunteerVacancies
	SelectByID(volunteerID int) *VolunteerVacancies
	Update(volunteer VolunteerVacancies) int64
	DeleteByID(volunteerID int) int64
}

type Usecase interface {
	FindAll(page, size int, skill string) []dtos.ResVolunteer
	FindByID(volunteerID int) *dtos.ResVolunteer
	Modify(volunteerData dtos.InputVolunteer, volunteerID int) bool
	Remove(volunteerID int) bool
}

type Handler interface {
	GetVolunteers() echo.HandlerFunc
	VolunteerDetails() echo.HandlerFunc
	UpdateVolunteer() echo.HandlerFunc
	DeleteVolunteer() echo.HandlerFunc
}
