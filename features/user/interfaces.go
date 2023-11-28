package user

import (
	"mime/multipart"
	"raihpeduli/features/user/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []User
	InsertUser(newUser *User) (*User, error)
	SelectByID(customerID int) *User
	SelectByEmail(email string) (*User, error)
	UpdateUser(user User) int64
	DeleteByID(customerID int) int64
	SendOTPByEmail(email string, otp string) error
	InsertVerification(email string, verificationKey string) error
	ValidateVerification(verificationKey string) string
	GetTotalData() int64
	UploadFile(file multipart.File, oldFilename string) (string, error)
}

type Usecase interface {
	FindAll(page, size int) ([]dtos.ResUser, int64)
	FindByID(customerID int) *dtos.ResUser
	Create(newUser dtos.InputUser) (*dtos.ResUser, []string, error)
	Modify(customerData dtos.InputUpdate, file multipart.File, oldData dtos.ResUser) (bool, []string)
	ModifyProfilePicture(file dtos.InputUpdateProfilePicture, oldData dtos.ResUser) (bool, []string)
	Remove(customerID int) bool
	ValidateVerification(verificationKey string) bool
	ForgetPassword(email dtos.ForgetPassword) error
	VerifyOTP(verificationKey string) string
	ResetPassword(newData dtos.ResetPassword) error
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
}
