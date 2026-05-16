package migration

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

var All = []Migration{
	{
		Name: "001_create_produk",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS produk (
					id         INTEGER PRIMARY KEY AUTOINCREMENT,
					nama       TEXT    NOT NULL,
					harga      INTEGER NOT NULL DEFAULT 0,
					stok       INTEGER NOT NULL DEFAULT 0,
					created_at TEXT    DEFAULT (datetime('now')),
					updated_at TEXT    DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
	{
		Name: "002_create_roles",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS roles (
					id          INTEGER PRIMARY KEY AUTOINCREMENT,
					name        TEXT    NOT NULL UNIQUE,
					permissions TEXT    NOT NULL DEFAULT '[]',
					created_at  TEXT    DEFAULT (datetime('now')),
					updated_at  TEXT    DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
	{
		Name: "003_create_users",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS users (
					id         INTEGER PRIMARY KEY AUTOINCREMENT,
					username   TEXT    NOT NULL UNIQUE,
					password   TEXT    NOT NULL,
					role_id    INTEGER NOT NULL REFERENCES roles(id),
					active     INTEGER NOT NULL DEFAULT 1,
					created_at TEXT    DEFAULT (datetime('now')),
					updated_at TEXT    DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
	{
		Name: "004_seed_defaults",
		Up: func(db *sql.DB) error {
			allPerms   := `["dashboard","kasir","produk","transaksi","laporan"]`
			kasirPerms := `["dashboard","kasir","transaksi"]`
			_, err := db.Exec(`INSERT OR IGNORE INTO roles (name, permissions) VALUES (?, ?)`, "admin", allPerms)
			if err != nil {
				return err
			}
			_, err = db.Exec(`INSERT OR IGNORE INTO roles (name, permissions) VALUES (?, ?)`, "kasir", kasirPerms)
			if err != nil {
				return err
			}

			var adminID int
			err = db.QueryRow(`SELECT id FROM roles WHERE name = 'admin'`).Scan(&adminID)
			if err != nil {
				return err
			}

			hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			_, err = db.Exec(`INSERT OR IGNORE INTO users (username, password, role_id) VALUES (?, ?, ?)`,
				"admin", string(hash), adminID)
			return err
		},
	},
	{
		Name: "005_add_pengaturan_permission",
		Up: func(db *sql.DB) error {
			var id int
			var perms string
			err := db.QueryRow(`SELECT id, permissions FROM roles WHERE name = 'admin'`).Scan(&id, &perms)
			if err != nil {
				return nil
			}
			if len(perms) > 0 && !containsStr(perms, `"pengaturan"`) {
				newPerms := perms[:len(perms)-1] + `,"pengaturan"]`
				_, err = db.Exec(`UPDATE roles SET permissions = ? WHERE id = ?`, newPerms, id)
			}
			return err
		},
	},
	{
		Name: "006_create_permissions",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS permissions (
					id         INTEGER PRIMARY KEY AUTOINCREMENT,
					key        TEXT NOT NULL UNIQUE,
					label      TEXT NOT NULL,
					icon_name  TEXT NOT NULL DEFAULT '',
					sort_order INTEGER NOT NULL DEFAULT 0,
					created_at TEXT DEFAULT (datetime('now')),
					updated_at TEXT DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
	{
		Name: "007_seed_permissions",
		Up: func(db *sql.DB) error {
			rows := []struct{ key, label, icon string; order int }{
				{"dashboard",  "Dashboard",  "LayoutDashboard", 1},
				{"kasir",      "Kasir",      "ShoppingCart",    2},
				{"produk",     "Produk",     "Package",         3},
				{"transaksi",  "Transaksi",  "ArrowLeftRight",  4},
				{"laporan",    "Laporan",    "BarChart2",       5},
				{"pengaturan", "Pengaturan", "Settings",        6},
			}
			for _, r := range rows {
				_, err := db.Exec(
					`INSERT OR IGNORE INTO permissions (key, label, icon_name, sort_order) VALUES (?, ?, ?, ?)`,
					r.key, r.label, r.icon, r.order,
				)
				if err != nil {
					return err
				}
			}
			return nil
		},
	},
	{
		Name: "008_role_permissions_pivot",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS role_permissions (
					role_id       INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
					permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
					PRIMARY KEY (role_id, permission_id)
				)
			`)
			if err != nil {
				return err
			}

			rows, err := db.Query(`SELECT id, permissions FROM roles`)
			if err != nil {
				return err
			}
			defer rows.Close()

			type roleRow struct {
				id   int
				keys []string
			}
			var roles []roleRow
			for rows.Next() {
				var id int
				var permsJSON string
				rows.Scan(&id, &permsJSON)
				keys := parseJSONStringArray(permsJSON)
				roles = append(roles, roleRow{id, keys})
			}

			for _, role := range roles {
				for _, key := range role.keys {
					db.Exec(`
						INSERT OR IGNORE INTO role_permissions (role_id, permission_id)
						SELECT ?, id FROM permissions WHERE key = ?
					`, role.id, key)
				}
			}

			db.Exec(`ALTER TABLE roles DROP COLUMN permissions`)
			return nil
		},
	},
	{
		Name: "009_create_toko",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS toko (
					id         INTEGER PRIMARY KEY CHECK (id = 1),
					nama       TEXT    NOT NULL DEFAULT '',
					telepon    TEXT    NOT NULL DEFAULT '',
					email      TEXT    NOT NULL DEFAULT '',
					alamat     TEXT    NOT NULL DEFAULT '',
					created_at TEXT    DEFAULT (datetime('now')),
					updated_at TEXT    DEFAULT (datetime('now'))
				)
			`)
			if err != nil {
				return err
			}
			_, err = db.Exec(`INSERT OR IGNORE INTO toko (id) VALUES (1)`)
			return err
		},
	},
	{
		Name: "010_create_kategori",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS kategori (
					id            INTEGER PRIMARY KEY AUTOINCREMENT,
					nama_kategori TEXT    NOT NULL,
					created_at    TEXT    DEFAULT (datetime('now')),
					updated_at    TEXT    DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
	{
		Name: "011_create_pemasok",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS pemasok (
					id           INTEGER PRIMARY KEY AUTOINCREMENT,
					nama_pemasok TEXT    NOT NULL,
					telepon      TEXT,
					alamat       TEXT,
					created_at   TEXT    DEFAULT (datetime('now')),
					updated_at   TEXT    DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
	{
		Name: "012_rebuild_produk",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`DROP TABLE IF EXISTS produk`)
			if err != nil {
				return err
			}
			_, err = db.Exec(`
				CREATE TABLE produk (
					id               INTEGER PRIMARY KEY AUTOINCREMENT,
					kode_produk      TEXT    NOT NULL UNIQUE,
					nama_produk      TEXT    NOT NULL,
					kategori_id      INTEGER REFERENCES kategori(id) ON DELETE SET NULL,
					tipe_kepemilikan TEXT    NOT NULL DEFAULT 'MILIK_SENDIRI',
					pemasok_id       INTEGER REFERENCES pemasok(id) ON DELETE SET NULL,
					deskripsi        TEXT,
					status_aktif     INTEGER NOT NULL DEFAULT 1,
					deleted_at       TEXT,
					created_at       TEXT    DEFAULT (datetime('now')),
					updated_at       TEXT    DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
	{
		Name: "013_create_satuan_produk",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS satuan_produk (
					id             INTEGER PRIMARY KEY AUTOINCREMENT,
					produk_id      INTEGER NOT NULL REFERENCES produk(id) ON DELETE CASCADE,
					nama_satuan    TEXT    NOT NULL,
					jenis_input    TEXT    NOT NULL DEFAULT 'PCS',
					nilai_konversi REAL    NOT NULL DEFAULT 1,
					satuan_utama   INTEGER NOT NULL DEFAULT 0,
					created_at     TEXT    DEFAULT (datetime('now')),
					updated_at     TEXT    DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
	{
		Name: "014_create_barcode_produk",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS barcode_produk (
					id               INTEGER PRIMARY KEY AUTOINCREMENT,
					satuan_produk_id INTEGER NOT NULL REFERENCES satuan_produk(id) ON DELETE CASCADE,
					kode_barcode     TEXT    NOT NULL UNIQUE,
					created_at       TEXT    DEFAULT (datetime('now')),
					updated_at       TEXT    DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
	{
		Name: "015_create_harga_produk",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS harga_produk (
					id               INTEGER PRIMARY KEY AUTOINCREMENT,
					satuan_produk_id INTEGER NOT NULL REFERENCES satuan_produk(id) ON DELETE CASCADE,
					harga_beli       REAL    NOT NULL DEFAULT 0,
					harga_jual       REAL    NOT NULL DEFAULT 0,
					berlaku_sejak    TEXT    NOT NULL DEFAULT (datetime('now')),
					created_at       TEXT    DEFAULT (datetime('now')),
					updated_at       TEXT    DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
	{
		Name: "016_create_stok_produk",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS stok_produk (
					id          INTEGER PRIMARY KEY AUTOINCREMENT,
					produk_id   INTEGER NOT NULL REFERENCES produk(id) ON DELETE CASCADE,
					jumlah_stok REAL    NOT NULL DEFAULT 0,
					created_at  TEXT    DEFAULT (datetime('now')),
					updated_at  TEXT    DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
	{
		Name: "017_create_konsinyasi_produk",
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS konsinyasi_produk (
					id                INTEGER PRIMARY KEY AUTOINCREMENT,
					satuan_produk_id  INTEGER NOT NULL REFERENCES satuan_produk(id) ON DELETE CASCADE,
					harga_hak_pemasok REAL    NOT NULL DEFAULT 0,
					komisi_toko       REAL    NOT NULL DEFAULT 0,
					created_at        TEXT    DEFAULT (datetime('now')),
					updated_at        TEXT    DEFAULT (datetime('now'))
				)
			`)
			return err
		},
	},
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func parseJSONStringArray(s string) []string {
	var result []string
	if len(s) < 2 {
		return result
	}
	s = s[1 : len(s)-1]
	if s == "" {
		return result
	}
	start := -1
	for i, c := range s {
		if c == '"' && start == -1 {
			start = i + 1
		} else if c == '"' && start != -1 {
			result = append(result, s[start:i])
			start = -1
		}
	}
	return result
}
