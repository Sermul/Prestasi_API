package service

import (
	"prestasi_api/app/repository"

	"github.com/gofiber/fiber/v2"
)

type ReportService struct {
	PostgresRepo repository.AchievementPostgresRepository
	StudentRepo  repository.StudentPostgresRepository
}

func NewReportService(
	pg repository.AchievementPostgresRepository,
	student repository.StudentPostgresRepository,
) *ReportService {
	return &ReportService{
		PostgresRepo: pg,
		StudentRepo:  student,
	}
}

// ===============================
// FR-011 Achievement Statistics
// ===============================
func (s *ReportService) Statistics(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	// =====================
	// ADMIN → SEMUA DATA
	// =====================
	if role == "Admin" {
		totalStudents, err := s.StudentRepo.CountAll()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to count students"})
		}

		totalAchievements, err := s.PostgresRepo.CountAll()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to count achievements"})
		}

		return c.JSON(fiber.Map{
			"scope":              "all",
			"total_students":     totalStudents,
			"total_achievements": totalAchievements,
		})
	}

	// =====================
	// DOSEN WALI → BIMBINGAN
	// =====================
	if role == "Dosen Wali" {
		lecturerID := c.Locals("lecturer_id").(string)

		studentIDs, err := s.StudentRepo.GetStudentIDsByAdvisor(lecturerID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to load advisees"})
		}

		totalAchievements := 0
		if len(studentIDs) > 0 {
			totalAchievements, err = s.PostgresRepo.CountByStudentIDs(studentIDs)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "failed to count achievements"})
			}
		}

		return c.JSON(fiber.Map{
			"scope":              "advisees",
			"total_students":     len(studentIDs),
			"total_achievements": totalAchievements,
		})
	}

	// =====================
	// MAHASISWA → DATA SENDIRI
	// =====================
	if role == "Mahasiswa" {
		studentID := c.Locals("student_id").(string)

		totalAchievements, err := s.PostgresRepo.CountByStudentID(studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to count achievements"})
		}

		return c.JSON(fiber.Map{
			"scope":              "self",
			"total_students":     1,
			"total_achievements": totalAchievements,
		})
	}

	return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
}

// ===============================
// Report per Student (detail)
// ===============================
func (s *ReportService) StudentReport(c *fiber.Ctx) error {
	id := c.Params("id")

	refs, err := s.PostgresRepo.GetByStudentID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(refs)
}
