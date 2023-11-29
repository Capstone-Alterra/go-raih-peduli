package usecase

import (
	"errors"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"raihpeduli/config"
	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/dtos"
	"raihpeduli/helpers"
	"time"

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
	if errorList, err := svc.validateInput(newFundraise, file); err != nil || len(errorList) > 0 {
		return nil, errorList, err
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
	if errorList, err := svc.validateInput(fundraiseData, file); len(errorList) > 0 || err != nil {
		return errorList, errors.New("error")
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

		if len(input.RejectedReason) < 20 {
			return []string{"rejected_reason must be at least 20 characters"}, errors.New("characters must be at least 20")
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

func (svc *service) validateInput(input dtos.InputFundraise, file multipart.File) ([]string, error) {
	var errorList []string 
	if errMap := svc.validation.ValidateRequest(input); errMap != nil {
		errorList = append(errorList, errMap...)
	}

	if len(input.Title) < 20 {
		errorList = append(errorList, "title must be at least 20 characters")
	}

	if len(input.Description) < 50 {
		errorList = append(errorList, "description must be at least 50 characters")
	}

	if input.Target < 100 {
		errorList = append(errorList, "target must be at less 100 rupiahs")
	}

	if file != nil {
		buffer := make([]byte, 512)
		
		if _, err := file.Read(buffer); err != nil {
			return nil, err
		}

		contentType := http.DetectContentType(buffer)
		isImage := contentType[:5] == "image"

		if !isImage {
			errorList = append(errorList, "photo file has to be an image (png, jpg, or jpeg)")
		}

		const maxFileSize = 5 * 1024 * 1024
		var fileSize int64

		buffer = make([]byte, 1024)
		for {
			n, err := file.Read(buffer)

			fileSize += int64(n)

			if err == io.EOF {
				break
			}

			if err != nil {
				errorList = append(errorList, "unknown file size")
			}
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}

		if fileSize > maxFileSize {
			errorList = append(errorList, "fize size exceeds the allowed limit (5MB)")
		}
	}

	if input.EndDate.Before(input.StartDate) {
		errorList = append(errorList, "end_date can not be earlier than start_date")
	}

	wibLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	startDate := input.StartDate.Truncate(24 * time.Hour)
	currentDate := time.Now().In(wibLocation).Truncate(24 * time.Hour)

	if startDate.Sub(currentDate).Milliseconds() < 0 {
		errorList = append(errorList, "start_date can not be earlier than current date")
	}

	return errorList, nil
}
