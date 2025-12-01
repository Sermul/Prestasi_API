package model

import "time"

type Student struct {
    ID           string     `json:"id"`
    UserID       string     `json:"user_id"`
    StudentID    string     `json:"student_id"`
    ProgramStudy string     `json:"program_study"`
    AcademicYear string     `json:"academic_year"`

    AdvisorID    *string     `json:"advisor_id"`   // ← WAJIB POINTER BIAR NULL BISA MASUK
    CreatedAt    *time.Time  `json:"created_at"`
    UpdatedAt    *time.Time  `json:"updated_at"`
}
