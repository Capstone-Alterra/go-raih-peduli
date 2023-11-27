package dtos

import "mime/multipart"

type InputUser struct {
	Fullname    string `json:"fullname" form:"fullname" validate:"required"`
	Address     string `json:"address" form:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required"`
	Gender      string `json:"gender" form:"gender" validate:"required"`
	Email       string `json:"email" form:"email" validate:"required"`
	Password    string `json:"password" form:"password" validate:"required"`
}

type InputUpdate struct {
	Fullname       string `json:"fullname" form:"fullname" validate:"required"`
	Address        string `json:"address" form:"address" validate:"required"`
	PhoneNumber    string `json:"phone_number" form:"phone_number" validate:"required"`
	Gender         string `json:"gender" form:"gender" validate:"required"`
	Email          string `json:"email" form:"email" validate:"required"`
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
	Email string `json:"email" form:"email" validate:"email"`
}

type ResetPassword struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}
