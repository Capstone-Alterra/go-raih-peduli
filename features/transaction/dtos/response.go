package dtos

type ResTransaction struct {
	PaymentType    string `json:"payment_type"`
	VirtualAccount string `json:"virtual_account,omitempty"`
	UrlCallback    string `json:"url_callback,omitempty"`
}
