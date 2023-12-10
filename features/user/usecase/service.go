package usecase

import (
	"errors"
	"mime/multipart"
	"raihpeduli/features/user"
	"raihpeduli/features/user/dtos"
	"raihpeduli/helpers"
	"strconv"
	"strings"

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
		return nil, errMap, errors.New("missing some data")
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

	return &resUser, nil, nil
}

func (svc *service) Modify(userData dtos.InputUpdate, file multipart.File, oldData dtos.ResUser) (error, []string) {
	errMap := svc.validation.ValidateRequest(userData)
	if errMap != nil {
		return nil, errMap
	}

	if userData.Nik != "" {
		_, err := strconv.Atoi(userData.Nik)
		if err != nil {
			return nil, []string{"NIK must be 16 digits of number"}
		}

		if len(userData.Nik) != 16 {
			return nil, []string{"NIK must be 16 digits of number"}
		}
	}

	url, err := svc.model.UploadFile(file, oldData.ProfilePicture)
	if err != nil {
		return errors.New("upload profile picture failed"), nil
	}

	var newUser user.User
	err = smapping.FillStruct(&newUser, smapping.MapFields(userData))
	if err != nil {
		return err, nil
	}

	newUser.ID = oldData.ID
	newUser.ProfilePicture = url

	rowsAffected := svc.model.UpdateUser(newUser)
	if rowsAffected == 0 {
		log.Error("There is No Customer Updated!")
		return errors.New("update user failed"), nil
	}

	return nil, nil
}

func (svc *service) ModifyProfilePicture(file dtos.InputUpdateProfilePicture, oldData dtos.ResUser) (error, []string) {
	errMap := svc.validation.ValidateRequest(file)
	if errMap != nil {
		return nil, errMap
	}

	url, err := svc.model.UploadFile(file.ProfilePicture, oldData.ProfilePicture)
	if err != nil {
		return errors.New("upload profile picture failed"), nil
	}

	var newUser user.User
	newUser.ID = oldData.ID
	newUser.ProfilePicture = url
	rowsAffected := svc.model.UpdateUser(newUser)

	if rowsAffected == 0 {
		log.Error("There is No Customer Updated!")
		return errors.New("update user failed"), nil
	}

	return nil, nil
}

func (svc *service) Remove(userID int) error {
	user := svc.model.SelectByID(userID)
	if user == nil {
		return errors.New("user not found")
	}

	err := svc.model.DeleteFile(user.ProfilePicture)
	if err != nil {
		return errors.New("delete profile picture failed")
	}

	rowsAffected := svc.model.DeleteByID(userID)
	if rowsAffected <= 0 {
		log.Error("There is No Customer Deleted!")
		return errors.New("delete user failed")
	}

	return nil
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

func (svc *service) MyProfile(userID int) *dtos.ResMyProfile {
	res := dtos.ResMyProfile{}
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

func (svc *service) CheckPassword(checkPassword dtos.CheckPassword, userID int) ([]string, error) {
	errMap := svc.validation.ValidateRequest(checkPassword)
	if errMap != nil {
		return errMap, nil
	}

	user := svc.model.SelectByID(userID)
	if user == nil {
		return nil, errors.New("user not found")
	}

	result := svc.hash.CompareHash(checkPassword.OldPassword, user.Password)
	if !result {
		return nil, errors.New("password not match")
	}

	return nil, nil
}

func (svc *service) ChangePassword(changePassword dtos.ChangePassword, userID int) ([]string, error) {
	errMap := svc.validation.ValidateRequest(changePassword)
	if errMap != nil {
		return errMap, nil
	}

	var user user.User

	user.ID = userID
	user.Password = svc.hash.HashPassword(changePassword.NewPassword)

	rowsAffected := svc.model.UpdateUser(user)
	if rowsAffected == 0 {
		return nil, errors.New("change password failed")
	}

	return nil, nil
}

func (svc *service) AddPersonalization(userID int, data dtos.InputPersonalization) error {
	user := svc.model.SelectByID(userID)
	if user == nil {
		return errors.New("user not found")
	}

	personalized := strings.Join(data.Personalization, ", ")
	user.Personalization = &personalized
	rowsAffected := svc.model.UpdateUser(*user)
	if rowsAffected == 0 {
		return errors.New("add personalization failed")
	}

	return nil
}
