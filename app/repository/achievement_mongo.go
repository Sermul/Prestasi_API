package repository

import (
	"context"
	"time"

	"app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementMongoRepository struct {
	Collection *mongo.Collection
}

func NewAchievementMongoRepository(col *mongo.Collection) *AchievementMongoRepository {
	return &AchievementMongoRepository{Collection: col}
}

func (r *AchievementMongoRepository) CreateAchievementMongo(a *model.AchievementMongo) (primitive.ObjectID, error) {
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	res, err := r.Collection.InsertOne(context.Background(), a)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *AchievementMongoRepository) SoftDeleteAchievementMongo(id primitive.ObjectID) error {
	now := time.Now()

	_, err := r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"deletedAt": now,
			"updatedAt": now,
		}},
	)
	return err
}
