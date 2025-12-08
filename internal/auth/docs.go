package auth

type RegisterResponseDoc struct {
	Message string          `json:"message"`
	Data    RegisterDataDoc `json:"data,omitempty"`
}

type RegisterDataDoc struct {
	Code int             `json:"code"`
	User userResponseDTO `json:"user"`
}
