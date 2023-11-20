package usecase

import (
	"errors"
	"raihpeduli/features/auth"
	"raihpeduli/features/auth/dtos"
	"raihpeduli/helpers"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type service struct {
	model     auth.Repository
	jwt       helpers.JWTInterface
	hash      helpers.HashInterface
	generator helpers.GeneratorInterface
	validator helpers.ValidationInterface
	converter helpers.ConverterInterface
}

func New(model auth.Repository, jwt helpers.JWTInterface, hash helpers.HashInterface, generator helpers.GeneratorInterface, validator helpers.ValidationInterface, converter helpers.ConverterInterface) auth.Usecase {
	return &service{
		model:     model,
		jwt:       jwt,
		hash:      hash,
		generator: generator,
		validator: validator,
		converter: converter,
	}
}

func (svc *service) Register(newData dtos.InputUser) (*dtos.ResUser, error) {
	newUser := auth.User{}

	err := svc.converter.Convert(&newUser, newData)
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

	err = svc.converter.Convert(&resCustomer, userModel)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	userID := strconv.Itoa(userModel.ID)
	roleID := strconv.Itoa(resCustomer.RoleID)
	tokenData := svc.jwt.GenerateJWT(userID, roleID)

	if tokenData == nil {
		log.Error("Token process failed")
	}

	// resCustomer.AccessToken = tokenData["access_token"].(string)
	// resCustomer.RefreshToken = tokenData["refresh_token"].(string)

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

	resUser := dtos.LoginResponse{}

	err = svc.converter.Convert(&resUser, user)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	userID := strconv.Itoa(user.ID)
	roleID := strconv.Itoa(resUser.RoleID)
	tokenData := svc.jwt.GenerateJWT(userID, roleID)

	if tokenData == nil {
		log.Error("Token process failed")
	}

	resUser.AccessToken = tokenData["access_token"].(string)
	resUser.RefreshToken = tokenData["refresh_token"].(string)

	return &resUser, nil
}

func (svc *service) ResendOTP(email string) bool {
	otp := svc.generator.GenerateRandomOTP()

	err := svc.model.SendOTPByEmail(email, otp)
	if err != nil {
		return false
	}

	return true
}
