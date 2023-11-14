package usecase

import (
	"mime/multipart"
	"raihpeduli/config"
	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/dtos"

	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
)

type service struct {
	model fundraise.Repository
}

func New(model fundraise.Repository) fundraise.Usecase {
	return &service {
		model: model,
	}
}

func (svc *service) FindAll(page int, size int, title string) []dtos.ResFundraise {
	var fundraises []dtos.ResFundraise

	entites, err := svc.model.Paginate(page, size, title)

	if err != nil {
		logrus.Error(err)
		return nil
	}

	for _, fundraise := range entites {
		var data dtos.ResFundraise

		if err := smapping.FillStruct(&data, smapping.MapFields(fundraise)); err != nil {
			logrus.Error(err)
		} 
		
		fundraises = append(fundraises, data)
	}

	return fundraises
}

func (svc *service) FindByID(fundraiseID int) *dtos.ResFundraise {
	var res dtos.ResFundraise
	fundraise, err := svc.model.SelectByID(fundraiseID)

	if err != nil {
		logrus.Error(err)
		return nil
	}
	
	if err := smapping.FillStruct(&res, smapping.MapFields(fundraise)); err != nil {
		logrus.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newFundraise dtos.InputFundraise, userID int, file multipart.File) (*dtos.ResFundraise, error) {
	var fundraise fundraise.Fundraise
	var url string = ""

	if file != nil {
		imageURL, err := svc.model.UploadFile(file, "")
	
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		url = imageURL
	} 
	
	fundraise.UserID = userID
	fundraise.Status = "pending"
	fundraise.Photo = url
	if err := smapping.FillStruct(&fundraise, smapping.MapFields(newFundraise)); err != nil {
		logrus.Error(err)
		return nil, err
	}

	if _, err := svc.model.Insert(fundraise); err != nil {
		return nil, err
	}

	var res dtos.ResFundraise
	
	res.Status = "pending"
	res.Photo = url
	res.UserID = userID
	if err := smapping.FillStruct(&res, smapping.MapFields(newFundraise)); err != nil {
		return nil, err
	}

	return &res, nil
}

func (svc *service) Modify(fundraiseData dtos.InputFundraise, file multipart.File, oldData dtos.ResFundraise) bool {
	var newFundraise fundraise.Fundraise
	var url string = ""
	var config = config.LoadCloudStorageConfig()
	var oldFilename string = oldData.Photo
	var urlLength int = len("https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/fundraises/")

	if file != nil {
		if len(oldFilename) > urlLength {
			oldFilename = oldFilename[urlLength:]
		}
		imageURL, err := svc.model.UploadFile(file, oldFilename)
		
		if err != nil {
			logrus.Error(err)
			return false
		}

		url = imageURL
	}

	if err := smapping.FillStruct(&newFundraise, smapping.MapFields(fundraiseData)); err != nil {
		logrus.Error(err)
		return false
	}
	
	newFundraise.Photo = url
	newFundraise.ID = oldData.ID
	newFundraise.UserID = oldData.UserID
	_, err := svc.model.Update(newFundraise)

	if err != nil {
		logrus.Error(err)
		return false
	}
	
	return true
}

func (svc *service) Remove(fundraiseID int) bool {
	_, err := svc.model.DeleteByID(fundraiseID)

	if err != nil {
		logrus.Error(err)
		return false
	}

	return true
}