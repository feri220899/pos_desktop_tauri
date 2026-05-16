package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"pos-desktop-tauri/backend/internal/helper"
	"pos-desktop-tauri/backend/internal/service"
)

func DiscoveryAdvertise(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
		Port int    `json:"port"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if body.Name == "" {
		hostname, _ := os.Hostname()
		body.Name = "POS Master - " + hostname
	}
	if body.Port == 0 {
		body.Port = 3001
	}

	if err := service.Advertise(body.Name, body.Port); err != nil {
		helper.ServerError(w, err.Error())
		return
	}
	helper.Success(w, map[string]any{"name": body.Name, "port": body.Port})
}

func DiscoveryScan(w http.ResponseWriter, r *http.Request) {
	masters, err := service.Scan(3 * time.Second)
	if err != nil {
		helper.ServerError(w, err.Error())
		return
	}
	helper.Success(w, masters)
}

func DiscoveryStop(w http.ResponseWriter, r *http.Request) {
	service.StopAdvertise()
	service.StopScan()
	helper.Success(w, nil)
}
