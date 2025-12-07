package users

import (
	"context"
	"errors"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	Logout(ctx context.Context, refreshToken string) error
	VerifyEmail(ctx context.Context, token string) error
}

type svc struct {
	//repository Querier interface
}

func NewService() Service {
	return &svc{}
}

func (s *svc) Register(ctx context.Context, req RegisterRequest) (UserResponse, error) {
	// TODO1: Insert User, hash password, create 1-time email token
	return UserResponse{
		ID:       "dummy-id",
		Email:    req.Email,
		Username: req.Username,
	}, nil
}
func (s *svc) Login(ctx context.Context, req LoginRequest) (Tokens, error) {
	// TODO2: verify user + password, return real token
	return Tokens{
		AccessToken:  "dummy-access-token",
		RefreshToken: "dummy-refresh-token",
	}, nil
}

func (s *svc) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	// TODO3: validate refresh tokens & issue new tokens

	if refreshToken == "" {
		return Tokens{}, errors.New("missing token")
	}

	return Tokens{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
	}, nil
}

func (s *svc) Logout(ctx context.Context, refreshToken string) error {
	// TODO4: revoke token in DB
	return nil
}

func (s *svc) VerifyEmail(ctx context.Context, token string) error {
	// TODO5: check token in DB, mark email_verified = true
	return nil
}
