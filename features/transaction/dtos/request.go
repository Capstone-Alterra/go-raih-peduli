package dtos

import "time"

type InputTransaction struct {
	FundraiseID int    `json:"fundraise_id" form:"fundraise_id" validate:"required"`
	PaymentType string `json:"payment_type" form:"payment_type" validate:"required"`
	Amount      int    `json:"amount" form:"amount" validate:"required"`
}

type InputToken struct {
	UserID    int       `json:"user_id" form:"user_id" validate:"required"`
	DeviceId  string    `json:"device_id" form:"device_id" validate:"required"`
	Timestamp time.Time `json:"timestamp" form:"timestamp" validate:"required"`
}

type Pagination struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}
