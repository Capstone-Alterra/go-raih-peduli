package dtos

type ResTransaction struct {
	ID             int    `json:"transaction_id"`
	UserID         int    `json:"user_id"`
	Email          string `json:"email"`
	Fullname       string `json:"fullname"`
	Address        string `json:"address"`
	PhoneNumber    string `json:"phone_number"`
	Photo          string `json:"photo"`
	FundraiseID    int    `json:"fundraise_id"`
	FundraiseName  string `json:"fundraise_name"`
	Amount         int    `json:"amount"`
	PaymentType    string `json:"payment_type"`
	VirtualAccount string `json:"virtual_account"`
	UrlCallback    string `json:"url_callback"`
	PaidAt         string `json:"paid_at"`
	ValidUntil     string `json:"valid_until"`
	Status         string `json:"status"`
}
