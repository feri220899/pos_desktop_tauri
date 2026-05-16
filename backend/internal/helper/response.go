package helper

import (
	"encoding/json"
	"net/http"
)

type envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func write(w http.ResponseWriter, status int, body envelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func Success(w http.ResponseWriter, data any) {
	write(w, http.StatusOK, envelope{Success: true, Data: data})
}

func Created(w http.ResponseWriter, data any) {
	write(w, http.StatusCreated, envelope{Success: true, Data: data})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func BadRequest(w http.ResponseWriter, msg string) {
	write(w, http.StatusBadRequest, envelope{Success: false, Message: msg})
}

func NotFound(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = "Data tidak ditemukan"
	}
	write(w, http.StatusNotFound, envelope{Success: false, Message: msg})
}

func Unauthorized(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = "Unauthorized"
	}
	write(w, http.StatusUnauthorized, envelope{Success: false, Message: msg})
}

func Forbidden(w http.ResponseWriter) {
	write(w, http.StatusForbidden, envelope{Success: false, Message: "Akses ditolak"})
}

func ServerError(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = "Terjadi kesalahan server"
	}
	write(w, http.StatusInternalServerError, envelope{Success: false, Message: msg})
}
