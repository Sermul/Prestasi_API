package repository

import (
    "prestasi_api/app/model"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementRepository interface {
    
    // MONGO FUNCTIONS
    CreateAchievementMongo(a *model.AchievementMongo) (primitive.ObjectID, error)
    GetAchievementMongoByID(id primitive.ObjectID) (*model.AchievementMongo, error)
    UpdateAchievementMongo(id primitive.ObjectID, a *model.AchievementMongo) error
    SoftDeleteAchievementMongo(id primitive.ObjectID) error

   
    // POSTGRES FUNCTIONS
    CreateReferencePostgres(ref *model.AchievementReference) error
    GetReferenceByID(refID string) (*model.AchievementReference, error)
    UpdateReferenceStatusPostgres(refID string, status string) error

    // Verify (FR-007)
    VerifyReference(refID string, verifiedBy string) error

    // Reject (FR-008)
    RejectReference(refID string, note string) error

    // List achievements  (FR-006)
    GetAdviseeAchievements(advisorID string) ([]model.AchievementReference, error)
}
