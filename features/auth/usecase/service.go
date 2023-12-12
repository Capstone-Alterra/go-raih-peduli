package usecase

import (
	"errors"
	"os"
	"raihpeduli/config"
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
	validator helpers.ValidationInterface
}

func New(model auth.Repository, jwt helpers.JWTInterface, hash helpers.HashInterface, generator helpers.GeneratorInterface, validator helpers.ValidationInterface) auth.Usecase {
	return &service{
		model:     model,
		jwt:       jwt,
		hash:      hash,
		generator: generator,
		validator: validator,
	}
}

func (svc *service) Register(newData dtos.InputUser) (*dtos.ResUser, []string, error) {
	config := config.LoadCloudStorageConfig()

	errMap := svc.validator.ValidateRequest(newData)
	if errMap != nil {
		return nil, errMap, nil
	}

	newUser := auth.User{}
	err := smapping.FillStruct(&newUser, smapping.MapFields(newData))
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	checkUser, err := svc.model.SelectByEmail(newUser.Email)
	if checkUser != nil {
		logrus.Print("User already exists")
		return nil, nil, errors.New("User already exists")
	}

	newUser.Password = svc.hash.HashPassword(newUser.Password)
	newUser.ProfilePicture = "https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/users/default"
	userModel, err := svc.model.Register(&newUser)
	if userModel == nil {
		return nil, nil, err
	}

	otp := svc.generator.GenerateRandomOTP()

	err = svc.model.SendOTPByEmail(userModel.Fullname, userModel.Email, otp, "1")
	if err != nil {
		return nil, nil, err
	}

	resCustomer := dtos.ResUser{}

	err = smapping.FillStruct(&resCustomer, smapping.MapFields(userModel))
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	return &resCustomer, nil, nil
}

func (svc *service) Login(data dtos.RequestLogin) (*dtos.LoginResponse, []string, error) {
	errMap := svc.validator.ValidateRequest(data)
	if errMap != nil {
		return nil, errMap, nil
	}

	user, err := svc.model.Login(data.Email)
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}

	if !svc.hash.CompareHash(data.Password, user.Password) {
		return nil, nil, errors.New("invalid password")
	}

	resUser := dtos.LoginResponse{}

	err = smapping.FillStruct(&resUser, smapping.MapFields(user))
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	userID := strconv.Itoa(user.ID)
	roleID := strconv.Itoa(resUser.RoleID)
	tokenData := svc.jwt.GenerateJWT(userID, roleID)

	if tokenData == nil {
		log.Error("Token process failed")
		return nil, nil, errors.New("generate token failed")
	}

	if err := svc.model.InsertToken(user.ID, data.FCMTokens); err != nil {
		return nil, nil, err
	}
	
	if user.Personalization == nil {
		resUser.PersonalizeUser = true
	}

	resUser.AccessToken = tokenData["access_token"].(string)
	resUser.RefreshToken = tokenData["refresh_token"].(string)

	return &resUser, nil, nil
}

func (svc *service) ResendOTP(email string) bool {
	otp := svc.generator.GenerateRandomOTP()

	data, err := svc.model.SelectByEmail(email)
	if err != nil {
		return false
	}
	err = svc.model.SendOTPByEmail(data.Fullname,email, otp, "1")
	if err != nil {
		return false
	}

	return true
}

func (svc *service) RefreshJWT(jwt dtos.RefreshJWT) (*dtos.ResJWT, error) {
	parsedRefreshToken, err := svc.jwt.ValidateToken(jwt.RefreshToken, os.Getenv("REFSECRET"))
	if err != nil {
		return nil, errors.New("validate token failed")
	}

	token := svc.jwt.RefereshJWT(parsedRefreshToken)
	if token == nil {
		return nil, errors.New("refresh jwt failed")
	}

	var resJWT dtos.ResJWT
	resJWT.AccessToken = token["access_token"].(string)
	resJWT.RefreshToken = token["refresh_token"].(string)

	return &resJWT, nil
}
