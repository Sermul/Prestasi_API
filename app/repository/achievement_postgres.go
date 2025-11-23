package repository

import (
	"context"
	"app/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AchievementPostgresRepository struct {
	DB *pgxpool.Pool
}

func NewAchievementPostgresRepository(db *pgxpool.Pool) *AchievementPostgresRepository {
	return &AchievementPostgresRepository{DB: db}
}

// FR-003: Create Reference in PostgreSQL
func (r *AchievementPostgresRepository) CreateReferencePostgres(ref *model.AchievementReference) error {
	query := `
	INSERT INTO achievement_references 
	(id, student_id, mongo_achievement_id, status, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW())
	`

	_, err := r.DB.Exec(context.Background(), query,
		ref.ID,
		ref.StudentID,
		ref.MongoID,
		ref.Status,
	)
	return err
}

// FR-004 & FR-005: Update status (submit / delete / verify)
func (r *AchievementPostgresRepository) UpdateReferenceStatusPostgres(refID string, status string) error {
	query := `
	UPDATE achievement_references
	SET status = $1, updated_at = NOW()
	WHERE id = $2
	`

	_, err := r.DB.Exec(context.Background(), query,
		status,
		refID,
	)

	return err
}
