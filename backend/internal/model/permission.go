package model

import (
	"database/sql"

	"pos-desktop-tauri/backend/internal/db"
)

type Permission struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Label     string `json:"label"`
	IconName  string `json:"icon_name"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func PermissionAll() ([]Permission, error) {
	rows, err := db.Get().Query(`SELECT id, key, label, icon_name, sort_order, created_at, updated_at FROM permissions ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Key, &p.Label, &p.IconName, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	if result == nil {
		result = []Permission{}
	}
	return result, nil
}

func PermissionFind(id int) (*Permission, error) {
	var p Permission
	err := db.Get().QueryRow(`SELECT id, key, label, icon_name, sort_order, created_at, updated_at FROM permissions WHERE id = ?`, id).
		Scan(&p.ID, &p.Key, &p.Label, &p.IconName, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}
