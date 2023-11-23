package volunteer

import (
	"mime/multipart"
	"raihpeduli/features/volunteer/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int, searchAndFilter dtos.SearchAndFilter) []VolunteerVacancies
	SelectByID(volunteerID int) *VolunteerVacancies
	Update(volunteer VolunteerVacancies) int64
	DeleteByID(volunteerID int) int64
	Insert(*VolunteerVacancies) (*VolunteerVacancies, error)
	UploadFile(file multipart.File, objectName string) (string, error)
	GetTotalDataVacancies() int64
	GetTotalDataVacanciesBySearchAndFilter(searchAndFilter dtos.SearchAndFilter) int64
	Register(registrar *VolunteerRelations) error
	UpdateStatusRegistrar(registrar VolunteerRelations) int64
	SelectRegistrarByID(registrarID int) *VolunteerRelations
	GetTotalVolunteerByVacancyID(vacancyID int) int64
	SelectVolunteerByVacancyID(vacancyID int, name string, page, size int) []Volunteer
	GetTotalVolunteer(vacancyID int, name string) int64
}

type Usecase interface {
	FindAll(page, size int, searchAndFilter dtos.SearchAndFilter) ([]dtos.ResVolunteer, int64)
	FindByID(volunteerID int) *dtos.ResVolunteer
	Modify(volunteerData dtos.InputVolunteer, volunteerID int) bool
	Remove(volunteerID int) bool
	Create(newVolunteer dtos.InputVolunteer, UserID int, file multipart.File) (*dtos.ResVolunteer, []string, error)
	Register(newApply dtos.ApplyVolunteer, userID int, file multipart.File) (bool, []string)
	UpdateStatusRegistrar(status string, registrarID int) bool
	FindAllVolunteerByVacancyID(page, size int, vacancyID int, name string) ([]dtos.ResRegistrantVacancy, int64)
}

type Handler interface {
	GetVolunteers() echo.HandlerFunc
	VolunteerDetails() echo.HandlerFunc
	UpdateVolunteer() echo.HandlerFunc
	DeleteVolunteer() echo.HandlerFunc
	CreateVolunteer() echo.HandlerFunc
	ApplyVacancies() echo.HandlerFunc
	UpdateStatusRegistrar() echo.HandlerFunc
	GetVolunteerByVacancyID() echo.HandlerFunc
}
