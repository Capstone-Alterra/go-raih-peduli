package dtos

type InputTransaction struct {
	FundraiseID int    `json:"fundraise_id" form:"fundraise_id" validate:"required"`
	PaymentType string `json:"payment_type" form:"payment_type" validate:"required"`
	Amount      int    `json:"amount" form:"amount" validate:"required"`
}

type Pagination struct {
	Page int `query:"page"`
	Size int `query:"size"`
}
