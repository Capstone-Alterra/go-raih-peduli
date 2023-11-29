package dtos

import "mime/multipart"

type InputUser struct {
	Fullname    string `json:"fullname" form:"fullname" validate:"required,alpha"`
	Address     string `json:"address" form:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required,min=10,max=13"`
	Gender      string `json:"gender" form:"gender" validate:"required"`
	Email       string `json:"email" form:"email" validate:"required,email"`
	Password    string `json:"password" form:"password" validate:"required,min=8"`
}

type InputUpdate struct {
	Fullname       string `json:"fullname" form:"fullname" validate:"required,alpha"`
	Address        string `json:"address" form:"address" validate:"required"`
	PhoneNumber    string `json:"phone_number" form:"phone_number" validate:"required,min=10,max=13"`
	Gender         string `json:"gender" form:"gender" validate:"required"`
	Email          string `json:"email" form:"email" validate:"required,email"`
	ProfilePicture string `json:"profile_picture" form:"profile_picture"`
}

type InputUpdateProfilePicture struct {
	ProfilePicture multipart.File `validate:"required"`
}

type Pagination struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

type VerifyOTP struct {
	OTP string `json:"otp" form:"otp"`
}

type ForgetPassword struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

type ResetPassword struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}
