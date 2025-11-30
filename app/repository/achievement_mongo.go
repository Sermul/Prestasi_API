package repository

import (
	"prestasi_api/database"
	"prestasi_api/app/model"
	"context"
	"errors"
	"time"
	 "io"
	 "os"
	  "path/filepath"
    "mime/multipart"
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
	UpdateAchievementMongo(id primitive.ObjectID, a *model.AchievementMongo) error
    AddAttachmentMongo(id primitive.ObjectID, file *multipart.FileHeader) (string, error)
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
// UPDATE ACHIEVEMENT (used by service.Update)
func (r *achievementMongoRepo) UpdateAchievementMongo(id primitive.ObjectID, a *model.AchievementMongo) error {
    ctx := context.TODO()
    a.UpdatedAt = time.Now()

    // build update doc: set fields from a (you may want to be selective)
    update := bson.M{
        "$set": bson.M{
            "achievementType": a.AchievementType,
            "title":           a.Title,
            "description":     a.Description,
            "details":         a.Details,
            "tags":            a.Tags,
            "points":          a.Points,
            "updatedAt":       a.UpdatedAt,
        },
    }

    _, err := r.collection.UpdateByID(ctx, id, update)
    return err
}

// ADD ATTACHMENT: saves file to ./uploads/<id>/ and returns URL/path
func (r *achievementMongoRepo) AddAttachmentMongo(id primitive.ObjectID, file *multipart.FileHeader) (string, error) {
    ctx := context.TODO()
    // ensure uploads directory
    uploadDir := filepath.Join(".", "uploads", id.Hex())
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        return "", err
    }

    // open uploaded file
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    // destination path
    dstPath := filepath.Join(uploadDir, file.Filename)
    dst, err := os.Create(dstPath)
    if err != nil {
        return "", err
    }
    defer dst.Close()

    // copy content
    if _, err := io.Copy(dst, src); err != nil {
        return "", err
    }

    // Optionally update the achievement document with attachment metadata
    _, err = r.collection.UpdateByID(ctx, id, bson.M{
        "$push": bson.M{
            "attachments": bson.M{
                "fileName":  file.Filename,
                "fileUrl":   dstPath, // change to actual public URL if needed
                "fileType":  file.Header.Get("Content-Type"),
                "uploadedAt": time.Now(),
            },
        },
        "$set": bson.M{"updatedAt": time.Now()},
    })
    if err != nil {
        return "", err
    }

    return dstPath, nil
}


