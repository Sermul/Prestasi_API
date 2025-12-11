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


func (s *AchievementService) Create(c *fiber.Ctx) error {
    var data model.AchievementMongo

    if err := c.BodyParser(&data); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
    }

    role, _ := c.Locals("role").(string)

    // === VALIDASI STUDENT ID ===

    // Mahasiswa tidak boleh kirim studentId di body
    if role == "Mahasiswa" && data.StudentID != "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "Mahasiswa tidak boleh mengirim studentId di body",
        })
    }

    // Admin wajib kirim studentId
    if role == "Admin" && data.StudentID == "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "Admin wajib menyertakan studentId",
        })
    }

    // === Ambil Student ID ===
    var studentID string
    if role == "Mahasiswa" {
        sid, ok := c.Locals("student_id").(string)
        if !ok || sid == "" {
            return c.Status(401).JSON(fiber.Map{"error": "student not found in token"})
        }
        studentID = sid
    } else if role == "Admin" {
        studentID = data.StudentID
    } else {
        return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
    }

    data.StudentID = studentID

    now := time.Now()
    data.CreatedAt = now
    data.UpdatedAt = now

    data.Points = calculatePoints(&data)

    mongoID, err := s.MongoRepo.CreateAchievementMongo(&data)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    ref := model.AchievementReference{
        ID:        uuid.New().String(),
        StudentID: studentID,
        MongoID:   mongoID.Hex(),
        Status:    "draft",
        CreatedAt: now,
        UpdatedAt: now,
    }

    if err := s.PostgresRepo.CreateReferencePostgres(&ref); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    // HISTORY: CREATED (status draft)
s.saveHistory(
    ref.ID,
    "",
    "draft",
    c.Locals("user_id").(string),
    role,
    "",
)


    return c.JSON(fiber.Map{
        "message":     "Achievement created",
        "referenceID": ref.ID,
        "mongoID":     mongoID.Hex(),
    })
}


func calculatePoints(a *model.AchievementMongo) int {
    switch a.AchievementType {
    case "competition":
        switch a.Details.CompetitionLevel {
        case "international":
            return 200
        case "national":
            return 100
        case "provincial":
            return 60
        case "city":
            return 40
        default:
            return 20
        }
    case "publication":
        return 150
    case "organization":
        return 50
    case "certification":
        return 80
    default:
        return 10
    }
}


//
func (s *AchievementService) Submit(c *fiber.Ctx) error {
    refID := c.Params("refId")
    role, _ := c.Locals("role").(string)

    ref, err := s.PostgresRepo.GetReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
    }

    // Mahasiswa: hanya boleh submit prestasinya sendiri
    if role == "Mahasiswa" {
        studentID := c.Locals("student_id").(string)
        if ref.StudentID != studentID {
            return c.Status(403).JSON(fiber.Map{"error": "not your achievement"})
        }
    }

    // Admin bypass semua aturan
    if ref.Status != "draft" {
        return c.Status(400).JSON(fiber.Map{"error": "only draft can be submitted"})
    }
// HISTORY: submit
s.saveHistory(
    ref.ID,
    ref.Status,
    "submitted",
    c.Locals("user_id").(string),
    role,
    "",
)

    now := time.Now()
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
    role, _ := c.Locals("role").(string)

    ref, err := s.PostgresRepo.GetReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
    }

    //=== Mahasiswa validation ===
    if role == "Mahasiswa" {
        studentID := c.Locals("student_id").(string)

        if ref.StudentID != studentID {
            return c.Status(403).JSON(fiber.Map{"error": "not your achievement"})
        }
        if ref.Status != "draft" {
            return c.Status(400).JSON(fiber.Map{"error": "cannot delete verified/submitted/rejected achievement"})
        }
    }

    //=== Admin bebas delete tanpa syarat ===

    oid, _ := primitive.ObjectIDFromHex(ref.MongoID)
    if err := s.MongoRepo.SoftDeleteAchievementMongo(oid); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
// HISTORY: delete
s.saveHistory(
    ref.ID,
    ref.Status,
    "deleted",
    c.Locals("user_id").(string),
    role,
    "",
)

    if err := s.PostgresRepo.UpdateReferenceStatusPostgres(refID, "deleted"); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Achievement deleted"})
}



// FR-006 — ADVISOR: LIST ACHIEVEMENTS
func (s *AchievementService) GetAdviseeAchievements(c *fiber.Ctx) error {
	advisorID := c.Locals("lecturer_id").(string)


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
    advisorLecturerID := c.Locals("lecturer_id").(string)
    advisorUserID := c.Locals("user_id").(string) // <-- pakai user_id
    role, _ := c.Locals("role").(string)

    refID := c.Params("refId")

    ref, err := s.PostgresRepo.GetReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference tidak ditemukan"})
    }

    // If role is Dosen Wali => validate advisee
    if role == "Dosen Wali" {
        studentIDs, _ := s.StudentRepo.GetStudentIDsByAdvisor(advisorLecturerID)
        valid := false
        for _, sid := range studentIDs {
            if sid == ref.StudentID {
                valid = true
                break
            }
        }
        if !valid {
            return c.Status(400).JSON(fiber.Map{"error": "bukan mahasiswa bimbingan"})
        }
    } else if role != "Admin" {
        // only Dosen Wali or Admin can verify
        return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
    }

    if ref.Status != "submitted" {
        return c.Status(400).JSON(fiber.Map{"error": "hanya status submitted yang bisa diverifikasi"})
    }
// HISTORY: verify
s.saveHistory(
    ref.ID,
    ref.Status,
    "verified",
    advisorUserID,
    role,
    "",
)

    // verifier id: for Admin use user_id claim as well
    if err := s.PostgresRepo.UpdateVerifyStatus(refID, advisorUserID); err != nil { // <-- save user_id
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Prestasi berhasil diverifikasi"})
}


// FR-008 — REJECT ACHIEVEMENT
func (s *AchievementService) Reject(c *fiber.Ctx) error {
    advisorLecturerID := c.Locals("lecturer_id").(string)
    advisorUserID := c.Locals("user_id").(string) // <-- pakai user_id
    role, _ := c.Locals("role").(string)

    refID := c.Params("refId")

    var body struct {
        Note string `json:"note"`
    }
    c.BodyParser(&body)

    ref, err := s.PostgresRepo.GetReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference tidak ditemukan"})
    }

    if ref.Status != "submitted" {
        return c.Status(400).JSON(fiber.Map{"error": "hanya status submitted yang bisa ditolak"})
    }

    // If role is Dosen Wali => validate advisee
    if role == "Dosen Wali" {
        studentIDs, _ := s.StudentRepo.GetStudentIDsByAdvisor(advisorLecturerID)
        valid := false
        for _, sid := range studentIDs {
            if sid == ref.StudentID {
                valid = true
                break
            }
        }
        if !valid {
            return c.Status(400).JSON(fiber.Map{"error": "bukan mahasiswa bimbingan"})
        }
    } else if role != "Admin" {
        // only Dosen Wali or Admin can reject
        return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
    }
// HISTORY: reject
s.saveHistory(
    ref.ID,
    ref.Status,
    "rejected",
    advisorUserID,
    role,
    body.Note,
)

    if err := s.PostgresRepo.RejectReference(refID, advisorUserID, body.Note); err != nil { // <-- save user_id
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Prestasi berhasil ditolak"})
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


// === LIST BY ROLE ===
func (s *AchievementService) List(c *fiber.Ctx) error {
    role, _ := c.Locals("role").(string)

    // === 1) Mahasiswa: lihat punya sendiri ===
    if role == "Mahasiswa" {
        return s.ListOwn(c)
    }

    // === 2) Dosen Wali: lihat prestasi mahasiswa bimbingan ===
  if role == "Dosen Wali" {
    advisorID, _ := c.Locals("lecturer_id").(string)

    studentIDs, err := s.StudentRepo.GetStudentIDsByAdvisor(advisorID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    // Tidak punya mahasiswa bimbingan
    if len(studentIDs) == 0 {
        return c.JSON(fiber.Map{"message": "Anda belum memiliki mahasiswa bimbingan"})
    }

    refs, err := s.PostgresRepo.GetByStudentIDs(studentIDs)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    // Ada mahasiswa tapi belum ada prestasi
    if len(refs) == 0 {
        return c.JSON(fiber.Map{"message": "Belum ada prestasi dari mahasiswa bimbingan"})
    }

    return buildAchievementResponse(c, refs, s)
}


    // === 3) Admin: lihat semua ===
    if role == "Admin" {
        refs, err := s.PostgresRepo.GetAllReferences()
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        return buildAchievementResponse(c, refs, s)
    }

    return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
}
 

// Helper to format response
func buildAchievementResponse(c *fiber.Ctx, refs []model.AchievementReference, s *AchievementService) error {
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
    role := c.Locals("role").(string)
    refID := c.Params("refId")

    // Ambil reference
    ref, err := s.PostgresRepo.GetReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
    }

    // ROLE: Mahasiswa → hanya miliknya sendiri
    if role == "Mahasiswa" {
        studentID := c.Locals("student_id").(string)
        if ref.StudentID != studentID {
            return c.Status(403).JSON(fiber.Map{"error": "not your achievement"})
        }
    }

    // ROLE: Dosen Wali → hanya mahasiswa bimbingan
    if role == "Dosen Wali" {
        lecturerID := c.Locals("lecturer_id").(string)

        studentIDs, _ := s.StudentRepo.GetStudentIDsByAdvisor(lecturerID)

        isAdvisee := false
        for _, sid := range studentIDs {
            if sid == ref.StudentID {
                isAdvisee = true
                break
            }
        }
        if !isAdvisee {
            return c.Status(403).JSON(fiber.Map{"error": "not your advisee"})
        }
    }

    // ROLE: Admin → bebas akses (tidak perlu validasi)

    // Ambil data dari Mongo
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
    role, _ := c.Locals("role").(string)

    // Ambil reference
    ref, err := s.PostgresRepo.GetReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
    }

    // === MAHASISWA VALIDATION ===
    if role == "Mahasiswa" {
        studentID := c.Locals("student_id").(string)

        // Mahasiswa hanya bisa update prestasi miliknya sendiri
        if ref.StudentID != studentID {
            return c.Status(403).JSON(fiber.Map{"error": "not your achievement"})
        }

        // Status harus draft
        if ref.Status != "draft" {
            return c.Status(400).JSON(fiber.Map{"error": "only draft can be updated"})
        }
    }

    // Parse body
    var body model.AchievementMongo
    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    // ❌ Mahasiswa tidak boleh kirim studentId
    if role == "Mahasiswa" && body.StudentID != "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "Mahasiswa tidak boleh mengirim studentId di body",
        })
    }

    // Admin bebas melakukan update → studentId diabaikan
    body.StudentID = "" // supaya tidak ke-set ulang di Mongo

    oid, _ := primitive.ObjectIDFromHex(ref.MongoID)
    if err := s.MongoRepo.UpdateAchievementMongo(oid, &body); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Achievement updated"})
}


// FR — UPLOAD ATTACHMENT
func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {
    refID := c.Params("refId")
    role, _ := c.Locals("role").(string)

    // ❌ Mahasiswa tidak boleh mengirim studentId di form-data
    if role == "Mahasiswa" && c.FormValue("studentId") != "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "Mahasiswa tidak boleh mengirim studentId",
        })
    }

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


// HISTORY — lihat perubahan status prestasi
func (s *AchievementService) History(c *fiber.Ctx) error {
    refID := c.Params("refId")
    role := c.Locals("role").(string)

    // Ambil reference dulu
    ref, err := s.PostgresRepo.GetReferenceByID(refID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "reference not found"})
    }

    // ===== VALIDASI ROLE =====

    // Mahasiswa → hanya miliknya sendiri
    if role == "Mahasiswa" {
        studentID := c.Locals("student_id").(string)
        if ref.StudentID != studentID {
            return c.Status(403).JSON(fiber.Map{
                "error": "Anda tidak boleh mengakses history prestasi orang lain",
            })
        }
    }

    // Dosen Wali → hanya mahasiswa bimbingan
    if role == "Dosen Wali" {
        lecturerID := c.Locals("lecturer_id").(string)

        adviseeIDs, _ := s.StudentRepo.GetStudentIDsByAdvisor(lecturerID)

        isAdvisee := false
        for _, sid := range adviseeIDs {
            if sid == ref.StudentID {
                isAdvisee = true
                break
            }
        }

        if !isAdvisee {
            return c.Status(403).JSON(fiber.Map{
                "error": "Prestasi ini bukan milik mahasiswa bimbingan Anda",
            })
        }
    }

    // Admin → bisa akses semua

    history, err := s.PostgresRepo.GetHistoryByReferenceID(refID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(history)
}


func (s *AchievementService) saveHistory(refID, oldStatus, newStatus, userID, role, note string) error {
    h := model.AchievementReferenceHistory{
        ID:            uuid.New().String(),
        ReferenceID:   refID,
        OldStatus:     oldStatus,
        NewStatus:     newStatus,
        Note:          note,
        ChangedBy:     userID,
        ChangedByRole: role,
        CreatedAt:     time.Now(),
    }

    return s.PostgresRepo.InsertHistory(&h)
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


