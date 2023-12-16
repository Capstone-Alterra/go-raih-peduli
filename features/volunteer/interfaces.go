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
	SelectByTittle(title string) error
	SelectBookmarkedVacancyID(ownerID int) (map[int]string, error)
	SelectBookmarkByVacancyAndOwnerID(vacancyID, ownerID int) string
	UpdateVacancy(vacancy VolunteerVacancies) int64
	DeleteVacancyByID(vacancyID int) error
	InsertVacancy(*VolunteerVacancies) (*VolunteerVacancies, error)
	UploadFile(file multipart.File, objectName string) (string, error)
	DeleteFile(filename string) error
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
	CheckUser(userID int) bool
	FindUserInVacancy(vacancyID, userID int) bool
	SelectAllSkills() ([]dtos.Skill, error)
	GetDeviceToken(userID int) string
}

type Usecase interface {
	FindAllVacancies(page, size int, searchAndFilter dtos.SearchAndFilter, ownerID int, suffix string) ([]dtos.ResVacancy, int64)
	FindVacancyByID(vacancyID, ownerID int) *dtos.ResVacancy
	ModifyVacancy(vacancyData dtos.InputVacancy, file multipart.File, oldData dtos.ResVacancy) ([]string, error)
	ModifyVacancyStatus(input dtos.StatusVacancies, oldData dtos.ResVacancy) (error, []string)
	RemoveVacancy(vacancyID int, oldData dtos.ResVacancy) error
	CreateVacancy(newVacancy dtos.InputVacancy, UserID int, file multipart.File) (*dtos.ResVacancy, []string, error)
	RegisterVacancy(newApply dtos.ApplyVacancy, userID int) (bool, []string)
	UpdateStatusRegistrar(input dtos.StatusRegistrar, registrarID int) (bool, []string)
	FindAllVolunteersByVacancyID(page, size int, vacancyID int, name string) ([]dtos.ResRegistrantVacancy, int64)
	FindDetailVolunteers(vacancyID, volunteerID int) *dtos.ResRegistrantVacancy
	CheckUser(userID int) bool
	FindUserInVacancy(vacancyID, userID int) bool
	FindAllSkills() ([]dtos.Skill, error)
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
	GetSkills() echo.HandlerFunc
}
