package usecase

import (
	"errors"
	user "raihpeduli/features/user"
	"raihpeduli/features/user/dtos"
	"raihpeduli/helpers"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
)

type service struct {
	model     user.Repository
	jwt       helpers.JWTInterface
	hash      helpers.HashInterface
	generator helpers.GeneratorInterface
}

func New(model user.Repository, jwt helpers.JWTInterface, hash helpers.HashInterface, generator helpers.GeneratorInterface) user.Usecase {
	return &service{
		model:     model,
		jwt:       jwt,
		hash:      hash,
		generator: generator,
	}
}

func (svc *service) FindAll(page, size int) []dtos.ResUser {
	var users []dtos.ResUser

	usersEnt := svc.model.Paginate(page, size)

	for _, user := range usersEnt {
		var data dtos.ResUser

		if err := smapping.FillStruct(&data, smapping.MapFields(user)); err != nil {
			log.Error(err.Error())
		}

		users = append(users, data)
	}

	return users
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

func (svc *service) Create(newData dtos.InputUser) (*dtos.ResUser, error) {
	newUser := user.User{}

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
	userModel, err := svc.model.InsertUser(&newUser)
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

func (svc *service) Modify(userData dtos.InputUser, userID int) bool {
	newUser := user.User{}

	err := smapping.FillStruct(&newUser, smapping.MapFields(userData))
	if err != nil {
		log.Error(err)
		return false
	}

	rowsAffected := svc.model.UpdateUser(newUser)

	if rowsAffected == 0 {
		log.Error("There is No Customer Updated!")
		return false
	}

	return true
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

	user.IsVerified = false
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

	if user.IsVerified {
		return errors.New("Cannot reset password")
	}

	user.Password = svc.hash.HashPassword(newData.Password)
	user.IsVerified = true
	rowsAffected := svc.model.UpdateUser(*user)

	if rowsAffected == 0 {
		log.Error("There is No Customer Updated!")
		return errors.New("There is No Customer Updated!")
	}

	return nil
}
