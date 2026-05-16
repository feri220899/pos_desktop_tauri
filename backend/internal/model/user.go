package model

import (
	"database/sql"
	"strings"

	"pos-desktop-tauri/backend/internal/db"
)

const userSelectWithRole = `
	SELECT u.id, u.username, u.active, u.role_id, u.created_at, u.updated_at,
	       r.name AS role_name,
	       COALESCE(
	           (SELECT '[' || GROUP_CONCAT('"' || p.key || '"') || ']'
	            FROM role_permissions rp
	            JOIN permissions p ON p.id = rp.permission_id
	            WHERE rp.role_id = r.id),
	           '[]'
	       ) AS permissions
	FROM users u
	JOIN roles r ON r.id = u.role_id
`

var userSortColumns = map[string]string{
	"username":  "u.username",
	"role_name": "r.name",
	"active":    "u.active",
	"id":        "u.id",
}

type User struct {
	ID          int      `json:"id"`
	Username    string   `json:"username"`
	Active      int      `json:"active"`
	RoleID      int      `json:"role_id"`
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

func scanUser(row interface{ Scan(...any) error }) (*User, error) {
	var u User
	var permsJSON string
	err := row.Scan(&u.ID, &u.Username, &u.Active, &u.RoleID, &u.CreatedAt, &u.UpdatedAt, &u.RoleName, &permsJSON)
	if err != nil {
		return nil, err
	}
	u.Permissions = parseJSONStringArray(permsJSON)
	return &u, nil
}

func UserAll() ([]User, error) {
	rows, err := db.Get().Query(userSelectWithRole + " ORDER BY u.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUsers(rows)
}

func UserPaginate(page, pageSize int, sortBy, sortOrder, search string) ([]User, int, error) {
	col := userSortColumns[sortBy]
	if col == "" {
		col = "u.id"
	}
	dir := "ASC"
	if strings.ToLower(sortOrder) == "desc" {
		dir = "DESC"
	}
	offset := (page - 1) * pageSize
	like   := "%" + search + "%"

	var total int
	err := db.Get().QueryRow(
		`SELECT COUNT(*) FROM users u JOIN roles r ON r.id = u.role_id WHERE u.username LIKE ?`,
		like,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.Get().Query(
		userSelectWithRole+` WHERE u.username LIKE ? ORDER BY `+col+` `+dir+` LIMIT ? OFFSET ?`,
		like, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	return users, total, err
}

func UserFind(id int) (*User, error) {
	row := db.Get().QueryRow(userSelectWithRole+" WHERE u.id = ?", id)
	u, err := scanUser(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func UserFindByUsername(username string) (*User, string, error) {
	var u User
	var permsJSON, hash string
	err := db.Get().QueryRow(`
		SELECT u.id, u.username, u.active, u.role_id, u.created_at, u.updated_at,
		       r.name AS role_name,
		       COALESCE(
		           (SELECT '[' || GROUP_CONCAT('"' || p.key || '"') || ']'
		            FROM role_permissions rp
		            JOIN permissions p ON p.id = rp.permission_id
		            WHERE rp.role_id = r.id),
		           '[]'
		       ) AS permissions,
		       u.password
		FROM users u JOIN roles r ON r.id = u.role_id
		WHERE u.username = ?
	`, username).Scan(&u.ID, &u.Username, &u.Active, &u.RoleID, &u.CreatedAt, &u.UpdatedAt, &u.RoleName, &permsJSON, &hash)
	if err == sql.ErrNoRows {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}
	u.Permissions = parseJSONStringArray(permsJSON)
	return &u, hash, nil
}

func UserCreate(username, password string, roleID, active int) (int64, error) {
	res, err := db.Get().Exec(
		`INSERT INTO users (username, password, role_id, active) VALUES (?, ?, ?, ?)`,
		username, password, roleID, active,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func UserUpdate(id int, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}
	var cols []string
	var vals []any
	for k, v := range fields {
		cols = append(cols, k+" = ?")
		vals = append(vals, v)
	}
	cols = append(cols, "updated_at = datetime('now')")
	vals = append(vals, id)
	_, err := db.Get().Exec(
		`UPDATE users SET `+strings.Join(cols, ", ")+` WHERE id = ?`,
		vals...,
	)
	return err
}

func UserDestroy(id int) error {
	_, err := db.Get().Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

func scanUsers(rows *sql.Rows) ([]User, error) {
	var users []User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	if users == nil {
		users = []User{}
	}
	return users, nil
}
