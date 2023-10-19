package dtos

type ResAdmin struct {
	ID          int            `json:"id"`
	Fullname    string         `json:"fullname"`
	NIK         string         `json:"nik"`
	Address     string         `json:"address"`
	PhoneNumber string         `json:"phone_number"`
	Gender      string         `json:"gender"`
	Email       string         `json:"email"`
	Token       map[string]any `json:"token"`
}
