package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"pos-desktop-tauri/backend/internal/helper"
	"pos-desktop-tauri/backend/internal/model"
)

func TokoShow(w http.ResponseWriter, r *http.Request) {
	t, err := model.TokoGet()
	if err != nil { helper.ServerError(w, ""); return }
	helper.Success(w, t)
}

func TokoUpdate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Nama    string `json:"nama"`
		Telepon string `json:"telepon"`
		Email   string `json:"email"`
		Alamat  string `json:"alamat"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if strings.TrimSpace(body.Nama) == "" {
		helper.BadRequest(w, "Nama toko wajib diisi")
		return
	}
	if err := model.TokoUpdate(
		strings.TrimSpace(body.Nama),
		strings.TrimSpace(body.Telepon),
		strings.TrimSpace(body.Email),
		strings.TrimSpace(body.Alamat),
	); err != nil {
		helper.ServerError(w, "")
		return
	}
	t, _ := model.TokoGet()
	helper.Success(w, t)
}
