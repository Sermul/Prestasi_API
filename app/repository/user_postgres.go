package repository

import (
	"context"
	"errors"
	"prestasi_api/app/model"
	"prestasi_api/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

// INTERFACE
type UserPostgresRepository interface {
	Create(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByID(id string) (*model.User, error)
	  // Tambahan untuk UserService
    GetAll() ([]model.User, error)
    Update(id string, user *model.User) error
    Delete(id string) error
    UpdateRole(id string, roleID string) error
}

type userPostgresRepo struct {
	pool *pgxpool.Pool
}


// CONSTRUCTOR
func NewUserPostgresRepository() UserPostgresRepository {
	return &userPostgresRepo{
		pool: database.Pg,
	}
}


// IMPLEMENTATION

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
func (r *userPostgresRepo) GetAll() ([]model.User, error) {
    rows, err := r.pool.Query(context.Background(),
        `SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at 
         FROM users`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []model.User
    for rows.Next() {
        var u model.User
        rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
        users = append(users, u)
    }
    return users, nil
}
func (r *userPostgresRepo) Update(id string, user *model.User) error {
    _, err := r.pool.Exec(context.Background(),
        `UPDATE users SET username=$1, email=$2, full_name=$3, role_id=$4, updated_at=$5
         WHERE id=$6`,
        user.Username, user.Email, user.FullName, user.RoleID, user.UpdatedAt, id,
    )
    return err
}
func (r *userPostgresRepo) Delete(id string) error {
    _, err := r.pool.Exec(context.Background(),
        `DELETE FROM users WHERE id=$1`, id)
    return err
}
func (r *userPostgresRepo) UpdateRole(id string, roleID string) error {
    _, err := r.pool.Exec(context.Background(),
        `UPDATE users SET role_id=$1, updated_at=NOW() WHERE id=$2`,
        roleID, id,
    )
    return err
}
