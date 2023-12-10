package dtos

type RequestLogin struct {
	Email     string `json:"email" form:"email" validate:"required,email"`
	Password  string `json:"password" form:"password" validate:"required,min=8"`
	FCMTokens string `json:"fcm_token" form:"fcm_token" validate:"required"`
}

type InputUser struct {
	Fullname    string `json:"fullname" form:"fullname" validate:"required,alphabetic"`
	Address     string `json:"address" form:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required,number,min=10,max=13"`
	Gender      string `json:"gender" form:"gender" validate:"required"`
	Email       string `json:"email" form:"email" validate:"required,email"`
	Password    string `json:"password" form:"password" validate:"required,min=8"`
}

type ResendOTP struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

type RefreshJWT struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required"`
}
