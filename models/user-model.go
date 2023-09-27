package models

type User struct {
	ID          int    `json:"id" form:"id" gorm:"type:int(11)"`
	Username    string `json:"username" form:"username" gorm:"type:varchar(255)"`
	Email       string `json:"email" form:"email" gorm:"type:varchar(255)"`
	PhoneNumber string `json:"phone-number" form:"phone-number" gorm:"type:varchar(20)"`
}