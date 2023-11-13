package dtos

type LoginResponse struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         int    `json:"role_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ResUser struct {
	ID           int    `json:"id"`
	RoleID       int    `json:"role_id"`
	Fullname     string `json:"fullname"`
	Address      string `json:"address"`
	PhoneNumber  string `json:"phone_number"`
	Gender       string `json:"gender"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
