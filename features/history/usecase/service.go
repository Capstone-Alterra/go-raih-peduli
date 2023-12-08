package usecase

import (
	"raihpeduli/features/history"
	"raihpeduli/features/history/dtos"
	"strings"

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
	if err != nil {
		return nil, err
	}

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
		data.PostType = "created_fundraises"
		fundraises = append(fundraises, data)
	}
	if err != nil {
		return nil, err
	}
	return fundraises, nil
}
func (svc *service) FindAllHistoryVolunteerVacanciesCreatedByUser(userID int) ([]dtos.ResVolunteersVacancyHistory, error) {
	var volunteers []dtos.ResVolunteersVacancyHistory
	var bookmarkIDs map[int]string
	var err error
	var entities []history.VolunteerVacancies

	entities, err = svc.model.HistoryVolunteerVacanciesCreatedByUser(userID)
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

		data.ID = volunteer.ID
		data.UserID = volunteer.UserID
		data.Title = volunteer.Title
		data.Description = volunteer.Description
		data.SkillsRequired = strings.Split(volunteer.SkillsRequired, ",")
		data.NumberOfVacancies = volunteer.NumberOfVacancies
		data.ApplicationDeadline = volunteer.ApplicationDeadline
		data.ContactEmail = volunteer.ContactEmail
		data.Province = volunteer.Province
		data.City = volunteer.City
		data.SubDistrict = volunteer.SubDistrict
		data.Photo = volunteer.Photo
		data.Status = volunteer.Status
		data.RejectedReason = volunteer.RejectedReason
		data.CreatedAt = volunteer.CreatedAt
		data.UpdatedAt = volunteer.UpdatedAt
		data.DeletedAt = volunteer.DeletedAt

		if bookmarkIDs != nil {
			bookmarkID, ok := bookmarkIDs[data.ID]

			if ok {
				data.BookmarkID = &bookmarkID
			}
		}
		data.TotalRegistrar = int(svc.model.GetTotalVolunteersByVacancyID(data.ID))
		data.PostType = "created_volunteer_vacancies"
		volunteers = append(volunteers, data)
	}

	return volunteers, nil
}
func (svc *service) FindAllHistoryVolunteerVacanciesRegisterByUser(userID int) ([]dtos.ResRegistrantVacancyHistory, error) {
	var volunteers []dtos.ResRegistrantVacancyHistory
	var err error
	var entities []history.Volunteer

	entities, err = svc.model.HistoryVolunteerVacanciesRegisterByUser(userID)
	if err != nil {
		return nil, err
	}

	for _, volunteer := range entities {
		var data dtos.ResRegistrantVacancyHistory
		data.ID = volunteer.ID
		data.Email = volunteer.Email
		data.Fullname = volunteer.Fullname
		data.Address = volunteer.Address
		data.PhoneNumber = volunteer.PhoneNumber
		data.Gender = volunteer.Gender
		data.Nik = volunteer.Nik
		data.Skills = strings.Split(volunteer.Skills, ", ")
		data.Resume = volunteer.Resume
		data.Reason = volunteer.Reason
		data.Photo = volunteer.Photo
		data.Status = volunteer.Status

		data.PostType = "registered_volunteer_vacancies"
		volunteers = append(volunteers, data)
	}

	return volunteers, nil
}
func (svc *service) FindAllHistoryUserTransaction(userID int) ([]dtos.ResTransactionHistory, error) {
	donations := []dtos.ResTransactionHistory{}
	entities, err := svc.model.HistoryUserTransaction(userID)

	if err != nil {
		return nil, err
	}
	for _, donation := range entities {
		var data dtos.ResTransactionHistory
		if err := smapping.FillStruct(&data, smapping.MapFields(donation)); err != nil {
			logrus.Error(err)
		}
		data.Fullname = donation.User.Fullname
		data.Address = donation.User.Address
		data.PhoneNumber = donation.User.PhoneNumber
		data.ProfilePicture = donation.User.ProfilePicture
		data.Email = donation.User.Email
		data.PostType = "donations"
		donations = append(donations, data)
	}
	return donations, nil

}
