package usecase

import (
	"errors"
	"raihpeduli/features/auth"
	"raihpeduli/features/auth/dtos"
	"raihpeduli/helpers"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
)

type service struct {
	model     auth.Repository
	jwt       helpers.JWTInterface
	hash      helpers.HashInterface
	generator helpers.GeneratorInterface
}

func New(model auth.Repository, jwt helpers.JWTInterface, hash helpers.HashInterface, generator helpers.GeneratorInterface) auth.Usecase {
	return &service{
		model:     model,
		jwt:       jwt,
		hash:      hash,
		generator: generator,
	}
}

func (svc *service) Register(newData dtos.InputUser) (*dtos.ResUser, error) {
	newUser := auth.User{}

	err := smapping.FillStruct(&newUser, smapping.MapFields(newData))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	checkUser, err := svc.model.SelectByEmail(newUser.Email)
	if checkUser != nil {
		logrus.Print("User already exists")
		return nil, errors.New("User already exists")
	}

	newUser.Password = svc.hash.HashPassword(newUser.Password)
	userModel, err := svc.model.Register(&newUser)
	if userModel == nil {
		return nil, err
	}

	otp := svc.generator.GenerateRandomOTP()

	err = svc.model.SendOTPByEmail(userModel.Email, otp)
	if err != nil {
		return nil, err
	}

	resCustomer := dtos.ResUser{}

	resCustomer.RoleID = userModel.RoleID
	resCustomer.Email = userModel.Email
	resCustomer.ID = userModel.ID
	resCustomer.Fullname = userModel.Fullname
	resCustomer.Address = userModel.Address
	resCustomer.PhoneNumber = userModel.PhoneNumber
	resCustomer.Gender = userModel.Gender

	userID := strconv.Itoa(userModel.ID)
	roleID := strconv.Itoa(resCustomer.RoleID)
	tokenData := svc.jwt.GenerateJWT(userID, roleID)

	if tokenData == nil {
		log.Error("Token process failed")
	}

	resCustomer.AccessToken = tokenData["access_token"].(string)
	resCustomer.RefreshToken = tokenData["refresh_token"].(string)

	return &resCustomer, nil
}

func (svc *service) Login(data dtos.RequestLogin) (*dtos.LoginResponse, error) {
	user, err := svc.model.Login(data.Email)
	if err != nil {
		return nil, err
	}

	if !svc.hash.CompareHash(data.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	tokenData := svc.jwt.GenerateJWT(strconv.Itoa(user.ID), strconv.Itoa(user.RoleID))
	return &dtos.LoginResponse{
		Name:         user.Fullname,
		Email:        user.Email,
		Role:         user.RoleID,
		AccessToken:  tokenData["access_token"].(string),
		RefreshToken: tokenData["refresh_token"].(string),
	}, nil
}

func (svc *service) InsertVerification(email string, verificationKey string) error {
	return svc.model.InsertVerification(email, verificationKey)
}
