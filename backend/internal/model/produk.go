package model

import (
	"database/sql"
	"fmt"
	"strings"

	"pos-desktop-tauri/backend/internal/db"
)

var produkSortColumns = map[string]string{
	"id":            "p.id",
	"kode_produk":   "p.kode_produk",
	"nama_produk":   "p.nama_produk",
	"nama_kategori": "k.nama_kategori",
	"jumlah_stok":   "sp.jumlah_stok",
	"harga_jual":    "h.harga_jual",
}

const produkBaseSelect = `
	SELECT
		p.id, p.kode_produk, p.nama_produk, p.tipe_kepemilikan,
		p.kategori_id, p.pemasok_id, p.deskripsi, p.status_aktif,
		p.deleted_at, p.created_at, p.updated_at,
		COALESCE(k.nama_kategori, ''),
		COALESCE(pe.nama_pemasok, ''),
		COALESCE(sp.jumlah_stok, 0),
		COALESCE(su.nama_satuan, 'PCS'),
		COALESCE(su.jenis_input, 'PCS'),
		COALESCE(su.nilai_konversi, 1),
		COALESCE(su.id, 0),
		COALESCE(bp.kode_barcode, ''),
		COALESCE(h.harga_beli, 0),
		COALESCE(h.harga_jual, 0),
		COALESCE(kp.harga_hak_pemasok, 0)
	FROM produk p
	LEFT JOIN kategori k         ON k.id = p.kategori_id
	LEFT JOIN pemasok pe         ON pe.id = p.pemasok_id
	LEFT JOIN stok_produk sp     ON sp.produk_id = p.id
	LEFT JOIN satuan_produk su   ON su.produk_id = p.id AND su.satuan_utama = 1
	LEFT JOIN barcode_produk bp  ON bp.satuan_produk_id = su.id
	LEFT JOIN (
		SELECT h1.*
		FROM harga_produk h1
		INNER JOIN (
			SELECT satuan_produk_id, MAX(berlaku_sejak) AS max_berlaku
			FROM harga_produk GROUP BY satuan_produk_id
		) h2 ON h1.satuan_produk_id = h2.satuan_produk_id AND h1.berlaku_sejak = h2.max_berlaku
	) h ON h.satuan_produk_id = su.id
	LEFT JOIN konsinyasi_produk kp ON kp.satuan_produk_id = su.id
`

type Produk struct {
	ID              int      `json:"id"`
	KodeProduk      string   `json:"kode_produk"`
	NamaProduk      string   `json:"nama_produk"`
	TipeKepemilikan string   `json:"tipe_kepemilikan"`
	KategoriID      *int     `json:"kategori_id"`
	PemasokID       *int     `json:"pemasok_id"`
	Deskripsi       *string  `json:"deskripsi"`
	StatusAktif     int      `json:"status_aktif"`
	DeletedAt       *string  `json:"deleted_at"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	NamaKategori    string   `json:"nama_kategori"`
	NamaPemasok     string   `json:"nama_pemasok"`
	JumlahStok      float64  `json:"jumlah_stok"`
	NamaSatuan      string   `json:"nama_satuan"`
	JenisInput      string   `json:"jenis_input"`
	NilaiKonversi   float64  `json:"nilai_konversi"`
	SatuanID        int      `json:"satuan_id"`
	KodeBarcode     string   `json:"kode_barcode"`
	HargaBeli       float64  `json:"harga_beli"`
	HargaJual       float64  `json:"harga_jual"`
	HargaTitip      float64  `json:"harga_titip"`
}

type SatuanProduk struct {
	ID            int     `json:"id"`
	NamaSatuan    string  `json:"nama_satuan"`
	JenisInput    string  `json:"jenis_input"`
	NilaiKonversi float64 `json:"nilai_konversi"`
	SatuanUtama   int     `json:"satuan_utama"`
	KodeBarcode   string  `json:"kode_barcode"`
	HargaBeli     float64 `json:"harga_beli"`
	HargaJual     float64 `json:"harga_jual"`
	HargaTitip    float64 `json:"harga_titip"`
}

func scanProduk(row interface{ Scan(...any) error }) (*Produk, error) {
	var p Produk
	err := row.Scan(
		&p.ID, &p.KodeProduk, &p.NamaProduk, &p.TipeKepemilikan,
		&p.KategoriID, &p.PemasokID, &p.Deskripsi, &p.StatusAktif,
		&p.DeletedAt, &p.CreatedAt, &p.UpdatedAt,
		&p.NamaKategori, &p.NamaPemasok,
		&p.JumlahStok, &p.NamaSatuan, &p.JenisInput, &p.NilaiKonversi, &p.SatuanID,
		&p.KodeBarcode, &p.HargaBeli, &p.HargaJual, &p.HargaTitip,
	)
	return &p, err
}

func ProdukPaginate(page, pageSize int, sortBy, sortOrder, search string) ([]Produk, int, error) {
	col := produkSortColumns[sortBy]
	if col == "" {
		col = "p.id"
	}
	dir    := sortDir(sortOrder)
	offset := (page - 1) * pageSize
	like   := "%" + search + "%"

	var total int
	err := db.Get().QueryRow(`
		SELECT COUNT(*) FROM produk p
		LEFT JOIN satuan_produk su ON su.produk_id = p.id AND su.satuan_utama = 1
		LEFT JOIN barcode_produk bp ON bp.satuan_produk_id = su.id
		WHERE (p.nama_produk LIKE ? OR p.kode_produk LIKE ? OR bp.kode_barcode LIKE ?) AND p.deleted_at IS NULL
	`, like, like, like).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.Get().Query(
		produkBaseSelect+` WHERE (p.nama_produk LIKE ? OR p.kode_produk LIKE ? OR bp.kode_barcode LIKE ?) AND p.deleted_at IS NULL ORDER BY `+col+` `+dir+` LIMIT ? OFFSET ?`,
		like, like, like, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []Produk
	for rows.Next() {
		p, err := scanProduk(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, *p)
	}
	if result == nil {
		result = []Produk{}
	}
	return result, total, nil
}

func ProdukFind(id int) (*Produk, error) {
	p, err := scanProduk(db.Get().QueryRow(produkBaseSelect+` WHERE p.id = ? AND p.deleted_at IS NULL`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func ProdukTrashed(page, pageSize int) ([]Produk, int, error) {
	offset := (page - 1) * pageSize
	var total int
	db.Get().QueryRow(`SELECT COUNT(*) FROM produk WHERE deleted_at IS NOT NULL`).Scan(&total)

	rows, err := db.Get().Query(
		produkBaseSelect+` WHERE p.deleted_at IS NOT NULL ORDER BY p.deleted_at DESC LIMIT ? OFFSET ?`,
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []Produk
	for rows.Next() {
		p, err := scanProduk(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, *p)
	}
	if result == nil {
		result = []Produk{}
	}
	return result, total, nil
}

func ProdukGenerateKode(prefix string) string {
	if len(prefix) == 0 {
		prefix = "A"
	}
	prefix = strings.ToUpper(string(prefix[0]))

	var maxNum sql.NullInt64
	db.Get().QueryRow(`
		SELECT MAX(CAST(SUBSTR(kode_produk, 2) AS INTEGER))
		FROM produk
		WHERE UPPER(SUBSTR(kode_produk, 1, 1)) = ? AND CAST(SUBSTR(kode_produk, 2) AS INTEGER) > 0
	`, prefix).Scan(&maxNum)

	next := int64(1)
	if maxNum.Valid {
		next = maxNum.Int64 + 1
	}
	return fmt.Sprintf("%s%06d", prefix, next)
}

type ProdukInput struct {
	KodeProduk      string   `json:"kode_produk"`
	NamaProduk      string   `json:"nama_produk"`
	KategoriID      *int     `json:"kategori_id"`
	TipeKepemilikan string   `json:"tipe_kepemilikan"`
	PemasokID       *int     `json:"pemasok_id"`
	Deskripsi       *string  `json:"deskripsi"`
	StatusAktif     int      `json:"status_aktif"`
	NamaSatuan      string   `json:"nama_satuan"`
	JenisInput      string   `json:"jenis_input"`
	NilaiKonversi   float64  `json:"nilai_konversi"`
	KodeBarcode     string   `json:"kode_barcode"`
	StokAwal        float64  `json:"stok_awal"`
	HargaBeli       float64  `json:"harga_beli"`
	HargaJual       float64  `json:"harga_jual"`
	HargaTitip      float64  `json:"harga_titip"`
}

func ProdukCreate(input ProdukInput) (int64, error) {
	tx, err := db.Get().Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	res, err := tx.Exec(`
		INSERT INTO produk (kode_produk, nama_produk, kategori_id, tipe_kepemilikan, pemasok_id, deskripsi, status_aktif)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, input.KodeProduk, input.NamaProduk, input.KategoriID, input.TipeKepemilikan, input.PemasokID, input.Deskripsi, input.StatusAktif)
	if err != nil {
		return 0, err
	}
	produkID, _ := res.LastInsertId()

	res, err = tx.Exec(`
		INSERT INTO satuan_produk (produk_id, nama_satuan, jenis_input, nilai_konversi, satuan_utama)
		VALUES (?, ?, ?, ?, 1)
	`, produkID, strings.ToUpper(input.NamaSatuan), input.JenisInput, input.NilaiKonversi)
	if err != nil {
		return 0, err
	}
	satuanID, _ := res.LastInsertId()

	if input.KodeBarcode != "" {
		_, err = tx.Exec(`INSERT INTO barcode_produk (satuan_produk_id, kode_barcode) VALUES (?, ?)`, satuanID, strings.TrimSpace(input.KodeBarcode))
		if err != nil {
			return 0, err
		}
	}

	hargaBeli := float64(0)
	if input.TipeKepemilikan == "MILIK_SENDIRI" {
		hargaBeli = input.HargaBeli
	}
	_, err = tx.Exec(`INSERT INTO harga_produk (satuan_produk_id, harga_beli, harga_jual, berlaku_sejak) VALUES (?, ?, ?, datetime('now'))`,
		satuanID, hargaBeli, input.HargaJual)
	if err != nil {
		return 0, err
	}

	if input.TipeKepemilikan == "TITIPAN" {
		komisi := input.HargaJual - input.HargaTitip
		if komisi < 0 {
			komisi = 0
		}
		_, err = tx.Exec(`INSERT INTO konsinyasi_produk (satuan_produk_id, harga_hak_pemasok, komisi_toko) VALUES (?, ?, ?)`,
			satuanID, input.HargaTitip, komisi)
		if err != nil {
			return 0, err
		}
	}

	_, err = tx.Exec(`INSERT INTO stok_produk (produk_id, jumlah_stok) VALUES (?, ?)`, produkID, input.StokAwal)
	if err != nil {
		return 0, err
	}

	return produkID, tx.Commit()
}

type ProdukUpdateInput struct {
	KodeProduk      string   `json:"kode_produk"`
	NamaProduk      string   `json:"nama_produk"`
	KategoriID      *int     `json:"kategori_id"`
	TipeKepemilikan string   `json:"tipe_kepemilikan"`
	PemasokID       *int     `json:"pemasok_id"`
	Deskripsi       *string  `json:"deskripsi"`
	StatusAktif     int      `json:"status_aktif"`
	NamaSatuan      string   `json:"nama_satuan"`
	JenisInput      string   `json:"jenis_input"`
	NilaiKonversi   float64  `json:"nilai_konversi"`
	KodeBarcode     string   `json:"kode_barcode"`
	JumlahStok      float64  `json:"jumlah_stok"`
	HargaBeli       float64  `json:"harga_beli"`
	HargaJual       float64  `json:"harga_jual"`
	HargaTitip      float64  `json:"harga_titip"`
}

func ProdukUpdate(id int, input ProdukUpdateInput) error {
	tx, err := db.Get().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		UPDATE produk SET kode_produk = ?, nama_produk = ?, kategori_id = ?, tipe_kepemilikan = ?,
			pemasok_id = ?, deskripsi = ?, status_aktif = ?, updated_at = datetime('now')
		WHERE id = ?
	`, input.KodeProduk, input.NamaProduk, input.KategoriID, input.TipeKepemilikan, input.PemasokID, input.Deskripsi, input.StatusAktif, id)
	if err != nil {
		return err
	}

	var satuanID int64
	err = tx.QueryRow(`SELECT id FROM satuan_produk WHERE produk_id = ? AND satuan_utama = 1`, id).Scan(&satuanID)
	if err == sql.ErrNoRows {
		res, e := tx.Exec(`
			INSERT INTO satuan_produk (produk_id, nama_satuan, jenis_input, nilai_konversi, satuan_utama)
			VALUES (?, ?, ?, ?, 1)
		`, id, strings.ToUpper(input.NamaSatuan), input.JenisInput, input.NilaiKonversi)
		if e != nil {
			err = e
			return err
		}
		satuanID, _ = res.LastInsertId()
	} else if err != nil {
		return err
	} else {
		_, err = tx.Exec(`
			UPDATE satuan_produk SET nama_satuan = ?, jenis_input = ?, nilai_konversi = ?, updated_at = datetime('now')
			WHERE id = ?
		`, strings.ToUpper(input.NamaSatuan), input.JenisInput, input.NilaiKonversi, satuanID)
		if err != nil {
			return err
		}
	}

	barcode := strings.TrimSpace(input.KodeBarcode)
	var barcodeRowID sql.NullInt64
	tx.QueryRow(`SELECT id FROM barcode_produk WHERE satuan_produk_id = ?`, satuanID).Scan(&barcodeRowID)
	if barcode != "" {
		if barcodeRowID.Valid {
			_, err = tx.Exec(`UPDATE barcode_produk SET kode_barcode = ?, updated_at = datetime('now') WHERE id = ?`, barcode, barcodeRowID.Int64)
		} else {
			_, err = tx.Exec(`INSERT INTO barcode_produk (satuan_produk_id, kode_barcode) VALUES (?, ?)`, satuanID, barcode)
		}
		if err != nil {
			return err
		}
	} else if barcodeRowID.Valid {
		_, err = tx.Exec(`DELETE FROM barcode_produk WHERE id = ?`, barcodeRowID.Int64)
		if err != nil {
			return err
		}
	}

	hargaBeli := float64(0)
	if input.TipeKepemilikan == "MILIK_SENDIRI" {
		hargaBeli = input.HargaBeli
	}
	_, err = tx.Exec(`INSERT INTO harga_produk (satuan_produk_id, harga_beli, harga_jual, berlaku_sejak) VALUES (?, ?, ?, datetime('now'))`,
		satuanID, hargaBeli, input.HargaJual)
	if err != nil {
		return err
	}

	if input.TipeKepemilikan == "TITIPAN" {
		komisi := input.HargaJual - input.HargaTitip
		if komisi < 0 {
			komisi = 0
		}
		var kpID sql.NullInt64
		tx.QueryRow(`SELECT id FROM konsinyasi_produk WHERE satuan_produk_id = ?`, satuanID).Scan(&kpID)
		if kpID.Valid {
			_, err = tx.Exec(`UPDATE konsinyasi_produk SET harga_hak_pemasok = ?, komisi_toko = ?, updated_at = datetime('now') WHERE satuan_produk_id = ?`,
				input.HargaTitip, komisi, satuanID)
		} else {
			_, err = tx.Exec(`INSERT INTO konsinyasi_produk (satuan_produk_id, harga_hak_pemasok, komisi_toko) VALUES (?, ?, ?)`,
				satuanID, input.HargaTitip, komisi)
		}
		if err != nil {
			return err
		}
	} else {
		_, err = tx.Exec(`DELETE FROM konsinyasi_produk WHERE satuan_produk_id = ?`, satuanID)
		if err != nil {
			return err
		}
	}

	var stokID sql.NullInt64
	tx.QueryRow(`SELECT id FROM stok_produk WHERE produk_id = ?`, id).Scan(&stokID)
	if stokID.Valid {
		_, err = tx.Exec(`UPDATE stok_produk SET jumlah_stok = ?, updated_at = datetime('now') WHERE produk_id = ?`, input.JumlahStok, id)
	} else {
		_, err = tx.Exec(`INSERT INTO stok_produk (produk_id, jumlah_stok) VALUES (?, ?)`, id, input.JumlahStok)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

func ProdukDestroy(id int) error {
	_, err := db.Get().Exec(`UPDATE produk SET deleted_at = datetime('now'), updated_at = datetime('now') WHERE id = ?`, id)
	return err
}

func ProdukRestore(id int) error {
	_, err := db.Get().Exec(`UPDATE produk SET deleted_at = NULL, updated_at = datetime('now') WHERE id = ?`, id)
	return err
}

func ProdukGetSatuan(produkID int) ([]SatuanProduk, error) {
	rows, err := db.Get().Query(`
		SELECT su.id, su.nama_satuan, su.jenis_input, su.nilai_konversi, su.satuan_utama,
		       COALESCE(bp.kode_barcode, ''),
		       COALESCE(h.harga_beli, 0), COALESCE(h.harga_jual, 0), COALESCE(kp.harga_hak_pemasok, 0)
		FROM satuan_produk su
		LEFT JOIN barcode_produk bp ON bp.satuan_produk_id = su.id
		LEFT JOIN (
			SELECT h1.* FROM harga_produk h1
			INNER JOIN (
				SELECT satuan_produk_id, MAX(berlaku_sejak) AS max_berlaku
				FROM harga_produk GROUP BY satuan_produk_id
			) h2 ON h1.satuan_produk_id = h2.satuan_produk_id AND h1.berlaku_sejak = h2.max_berlaku
		) h ON h.satuan_produk_id = su.id
		LEFT JOIN konsinyasi_produk kp ON kp.satuan_produk_id = su.id
		WHERE su.produk_id = ? ORDER BY su.satuan_utama DESC, su.nama_satuan
	`, produkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SatuanProduk
	for rows.Next() {
		var s SatuanProduk
		if err := rows.Scan(&s.ID, &s.NamaSatuan, &s.JenisInput, &s.NilaiKonversi, &s.SatuanUtama,
			&s.KodeBarcode, &s.HargaBeli, &s.HargaJual, &s.HargaTitip); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	if result == nil {
		result = []SatuanProduk{}
	}
	return result, nil
}

type SatuanInput struct {
	NamaSatuan    string  `json:"nama_satuan"`
	JenisInput    string  `json:"jenis_input"`
	NilaiKonversi float64 `json:"nilai_konversi"`
	SatuanUtama   int     `json:"satuan_utama"`
	KodeBarcode   string  `json:"kode_barcode"`
	HargaBeli     float64 `json:"harga_beli"`
	HargaJual     float64 `json:"harga_jual"`
	HargaTitip    float64 `json:"harga_titip"`
}

func ProdukStoreSatuan(produkID int, input SatuanInput) (int64, error) {
	var tipe string
	db.Get().QueryRow(`SELECT tipe_kepemilikan FROM produk WHERE id = ?`, produkID).Scan(&tipe)

	tx, err := db.Get().Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if input.SatuanUtama == 1 {
		_, err = tx.Exec(`UPDATE satuan_produk SET satuan_utama = 0, updated_at = datetime('now') WHERE produk_id = ?`, produkID)
		if err != nil {
			return 0, err
		}
	}

	res, err := tx.Exec(`
		INSERT INTO satuan_produk (produk_id, nama_satuan, jenis_input, nilai_konversi, satuan_utama)
		VALUES (?, ?, ?, ?, ?)
	`, produkID, strings.ToUpper(input.NamaSatuan), input.JenisInput, input.NilaiKonversi, input.SatuanUtama)
	if err != nil {
		return 0, err
	}
	satuanID, _ := res.LastInsertId()

	if bc := strings.TrimSpace(input.KodeBarcode); bc != "" {
		_, err = tx.Exec(`INSERT INTO barcode_produk (satuan_produk_id, kode_barcode) VALUES (?, ?)`, satuanID, bc)
		if err != nil {
			return 0, err
		}
	}

	hargaBeli := float64(0)
	if tipe == "MILIK_SENDIRI" {
		hargaBeli = input.HargaBeli
	}
	_, err = tx.Exec(`INSERT INTO harga_produk (satuan_produk_id, harga_beli, harga_jual, berlaku_sejak) VALUES (?, ?, ?, datetime('now'))`,
		satuanID, hargaBeli, input.HargaJual)
	if err != nil {
		return 0, err
	}

	if tipe == "TITIPAN" {
		komisi := input.HargaJual - input.HargaTitip
		if komisi < 0 {
			komisi = 0
		}
		_, err = tx.Exec(`INSERT INTO konsinyasi_produk (satuan_produk_id, harga_hak_pemasok, komisi_toko) VALUES (?, ?, ?)`,
			satuanID, input.HargaTitip, komisi)
		if err != nil {
			return 0, err
		}
	}

	return satuanID, tx.Commit()
}

func ProdukUpdateSatuan(produkID, satuanID int, input SatuanInput) error {
	var tipe string
	db.Get().QueryRow(`SELECT tipe_kepemilikan FROM produk WHERE id = ?`, produkID).Scan(&tipe)

	tx, err := db.Get().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if input.SatuanUtama == 1 {
		_, err = tx.Exec(`UPDATE satuan_produk SET satuan_utama = 0, updated_at = datetime('now') WHERE produk_id = ?`, produkID)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`
		UPDATE satuan_produk SET nama_satuan = ?, jenis_input = ?, nilai_konversi = ?, satuan_utama = ?, updated_at = datetime('now')
		WHERE id = ?
	`, strings.ToUpper(input.NamaSatuan), input.JenisInput, input.NilaiKonversi, input.SatuanUtama, satuanID)
	if err != nil {
		return err
	}

	barcode := strings.TrimSpace(input.KodeBarcode)
	var barcodeRowID sql.NullInt64
	tx.QueryRow(`SELECT id FROM barcode_produk WHERE satuan_produk_id = ?`, satuanID).Scan(&barcodeRowID)
	if barcode != "" {
		if barcodeRowID.Valid {
			_, err = tx.Exec(`UPDATE barcode_produk SET kode_barcode = ?, updated_at = datetime('now') WHERE id = ?`, barcode, barcodeRowID.Int64)
		} else {
			_, err = tx.Exec(`INSERT INTO barcode_produk (satuan_produk_id, kode_barcode) VALUES (?, ?)`, satuanID, barcode)
		}
		if err != nil {
			return err
		}
	} else if barcodeRowID.Valid {
		_, err = tx.Exec(`DELETE FROM barcode_produk WHERE id = ?`, barcodeRowID.Int64)
		if err != nil {
			return err
		}
	}

	hargaBeli := float64(0)
	if tipe == "MILIK_SENDIRI" {
		hargaBeli = input.HargaBeli
	}
	_, err = tx.Exec(`INSERT INTO harga_produk (satuan_produk_id, harga_beli, harga_jual, berlaku_sejak) VALUES (?, ?, ?, datetime('now'))`,
		satuanID, hargaBeli, input.HargaJual)
	if err != nil {
		return err
	}

	if tipe == "TITIPAN" {
		komisi := input.HargaJual - input.HargaTitip
		if komisi < 0 {
			komisi = 0
		}
		var kpID sql.NullInt64
		tx.QueryRow(`SELECT id FROM konsinyasi_produk WHERE satuan_produk_id = ?`, satuanID).Scan(&kpID)
		if kpID.Valid {
			_, err = tx.Exec(`UPDATE konsinyasi_produk SET harga_hak_pemasok = ?, komisi_toko = ?, updated_at = datetime('now') WHERE satuan_produk_id = ?`,
				input.HargaTitip, komisi, satuanID)
		} else {
			_, err = tx.Exec(`INSERT INTO konsinyasi_produk (satuan_produk_id, harga_hak_pemasok, komisi_toko) VALUES (?, ?, ?)`,
				satuanID, input.HargaTitip, komisi)
		}
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func ProdukDestroySatuan(produkID, satuanID int) error {
	var count int
	db.Get().QueryRow(`SELECT COUNT(*) FROM satuan_produk WHERE produk_id = ?`, produkID).Scan(&count)
	if count <= 1 {
		return fmt.Errorf("Minimal harus ada 1 satuan")
	}
	_, err := db.Get().Exec(`DELETE FROM satuan_produk WHERE id = ? AND produk_id = ?`, satuanID, produkID)
	return err
}
