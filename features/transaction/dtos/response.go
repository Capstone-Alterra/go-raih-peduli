package dtos

type ResTransaction struct {
	ID             int    `json:"transaction_id"`
	UserID         int    `json:"user_id"`
	FundraiseID    int    `json:"fundraise_id"`
	Amount         int    `json:"amount"`
	PaymentType    string `json:"payment_type"`
	VirtualAccount string `json:"virtual_account"`
	UrlCallback    string `json:"url_callback"`
	PaidAt         string `json:"paid_at"`
	Status         string `json:"status"`
}

type PaginationResponse struct {
	TotalData    int64 `json:"total_data"`
	CurrentPage  int   `json:"current_page"`
	PreviousPage int   `json:"previous_page"`
	NextPage     int   `json:"next_page"`
	TotalPage    int   `json:"total_page"`
}
