package middleware

import (
	"net/http"

	"pos-desktop-tauri/backend/internal/helper"
)

func Can(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUser(r)
			if user == nil {
				helper.Forbidden(w)
				return
			}
			for _, p := range user.Permissions {
				if p == permission {
					next.ServeHTTP(w, r)
					return
				}
			}
			helper.Forbidden(w)
		})
	}
}
