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
type UserPostgresRepository interface {
	Create(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByID(id string) (*model.User, error)
}

type userPostgresRepo struct {
	pool *pgxpool.Pool
}

// =======================
// CONSTRUCTOR
// =======================
func NewUserPostgresRepository() UserPostgresRepository {
	return &userPostgresRepo{
		pool: database.Pg,
	}
}

// =======================
// IMPLEMENTATION
// =======================

func (r *userPostgresRepo) Create(u *model.User) error {
	_, err := r.pool.Exec(
		context.Background(),
		`INSERT INTO users 
		(id, username, email, password_hash, full_name, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		u.ID, u.Username, u.Email, u.PasswordHash, u.FullName, u.RoleID, u.IsActive,
	)
	return err
}

func (r *userPostgresRepo) GetByEmail(email string) (*model.User, error) {
	var u model.User

	err := r.pool.QueryRow(
		context.Background(),
		`SELECT id, username, email, password_hash, full_name, role_id, is_active,
		        created_at, updated_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &u, nil
}

func (r *userPostgresRepo) GetByUsername(username string) (*model.User, error) {
	var u model.User

	err := r.pool.QueryRow(
		context.Background(),
		`SELECT id, username, email, password_hash, full_name, role_id, is_active,
		        created_at, updated_at
		 FROM users WHERE username = $1`,
		username,
	).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &u, nil
}

func (r *userPostgresRepo) GetByID(id string) (*model.User, error) {
	var u model.User

	err := r.pool.QueryRow(
		context.Background(),
		`SELECT id, username, email, password_hash, full_name, role_id, is_active,
		        created_at, updated_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &u, nil
}
