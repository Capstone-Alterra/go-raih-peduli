package dtos

type ResCustomer struct {
	ID          int            `json:"id"`
	RoleID      int            `json:"role_id"`
	Fullname    string         `json:"fullname"`
	Address     string         `json:"address"`
	PhoneNumber string         `json:"phone_number"`
	Gender      string         `json:"gender"`
	Email       string         `json:"email"`
	Token       map[string]any `json:"token"`
}

type ResLogin struct {
	Name  string         `json:"name"`
	Email string         `json:"email"`
	Role  string         `json:"role_id"`
	Token map[string]any `json:"token"`
}
