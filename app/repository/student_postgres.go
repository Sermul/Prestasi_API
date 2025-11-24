package repository

import (
	"context"
	"prestasi_api/database"
)

type StudentPostgresRepository interface {
	GetStudentIDsByAdvisor(advisorID string) ([]string, error)
}

type studentPostgresRepo struct{}

func NewStudentPostgresRepository() StudentPostgresRepository {
	return &studentPostgresRepo{}
}

func (r *studentPostgresRepo) GetStudentIDsByAdvisor(advisorID string) ([]string, error) {
	rows, err := database.Pg.Query(
		context.Background(),
		`SELECT id FROM students WHERE advisor_id = $1`,
		advisorID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
