package users

import (
	"contacts/internal/json"
	"contacts/internal/transport"
	"errors"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{service: s}
}

// // Register godoc
// @Summary      Register new user
// @Description  Create a new user account and send verification email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      registerRequestDTO  true  "register data"
// @Success      201    {object}  users.RegisterResponseDoc
// @Failure      400    {string}  string  "invalid body"
// @Failure      409    {string}  string  "email or username already exists"
// @Failure      500    {string}  string  "server error. registration failed"
// @Router       /auth/register [post]
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
			http.Error(w, "email or username already exists", http.StatusConflict)
			return
		}
		http.Error(w, "server error. registration failed", http.StatusInternalServerError)
		return
	}

	userResp := newUserDTO(user)

	resp := transport.Envelope{
		Message: "user registered successfully",
		Data: map[string]interface{}{
			"code": http.StatusCreated,
			"user": userResp,
		},
	}

	json.Write(w, http.StatusCreated, resp)
}

// // Login godoc
// @Summary      Login with validated user
// @Description  Login with a validated user and generate then return the token pairs
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      loginRequestDTO  true  "login data"
// @Success      200    {object}  TokenPair
// @Failure      400    {string}  string  "invalid body"
// @Failure      403    {string}  string  "email not verified - please check your inbox"
// @Failure      401    {string}  string  "invalid credentials"
// @Failure      500    {string}  string  "server error. login failed failed"
// @Router       /auth/login [post]
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
		if errors.Is(err, ErrEmailNotVerified) {
			http.Error(w, "email not verified", http.StatusForbidden)
			return
		}
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "server error. login failed", http.StatusInternalServerError)
		return
	}

	tokenPairResp := newTokensDTO(tokens)
	resp := transport.Envelope{
		Message: "logged in successfully",
		Data: map[string]interface{}{
			"code": http.StatusOK,
			"user": tokenPairResp,
		},
	}

	json.Write(w, http.StatusOK, resp)
}

// // Refresh godoc
// @Summary      Generate new refresh token
// @Description  If Refresh token is still valid, refresh the refresh token through this endpoint and regenerate token pairs
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      refreshTokenRequestDTO  true  "refresh token data"
// @Success      200    {object}  tokensDTO
// @Failure      400    {string}  string  "invalid body"
// @Failure      401    {string}  string  "invalid refresh token"
// @Failure      500    {string}  string  "server error. Refresh failed, please try again later"
// @Router       /auth/refresh [post]
func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {

	// struct to hold refresh token
	var body refreshTokenRequestDTO

	if err := json.Read(r, &body); err != nil || body.RefreshToken == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	newTokens, err := h.service.RefreshTokens(r.Context(), body.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "server error. Refresh failed, please try again later", http.StatusInternalServerError)
		return
	}

	// Send the new Access Token with Refresh Token to frontend
	json.Write(w, http.StatusOK, newTokensDTO(newTokens))
}

// // Logout godoc
// @Summary      logout
// @Description  Revoke existing refresh token and logout
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      refreshTokenRequestDTO  true  "refresh token data"
// @Success      200    {object}  LogoutResponse
// @Failure      400    {string}  string  "invalid body"
// @Failure      500    {string}  string  "server error. logout failed."
// @Router       /auth/logout [post]
func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	var body refreshTokenRequestDTO

	// Marshall JSON and give Bad Request error if body is invalid
	if err := json.Read(r, &body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.service.Logout(r.Context(), body.RefreshToken); err != nil {
		http.Error(w, "server error. logout failed", http.StatusInternalServerError)
		return
	}

	// Send dummy logout data
	json.Write(w, http.StatusOK, map[string]string{
		"message": "logged out",
	})
}

// VerifyEmail godoc
// @Summary     Email verification
// @Description Verify user email using a one-time token valid for 5 minutes after registration.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       token query string true "email verification token"
// @Success     200 {object} VerifyEmailResponse
// @Failure     400 {string} string "missing, expired, or invalid email verification token"
// @Failure     409 {string} string "token already used; verification already completed"
// @Failure     500 {string} string "server error"
// @Router      /auth/verify-email [get]
func (h *handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, "missing email verification token", http.StatusBadRequest)
		return
	}

	if err := h.service.VerifyEmail(r.Context(), token); err != nil {
		if errors.Is(err, ErrUsedVerifyToken) {
			http.Error(w, "token already used verificaiton is complete, please log in instead", http.StatusConflict)
			return
		}
		if errors.Is(err, ErrInvalidVerifyToken) {
			http.Error(w, "expired email verification token", http.StatusBadRequest)
		}
		http.Error(w, "invalid email verification token", http.StatusBadRequest)
		return
	}

	// Send back response after verifying
	json.Write(w, http.StatusOK, map[string]string{
		"message": "email verified",
	})
}
