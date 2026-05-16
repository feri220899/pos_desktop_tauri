package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"pos-desktop-tauri/backend/internal/helper"
	"pos-desktop-tauri/backend/internal/model"
)

func ProdukIndex(w http.ResponseWriter, r *http.Request) {
	page, _     := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 10 }

	data, total, err := model.ProdukPaginate(page, pageSize,
		r.URL.Query().Get("sortBy"),
		r.URL.Query().Get("sortOrder"),
		r.URL.Query().Get("search"),
	)
	if err != nil { helper.ServerError(w, ""); return }
	helper.Success(w, map[string]any{"data": data, "total": total})
}

func ProdukShow(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := model.ProdukFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if p == nil   { helper.NotFound(w, ""); return }
	helper.Success(w, p)
}

func ProdukStore(w http.ResponseWriter, r *http.Request) {
	var input model.ProdukInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.BadRequest(w, "Format data tidak valid")
		return
	}
	if input.TipeKepemilikan == "MILIK_SENDIRI" && input.HargaBeli == 0 {
		helper.BadRequest(w, "Harga beli wajib untuk produk MILIK_SENDIRI")
		return
	}
	if input.TipeKepemilikan == "TITIPAN" && input.HargaTitip == 0 {
		helper.BadRequest(w, "Harga titip wajib untuk produk TITIPAN")
		return
	}
	if input.HargaJual == 0 {
		helper.BadRequest(w, "Harga jual wajib diisi")
		return
	}
	id, err := model.ProdukCreate(input)
	if err != nil { helper.ServerError(w, err.Error()); return }
	helper.Created(w, map[string]any{"produk_id": id})
}

func ProdukUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := model.ProdukFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if p == nil   { helper.NotFound(w, ""); return }

	var input model.ProdukUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.BadRequest(w, "Format data tidak valid")
		return
	}
	if input.TipeKepemilikan == "MILIK_SENDIRI" && input.HargaBeli == 0 {
		helper.BadRequest(w, "Harga beli wajib untuk produk MILIK_SENDIRI")
		return
	}
	if input.TipeKepemilikan == "TITIPAN" && input.HargaTitip == 0 {
		helper.BadRequest(w, "Harga titip wajib untuk produk TITIPAN")
		return
	}
	if input.HargaJual == 0 {
		helper.BadRequest(w, "Harga jual wajib diisi")
		return
	}
	if err := model.ProdukUpdate(id, input); err != nil {
		helper.ServerError(w, err.Error())
		return
	}
	helper.Success(w, nil)
}

func ProdukDestroy(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := model.ProdukFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if p == nil   { helper.NotFound(w, ""); return }
	if err := model.ProdukDestroy(id); err != nil { helper.ServerError(w, ""); return }
	helper.NoContent(w)
}

func ProdukRestore(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if err := model.ProdukRestore(id); err != nil { helper.ServerError(w, ""); return }
	helper.Success(w, nil)
}

func ProdukTrashed(w http.ResponseWriter, r *http.Request) {
	page, _     := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 10 }
	data, total, err := model.ProdukTrashed(page, pageSize)
	if err != nil { helper.ServerError(w, ""); return }
	helper.Success(w, map[string]any{"data": data, "total": total})
}

func ProdukNextKode(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	if prefix == "" { prefix = "A" }
	helper.Success(w, map[string]any{"kode_produk": model.ProdukGenerateKode(prefix)})
}

func ProdukGetSatuan(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := model.ProdukFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if p == nil   { helper.NotFound(w, ""); return }
	data, err := model.ProdukGetSatuan(id)
	if err != nil { helper.ServerError(w, ""); return }
	helper.Success(w, data)
}

func ProdukStoreSatuan(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := model.ProdukFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if p == nil   { helper.NotFound(w, ""); return }

	var input model.SatuanInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.BadRequest(w, "Format tidak valid")
		return
	}
	satuanID, err := model.ProdukStoreSatuan(id, input)
	if err != nil { helper.ServerError(w, err.Error()); return }
	helper.Created(w, map[string]any{"satuan_id": satuanID})
}

func ProdukUpdateSatuan(w http.ResponseWriter, r *http.Request) {
	id, _       := strconv.Atoi(chi.URLParam(r, "id"))
	satuanID, _ := strconv.Atoi(chi.URLParam(r, "satuanId"))
	p, err := model.ProdukFind(id)
	if err != nil { helper.ServerError(w, ""); return }
	if p == nil   { helper.NotFound(w, ""); return }

	var input model.SatuanInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.BadRequest(w, "Format tidak valid")
		return
	}
	if err := model.ProdukUpdateSatuan(id, satuanID, input); err != nil {
		helper.ServerError(w, err.Error())
		return
	}
	helper.Success(w, nil)
}

func ProdukDestroySatuan(w http.ResponseWriter, r *http.Request) {
	id, _       := strconv.Atoi(chi.URLParam(r, "id"))
	satuanID, _ := strconv.Atoi(chi.URLParam(r, "satuanId"))
	err := model.ProdukDestroySatuan(id, satuanID)
	if err != nil {
		if err.Error() == "Minimal harus ada 1 satuan" {
			helper.BadRequest(w, err.Error())
		} else {
			helper.ServerError(w, err.Error())
		}
		return
	}
	helper.NoContent(w)
}
