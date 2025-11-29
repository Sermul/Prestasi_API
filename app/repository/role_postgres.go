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
type RolePostgresRepository interface {
	GetByID(id string) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	GetAll() ([]model.Role, error)
}

type rolePostgresRepo struct {
	pool *pgxpool.Pool
}

// =======================
// CONSTRUCTOR
// =======================
func NewRolePostgresRepository() RolePostgresRepository {
	return &rolePostgresRepo{
		pool: database.Pg,
	}
}

// =======================
// IMPLEMENTATION
// =======================
func (r *rolePostgresRepo) GetByID(id string) (*model.Role, error) {
	var role model.Role

	err := r.pool.QueryRow(context.Background(),
		`SELECT id, name, description, created_at 
		 FROM roles WHERE id = $1`,
		id,
	).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)

	if err != nil {
		return nil, errors.New("role not found")
	}

	return &role, nil
}

func (r *rolePostgresRepo) GetByName(name string) (*model.Role, error) {
	var role model.Role

	err := r.pool.QueryRow(context.Background(),
		`SELECT id, name, description, created_at 
		 FROM roles WHERE name = $1`,
		name,
	).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)

	if err != nil {
		return nil, errors.New("role not found")
	}

	return &role, nil
}

func (r *rolePostgresRepo) GetAll() ([]model.Role, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT id, name, description, created_at FROM roles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Role
	for rows.Next() {
		var role model.Role
		rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
		list = append(list, role)
	}

	return list, nil
}
