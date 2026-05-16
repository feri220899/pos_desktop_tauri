package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	appDB        "pos-desktop-tauri/backend/internal/db"
	"pos-desktop-tauri/backend/internal/migration"
	"pos-desktop-tauri/backend/internal/router"
	"pos-desktop-tauri/backend/internal/service"
)

func main() {
	var port, dbPath, appMode string
	flag.StringVar(&port,    "port", "", "HTTP port (default 3001)")
	flag.StringVar(&dbPath,  "db",   "", "SQLite database path")
	flag.StringVar(&appMode, "mode", "", "app_mode: master or client")
	flag.Parse()

	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "3001"
	}
	if dbPath == "" {
		dbPath = os.Getenv("DB_PATH")
	}
	if dbPath == "" {
		home, _ := os.UserHomeDir()
		dbPath = filepath.Join(home, ".local", "share", "pos", "data.db")
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("cannot create db dir: %v", err)
	}

	if err := appDB.Init(dbPath); err != nil {
		log.Fatalf("db init: %v", err)
	}
	defer appDB.Close()

	if err := migration.Run(appDB.Get(), migration.All); err != nil {
		log.Fatalf("migration: %v", err)
	}

	// Auto-advertise jika mode master
	if appMode == "master" {
		hostname, _ := os.Hostname()
		portNum, _  := strconv.Atoi(port)
		go service.Advertise("POS Master - "+hostname, portNum)
	}

	log.Printf("[server] listening on :%s (mode: %s)", port, appMode)
	if err := http.ListenAndServe(":"+port, router.New()); err != nil {
		log.Fatalf("server: %v", err)
	}
}
