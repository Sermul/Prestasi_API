package service

import (
	"prestasi_api/app/model"
	"prestasi_api/app/repository"
	"time"
 "errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService struct {
	MongoRepo    repository.AchievementMongoRepository
	PostgresRepo repository.AchievementPostgresRepository
	StudentRepo  repository.StudentPostgresRepository
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
    userID := c.Locals("student_id").(string) // dari JWT middleware

    // 1. Ambil reference lengkap
    ref, err := s.PostgresRepo.GetReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Reference not found"})
    }

    // 2. Cek apakah pemilik data
    if ref.StudentID != userID {
        return c.Status(403).JSON(fiber.Map{"error": "Not your achievement"})
    }

    // 3. Cek status harus draft
    if ref.Status != "draft" {
        return c.Status(400).JSON(fiber.Map{"error": "Only draft achievements can be deleted"})
    }

    // 4. Konversi mongo ID
    oid, err := primitive.ObjectIDFromHex(ref.MongoID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid MongoID"})
    }

    // 5. Soft delete MongoDB
    if err := s.MongoRepo.SoftDeleteAchievementMongo(oid); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    // 6. Update status Postgres → deleted
    if err := s.PostgresRepo.UpdateReferenceStatusPostgres(refID, "deleted"); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Achievement deleted successfully"})
}
// FR-006 — Get Achievements for Advisee Students (Dosen Wali)
func (s *AchievementService) GetAdviseeAchievements(advisorID string) ([]map[string]interface{}, error) {

	// 1. Ambil semua mahasiswa bimbingan dosen
	studentIDs, err := s.StudentRepo.GetStudentIDsByAdvisor(advisorID)
	if err != nil {
		return nil, err
	}

	// 2. Ambil semua reference prestasi milik mereka
	refs, err := s.PostgresRepo.GetByStudentIDs(studentIDs)
	if err != nil {
		return nil, err
	}

	// 3. Ambil detail prestasi dari MongoDB
	var results []map[string]interface{}

	for _, ref := range refs {
		oid, _ := primitive.ObjectIDFromHex(ref.MongoID)
		ach, _ := s.MongoRepo.GetByID(oid)

		result := map[string]interface{}{
			"reference":   ref,
			"achievement": ach,
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *AchievementService) Verify(advisorID string, refID string) error {

	// 1. Ambil reference
	ref, err := s.PostgresRepo.GetReferenceByID(refID)
	if err != nil {
		return errors.New("reference not found")
	}

	// 2. Ambil list mahasiswa bimbingan
	studentIDs, err := s.StudentRepo.GetStudentIDsByAdvisor(advisorID)
	if err != nil {
		return err
	}

	// 3. Validasi apakah prestasi milik mahasiswa bimbingan
	isAdvisee := false
	for _, sid := range studentIDs {
		if sid == ref.StudentID {
			isAdvisee = true
			break
		}
	}
	if !isAdvisee {
		return errors.New("not your advisee")
	}

	// 4. Validasi status
	if ref.Status != "submitted" {
		return errors.New("only submitted achievements can be verified")
	}

	// 5. Update status
	return s.PostgresRepo.UpdateVerifyStatus(refID, advisorID)
}


func (s *AchievementService) Reject(advisorID string, refID string, note string) error {

	ref, err := s.PostgresRepo.GetReferenceByID(refID)
	if err != nil {
		return err
	}

	// Status harus submitted
	if ref.Status != "submitted" {
		return errors.New("only submitted achievements can be rejected")
	}

	// Validasi mahasiswa bimbingan
	students, err := s.StudentRepo.GetStudentIDsByAdvisor(advisorID)
	if err != nil {
		return err
	}

	isMine := false
	for _, sid := range students {
		if sid == ref.StudentID {
			isMine = true
			break
		}
	}
	if !isMine {
		return errors.New("not your advisee")
	}

	// Update status
	return s.PostgresRepo.RejectReference(refID, advisorID, note)
}
