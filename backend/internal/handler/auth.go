package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"pos-desktop-tauri/backend/internal/helper"
	"pos-desktop-tauri/backend/internal/middleware"
	"pos-desktop-tauri/backend/internal/model"
)

func jwtSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "p0s-d3skt0p-jwt-s3cr3t-k3y-2024"
	}
	return []byte(s)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Username == "" || body.Password == "" {
		helper.BadRequest(w, "Username dan password wajib diisi")
		return
	}

	user, hash, err := model.UserFindByUsername(body.Username)
	if err != nil {
		helper.ServerError(w, "")
		return
	}
	if user == nil || user.Active == 0 {
		helper.BadRequest(w, "Username atau password salah")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.Password)) != nil {
		helper.BadRequest(w, "Username atau password salah")
		return
	}

	claims := middleware.UserClaims{
		ID:          user.ID,
		Username:    user.Username,
		Role:        user.RoleName,
		Permissions: user.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret())
	if err != nil {
		helper.ServerError(w, "")
		return
	}

	payload := map[string]any{
		"id":          user.ID,
		"username":    user.Username,
		"role":        user.RoleName,
		"permissions": user.Permissions,
	}
	helper.Success(w, map[string]any{"token": token, "user": payload})
}

func Me(w http.ResponseWriter, r *http.Request) {
	helper.Success(w, middleware.GetUser(r))
}
