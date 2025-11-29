package repository

import (
	"context"
	"errors"
	"prestasi_api/app/model"
	"prestasi_api/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

// =======================
// INTERFACE
// =======================
type PermissionPostgresRepository interface {
	GetByID(id string) (*model.Permission, error)
	GetByName(name string) (*model.Permission, error)
	GetByRoleID(roleID string) ([]model.Permission, error)
}

type permissionPostgresRepo struct {
	pool *pgxpool.Pool
}

// =======================
// CONSTRUCTOR
// =======================
func NewPermissionPostgresRepository() PermissionPostgresRepository {
	return &permissionPostgresRepo{
		pool: database.Pg,
	}
}

// =======================
// IMPLEMENTATION
// =======================
func (r *permissionPostgresRepo) GetByID(id string) (*model.Permission, error) {
	var p model.Permission

	err := r.pool.QueryRow(context.Background(),
		`SELECT id, name, resource, action, description 
		 FROM permissions WHERE id = $1`,
		id,
	).Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description)

	if err != nil {
		return nil, errors.New("permission not found")
	}

	return &p, nil
}

func (r *permissionPostgresRepo) GetByName(name string) (*model.Permission, error) {
	var p model.Permission

	err := r.pool.QueryRow(context.Background(),
		`SELECT id, name, resource, action, description
		 FROM permissions WHERE name = $1`,
		name,
	).Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description)

	if err != nil {
		return nil, errors.New("permission not found")
	}

	return &p, nil
}

func (r *permissionPostgresRepo) GetByRoleID(roleID string) ([]model.Permission, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT p.id, p.name, p.resource, p.action, p.description
		 FROM role_permissions rp
		 JOIN permissions p ON p.id = rp.permission_id
		 WHERE rp.role_id = $1`,
		roleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Permission

	for rows.Next() {
		var p model.Permission
		err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}
