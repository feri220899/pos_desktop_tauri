package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"pos-desktop-tauri/backend/internal/helper"
	"pos-desktop-tauri/backend/internal/model"
)

func RoleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := model.RoleAll()
	if err != nil { helper.ServerError(w, ""); return }
	helper.Success(w, data)
}

func RoleStore(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if strings.TrimSpace(body.Name) == "" {
		helper.BadRequest(w, "Nama role wajib diisi")
		return
	}
	id, err := model.RoleCreate(strings.TrimSpace(body.Name), body.Permissions)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			helper.BadRequest(w, "Nama role sudah digunakan")
		} else {
			helper.ServerError(w, "")
		}
		return
	}
	helper.Created(w, map[string]any{"id": id})
}

func RoleUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	role, err := model.RoleFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if role == nil { helper.NotFound(w, ""); return }

	var body struct {
		Permissions []string `json:"permissions"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if body.Permissions == nil {
		helper.BadRequest(w, "Permissions harus berupa array")
		return
	}
	if err := model.RoleUpdate(id, body.Permissions); err != nil {
		helper.ServerError(w, "")
		return
	}
	helper.Success(w, nil)
}
