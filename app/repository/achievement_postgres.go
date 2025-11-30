package repository

import (
	"context"
	"errors"
	"time"

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
	SaveSubmittedAt(refID string, t time.Time) error
	GetByStudentID(studentID string) ([]model.AchievementReference, error)
	GetAllReferences() ([]model.AchievementReference, error)
	GetHistoryByReferenceID(refID string) ([]map[string]interface{}, error)
}

type achievementPostgresRepo struct {
	pool *pgxpool.Pool
}

func NewAchievementPostgresRepository() AchievementPostgresRepository {
	return &achievementPostgresRepo{
		pool: database.Pg,
	}
}

// CREATE REFERENCE
func (r *achievementPostgresRepo) CreateReferencePostgres(ref *model.AchievementReference) error {
	_, err := r.pool.Exec(
		context.Background(),
		`INSERT INTO achievement_references
        (id, student_id, mongo_achievement_id, status,
         submitted_at, verified_at, verified_by, rejection_note,
         created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		ref.ID, ref.StudentID, ref.MongoID, ref.Status,
		ref.SubmittedAt, ref.VerifiedAt, ref.VerifiedBy, ref.RejectionNote,
		ref.CreatedAt, ref.UpdatedAt)
	return err
}

// UPDATE STATUS
func (r *achievementPostgresRepo) UpdateReferenceStatusPostgres(refID string, status string) error {
	_, err := r.pool.Exec(
		context.Background(),
		`UPDATE achievement_references
		 SET status=$1, updated_at=NOW()
		 WHERE id=$2`,
		status, refID,
	)
	return err
}

// GET MONGO ID
func (r *achievementPostgresRepo) GetMongoID(refID string) (string, error) {
	var mongoID *string

	err := r.pool.QueryRow(
		context.Background(),
		`SELECT mongo_achievement_id FROM achievement_references WHERE id=$1`,
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

// GET REFERENCE BY ID
func (r *achievementPostgresRepo) GetReferenceByID(refID string) (*model.AchievementReference, error) {
	var ref model.AchievementReference

	err := r.pool.QueryRow(
		context.Background(),
		`SELECT id, student_id, mongo_achievement_id, status,
		        submitted_at, verified_at, verified_by, rejection_note,
		        created_at, updated_at
		 FROM achievement_references WHERE id=$1`,
		refID,
	).Scan(
		&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
		&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
		&ref.CreatedAt, &ref.UpdatedAt,
	)

	if err != nil {
		return nil, errors.New("reference not found")
	}

	return &ref, nil
}

// GET BY MULTIPLE STUDENT IDs (Advisor)
func (r *achievementPostgresRepo) GetByStudentIDs(studentIDs []string) ([]model.AchievementReference, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT id, student_id, mongo_achievement_id, status, created_at, updated_at
         FROM achievement_references
         WHERE student_id = ANY($1)`,
		studentIDs)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status, &ref.CreatedAt, &ref.UpdatedAt)
		refs = append(refs, ref)
	}
	return refs, nil
}

// VERIFY
func (r *achievementPostgresRepo) UpdateVerifyStatus(refID string, verifierID string) error {
	_, err := r.pool.Exec(
		context.Background(),
		`UPDATE achievement_references
		 SET status='verified', verified_by=$1, verified_at=NOW(), updated_at=NOW()
		 WHERE id=$2`,
		verifierID, refID)
	return err
}

// REJECT
func (r *achievementPostgresRepo) RejectReference(refID string, advisorID string, note string) error {
	_, err := r.pool.Exec(
		context.Background(),
		`UPDATE achievement_references
		 SET status='rejected', rejection_note=$1,
		     verified_by=$2, verified_at=NOW(),
		     updated_at=NOW()
		 WHERE id=$3`,
		note, advisorID, refID)
	return err
}

// SAVE SUBMITTED_AT
func (r *achievementPostgresRepo) SaveSubmittedAt(refID string, t time.Time) error {
	_, err := r.pool.Exec(
		context.Background(),
		`UPDATE achievement_references
         SET submitted_at=$1, updated_at=NOW()
         WHERE id=$2`,
		t, refID)
	return err
}

// LIST OWN (Mahasiswa)
func (r *achievementPostgresRepo) GetByStudentID(studentID string) ([]model.AchievementReference, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT id, student_id, mongo_achievement_id, status,
		        submitted_at, verified_at, verified_by, rejection_note,
		        created_at, updated_at
		 FROM achievement_references
		 WHERE student_id=$1`,
		studentID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
			&ref.CreatedAt, &ref.UpdatedAt,
		)
		refs = append(refs, ref)
	}
	return refs, nil
}

// ADMIN LIST ALL
func (r *achievementPostgresRepo) GetAllReferences() ([]model.AchievementReference, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT id, student_id, mongo_achievement_id, status,
		        submitted_at, verified_at, verified_by, rejection_note,
		        created_at, updated_at
		 FROM achievement_references
		 ORDER BY created_at DESC`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoID, &ref.Status,
			&ref.SubmittedAt, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote,
			&ref.CreatedAt, &ref.UpdatedAt,
		)
		refs = append(refs, ref)
	}
	return refs, nil
}

// HISTORY
func (r *achievementPostgresRepo) GetHistoryByReferenceID(refID string) ([]map[string]interface{}, error) {
	rows, err := r.pool.Query(context.Background(),
		`SELECT status, note, changed_by, changed_at
		 FROM achievement_reference_history
		 WHERE reference_id=$1
		 ORDER BY changed_at ASC`,
		refID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var status, note, changedBy *string
		var changedAt *time.Time

		rows.Scan(&status, &note, &changedBy, &changedAt)

		entry := map[string]interface{}{
			"status":     nil,
			"note":       nil,
			"changed_by": nil,
			"changed_at": nil,
		}

		if status != nil {
			entry["status"] = *status
		}
		if note != nil {
			entry["note"] = *note
		}
		if changedBy != nil {
			entry["changed_by"] = *changedBy
		}
		if changedAt != nil {
			entry["changed_at"] = *changedAt
		}

		history = append(history, entry)
	}

	return history, nil
}
