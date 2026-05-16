package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var instance *sql.DB

func Init(dbPath string) error {
	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on", dbPath)
	d, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	d.SetMaxOpenConns(1)
	instance = d
	log.Printf("[db] connected: %s", dbPath)
	return nil
}

func Get() *sql.DB {
	if instance == nil {
		panic("db not initialized — call db.Init first")
	}
	return instance
}

func Close() {
	if instance != nil {
		instance.Close()
	}
}
