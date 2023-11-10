package dtos

type InputCustomer struct {
	Fullname    string `json:"fullname" form:"fullname"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Gender      string `json:"gender" form:"gender"`
	Email       string `json:"email" form:"email"`
	Password    string `json:"password" form:"password"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
}

type VerifyOTP struct {
	Email string `json:"email" form:"email"`
	OTP   string `json:"otp" form:"otp"`
}
