package repository

import (
	"prestasi_api/database"
	"prestasi_api/app/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Interface sesuai kebutuhan modul
type AchievementPostgresRepository interface {
	CreateReferencePostgres(ref *model.AchievementReference) error
	UpdateReferenceStatusPostgres(refID string, status string) error
	GetMongoID(refID string) (string, error)
}

// Implementasi repo
type achievementPostgresRepo struct {
	pool *pgxpool.Pool
}

// Constructor
func NewAchievementPostgresRepository() AchievementPostgresRepository {
	return &achievementPostgresRepo{
		pool: database.Pg,
	}
}


// CREATE REFERENCE (FR-003)

func (r *achievementPostgresRepo) CreateReferencePostgres(ref *model.AchievementReference) error {
	_, err := r.pool.Exec(
		context.Background(),
		`INSERT INTO achievement_reference 
		 (id, student_id, mongo_id, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		ref.ID, ref.StudentID, ref.MongoID, ref.Status, ref.CreatedAt, ref.UpdatedAt,
	)
	return err
}


// UPDATE STATUS (FR-004 & FR-005)

func (r *achievementPostgresRepo) UpdateReferenceStatusPostgres(refID string, status string) error {
	_, err := r.pool.Exec(
		context.Background(),
		`UPDATE achievement_reference
		 SET status = $1, updated_at = NOW()
		 WHERE id = $2`,
		status, refID,
	)
	return err
}


// GET mongo_id (FR-005)
func (r *achievementPostgresRepo) GetMongoID(refID string) (string, error) {
	var mongoID string

	err := r.pool.QueryRow(
		context.Background(),
		`SELECT mongo_id 
		 FROM achievement_reference
		 WHERE id = $1`,
		refID,
	).Scan(&mongoID)

	if err != nil {
		return "", errors.New("reference not found")
	}

	return mongoID, nil
}
