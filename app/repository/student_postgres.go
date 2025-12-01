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
}

type studentPostgresRepo struct{}

func NewStudentPostgresRepository() StudentPostgresRepository {
    return &studentPostgresRepo{}
}

// ================================
// GetStudentIDsByAdvisor
// ================================
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

// ================================
// GetByID
// ================================
func (r *studentPostgresRepo) GetByID(studentID string) (*model.Student, error) {
    var s model.Student

    err := database.Pg.QueryRow(
        context.Background(),
        `SELECT id, user_id, student_id, program_study, academic_year, advisor_id,
                created_at, updated_at
         FROM students
         WHERE id = $1 OR student_id = $1
         LIMIT 1`,
        studentID,
    ).Scan(
        &s.ID,
        &s.UserID,
        &s.StudentID,
        &s.ProgramStudy,
        &s.AcademicYear,
        &s.AdvisorID,
        &s.CreatedAt,
        &s.UpdatedAt,
    )

    if err != nil {
        return nil, errors.New("student not found")
    }

    return &s, nil
}

// ================================
// GetByUserID
// ================================
func (r *studentPostgresRepo) GetByUserID(userID string) (*model.Student, error) {
    var s model.Student

    err := database.Pg.QueryRow(
        context.Background(),
        `SELECT id, user_id, student_id, program_study, academic_year, advisor_id,
                created_at, updated_at
         FROM students
         WHERE user_id = $1
         LIMIT 1`,
        userID,
    ).Scan(
        &s.ID,
        &s.UserID,
        &s.StudentID,
        &s.ProgramStudy,
        &s.AcademicYear,
        &s.AdvisorID,
        &s.CreatedAt,
        &s.UpdatedAt,
    )

    if err != nil {
        return nil, errors.New("student not found")
    }

    return &s, nil
}
