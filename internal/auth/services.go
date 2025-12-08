package auth

import (
	repo "contacts/internal/adapters/postgresql/sqlc"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrConflict            = errors.New("user with that email or username already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrEmailNotVerified    = errors.New("email not verified")
	ErrInvalidVerifyToken  = errors.New("invalid email verification token")
	ErrUsedVerifyToken     = errors.New("email already verified, please log in instead")
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

func (s *svc) Register(ctx context.Context, cmd RegisterCommand) (UserResponse, error) {
	// TODO1: Insert User, hash password, create 1-time email token

	// password hash
	pwHash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, nil
	}

	// Mandatory Fields for Registration
	arg := repo.CreateUserParams{
		Email:        cmd.Email,
		Username:     cmd.Username,
		PasswordHash: string(pwHash),
		Name:         cmd.Name,
	}

	user, err := s.repo.CreateUser(ctx, arg)
	if err != nil {
		return UserResponse{}, ErrConflict
	}

	// Generate Email Verification Token
	token, expiresAt, err := generateEmailVerificationToken()
	if err != nil {
		return UserResponse{}, err
	}

	// Store email verification token into DB
	_, err = s.repo.CreateEmailVerificationToken(ctx, repo.CreateEmailVerificationTokenParams{
		UserID: user.ID,
		Token:  token,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return UserResponse{}, err
	}

	// TASK0 Sending Email ---  Note by Saad: Please Check Again

	go SendEmailVerification(
		cmd.Email,
		"Verify your email",
		fmt.Sprintf("Click here to verify your registration to Contacts: %s%s?token=%s", appBaseURL, emaiLVerificationPath, token),
	)
	// Final return if everything went right
	return UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

// Login Method

func (s *svc) Login(ctx context.Context, cmd LoginCommand) (TokenPair, error) {
	// TODO2: verify user + password, return real token

	user, err := s.repo.GetUserByEmailOrUsername(ctx, cmd.Identifier)

	if err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}

	// compare hashed password and request password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(cmd.Password),
	); err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}

	if !user.EmailVerified {
		return TokenPair{}, ErrEmailNotVerified
	}

	// Generate access token with the skey
	access, err := generateAccessToken(user.ID.String())
	if err != nil {
		return TokenPair{}, err
	}

	// generate refresh token
	plainRefresh, refreshHash, refreshExp, err := generateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}

	// Store refresh token to repo
	_, err = s.repo.CreateRefreshToken(ctx, repo.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  refreshExp,
			Valid: true,
		},
	})
	if err != nil {
		return TokenPair{}, err
	}

	// return access and plain refresh token
	return TokenPair{
		AccessToken:  access,
		RefreshToken: plainRefresh,
	}, nil
}

// Access Token Refresh Method

func (s *svc) RefreshTokens(ctx context.Context, refreshToken string) (TokenPair, error) {
	// TODO3: validate refresh tokens & issue new tokens

	if refreshToken == "" {
		return TokenPair{}, ErrInvalidRefreshToken
	}

	hashed := hashRefreshToken(refreshToken)

	dbToken, err := s.repo.GetRefreshTokenByHash(ctx, hashed)

	if err != nil {
		return TokenPair{}, ErrInvalidRefreshToken
	}

	now := time.Now()

	// Check if refresh token was revoked or has already expired then we reject and send error message
	if dbToken.RevokedAt.Valid || dbToken.ExpiresAt.Time.Before(now) {
		return TokenPair{}, ErrInvalidRefreshToken
	}

	access, err := generateAccessToken(dbToken.UserID.String())
	if err != nil {
		return TokenPair{}, err
	}

	// rotate: revoke old token, create new one
	if err := s.repo.RevokeRefreshToken(ctx, dbToken.ID); err != nil {
		return TokenPair{}, err
	}

	// Generate a new refresh token
	newPlain, newHash, newExp, err := generateRefreshToken()
	if err != nil {
		return TokenPair{}, err
	}

	// Store the new refresh token in the db
	_, err = s.repo.CreateRefreshToken(ctx, repo.CreateRefreshTokenParams{
		UserID:    dbToken.UserID,
		TokenHash: newHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  newExp,
			Valid: true,
		},
	})
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  access,
		RefreshToken: newPlain,
	}, nil
}

// Logout Method

func (s *svc) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil // idempotent
	}

	hashed := hashRefreshToken(refreshToken)

	dbToken, err := s.repo.GetRefreshTokenByHash(ctx, hashed)
	if err != nil {
		return nil // already invalid, so treated as success
	}
	// TODO4: revoke token in DB
	return s.repo.RevokeRefreshToken(ctx, dbToken.ID)
}

// Email verification method

func (s *svc) VerifyEmail(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidVerifyToken
	}
	// get the email verification token from the db
	dbToken, err := s.repo.GetEmailVerificationToken(ctx, token)
	if err != nil {
		return ErrInvalidVerifyToken
	}

	now := time.Now()

	if dbToken.UsedAt.Valid {
		return ErrUsedVerifyToken
	} else if dbToken.ExpiresAt.Time.Before(now) {
		return ErrInvalidVerifyToken
	}

	// Mark in auth table that email has been verified after registration
	if err := s.repo.MarkUserEmailVerified(ctx, dbToken.UserID); err != nil {
		return err
	}

	// Mark that the email verification token has already been used
	if err := s.repo.MarkEmailVerificationTokenUsed(ctx, dbToken.ID); err != nil {
		return err
	}
	return nil
}
