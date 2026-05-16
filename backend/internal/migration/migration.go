package migration

import (
	"database/sql"
	"log"
)

type Migration struct {
	Name string
	Up   func(db *sql.DB) error
}

func Run(db *sql.DB, migrations []Migration) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT    NOT NULL UNIQUE,
			created_at TEXT    DEFAULT (datetime('now'))
		)
	`)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		var exists int
		db.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = ?`, m.Name).Scan(&exists)
		if exists > 0 {
			continue
		}

		log.Printf("[migration] running: %s", m.Name)
		if err := m.Up(db); err != nil {
			return err
		}

		_, err = db.Exec(`INSERT INTO migrations (name) VALUES (?)`, m.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
