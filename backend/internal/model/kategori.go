package model

import (
	"database/sql"
	"strings"

	"pos-desktop-tauri/backend/internal/db"
)

var kategoriSortColumns = map[string]string{
	"id":            "id",
	"nama_kategori": "nama_kategori",
}

type Kategori struct {
	ID           int    `json:"id"`
	NamaKategori string `json:"nama_kategori"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func KategoriAll() ([]Kategori, error) {
	rows, err := db.Get().Query(`SELECT id, nama_kategori, created_at, updated_at FROM kategori ORDER BY nama_kategori`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanKategoris(rows)
}

func KategoriPaginate(page, pageSize int, sortBy, sortOrder, search string) ([]Kategori, int, error) {
	col := kategoriSortColumns[sortBy]
	if col == "" {
		col = "id"
	}
	dir    := sortDir(sortOrder)
	offset := (page - 1) * pageSize
	like   := "%" + search + "%"

	var total int
	err := db.Get().QueryRow(`SELECT COUNT(*) FROM kategori WHERE nama_kategori LIKE ?`, like).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.Get().Query(
		`SELECT id, nama_kategori, created_at, updated_at FROM kategori WHERE nama_kategori LIKE ? ORDER BY `+col+` `+dir+` LIMIT ? OFFSET ?`,
		like, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	data, err := scanKategoris(rows)
	return data, total, err
}

func KategoriFind(id int) (*Kategori, error) {
	var k Kategori
	err := db.Get().QueryRow(`SELECT id, nama_kategori, created_at, updated_at FROM kategori WHERE id = ?`, id).
		Scan(&k.ID, &k.NamaKategori, &k.CreatedAt, &k.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &k, err
}

func KategoriCreate(namaKategori string) (int64, error) {
	res, err := db.Get().Exec(`INSERT INTO kategori (nama_kategori) VALUES (?)`, namaKategori)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func KategoriUpdate(id int, namaKategori string) error {
	_, err := db.Get().Exec(`UPDATE kategori SET nama_kategori = ?, updated_at = datetime('now') WHERE id = ?`, namaKategori, id)
	return err
}

func KategoriDestroy(id int) error {
	_, err := db.Get().Exec(`DELETE FROM kategori WHERE id = ?`, id)
	return err
}

func scanKategoris(rows *sql.Rows) ([]Kategori, error) {
	var result []Kategori
	for rows.Next() {
		var k Kategori
		if err := rows.Scan(&k.ID, &k.NamaKategori, &k.CreatedAt, &k.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, k)
	}
	if result == nil {
		result = []Kategori{}
	}
	return result, nil
}

var _ = strings.ToLower
