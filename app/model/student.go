package model

import "time"

type Student struct {
    ID           string    `db:"id" json:"id"`
    UserID       string    `db:"user_id" json:"userId"`
    StudentID    string    `db:"student_id" json:"studentId"`
    ProgramStudy string    `db:"program_study" json:"programStudy"`
    AcademicYear string    `db:"academic_year" json:"academicYear"`
    AdvisorID    string    `db:"advisor_id" json:"advisorId"`
    CreatedAt    time.Time `db:"created_at" json:"createdAt"`
    UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}
