package usecase

import (
	"errors"
	"math"
	"mime/multipart"
	"raihpeduli/config"
	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/dtos"
	"raihpeduli/helpers"

	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
)

type service struct {
	model      fundraise.Repository
	validation helpers.ValidationInterface
}

func New(model fundraise.Repository, validation helpers.ValidationInterface) fundraise.Usecase {
	return &service{
		model:      model,
		validation: validation,
	}
}

func (svc *service) FindAll(pagination dtos.Pagination, searchAndFilter dtos.SearchAndFilter, ownerID int, suffix string) ([]dtos.ResFundraise, int64) {
	var fundraises []dtos.ResFundraise
	var bookmarkIDs map[int]string

	if pagination.Page == 0 {
		pagination.Page = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}

	if searchAndFilter.MaxTarget == 0 {
		searchAndFilter.MaxTarget = math.MaxInt32
	}

	var entities []fundraise.Fundraise
	var err error

	if suffix == "mobile" {
		entities, err = svc.model.PaginateMobile(pagination, searchAndFilter)
	} else {
		entities, err = svc.model.Paginate(pagination, searchAndFilter)
	}

	if err != nil {
		logrus.Error(err)
		return nil, 0
	}

	if ownerID != 0 {
		bookmarkIDs, err = svc.model.SelectBookmarkedFundraiseID(ownerID)

		if err != nil {
			logrus.Error(err)
			return nil, 0
		}
	}

	if err != nil {
		logrus.Error(err)
		return nil, 0
	}

	for _, fundraise := range entities {
		var data dtos.ResFundraise

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

	var totalData int64 = 0

	if searchAndFilter.Title != "" || searchAndFilter.MinTarget != 0 || searchAndFilter.MaxTarget != math.MaxInt32 {
		if suffix == "mobile" {
			totalData = svc.model.GetTotalDataBySearchAndFilterMobile(searchAndFilter)
		} else {
			totalData = svc.model.GetTotalDataBySearchAndFilter(searchAndFilter)
		}
	} else {
		if suffix == "mobile" {
			totalData = svc.model.GetTotalDataMobile()
		} else {
			totalData = svc.model.GetTotalData()
		}
	}

	return fundraises, totalData
}

func (svc *service) FindByID(fundraiseID, ownerID int) *dtos.ResFundraise {
	var res dtos.ResFundraise
	fundraise, err := svc.model.SelectByID(fundraiseID)

	if err != nil {
		logrus.Error(err)
		return nil
	}

	var bookmarkID string

	if ownerID != 0 {
		bookmarkID, err = svc.model.SelectBookmarkByFundraiseAndOwnerID(fundraiseID, ownerID)

		if bookmarkID != "" {
			res.BookmarkID = &bookmarkID
		}
	}

	if err := smapping.FillStruct(&res, smapping.MapFields(fundraise)); err != nil {
		logrus.Error(err)
		return nil
	}

	if res.FundAcquired, err = svc.model.TotalFundAcquired(fundraiseID); err != nil {
		logrus.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newFundraise dtos.InputFundraise, userID int, file multipart.File) (*dtos.ResFundraise, []string, error) {
	if errMap := svc.validation.ValidateRequest(newFundraise); errMap != nil {
		return nil, errMap, errors.New("error")
	}

	var fundraise fundraise.Fundraise
	var url string = ""

	if file != nil {
		imageURL, err := svc.model.UploadFile(file)

		if err != nil {
			logrus.Error(err)
			return nil, nil, err
		}

		url = imageURL
	} else {
		config := config.LoadCloudStorageConfig()
		url = "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/fundraises/default"
	}

	fundraise.UserID = userID
	fundraise.Status = "pending"
	fundraise.Photo = url
	if err := smapping.FillStruct(&fundraise, smapping.MapFields(newFundraise)); err != nil {
		logrus.Error(err)
		return nil, nil, err
	}

	inserted, err := svc.model.Insert(fundraise)

	if err != nil {
		return nil, nil, err
	}

	var res dtos.ResFundraise

	if err := smapping.FillStruct(&res, smapping.MapFields(inserted)); err != nil {
		return nil, nil, err
	}

	return &res, nil, nil
}

func (svc *service) Modify(fundraiseData dtos.InputFundraise, file multipart.File, oldData dtos.ResFundraise) ([]string, error) {
	if errMap := svc.validation.ValidateRequest(fundraiseData); errMap != nil {
		return errMap, errors.New("error")
	}

	var newFundraise fundraise.Fundraise
	var url string = ""
	var config = config.LoadCloudStorageConfig()
	var oldFilename string = oldData.Photo
	var urlLength int = len("https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/fundraises/")

	if file != nil {
		if len(oldFilename) > urlLength {
			oldFilename = oldFilename[urlLength:]
		}

		if err := svc.model.DeleteFile(oldFilename); err != nil {
			return nil, err
		}

		imageURL, err := svc.model.UploadFile(file)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		url = imageURL
	}

	if err := smapping.FillStruct(&newFundraise, smapping.MapFields(fundraiseData)); err != nil {
		logrus.Error(err)
		return nil, err
	}

	newFundraise.Photo = url
	newFundraise.ID = oldData.ID
	newFundraise.UserID = oldData.UserID
	newFundraise.Status = oldData.Status

	if err := svc.model.Update(newFundraise); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return nil, nil
}

func (svc *service) ModifyStatus(input dtos.InputFundraiseStatus, oldData dtos.ResFundraise) ([]string, error) {
	if errMap := svc.validation.ValidateRequest(input); errMap != nil {
		return errMap, errors.New("error")
	}

	var newFundraise fundraise.Fundraise

	if err := smapping.FillStruct(&newFundraise, smapping.MapFields(oldData)); err != nil {
		logrus.Error(err)
		return nil, err
	}

	newFundraise.Status = input.Status
	if input.Status == "rejected" {
		if input.RejectedReason == "" {
			return []string{"rejected_reason field is required when the status is rejected"}, errors.New("error reason empty string")
		}
		newFundraise.RejectedReason = input.RejectedReason
	}

	if err := svc.model.Update(newFundraise); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return nil, nil
}

func (svc *service) Remove(fundraiseID int, oldData dtos.ResFundraise) error {
	var config = config.LoadCloudStorageConfig()
	var oldFilename string = oldData.Photo
	var urlLength int = len("https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/fundraises/")

	if len(oldFilename) > urlLength {
		oldFilename = oldFilename[urlLength:]
	}

	if oldFilename != "default" {
		svc.model.DeleteFile(oldFilename)
	}

	if err := svc.model.DeleteByID(fundraiseID); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
