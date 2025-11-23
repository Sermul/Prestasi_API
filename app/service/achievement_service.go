package service

import (
	"app/model"
	"app/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService struct {
	MongoRepo repository.AchievementRepository
	PgRepo    repository.AchievementRepository
}

// =========================================
// CREATE (FR-003)
// =========================================
func CreateAchievementService(c *fiber.Ctx) error {
	var a model.AchievementMongo

	if err := c.BodyParser(&a); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	// create mongo
	mongoID, err := Mongo.Achievement.CreateAchievementMongo(&a)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	ref := model.AchievementReference{
		ID:        uuid.New().String(),
		StudentID: a.StudentID,
		MongoID:   mongoID.Hex(),
		Status:    "draft",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = Postgres.Achievement.CreateReferencePostgres(&ref)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":        "Achievement created",
		"referenceID":    ref.ID,
		"mongoID":        mongoID.Hex(),
	})
}

// =========================================
// SUBMIT (FR-004)
// =========================================
func SubmitAchievementService(c *fiber.Ctx) error {
	refID := c.Params("refId")

	err := Postgres.Achievement.UpdateReferenceStatusPostgres(refID, "submitted")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Achievement submitted"})
}

// =========================================
// SOFT DELETE (FR-005)
// =========================================
func DeleteAchievementService(c *fiber.Ctx) error {
	refID := c.Params("refId")

	// update postgres status
	err := Postgres.Achievement.UpdateReferenceStatusPostgres(refID, "deleted")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// convert mongo id
	mongoID, err := primitive.ObjectIDFromHex(refID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid mongo ID"})
	}

	err = Mongo.Achievement.SoftDeleteAchievementMongo(mongoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Achievement soft deleted"})
}
