package handler

import (
	"net/http"

	"pos-desktop-tauri/backend/internal/helper"
	"pos-desktop-tauri/backend/internal/model"
)

func PermissionIndex(w http.ResponseWriter, r *http.Request) {
	data, err := model.PermissionAll()
	if err != nil { helper.ServerError(w, ""); return }
	helper.Success(w, data)
}
