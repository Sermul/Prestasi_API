package repository

import (
	"context"
	"prestasi_api/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

// =======================
// INTERFACE
// =======================
type RolePermissionPostgresRepository interface {
	Assign(roleID string, permissionID string) error
	GetPermissionsByRoleID(roleID string) ([]string, error)
}

type rolePermissionPostgresRepo struct {
	pool *pgxpool.Pool
}

// =======================
// CONSTRUCTOR
// =======================
func NewRolePermissionPostgresRepository() RolePermissionPostgresRepository {
	return &rolePermissionPostgresRepo{
		pool: database.Pg,
	}
}

// =======================
// IMPLEMENTATION
// =======================
func (r *rolePermissionPostgresRepo) Assign(roleID string, permissionID string) error {
	_, err := r.pool.Exec(
		context.Background(),
		`INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)`,
		roleID, permissionID,
	)
	return err
}

func (r *rolePermissionPostgresRepo) GetPermissionsByRoleID(roleID string) ([]string, error) {
	rows, err := r.pool.Query(
		context.Background(),
		`SELECT permission_id FROM role_permissions WHERE role_id = $1`,
		roleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var permID string
		rows.Scan(&permID)
		list = append(list, permID)
	}

	return list, nil
}
