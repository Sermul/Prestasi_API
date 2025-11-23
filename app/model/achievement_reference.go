package model

import "time"

type AchievementReference struct {
	ID               string     `db:"id"`
	StudentID        string     `db:"student_id"`
	MongoID          string     `db:"mongo_achievement_id"`
	Status           string     `db:"status"`
	SubmittedAt      *time.Time `db:"submitted_at"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
}
