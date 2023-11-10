package dtos

type InputAdmin struct {
	Fullname    string `json:"fullname" form:"fullname"`
	NIK         string `json:"nik" form:"nik"`
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

type LoginAdmin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}