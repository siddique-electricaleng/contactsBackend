package users

type RegisterResponse struct {
	User    userResponseDTO `json:"user"`
	Message string          `json:"message"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}
type VerifyEmailResponse struct {
	Message string `json:"message"`
}
