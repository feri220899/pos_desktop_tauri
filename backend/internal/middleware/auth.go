package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"pos-desktop-tauri/backend/internal/helper"
)

type contextKey string

const UserKey contextKey = "user"

type UserClaims struct {
	ID          int      `json:"id"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func jwtSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "p0s-d3skt0p-jwt-s3cr3t-k3y-2024"
	}
	return []byte(s)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			helper.Unauthorized(w, "")
			return
		}

		tokenStr := header[7:]
		claims   := &UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return jwtSecret(), nil
		})
		if err != nil || !token.Valid {
			helper.Unauthorized(w, "Token tidak valid atau kadaluarsa")
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(r *http.Request) *UserClaims {
	u, _ := r.Context().Value(UserKey).(*UserClaims)
	return u
}
