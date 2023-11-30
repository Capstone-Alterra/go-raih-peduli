package usecase

import (
	"raihpeduli/features/history"
	"raihpeduli/features/history/dtos"

	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
)

type service struct {
	model history.Repository
}

func New(model history.Repository) history.Usecase {
	return &service{
		model: model,
	}
}

func (svc *service) FindAllHistoryFundraiseCreatedByUser(userID int) ([]dtos.ResFundraisesHistory, error) {
	fundraises := []dtos.ResFundraisesHistory{}
	var bookmarkIDs map[int]string

	entities, err := svc.model.HistoryFundraiseCreatedByUser(userID)

	if userID != 0 {
		bookmarkIDs, err = svc.model.SelectBookmarkedFundraiseID(userID)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	for _, fundraise := range entities {
		var data dtos.ResFundraisesHistory

		if err := smapping.FillStruct(&data, smapping.MapFields(fundraise)); err != nil {
			logrus.Error(err)
		}

		if bookmarkIDs != nil {
			bookmarkID, ok := bookmarkIDs[data.ID]

			if ok {
				data.BookmarkID = &bookmarkID
			}
		}

		if data.FundAcquired, err = svc.model.TotalFundAcquired(data.ID); err != nil {
			logrus.Error(err)
		}
		fundraises = append(fundraises, data)
	}
	if err != nil {
		return nil, err
	}
	return fundraises, nil
}
func (svc *service) FindAllHistoryVolunteerVacanciesCreatedByUser(userID int) ([]dtos.ResVolunteersVacancyHistory, error) {
	return []dtos.ResVolunteersVacancyHistory{}, nil
}
func (svc *service) FindAllHistoryVolunteerVacanciewsRegisterByUser(userID int) ([]dtos.ResVolunteersVacancyHistory, error) {
	var volunteers []dtos.ResVolunteersVacancyHistory
	var bookmarkIDs map[int]string
	var err error
	var entities []history.VolunteerVacancies

	entities, err = svc.model.HistoryVolunteerVacanciewsRegisterByUser(userID)
	if err != nil {
		return nil, err
	}
	if userID != 0 {
		bookmarkIDs, err = svc.model.SelectBookmarkedVacancyID(userID)
		if err != nil {
			return nil, err
		}
	}

	for _, volunteer := range entities {
		var data dtos.ResVolunteersVacancyHistory

		if err := smapping.FillStruct(&data, smapping.MapFields(volunteer)); err != nil {
			logrus.Error(err)
		}

		if bookmarkIDs != nil {
			bookmarkID, ok := bookmarkIDs[data.ID]

			if ok {
				data.BookmarkID = &bookmarkID
			}
		}
		data.TotalRegistrar = int(svc.model.GetTotalVolunteersByVacancyID(data.ID))
		volunteers = append(volunteers, data)
	}

	return volunteers, nil
}
func (svc *service) FindAllHistoryUserTransaction(userID int) ([]dtos.ResTransactionHistory, error) {
	return []dtos.ResTransactionHistory{}, nil

}
