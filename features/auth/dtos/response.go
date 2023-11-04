package dtos

type LoginResponse struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         int    `json:"role_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
