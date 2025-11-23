package repository

import (
	"app/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementRepository interface {
	// Mongo
	CreateAchievementMongo(a *model.AchievementMongo) (primitive.ObjectID, error)
	SoftDeleteAchievementMongo(id primitive.ObjectID) error

	// Postgres
	CreateReferencePostgres(ref *model.AchievementReference) error
	UpdateReferenceStatusPostgres(refID string, status string) error
}
