package users

/*
	Token Helper/Generator Functions
*/

import (
	"contacts/internal/env"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret            = []byte(env.GetStringNoFallback("JWT_SECRET"))
	accessTokenTTL       = time.Duration(env.GetInt("ACCESS_TOKEN_TTL", 15)) * time.Minute
	refreshTokenTTL      = time.Duration(env.GetInt("REFRESH_TOKEN_TTL", 30)) * time.Hour * 24
	emailVerificationTTL = time.Duration(env.GetInt("EMAIL_TOKEN_EXPIRE_MINUTES", 5)) * time.Minute
)

func generateAccessToken(userID string) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(accessTokenTTL).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}

func generateRefreshToken() (plain string, hash string, expiresAt time.Time, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", time.Time{}, err
	}

	plain = hex.EncodeToString(b)
	expiresAt = time.Now().Add(refreshTokenTTL)

	h := sha256.Sum256([]byte(plain + string(jwtSecret)))
	hash = hex.EncodeToString(h[:])

	return
}

func hashRefreshToken(plain string) string {
	h := sha256.Sum256([]byte(plain + string(jwtSecret)))
	return hex.EncodeToString(h[:])
}

// generate email verification token
func generateEmailVerificationToken() (string, time.Time, error) {

	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", time.Time{}, err
	}

	token := hex.EncodeToString(b)
	expiresAt := time.Now().Add(emailVerificationTTL)
	return token, expiresAt, nil
}
