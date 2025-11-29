package repository

import (
	"prestasi_api/database"
	"prestasi_api/app/model"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementMongoRepository interface {
	CreateAchievementMongo(data *model.AchievementMongo) (primitive.ObjectID, error)
	SoftDeleteAchievementMongo(id primitive.ObjectID) error
	RestoreAchievementMongo(id primitive.ObjectID) error
	GetByID(id primitive.ObjectID) (*model.AchievementMongo, error)
	GetAll() ([]model.AchievementMongo, error)
}

type achievementMongoRepo struct {
	collection *mongo.Collection
}

func NewAchievementMongoRepository() AchievementMongoRepository {
	return &achievementMongoRepo{
		collection: database.Mongo.Collection("achievements"),
	}
}


// CREATE (FR-003)
func (r *achievementMongoRepo) CreateAchievementMongo(data *model.AchievementMongo) (primitive.ObjectID, error) {
	ctx := context.TODO()

	data.ID = primitive.NewObjectID()
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	data.DeletedAt = nil // soft delete default

	_, err := r.collection.InsertOne(ctx, data)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return data.ID, nil
}


// SOFT DELETE (FR-005)

func (r *achievementMongoRepo) SoftDeleteAchievementMongo(id primitive.ObjectID) error {
	ctx := context.TODO()

	now := time.Now()

	_, err := r.collection.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{
			"deletedAt": now,
		},
	})

	return err
}


// RESTORE (SoftDelete)

func (r *achievementMongoRepo) RestoreAchievementMongo(id primitive.ObjectID) error {
	ctx := context.TODO()

	_, err := r.collection.UpdateByID(ctx, id, bson.M{
		"$unset": bson.M{
			"deletedAt": "",
		},
	})

	return err
}


// GET BY ID

func (r *achievementMongoRepo) GetByID(id primitive.ObjectID) (*model.AchievementMongo, error) {
    ctx := context.TODO()

    var result model.AchievementMongo

    err := r.collection.FindOne(ctx, bson.M{
        "_id": id,
        "deletedAt": bson.M{"$exists": false},
    }).Decode(&result)
    if err != nil {
        return nil, errors.New("data not found")
    }

    return &result, nil
}


// GET ALL

func (r *achievementMongoRepo) GetAll() ([]model.AchievementMongo, error) {
    ctx := context.TODO()

    cursor, err := r.collection.Find(ctx, bson.M{
        "deletedAt": bson.M{"$exists": false},
    })
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []model.AchievementMongo
    err = cursor.All(ctx, &results)

    return results, err
}

