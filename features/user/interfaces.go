package user

import (
	"mime/multipart"
	"raihpeduli/features/user/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(searchAndFilter dtos.SearchAndFilter) []User
	InsertUser(newUser *User) (*User, error)
	SelectByID(customerID int) *User
	SelectByEmail(email string) (*User, error)
	UpdateUser(user User) int64
	DeleteByID(customerID int) int64
	SendOTPByEmail(fullname string, email string, otp string, status string) error
	InsertVerification(email string, verificationKey string) error
	ValidateVerification(verificationKey string) string
	GetTotalData() int64
	GetTotalDataByName(name string) int64
	UploadFile(file multipart.File, oldFilename string) (string, error)
	DeleteFile(fileName string) error
}

type Usecase interface {
	FindAll(searchAndFilter dtos.SearchAndFilter) ([]dtos.ResUser, int64)
	FindByID(customerID int) *dtos.ResUser
	Create(newUser dtos.InputUser) (*dtos.ResUser, []string, error)
	Modify(customerData dtos.InputUpdate, file multipart.File, oldData dtos.ResUser) (error, []string)
	ModifyProfilePicture(file dtos.InputUpdateProfilePicture, oldData dtos.ResUser) (error, []string)
	Remove(customerID int) error
	ValidateVerification(verificationKey string) bool
	ForgetPassword(email dtos.ForgetPassword) error
	VerifyOTP(verificationKey string) string
	ResetPassword(newData dtos.ResetPassword) error
	MyProfile(userID int) *dtos.ResMyProfile
	CheckPassword(checkPassword dtos.CheckPassword, userID int) ([]string, error)
	ChangePassword(changePassword dtos.ChangePassword, userID int) ([]string, error)
	AddPersonalization(userID int, data dtos.InputPersonalization) error
}

type Handler interface {
	GetUsers() echo.HandlerFunc
	UserDetails() echo.HandlerFunc
	CreateUser() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
	UpdateProfilePicture() echo.HandlerFunc
	DeleteUser() echo.HandlerFunc
	VerifyEmail() echo.HandlerFunc
	ForgetPassword() echo.HandlerFunc
	VerifyOTP() echo.HandlerFunc
	ResetPassword() echo.HandlerFunc
	MyProfile() echo.HandlerFunc
	CheckPassword() echo.HandlerFunc
	ChangePassword() echo.HandlerFunc
	AddPersonalization() echo.HandlerFunc
}
