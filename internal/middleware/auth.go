package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "userID"

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

func Auth(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			log.Printf("Printing key userID: %v", userIDKey)
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {

				// log.Printf("Authentication Header: %q", authHeader)

				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}

			authHeaderParts := strings.SplitN(authHeader, " ", 2)
			if len(authHeaderParts) != 2 || !strings.EqualFold(authHeaderParts[0], "Bearer") {

				// log.Printf("AUTH: missing or bad prefix")

				http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
				return
			}

			// DEV Comments: Take token then parse and validate

			authTokenStr := authHeaderParts[1]

			token, err := jwt.Parse(authTokenStr, func(token *jwt.Token) (interface{}, error) {
				// Security Check - Verify the Signing Method matches from the payload in the header

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

					// log.Printf("Invalid signature")

					return nil, jwt.ErrSignatureInvalid
				}

				// return the jwtsecret key
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				// log.Printf("Token invalid. Token information: %v", token)
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				// log.Printf("Wrong or invalid claims: %v", token.Claims)
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			sub, ok := claims["sub"].(string)
			if !ok || sub == "" {
				// log.Printf("Missing User ID in token: %v", sub)
				http.Error(w, "missing user id in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
