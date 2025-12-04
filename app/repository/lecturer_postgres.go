package repository

import (
	"context"
	"errors"
	"prestasi_api/app/model"
	"prestasi_api/database"
)

type LecturerPostgresRepository interface {
	Create(l *model.Lecturer) error
	GetByUserID(userID string) (*model.Lecturer, error)
	GetByLecturerID(lecturerID string) (*model.Lecturer, error)
}

type lecturerPostgresRepo struct{}

func NewLecturerPostgresRepository() LecturerPostgresRepository {
	return &lecturerPostgresRepo{}
}

func (r *lecturerPostgresRepo) Create(l *model.Lecturer) error {
	_, err := database.Pg.Exec(
		context.Background(),
		`INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
		 VALUES ($1, $2, $3, $4, NOW())`,
		l.ID, l.UserID, l.LecturerID, l.Department,
	)
	return err
}

func (r *lecturerPostgresRepo) GetByUserID(userID string) (*model.Lecturer, error) {
	var l model.Lecturer

	err := database.Pg.QueryRow(
		context.Background(),
		`SELECT id, user_id, lecturer_id, department, created_at
		 FROM lecturers
		 WHERE user_id = $1 LIMIT 1`,
		userID,
	).Scan(
		&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt,
	)

	if err != nil {
		return nil, errors.New("lecturer not found")
	}

	return &l, nil
}

func (r *lecturerPostgresRepo) GetByLecturerID(lecturerID string) (*model.Lecturer, error) {
	var l model.Lecturer

	err := database.Pg.QueryRow(
		context.Background(),
		`SELECT id, user_id, lecturer_id, department, created_at
		 FROM lecturers
		 WHERE lecturer_id = $1 LIMIT 1`,
		lecturerID,
	).Scan(
		&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt,
	)

	if err != nil {
		return nil, errors.New("lecturer not found")
	}

	return &l, nil
}
