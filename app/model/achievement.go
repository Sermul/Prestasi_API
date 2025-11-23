package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementMongo struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	StudentID       string             `bson:"studentId"`
	AchievementType string             `bson:"achievementType"`
	Title           string             `bson:"title"`
	Description     string             `bson:"description"`
	Details         interface{}        `bson:"details"`
	Tags            []string           `bson:"tags"`
	Points          int                `bson:"points"`
	CreatedAt       time.Time          `bson:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt"`
	DeletedAt       *time.Time         `bson:"deletedAt,omitempty"`
}
