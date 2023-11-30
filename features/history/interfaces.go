package history

import (
	"raihpeduli/features/history/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	HistoryFundraiseCreatedByUser(userID int) ([]Fundraise, error)
	TotalFundAcquired(fundraiseID int) (int32, error)
	SelectBookmarkedFundraiseID(ownerID int) (map[int]string, error)
	HistoryVolunteerVacanciesCreatedByUser(userID int) ([]VolunteerVacancies, error)
	HistoryVolunteerVacanciewsRegisterByUser(userID int) ([]VolunteerVacancies, error)
	SelectBookmarkedVacancyID(ownerID int) (map[int]string, error)
	GetTotalVolunteersByVacancyID(vacancyID int) int64
	HistoryUserTransaction(userID int) ([]Transaction, error)
}

type Usecase interface {
	FindAllHistoryFundraiseCreatedByUser(userID int) ([]dtos.ResFundraisesHistory, error)
	FindAllHistoryVolunteerVacanciesCreatedByUser(userID int) ([]dtos.ResVolunteersVacancyHistory, error)
	FindAllHistoryVolunteerVacanciewsRegisterByUser(userID int) ([]dtos.ResVolunteersVacancyHistory, error)
	FindAllHistoryUserTransaction(userID int) ([]dtos.ResTransactionHistory, error)
}

type Handler interface {
	GetHistoryFundraiseCreatedByUser() echo.HandlerFunc
	GetHistoryVolunteerVacanciesCreatedByUser() echo.HandlerFunc
	GetHistoryVolunteerVacanciewsRegisterByUser() echo.HandlerFunc
	GetHistoryUserTransaction() echo.HandlerFunc
}
