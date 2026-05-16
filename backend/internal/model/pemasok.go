package model

import (
	"database/sql"

	"pos-desktop-tauri/backend/internal/db"
)

var pemasokSortColumns = map[string]string{
	"id":           "id",
	"nama_pemasok": "nama_pemasok",
	"telepon":      "telepon",
}

type Pemasok struct {
	ID          int     `json:"id"`
	NamaPemasok string  `json:"nama_pemasok"`
	Telepon     *string `json:"telepon"`
	Alamat      *string `json:"alamat"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func PemasokAll() ([]Pemasok, error) {
	rows, err := db.Get().Query(`SELECT id, nama_pemasok, telepon, alamat, created_at, updated_at FROM pemasok ORDER BY nama_pemasok`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPemasoks(rows)
}

func PemasokPaginate(page, pageSize int, sortBy, sortOrder, search string) ([]Pemasok, int, error) {
	col := pemasokSortColumns[sortBy]
	if col == "" {
		col = "id"
	}
	dir    := sortDir(sortOrder)
	offset := (page - 1) * pageSize
	like   := "%" + search + "%"

	var total int
	err := db.Get().QueryRow(`SELECT COUNT(*) FROM pemasok WHERE nama_pemasok LIKE ?`, like).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.Get().Query(
		`SELECT id, nama_pemasok, telepon, alamat, created_at, updated_at FROM pemasok WHERE nama_pemasok LIKE ? ORDER BY `+col+` `+dir+` LIMIT ? OFFSET ?`,
		like, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	data, err := scanPemasoks(rows)
	return data, total, err
}

func PemasokFind(id int) (*Pemasok, error) {
	var p Pemasok
	err := db.Get().QueryRow(`SELECT id, nama_pemasok, telepon, alamat, created_at, updated_at FROM pemasok WHERE id = ?`, id).
		Scan(&p.ID, &p.NamaPemasok, &p.Telepon, &p.Alamat, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func PemasokCreate(nama string, telepon, alamat *string) (int64, error) {
	res, err := db.Get().Exec(`INSERT INTO pemasok (nama_pemasok, telepon, alamat) VALUES (?, ?, ?)`, nama, telepon, alamat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func PemasokUpdate(id int, nama string, telepon, alamat *string) error {
	_, err := db.Get().Exec(
		`UPDATE pemasok SET nama_pemasok = ?, telepon = ?, alamat = ?, updated_at = datetime('now') WHERE id = ?`,
		nama, telepon, alamat, id,
	)
	return err
}

func PemasokDestroy(id int) error {
	_, err := db.Get().Exec(`DELETE FROM pemasok WHERE id = ?`, id)
	return err
}

func scanPemasoks(rows *sql.Rows) ([]Pemasok, error) {
	var result []Pemasok
	for rows.Next() {
		var p Pemasok
		if err := rows.Scan(&p.ID, &p.NamaPemasok, &p.Telepon, &p.Alamat, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	if result == nil {
		result = []Pemasok{}
	}
	return result, nil
}
