package users

/*
	Internal Domain Model for Users Package
*/

type RegisterCommand struct {
	Email        string
	Username     string
	Password     string
	Name         string
	MobileNumber string
}

type LoginCommand struct {
	Identifier string // needs to be either email or username
	Password   string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type UserResponse struct {
	ID       string
	Email    string
	Username string
}
