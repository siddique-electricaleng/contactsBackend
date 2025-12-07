package users

/*
	Internal Domain Model for Users Package
*/

type RegisterCommand struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	MobileNumber string `json:"mobileNumber"`
}

type LoginCommand struct {
	EmailOrUsername string `json:"verify"`
	Password        string `json:"password"`
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
