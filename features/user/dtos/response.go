package dtos

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

type ResLogin struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type PaginationResponse struct {
	TotalData    int64 `json:"total_data"`
	CurrentPage  int   `json:"current_page"`
	PreviousPage int   `json:"previous_page"`
	NextPage     int   `json:"next_page"`
	TotalPage    int   `json:"total_page"`
}
