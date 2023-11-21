package dtos

type InputUser struct {
	Fullname    string `json:"fullname" form:"fullname"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Gender      string `json:"gender" form:"gender"`
	Email       string `json:"email" form:"email"`
	Password    string `json:"password" form:"password"`
}

type InputUpdate struct {
	Fullname       string `json:"fullname" form:"fullname"`
	Address        string `json:"address" form:"address"`
	PhoneNumber    string `json:"phone_number" form:"phone_number"`
	Gender         string `json:"gender" form:"gender"`
	Email          string `json:"email" form:"email"`
	ProfilePicture string `json:"profile_picture" form:"profile_picture"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
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
