package service

import (
	"prestasi_api/app/model"
	"prestasi_api/app/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService struct {
	MongoRepo    repository.AchievementMongoRepository
	PostgresRepo repository.AchievementPostgresRepository
}

// FR-003 — Create Achievement
func (s *AchievementService) Create(c *fiber.Ctx) error {
	var data model.AchievementMongo

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	mongoID, err := s.MongoRepo.CreateAchievementMongo(&data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	ref := model.AchievementReference{
		ID:        uuid.New().String(),
		StudentID: data.StudentID,
		MongoID:   mongoID.Hex(),
		Status:    "draft",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.PostgresRepo.CreateReferencePostgres(&ref); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":     "Achievement created",
		"referenceID": ref.ID,
		"mongoID":     mongoID.Hex(),
	})
}

// FR-004 — Submit Achievement
func (s *AchievementService) Submit(c *fiber.Ctx) error {
	refID := c.Params("refId")

	err := s.PostgresRepo.UpdateReferenceStatusPostgres(refID, "submitted")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Achievement submitted"})
}

// FR-005 — Soft Delete Achievement
func (s *AchievementService) Delete(c *fiber.Ctx) error {
	refID := c.Params("refId")

	// Ambil mongoID dari postgres
	mongoIDStr, err := s.PostgresRepo.GetMongoID(refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Reference not found"})
	}

	oid, err := primitive.ObjectIDFromHex(mongoIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Mongo ID"})
	}

	// Soft delete postgres
	if err := s.PostgresRepo.UpdateReferenceStatusPostgres(refID, "deleted"); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Soft delete mongo
	if err := s.MongoRepo.SoftDeleteAchievementMongo(oid); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Achievement soft deleted"})
}
