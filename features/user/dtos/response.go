package dtos

type ResUser struct {
	ID             int    `json:"id"`
	RoleID         int    `json:"role_id"`
	Fullname       string `json:"fullname"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	PhoneNumber    string `json:"phone_number"`
	Gender         string `json:"gender"`
	ProfilePicture string `json:"profile_picture"`
}

type ResLogin struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ResMyProfile struct {
	ID             int    `json:"id"`
	RoleID         int    `json:"role_id"`
	Fullname       string `json:"fullname"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	PhoneNumber    string `json:"phone_number"`
	Gender         string `json:"gender"`
	Nik            string `json:"nik"`
	ProfilePicture string `json:"profile_picture"`
}
