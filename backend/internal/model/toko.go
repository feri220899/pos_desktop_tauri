package model

import "pos-desktop-tauri/backend/internal/db"

type Toko struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama"`
	Telepon   string `json:"telepon"`
	Email     string `json:"email"`
	Alamat    string `json:"alamat"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TokoGet() (*Toko, error) {
	var t Toko
	err := db.Get().QueryRow(`SELECT id, nama, telepon, email, alamat, created_at, updated_at FROM toko WHERE id = 1`).
		Scan(&t.ID, &t.Nama, &t.Telepon, &t.Email, &t.Alamat, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}

func TokoUpdate(nama, telepon, email, alamat string) error {
	_, err := db.Get().Exec(`
		UPDATE toko SET nama = ?, telepon = ?, email = ?, alamat = ?, updated_at = datetime('now')
		WHERE id = 1
	`, nama, telepon, email, alamat)
	return err
}
