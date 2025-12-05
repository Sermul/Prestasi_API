package repository

import (
    "context"
    "errors"
    "prestasi_api/app/model"
    "prestasi_api/database"
)

type StudentPostgresRepository interface {
    GetStudentIDsByAdvisor(advisorID string) ([]string, error)
    GetByID(studentID string) (*model.Student, error)
    GetByUserID(userID string) (*model.Student, error)
    Create(student *model.Student) error
    UpdateAdvisor(studentID, lecturerID string) error
    GetAll() ([]*model.Student, error)
    GetStudentAchievements(studentID string) ([]*model.AchievementReference, error)
}

type studentPostgresRepo struct{}

func NewStudentPostgresRepository() StudentPostgresRepository {
    return &studentPostgresRepo{}
}

/* ================================
   GetStudentIDsByAdvisor
================================ */
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

/* ================================
   GetByID  (FIXED)
================================ */
func (r *studentPostgresRepo) GetByID(id string) (*model.Student, error) {
    var s model.Student

    err := database.Pg.QueryRow(
        context.Background(),
        `SELECT 
            id, user_id, student_id, program_study, academic_year,
            COALESCE(advisor_id, '') AS advisor_id,
            created_at, updated_at
         FROM students
         WHERE id = $1 OR student_id = $1
         LIMIT 1`,
        id,
    ).Scan(
        &s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy,
        &s.AcademicYear, &s.AdvisorID, &s.CreatedAt, &s.UpdatedAt,
    )

    if err != nil {
        return nil, errors.New("student not found")
    }

    return &s, nil
}

/* ================================
   GetByUserID
================================ */
func (r *studentPostgresRepo) GetByUserID(userID string) (*model.Student, error) {
    var s model.Student

    err := database.Pg.QueryRow(
        context.Background(),
        `SELECT id, user_id, student_id, program_study, academic_year,
                advisor_id, created_at, updated_at
         FROM students
         WHERE user_id = $1
         LIMIT 1`,
        userID,
    ).Scan(
        &s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy,
        &s.AcademicYear, &s.AdvisorID, &s.CreatedAt, &s.UpdatedAt,
    )

    if err != nil {
        return nil, errors.New("student not found")
    }

    return &s, nil
}

/* ================================
   Create
================================ */
func (r *studentPostgresRepo) Create(s *model.Student) error {
    _, err := database.Pg.Exec(
        context.Background(),
        `INSERT INTO students 
        (id, user_id, student_id, program_study, academic_year, advisor_id, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
        s.ID, s.UserID, s.StudentID, s.ProgramStudy,
        s.AcademicYear, s.AdvisorID, s.CreatedAt, s.UpdatedAt,
    )
    return err
}

/* ================================
   UpdateAdvisor
================================ */
func (r *studentPostgresRepo) UpdateAdvisor(studentID, lecturerID string) error {
    _, err := database.Pg.Exec(
        context.Background(),
        `UPDATE students 
         SET advisor_id = $1, updated_at = NOW()
         WHERE id = $2`,
        lecturerID, studentID,
    )
    return err
}

/* ================================
   GetAll
================================ */
func (r *studentPostgresRepo) GetAll() ([]*model.Student, error) {
    rows, err := database.Pg.Query(context.Background(),
        `SELECT id, user_id, student_id, program_study, academic_year,
                advisor_id, created_at, updated_at
         FROM students
         ORDER BY created_at DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []*model.Student

    for rows.Next() {
        var s model.Student
        if err := rows.Scan(
            &s.ID, &s.UserID, &s.StudentID,
            &s.ProgramStudy, &s.AcademicYear,
            &s.AdvisorID, &s.CreatedAt, &s.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        result = append(result, &s)
    }

    return result, nil
}

/* ================================
   GetStudentAchievements (FIXED TYPO)
================================ */
func (r *studentPostgresRepo) GetStudentAchievements(studentID string) ([]*model.AchievementReference, error) {
    rows, err := database.Pg.Query(context.Background(),
        `SELECT id, student_id, mongo_achievement_id, status,
                submitted_at, verified_at, verified_by,
                rejection_note, created_at, updated_at
         FROM achievement_references
         WHERE student_id = $1
         ORDER BY created_at DESC`,
        studentID,
    )
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var list []*model.AchievementReference

    for rows.Next() {
        var a model.AchievementReference
        if err := rows.Scan(
            &a.ID, &a.StudentID, &a.MongoID,
            &a.Status, &a.SubmittedAt, &a.VerifiedAt,
            &a.VerifiedBy, &a.RejectionNote,
            &a.CreatedAt, &a.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        list = append(list, &a)
    }

    return list, nil
}
