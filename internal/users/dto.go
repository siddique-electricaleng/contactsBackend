package users

/*
	DTOs for reponses: Users Package
*/

type registerRequestDTO struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	MobileNumber string `json:"mobileNumber"` // optional
}

type loginRequestDTO struct {
	EmailOrUsername string `json:"verify"`
	Password        string `json:"password"`
}

type tokensDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type userResponseDTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (d registerRequestDTO) toCommand() RegisterCommand {
	return RegisterCommand(d)
}

func (d loginRequestDTO) toCommand() LoginCommand {
	return LoginCommand(d)
}

func newTokensDTO(t TokenPair) tokensDTO {
	return tokensDTO(t)
}

func newUserDTO(u UserResponse) userResponseDTO {
	return userResponseDTO(u)
}
