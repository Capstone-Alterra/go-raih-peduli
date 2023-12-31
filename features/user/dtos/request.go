package dtos

import "mime/multipart"

type InputUser struct {
	Fullname    string `json:"fullname" form:"fullname" validate:"required,alphabetic"`
	Address     string `json:"address" form:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required,number,min=10,max=13"`
	Gender      string `json:"gender" form:"gender" validate:"required"`
	Email       string `json:"email" form:"email" validate:"required,email"`
	Password    string `json:"password" form:"password" validate:"required,min=8"`
}

type InputUpdate struct {
	Email          string `json:"email" form:"email" validate:"required,email"`
	Fullname       string `json:"fullname" form:"fullname" validate:"required,alphabetic"`
	Address        string `json:"address" form:"address" validate:"required"`
	PhoneNumber    string `json:"phone_number" form:"phone_number" validate:"required,number,min=10,max=13"`
	Gender         string `json:"gender" form:"gender"`
	ProfilePicture string `json:"profile_picture" form:"profile_picture"`
	Nik            string `json:"nik" form:"nik"`
}

type InputUpdateProfilePicture struct {
	ProfilePicture multipart.File `validate:"required"`
}

type InputPersonalization struct {
	Personalization []string `json:"personalization" form:"personalization"`
}

type SearchAndFilter struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Name     string `query:"name"`
}

type VerifyOTP struct {
	OTP string `json:"otp" form:"otp"`
}

type ForgetPassword struct {
	Email string `json:"email" form:"email" validate:"email"`
}

type ResetPassword struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

type CheckPassword struct {
	OldPassword string `json:"old_password" form:"old_password" validate:"required"`
}

type ChangePassword struct {
	NewPassword string `json:"new_password" form:"new_password" validate:"required,min=8"`
}
