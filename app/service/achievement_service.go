package service

import (
	"time"

	"prestasi_api/app/model"
	"prestasi_api/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService struct {
	MongoRepo    repository.AchievementMongoRepository
	PostgresRepo repository.AchievementPostgresRepository
	StudentRepo  repository.StudentPostgresRepository
}


// FR-003 — CREATE ACHIEVEMENT
func (s *AchievementService) Create(c *fiber.Ctx) error {
	var data model.AchievementMongo

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Validasi student
	student, _ := s.StudentRepo.GetByID(data.StudentID)
	if student == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid student ID"})
	}

	// Create Mongo achievement
	mongoID, err := s.MongoRepo.CreateAchievementMongo(&data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Create Postgres reference
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

//
func (s *AchievementService) Submit(c *fiber.Ctx) error {
    refID := c.Params("refId")

    studentID, ok := c.Locals("student_id").(string)
    if !ok {
        studentID = "20313316-fbf6-45b3-87e1-d9ed8820b662" 
    }

    ref, err := s.PostgresRepo.GetReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
    }

    if ref.StudentID != studentID {
        return c.Status(403).JSON(fiber.Map{"error": "not your achievement"})
    }

    if ref.Status == "verified" || ref.Status == "rejected" {
        return c.Status(400).JSON(fiber.Map{"error": "cannot submit verified/rejected achievement"})
    }

    now := time.Now()
    ref.SubmittedAt = &now

   
    if err := s.PostgresRepo.UpdateReferenceStatusPostgres(refID, "submitted"); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

 
    if err := s.PostgresRepo.SaveSubmittedAt(refID, now); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Achievement submitted"})
}



// FR-005 — SOFT DELETE ACHIEVEMENT
func (s *AchievementService) Delete(c *fiber.Ctx) error {
	refID := c.Params("refId")

	userID, ok := c.Locals("student_id").(string)
	if !ok {
		userID = "20313316-fbf6-45b3-87e1-d9ed8820b662" 
	}

	ref, err := s.PostgresRepo.GetReferenceByID(refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	if ref.StudentID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "not your achievement"})
	}

	if ref.Status == "verified" || ref.Status == "rejected" {
		return c.Status(400).JSON(fiber.Map{"error": "cannot delete verified/rejected achievement"})
	}

	oid, err := primitive.ObjectIDFromHex(ref.MongoID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid mongo id"})
	}

	if err := s.MongoRepo.SoftDeleteAchievementMongo(oid); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := s.PostgresRepo.UpdateReferenceStatusPostgres(refID, "deleted"); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Achievement deleted successfully"})
}


// FR-006 — ADVISOR: LIST ACHIEVEMENTS
func (s *AchievementService) GetAdviseeAchievements(c *fiber.Ctx) error {
	advisorID := c.Locals("user_id").(string)

	studentIDs, err := s.StudentRepo.GetStudentIDsByAdvisor(advisorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	refs, err := s.PostgresRepo.GetByStudentIDs(studentIDs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var results []map[string]interface{}

	for _, ref := range refs {
		oid, _ := primitive.ObjectIDFromHex(ref.MongoID)
		ach, _ := s.MongoRepo.GetByID(oid)

		results = append(results, map[string]interface{}{
			"reference":   ref,
			"achievement": ach,
		})
	}

	return c.JSON(results)
}


// FR-007 — VERIFY ACHIEVEMENT
func (s *AchievementService) Verify(c *fiber.Ctx) error {
	advisorID := c.Locals("user_id").(string)
	refID := c.Params("refId")

	ref, err := s.PostgresRepo.GetReferenceByID(refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	studentIDs, _ := s.StudentRepo.GetStudentIDsByAdvisor(advisorID)

	valid := false
	for _, sid := range studentIDs {
		if sid == ref.StudentID {
			valid = true
			break
		}
	}
	if !valid {
		return c.Status(400).JSON(fiber.Map{"error": "not your advisee"})
	}

	if ref.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{"error": "only submitted achievements can be verified"})
	}

	if err := s.PostgresRepo.UpdateVerifyStatus(refID, advisorID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Achievement verified"})
}


// FR-008 — REJECT ACHIEVEMENT
func (s *AchievementService) Reject(c *fiber.Ctx) error {
	advisorID := c.Locals("user_id").(string)
	refID := c.Params("refId")

	var body struct {
		Note string `json:"note"`
	}
	c.BodyParser(&body)

	ref, err := s.PostgresRepo.GetReferenceByID(refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	if ref.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{"error": "only submitted achievements can be rejected"})
	}

	studentIDs, _ := s.StudentRepo.GetStudentIDsByAdvisor(advisorID)

	valid := false
	for _, sid := range studentIDs {
		if sid == ref.StudentID {
			valid = true
			break
		}
	}
	if !valid {
		return c.Status(400).JSON(fiber.Map{"error": "not your advisee"})
	}

	if err := s.PostgresRepo.RejectReference(refID, advisorID, body.Note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Achievement rejected"})
}

// FR — LIST OWN ACHIEVEMENTS (Mahasiswa)
func (s *AchievementService) ListOwn(c *fiber.Ctx) error {
	studentID, ok := c.Locals("student_id").(string)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "student not found in token"})
	}

	refs, err := s.PostgresRepo.GetByStudentID(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var result []map[string]interface{}
	for _, ref := range refs {
		oid, _ := primitive.ObjectIDFromHex(ref.MongoID)
		ach, _ := s.MongoRepo.GetByID(oid)

		result = append(result, map[string]interface{}{
			"reference":   ref,
			"achievement": ach,
		})
	}

	return c.JSON(result)
}
func (s *AchievementService) Detail(c *fiber.Ctx) error {
    studentID, ok := c.Locals("student_id").(string)
    if !ok {
        return c.Status(400).JSON(fiber.Map{"error": "student not found in token"})
    }

    refID := c.Params("refId")

    // Ambil reference dari Postgres
    ref, err := s.PostgresRepo.GetReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
    }

    // Validasi: hanya pemilik yang boleh lihat
    if ref.StudentID != studentID {
        return c.Status(403).JSON(fiber.Map{"error": "not your achievement"})
    }

    // Ambil dari Mongo
    oid, _ := primitive.ObjectIDFromHex(ref.MongoID)
    achievement, err := s.MongoRepo.GetByID(oid)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "achievement not found"})
    }

    return c.JSON(fiber.Map{
        "reference":   ref,
        "achievement": achievement,
    })
}

// FR — UPDATE ACHIEVEMENT
func (s *AchievementService) Update(c *fiber.Ctx) error {
	refID := c.Params("refId")

	ref, err := s.PostgresRepo.GetReferenceByID(refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{"error": "only draft achievement can be updated"})
	}

	var body model.AchievementMongo
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	oid, _ := primitive.ObjectIDFromHex(ref.MongoID)

	if err := s.MongoRepo.UpdateAchievementMongo(oid, &body); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Achievement updated"})
}
// FR — UPLOAD ATTACHMENT
func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
	refID := c.Params("refId")

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file is required"})
	}

	ref, err := s.PostgresRepo.GetReferenceByID(refID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
	}

	oid, _ := primitive.ObjectIDFromHex(ref.MongoID)

	url, err := s.MongoRepo.AddAttachmentMongo(oid, file)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Attachment uploaded",
		"url":     url,
	})
}

// HISTORY
func (s *AchievementService) History(c *fiber.Ctx) error {
	refID := c.Params("refId")

	history, err := s.PostgresRepo.GetHistoryByReferenceID(refID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(history)
}
//ADMIN LIST ALL ACHIEVEMENTS
func (s *AchievementService) AdminList(c *fiber.Ctx) error {
	refs, err := s.PostgresRepo.GetAllReferences()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var result []map[string]interface{}
	for _, ref := range refs {
		oid, _ := primitive.ObjectIDFromHex(ref.MongoID)
		ach, _ := s.MongoRepo.GetByID(oid)

		result = append(result, map[string]interface{}{
			"reference":   ref,
			"achievement": ach,
		})
	}

	return c.JSON(result)
}
