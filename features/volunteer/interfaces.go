package volunteer

import (
	"mime/multipart"
	"raihpeduli/features/volunteer/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int, searchAndFilter dtos.SearchAndFilter) []VolunteerVacancies
	PaginateMobile(page, size int, searchAndFilter dtos.SearchAndFilter) []VolunteerVacancies
	SelectVacancyByID(vacancyID int) *VolunteerVacancies
	SelectBookmarkedVacancyID(ownerID int) (map[int]string, error)
	SelectBookmarkByVacancyAndOwnerID(vacancyID, ownerID int) string
	UpdateVacancy(vacancy VolunteerVacancies) int64
	DeleteVacancyByID(vacancyID int) int64
	InsertVacancy(*VolunteerVacancies) (*VolunteerVacancies, error)
	UploadFile(file multipart.File, objectName string) (string, error)
	GetTotalDataVacancies() int64
	GetTotalDataVacanciesMobile() int64
	GetTotalDataVacanciesBySearchAndFilter(searchAndFilter dtos.SearchAndFilter) int64
	GetTotalDataVacanciesBySearchAndFilterMobile(searchAndFilter dtos.SearchAndFilter) int64
	RegisterVacancy(registrar *VolunteerRelations) error
	UpdateStatusRegistrar(registrar VolunteerRelations) int64
	SelectRegistrarByID(registrarID int) *VolunteerRelations
	GetTotalVolunteersByVacancyID(vacancyID int) int64
	SelectVolunteersByVacancyID(vacancyID int, name string, page, size int) []Volunteer
	GetTotalVolunteers(vacancyID int, name string) int64
	SelectVolunteerDetails(vacancyID int, volunteerID int) *Volunteer
}

type Usecase interface {
	FindAllVacancies(page, size int, searchAndFilter dtos.SearchAndFilter, ownerID int, status string) ([]dtos.ResVacancy, int64)
	FindVacancyByID(vacancyID, ownerID int) *dtos.ResVacancy
	ModifyVacancy(vacancyData dtos.InputVacancy, file multipart.File, oldData dtos.ResVacancy) (bool, []string)
	ModifyVacancyStatus(input dtos.StatusVacancies, oldData dtos.ResVacancy) (bool, []string)
	RemoveVacancy(vacancyID int) bool
	CreateVacancy(newVacancy dtos.InputVacancy, UserID int, file multipart.File) (*dtos.ResVacancy, []string, error)
	RegisterVacancy(newApply dtos.ApplyVacancy, userID int, file multipart.File) (bool, []string)
	UpdateStatusRegistrar(status string, registrarID int) bool
	FindAllVolunteersByVacancyID(page, size int, vacancyID int, name string) ([]dtos.ResRegistrantVacancy, int64)
	FindDetailVolunteers(vacancyID, volunteerID int) *dtos.ResRegistrantVacancy
}

type Handler interface {
	GetVacancies(suffix string) echo.HandlerFunc
	VacancyDetails() echo.HandlerFunc
	UpdateVacancy() echo.HandlerFunc
	UpdateStatusVacancy() echo.HandlerFunc
	DeleteVacancy() echo.HandlerFunc
	CreateVacancy() echo.HandlerFunc
	ApplyVacancy() echo.HandlerFunc
	UpdateStatusRegistrar() echo.HandlerFunc
	GetVolunteersByVacancyID() echo.HandlerFunc
	GetVolunteer() echo.HandlerFunc
}
