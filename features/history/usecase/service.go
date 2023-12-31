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
	var entities []history.VolunteerRegistered

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
		data.Status = volunteer.Status
		data.VacancyID = volunteer.VolunteerID
		data.VacancyName = volunteer.VolunteerName
		data.VacancyPhoto = volunteer.VolunteerPhoto
		data.PostType = "registered_volunteer_vacancies"
		data.CreatedAt = volunteer.CreatedAt
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
		data.Email = donation.User.Email
		data.FundraiseName = donation.Fundraise.Title
		data.FundraisePhoto = donation.Fundraise.Photo
		data.PostType = "donations"
		data.CreatedAt = donation.CreatedAt
		switch donation.Status {
		case "2":
			data.Status = "Waiting For Payment"
		case "3":
			data.Status = "Failed / Cancelled"
		case "4":
			data.Status = "Transaction Success"
		case "5":
			data.Status = "Paid"
		default:
			data.Status = "Created"
		}
		switch donation.PaymentType {
		case "4":
			data.PaymentType = "Bank Permata"
		case "5":
			data.PaymentType = "Bank CIMB"
		case "6":
			data.PaymentType = "Bank BCA"
		case "7":
			data.PaymentType = "Bank BRI"
		case "8":
			data.PaymentType = "Bank BNI"
		case "10":
			data.PaymentType = "Gopay"
		case "11":
			data.PaymentType = "Qris"
		default:
			data.PaymentType = "Other"
		}
		donations = append(donations, data)
	}

	return donations, nil

}
