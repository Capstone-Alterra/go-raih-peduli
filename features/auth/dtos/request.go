package dtos

type RequestLogin struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type InputUser struct {
	Fullname    string `json:"fullname" form:"fullname" validate:"required"`
	Address     string `json:"address" form:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required"`
	Gender      string `json:"gender" form:"gender" validate:"required"`
	Email       string `json:"email" form:"email" validate:"required"`
	Password    string `json:"password" form:"password" validate:"required"`
}

type ResendOTP struct {
	Email string `json:"email" form:"email" validate:"required"`
}
