package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"pos-desktop-tauri/backend/internal/helper"
	"pos-desktop-tauri/backend/internal/model"
)

func PemasokIndex(w http.ResponseWriter, r *http.Request) {
	page, _     := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 10 }
	sortBy    := r.URL.Query().Get("sortBy")
	sortOrder := r.URL.Query().Get("sortOrder")
	search    := r.URL.Query().Get("search")

	data, total, err := model.PemasokPaginate(page, pageSize, sortBy, sortOrder, search)
	if err != nil {
		helper.ServerError(w, "")
		return
	}
	helper.Success(w, map[string]any{"data": data, "total": total})
}

func PemasokStore(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NamaPemasok string  `json:"nama_pemasok"`
		Telepon     *string `json:"telepon"`
		Alamat      *string `json:"alamat"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if body.NamaPemasok == "" {
		helper.BadRequest(w, "nama_pemasok wajib diisi")
		return
	}
	id, err := model.PemasokCreate(body.NamaPemasok, body.Telepon, body.Alamat)
	if err != nil {
		helper.ServerError(w, "")
		return
	}
	helper.Created(w, map[string]any{"id": id})
}

func PemasokUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := model.PemasokFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if p == nil   { helper.NotFound(w, ""); return }

	var body struct {
		NamaPemasok string  `json:"nama_pemasok"`
		Telepon     *string `json:"telepon"`
		Alamat      *string `json:"alamat"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if body.NamaPemasok == "" {
		helper.BadRequest(w, "nama_pemasok wajib diisi")
		return
	}
	if err := model.PemasokUpdate(id, body.NamaPemasok, body.Telepon, body.Alamat); err != nil {
		helper.ServerError(w, "")
		return
	}
	helper.Success(w, nil)
}

func PemasokDestroy(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := model.PemasokFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if p == nil   { helper.NotFound(w, ""); return }

	if err := model.PemasokDestroy(id); err != nil {
		helper.ServerError(w, "")
		return
	}
	helper.NoContent(w)
}
