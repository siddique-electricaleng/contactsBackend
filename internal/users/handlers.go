package users

import (
	"contacts/internal/json"
	"errors"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{service: s}
}

// ENDPOINT 1: POST /api/v1/auth/register
func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var dto registerRequestDTO

	// Marshall the json with error checking
	if err := json.Read(r, &dto); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), dto.toCommand())
	if err != nil {
		if errors.Is(err, ErrConflict) {
			http.Error(w, "eamil or username already exists", http.StatusConflict)
			return
		}
		http.Error(w, "registration failed", http.StatusInternalServerError)
		return
	}

	// Dummy map data
	json.Write(w, http.StatusCreated, map[string]any{
		"user":    newUserDTO(user),
		"message": "verification email will be sent",
	})
}

// ENDPOINT 2: POST /api/v1/auth/login
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var dto loginRequestDTO

	// Marshall the json
	if err := json.Read(r, &dto); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Send request body and receive access + refresh tokens
	tokens, err := h.service.Login(r.Context(), dto.toCommand())
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, newTokensDTO(tokens))
}

// ENDPOINT 3: POST /api/v1/auth/refresh
func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {

	// struct to hold refresh token
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := json.Read(r, &body); err != nil || body.RefreshToken == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	newTokens, err := h.service.RefreshTokens(r.Context(), body.RefreshToken)
	if err != nil {
		http.Error(w, "refresh failed", http.StatusUnauthorized)
		return
	}

	// Send the new Access Token with Refresh Token to frontend
	json.Write(w, http.StatusOK, newTokensDTO(newTokens))
}

// ENDPOINT 4: POST /api/v1/auth/logout
func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}

	// Marshall JSON and give Bad Request error if body is invalid
	if err := json.Read(r, &body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.service.Logout(r.Context(), body.RefreshToken); err != nil {
		http.Error(w, "logout failed", http.StatusBadRequest)
		return
	}

	// Send dummy logout data
	json.Write(w, http.StatusOK, map[string]string{
		"message": "logged out",
	})
}

// ENDPOINT 5: GET /api/v1/auth/verify?token=123
func (h *handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}

	if err := h.service.VerifyEmail(r.Context(), token); err != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	// Send back response after verifying
	json.Write(w, http.StatusOK, map[string]string{
		"message": "email verified",
	})
}
