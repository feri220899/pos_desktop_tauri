package model

import (
	"database/sql"

	"pos-desktop-tauri/backend/internal/db"
)

const rolePermsSubquery = `
	COALESCE(
		(SELECT '[' || GROUP_CONCAT('"' || p.key || '"') || ']'
		 FROM role_permissions rp
		 JOIN permissions p ON p.id = rp.permission_id
		 WHERE rp.role_id = roles.id),
		'[]'
	) AS permissions
`

type Role struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

func scanRole(row interface{ Scan(...any) error }) (*Role, error) {
	var r Role
	var permsJSON string
	err := row.Scan(&r.ID, &r.Name, &r.CreatedAt, &r.UpdatedAt, &permsJSON)
	if err != nil {
		return nil, err
	}
	r.Permissions = parseJSONStringArray(permsJSON)
	return &r, nil
}

func RoleAll() ([]Role, error) {
	rows, err := db.Get().Query(`SELECT id, name, created_at, updated_at, ` + rolePermsSubquery + ` FROM roles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Role
	for rows.Next() {
		r, err := scanRole(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *r)
	}
	if result == nil {
		result = []Role{}
	}
	return result, nil
}

func RoleFind(id int) (*Role, error) {
	row := db.Get().QueryRow(`SELECT id, name, created_at, updated_at, `+rolePermsSubquery+` FROM roles WHERE id = ?`, id)
	r, err := scanRole(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return r, err
}

func RoleCreate(name string, permissions []string) (int64, error) {
	res, err := db.Get().Exec(`INSERT INTO roles (name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, setRolePermissions(int(id), permissions)
}

func RoleUpdate(id int, permissions []string) error {
	return setRolePermissions(id, permissions)
}

func RoleDestroy(id int) error {
	_, err := db.Get().Exec(`DELETE FROM roles WHERE id = ?`, id)
	return err
}

func setRolePermissions(roleID int, keys []string) error {
	d := db.Get()
	_, err := d.Exec(`DELETE FROM role_permissions WHERE role_id = ?`, roleID)
	if err != nil {
		return err
	}
	for _, key := range keys {
		d.Exec(`
			INSERT OR IGNORE INTO role_permissions (role_id, permission_id)
			SELECT ?, id FROM permissions WHERE key = ?
		`, roleID, key)
	}
	return nil
}
