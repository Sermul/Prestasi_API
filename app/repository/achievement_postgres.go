package repository

import (
	"context"
	"errors"
	"prestasi_api/app/model"
	"prestasi_api/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AchievementPostgresRepository interface {
	CreateReferencePostgres(ref *model.AchievementReference) error
	UpdateReferenceStatusPostgres(refID string, status string) error
	GetMongoID(refID string) (string, error)
	GetReferenceByID(refID string) (*model.AchievementReference, error)
	GetByStudentIDs(studentIDs []string) ([]model.AchievementReference, error)
	UpdateVerifyStatus(refID string, verifierID string) error
RejectReference(refID string, advisorID string, note string) error

}

type achievementPostgresRepo struct {
	pool *pgxpool.Pool
}

func NewAchievementPostgresRepository() AchievementPostgresRepository {
	return &achievementPostgresRepo{
		pool: database.Pg,
	}
}

// CREATE
func (r *achievementPostgresRepo) CreateReferencePostgres(ref *model.AchievementReference) error {
	_, err := r.pool.Exec(
		context.Background(),
		 `INSERT INTO achievement_references
        (id, student_id, mongo_achievement_id, status,
         submitted_at, verified_at, verified_by, rejection_note,
         created_at, updated_at)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
    ref.ID,
    ref.StudentID,
    ref.MongoID,
    ref.Status,
    ref.SubmittedAt,
    ref.VerifiedAt,
    ref.VerifiedBy,
    ref.RejectionNote,
    ref.CreatedAt,
    ref.UpdatedAt,
)
	return err
}

// UPDATE STATUS
func (r *achievementPostgresRepo) UpdateReferenceStatusPostgres(refID string, status string) error {
	_, err := r.pool.Exec(
		context.Background(),
		`UPDATE achievement_references
		 SET status = $1, updated_at = NOW()
		 WHERE id = $2`,
		status, refID,
	)
	return err
}

// GET MONGO ACHIEVEMENT ID
func (r *achievementPostgresRepo) GetMongoID(refID string) (string, error) {
	var mongoID *string
err := r.pool.QueryRow(
    context.Background(),
    `SELECT mongo_achievement_id FROM achievement_references WHERE id = $1`,
    refID,
).Scan(&mongoID)

if err != nil {
    return "", errors.New("reference not found")
}

if mongoID == nil {
    return "", nil
}

return *mongoID, nil

}

// GET FULL REFERENCE
func (r *achievementPostgresRepo) GetReferenceByID(refID string) (*model.AchievementReference, error) {
	var ref model.AchievementReference

	err := r.pool.QueryRow(
		context.Background(),
		`SELECT id, student_id, mongo_achievement_id, status,
            submitted_at, verified_at, verified_by, rejection_note,
            created_at, updated_at
     FROM achievement_references
     WHERE id = $1`,
    refID,
).Scan(
    &ref.ID,
    &ref.StudentID,
    &ref.MongoID,
    &ref.Status,
    &ref.SubmittedAt,
    &ref.VerifiedAt,
    &ref.VerifiedBy,
    &ref.RejectionNote,
    &ref.CreatedAt,
    &ref.UpdatedAt,
)

	if err != nil {
		return nil, errors.New("reference not found")
	}

	return &ref, nil
}
func (r *achievementPostgresRepo) GetByStudentIDs(studentIDs []string) ([]model.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, created_at, updated_at
		FROM achievement_references
        WHERE student_id = ANY($1)
	`

	rows, err := r.pool.Query(context.Background(), query, studentIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference

	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoID,
			&ref.Status,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		refs = append(refs, ref)
	}

	return refs, nil
}
func (r *achievementPostgresRepo) UpdateVerifyStatus(refID string, verifierID string) error {
	_, err := r.pool.Exec(
		context.Background(),
		`UPDATE achievement_references
		 SET status = 'verified',
		     verified_by = $1,
		     verified_at = NOW(),
		     updated_at = NOW()
		 WHERE id = $2`,
		verifierID, refID,
	)
	return err
}
func (r *achievementPostgresRepo) RejectReference(refID string, advisorID string, note string) error {
	_, err := r.pool.Exec(
		context.Background(),
		`UPDATE achievement_references
		 SET status = 'rejected',
		     rejection_note = $1,
		     verified_by = $2,
		     verified_at = NOW(),
		     updated_at = NOW()
		 WHERE id = $3`,
		note,
		advisorID,
		refID,
	)

	return err
}
