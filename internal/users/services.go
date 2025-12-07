package users

import (
	repo "contacts/internal/adapters/postgresql/sqlc"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrConflict           = errors.New("user with that email or username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service interface {
	Register(ctx context.Context, req RegisterCommand) (UserResponse, error)
	Login(ctx context.Context, req LoginCommand) (TokenPair, error)
	RefreshTokens(ctx context.Context, refreshToken string) (TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	VerifyEmail(ctx context.Context, token string) error
}

type svc struct {
	repo repo.Querier
}

func NewService(r *repo.Queries) Service {
	return &svc{repo: r}
}

// Registration Method

func (s *svc) Register(ctx context.Context, req RegisterCommand) (UserResponse, error) {
	// TODO1: Insert User, hash password, create 1-time email token

	// password hash
	pwHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, nil
	}

	arg := repo.CreateUserParams{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(pwHash),
		Name:         req.Name,
	}

	user, err := s.repo.CreateUser(ctx, arg)
	if err != nil {
		return UserResponse{}, ErrConflict
	}

	// Final return if everything went right
	return UserResponse{
		ID:       "",
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

// Login Method

func (s *svc) Login(ctx context.Context, req LoginCommand) (TokenPair, error) {
	// TODO2: verify user + password, return real token

	user, err := s.repo.GetUserByEmailOrUsername(ctx, req.EmailOrUsername)

	if err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}

	return TokenPair{
		AccessToken:  "dummy-access-token",
		RefreshToken: "dummy-refresh-token",
	}, nil
}

// Access Token Refresh Method

func (s *svc) RefreshTokens(ctx context.Context, refreshToken string) (TokenPair, error) {
	// TODO3: validate refresh tokens & issue new tokens

	if refreshToken == "" {
		return TokenPair{}, errors.New("missing token")
	}

	return TokenPair{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
	}, nil
}

// Logout Method

func (s *svc) Logout(ctx context.Context, refreshToken string) error {
	// TODO4: revoke token in DB
	return nil
}

// Email verification method

func (s *svc) VerifyEmail(ctx context.Context, token string) error {
	// TODO5: check token in DB, mark email_verified = true
	return nil
}
