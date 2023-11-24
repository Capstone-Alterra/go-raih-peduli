package volunteer

import (
	"mime/multipart"
	"raihpeduli/features/volunteer/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int, searchAndFilter dtos.SearchAndFilter) []VolunteerVacancies
	SelectVacancyByID(volunteerID int) *VolunteerVacancies
	UpdateVacancy(volunteer VolunteerVacancies) int64
	DeleteVacancyByID(volunteerID int) int64
	InsertVacancy(*VolunteerVacancies) (*VolunteerVacancies, error)
	UploadFile(file multipart.File, objectName string) (string, error)
	GetTotalDataVacancies() int64
	GetTotalDataVacanciesBySearchAndFilter(searchAndFilter dtos.SearchAndFilter) int64
	RegisterVacancy(registrar *VolunteerRelations) error
	UpdateStatusRegistrar(registrar VolunteerRelations) int64
	SelectRegistrarByID(registrarID int) *VolunteerRelations
	GetTotalVolunteersByVacancyID(vacancyID int) int64
	SelectVolunteersByVacancyID(vacancyID int, name string, page, size int) []Volunteer
	GetTotalVolunteers(vacancyID int, name string) int64
}

type Usecase interface {
	FindAllVacancies(page, size int, searchAndFilter dtos.SearchAndFilter) ([]dtos.ResVacancy, int64)
	FindVacancyByID(volunteerID int) *dtos.ResVacancy
	ModifyVacancy(volunteerData dtos.InputVacancy, volunteerID int) bool
	RemoveVacancy(volunteerID int) bool
	CreateVacancy(newVolunteer dtos.InputVacancy, UserID int, file multipart.File) (*dtos.ResVacancy, []string, error)
	RegisterVacancy(newApply dtos.ApplyVacancy, userID int, file multipart.File) (bool, []string)
	UpdateStatusRegistrar(status string, registrarID int) bool
	FindAllVolunteersByVacancyID(page, size int, vacancyID int, name string) ([]dtos.ResRegistrantVacancy, int64)
}

type Handler interface {
	GetVacancies() echo.HandlerFunc
	VacancyDetails() echo.HandlerFunc
	UpdateVacancy() echo.HandlerFunc
	DeleteVacancy() echo.HandlerFunc
	CreateVacancy() echo.HandlerFunc
	ApplyVacancy() echo.HandlerFunc
	UpdateStatusRegistrar() echo.HandlerFunc
	GetVolunteersByVacancyID() echo.HandlerFunc
}
