package dtos

type ResTransaction struct {
	ID             int    `json:"transaction_id"`
	UserID         int    `json:"user_id"`
	Fullname       string `json:"fullname"`
	Address        string `json:"address"`
	PhoneNumber    string `json:"phone_number"`
	ProfilePicture string `json:"profile_picture"`
	FundraiseID    int    `json:"fundraise_id"`
	Amount         int    `json:"amount"`
	PaymentType    string `json:"payment_type"`
	VirtualAccount string `json:"virtual_account"`
	UrlCallback    string `json:"url_callback"`
	PaidAt         string `json:"paid_at"`
	Status         string `json:"status"`
}
