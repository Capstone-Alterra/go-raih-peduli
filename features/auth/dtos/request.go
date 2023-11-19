package dtos

type RequestLogin struct {
	Email    string `json:"email" form:"email" validation:"required"`
	Password string `json:"password" form:"password" validation:"required"`
}

type InputUser struct {
	Fullname    string `json:"fullname" form:"fullname" validation:"required"`
	Address     string `json:"address" form:"address" validation:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validation:"required"`
	Gender      string `json:"gender" form:"gender" validation:"required"`
	Email       string `json:"email" form:"email" validation:"required"`
	Password    string `json:"password" form:"password" validation:"required"`
}

type ResendOTP struct {
	Email string `json:"email" form:"email" validation:"required"`
}
