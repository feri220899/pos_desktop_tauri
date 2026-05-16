package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"pos-desktop-tauri/backend/internal/helper"
	"pos-desktop-tauri/backend/internal/model"
)

func KategoriIndex(w http.ResponseWriter, r *http.Request) {
	page, _     := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 10 }
	sortBy    := r.URL.Query().Get("sortBy")
	sortOrder := r.URL.Query().Get("sortOrder")
	search    := r.URL.Query().Get("search")

	data, total, err := model.KategoriPaginate(page, pageSize, sortBy, sortOrder, search)
	if err != nil {
		helper.ServerError(w, "")
		return
	}
	helper.Success(w, map[string]any{"data": data, "total": total})
}

func KategoriStore(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NamaKategori string `json:"nama_kategori"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if body.NamaKategori == "" {
		helper.BadRequest(w, "nama_kategori wajib diisi")
		return
	}
	id, err := model.KategoriCreate(body.NamaKategori)
	if err != nil {
		helper.ServerError(w, "")
		return
	}
	helper.Created(w, map[string]any{"id": id})
}

func KategoriUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	k, err := model.KategoriFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if k == nil   { helper.NotFound(w, ""); return }

	var body struct {
		NamaKategori string `json:"nama_kategori"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if body.NamaKategori == "" {
		helper.BadRequest(w, "nama_kategori wajib diisi")
		return
	}
	if err := model.KategoriUpdate(id, body.NamaKategori); err != nil {
		helper.ServerError(w, "")
		return
	}
	helper.Success(w, nil)
}

func KategoriDestroy(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	k, err := model.KategoriFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if k == nil   { helper.NotFound(w, ""); return }

	if err := model.KategoriDestroy(id); err != nil {
		helper.ServerError(w, "")
		return
	}
	helper.NoContent(w)
}
