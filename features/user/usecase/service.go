package usecase

import (
	"errors"
	"mime/multipart"
	"os"
	"raihpeduli/config"
	user "raihpeduli/features/user"
	"raihpeduli/features/user/dtos"
	"raihpeduli/helpers"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
)

type service struct {
	model      user.Repository
	jwt        helpers.JWTInterface
	hash       helpers.HashInterface
	generator  helpers.GeneratorInterface
	validation helpers.ValidationInterface
}

func New(model user.Repository, jwt helpers.JWTInterface, hash helpers.HashInterface, generator helpers.GeneratorInterface, validation helpers.ValidationInterface) user.Usecase {
	return &service{
		model:      model,
		jwt:        jwt,
		hash:       hash,
		generator:  generator,
		validation: validation,
	}
}

func (svc *service) FindAll(page, size int) ([]dtos.ResUser, int64) {
	var users []dtos.ResUser

	usersEnt := svc.model.Paginate(page, size)
	totalData := svc.model.GetTotalData()

	for _, user := range usersEnt {
		var data dtos.ResUser

		if err := smapping.FillStruct(&data, smapping.MapFields(user)); err != nil {
			log.Error(err.Error())
		}

		users = append(users, data)
	}

	return users, totalData
}

func (svc *service) FindByID(userID int) *dtos.ResUser {
	res := dtos.ResUser{}
	user := svc.model.SelectByID(userID)

	if user == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(user))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newData dtos.InputUser) (*dtos.ResUser, []string, error) {
	if errMap := svc.validation.ValidateRequest(newData); errMap != nil {
		return nil, errMap, nil
	}

	newUser := user.User{}

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
	userModel, err := svc.model.InsertUser(&newUser)
	if userModel == nil {
		return nil, nil, err
	}

	otp := svc.generator.GenerateRandomOTP()

	err = svc.model.SendOTPByEmail(userModel.Email, otp)
	if err != nil {
		return nil, nil, err
	}

	resUser := dtos.ResUser{}

	err = smapping.FillStruct(&resUser, smapping.MapFields(userModel))
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	userID := strconv.Itoa(userModel.ID)
	roleID := strconv.Itoa(resUser.RoleID)
	tokenData := svc.jwt.GenerateJWT(userID, roleID)

	if tokenData == nil {
		log.Error("Token process failed")
	}

	resUser.AccessToken = tokenData["access_token"].(string)
	resUser.RefreshToken = tokenData["refresh_token"].(string)

	return &resUser, nil, nil
}

func (svc *service) Modify(userData dtos.InputUpdate, file multipart.File, oldData dtos.ResUser) (bool, []string) {
	errMap := svc.validation.ValidateRequest(userData)
	if errMap != nil {
		return false, errMap
	}

	var newUser user.User
	err := smapping.FillStruct(&newUser, smapping.MapFields(userData))
	if err != nil {
		log.Error(err)
		return false, nil
	}

	newUser.ID = oldData.ID
	rowsAffected := svc.model.UpdateUser(newUser)

	if rowsAffected == 0 {
		log.Error("There is No Customer Updated!")
		return false, nil
	}

	return true, nil
}

func (svc *service) ModifyProfilePicture(file dtos.InputUpdateProfilePicture, oldData dtos.ResUser) (bool, []string) {
	errMap := svc.validation.ValidateRequest(file)
	if errMap != nil {
		return false, errMap
	}

	var newUser user.User
	var config = config.LoadCloudStorageConfig()
	var oldFilename string = oldData.ProfilePicture
	var urlLength int = len("https://storage.googleapis.com/" + config.CLOUD_BUCKET_NAME + "/users/")

	if file.ProfilePicture != nil {
		if oldFilename == "https://storage.googleapis.com/raih_peduli/users/user.png" {
			oldFilename = ""
		} else if len(oldFilename) > urlLength {
			oldFilename = oldFilename[urlLength:]
		}
		imageURL, err := svc.model.UploadFile(file.ProfilePicture, oldFilename)

		if err != nil {
			logrus.Error(err)
			return false, nil
		}

		newUser.ProfilePicture = imageURL
	}

	newUser.ID = oldData.ID
	rowsAffected := svc.model.UpdateUser(newUser)

	if rowsAffected == 0 {
		log.Error("There is No Customer Updated!")
		return false, nil
	}

	return true, nil
}

func (svc *service) Remove(userID int) bool {
	rowsAffected := svc.model.DeleteByID(userID)

	if rowsAffected <= 0 {
		log.Error("There is No Customer Deleted!")
		return false
	}

	return true
}

func (svc *service) ValidateVerification(verificationKey string) bool {
	email := svc.model.ValidateVerification(verificationKey)
	if email == "" {
		return false
	}

	user, err := svc.model.SelectByEmail(email)
	if err != nil {
		return false
	}

	user.IsVerified = true
	rowsAffected := svc.model.UpdateUser(*user)
	if rowsAffected <= 0 {
		log.Error("There is No Customer Deleted!")
		return false
	}

	return true
}

func (svc *service) ForgetPassword(data dtos.ForgetPassword) error {
	user, err := svc.model.SelectByEmail(data.Email)

	if err != nil {
		return err
	}

	rowsAffected := svc.model.UpdateUser(*user)

	if rowsAffected == 0 {
		log.Error("There is No Customer Updated!")
		return errors.New("There is No Customer Updated!")
	}

	otp := svc.generator.GenerateRandomOTP()

	err = svc.model.SendOTPByEmail(user.Email, otp)
	if err != nil {
		return err
	}

	return nil
}

func (svc *service) VerifyOTP(verificationKey string) string {
	email := svc.model.ValidateVerification(verificationKey)
	if email == "" {
		return ""
	}

	user, err := svc.model.SelectByEmail(email)
	if err != nil {
		return ""
	}

	userID := strconv.Itoa(user.ID)
	roleID := strconv.Itoa(user.RoleID)
	token := svc.jwt.GenerateTokenResetPassword(userID, roleID)

	return token
}

func (svc *service) ResetPassword(newData dtos.ResetPassword) error {
	user, err := svc.model.SelectByEmail(newData.Email)

	if err != nil {
		return err
	}

	user.Password = svc.hash.HashPassword(newData.Password)
	rowsAffected := svc.model.UpdateUser(*user)

	if rowsAffected == 0 {
		log.Error("There is No Customer Updated!")
		return errors.New("There is No Customer Updated!")
	}

	return nil
}

func (svc *service) RefreshJWT(jwt dtos.RefreshJWT) (*dtos.ResJWT, error) {
	parsedToken, err := svc.jwt.ValidateToken(jwt.RefreshToken, os.Getenv("SECRET"))
	if err != nil {
		return nil, errors.New("validate token failed")
	}

	token := svc.jwt.RefereshJWT(jwt.AccessToken, parsedToken)
	if token == nil {
		return nil, errors.New("refresh jwt failed")
	}

	var resJWT dtos.ResJWT
	resJWT.AccessToken = token["access_token"].(string)
	resJWT.RefreshToken = token["refresh_token"].(string)

	return &resJWT, nil
}
