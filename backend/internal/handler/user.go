package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"pos-desktop-tauri/backend/internal/helper"
	"pos-desktop-tauri/backend/internal/middleware"
	"pos-desktop-tauri/backend/internal/model"
)

func UserIndex(w http.ResponseWriter, r *http.Request) {
	page, _     := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 10 }

	data, total, err := model.UserPaginate(page, pageSize,
		r.URL.Query().Get("sortBy"),
		r.URL.Query().Get("sortOrder"),
		r.URL.Query().Get("search"),
	)
	if err != nil { helper.ServerError(w, ""); return }
	helper.Success(w, map[string]any{"data": data, "total": total})
}

func UserStore(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		RoleID   int    `json:"role_id"`
		Active   *int   `json:"active"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if body.Username == "" || body.Password == "" || body.RoleID == 0 {
		helper.BadRequest(w, "Username, password, dan role wajib diisi")
		return
	}
	active := 1
	if body.Active != nil { active = *body.Active }

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil { helper.ServerError(w, ""); return }

	id, err := model.UserCreate(body.Username, string(hash), body.RoleID, active)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			helper.BadRequest(w, "Username sudah digunakan")
		} else {
			helper.ServerError(w, "")
		}
		return
	}
	helper.Created(w, map[string]any{"id": id})
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	u, err := model.UserFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if u == nil   { helper.NotFound(w, ""); return }

	var body struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
		RoleID   *int    `json:"role_id"`
		Active   *int    `json:"active"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	fields := map[string]any{}
	if body.Username != nil { fields["username"] = *body.Username }
	if body.RoleID   != nil { fields["role_id"]  = *body.RoleID }
	if body.Active   != nil { fields["active"]   = *body.Active }
	if body.Password != nil && *body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*body.Password), bcrypt.DefaultCost)
		if err != nil { helper.ServerError(w, ""); return }
		fields["password"] = string(hash)
	}

	if err := model.UserUpdate(id, fields); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			helper.BadRequest(w, "Username sudah digunakan")
		} else {
			helper.ServerError(w, "")
		}
		return
	}
	helper.Success(w, nil)
}

func UserDestroy(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	caller := middleware.GetUser(r)
	if caller != nil && caller.ID == id {
		helper.BadRequest(w, "Tidak bisa menghapus akun sendiri")
		return
	}
	u, err := model.UserFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if u == nil   { helper.NotFound(w, ""); return }
	if err := model.UserDestroy(id); err != nil { helper.ServerError(w, ""); return }
	helper.NoContent(w)
}
