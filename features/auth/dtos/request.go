package dtos

type RequestLogin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
