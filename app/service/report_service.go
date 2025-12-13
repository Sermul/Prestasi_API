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

// ======================================================
// FR-011 Achievement Statistics
// ======================================================
func (s *ReportService) Statistics(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	// =====================
	// ADMIN → SEMUA PRESTASI
	// =====================
	if role == "Admin" {
		refs, err := s.PostgresRepo.GetAllReferences()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		stats := map[string]int{
			"draft":     0,
			"submitted": 0,
			"verified":  0,
			"rejected":  0,
			"deleted":   0,
		}

		for _, r := range refs {
			stats[r.Status]++
		}

		return c.JSON(fiber.Map{
			"scope":          "all",
			"total_records": len(refs),
			"by_status":     stats,
		})
	}

	// =====================
	// DOSEN WALI → MAHASISWA BIMBINGAN
	// =====================
	if role == "Dosen Wali" {
		lecturerID := c.Locals("lecturer_id").(string)

		studentIDs, err := s.StudentRepo.GetStudentIDsByAdvisor(lecturerID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		refs, err := s.PostgresRepo.GetByStudentIDs(studentIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		stats := map[string]int{
			"draft":     0,
			"submitted": 0,
			"verified":  0,
			"rejected":  0,
			"deleted":   0,
		}

		for _, r := range refs {
			stats[r.Status]++
		}

		return c.JSON(fiber.Map{
			"scope":          "advisees",
			"total_records": len(refs),
			"by_status":     stats,
		})
	}

	// =====================
	// MAHASISWA → PRESTASI SENDIRI
	// =====================
	if role == "Mahasiswa" {
		studentID := c.Locals("student_id").(string)

		refs, err := s.PostgresRepo.GetByStudentID(studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		stats := map[string]int{
			"draft":     0,
			"submitted": 0,
			"verified":  0,
			"rejected":  0,
			"deleted":   0,
		}

		for _, r := range refs {
			stats[r.Status]++
		}

		return c.JSON(fiber.Map{
			"scope":          "self",
			"total_records": len(refs),
			"by_status":     stats,
		})
	}

	return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
}

// ======================================================
// Report Detail per Student
// ======================================================
func (s *ReportService) StudentReport(c *fiber.Ctx) error {
	targetStudentID := c.Params("id")
	role := c.Locals("role").(string)

	// =====================
	// RBAC
	// =====================
	if role == "Mahasiswa" {
		if c.Locals("student_id") != targetStudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if role == "Dosen Wali" {
		lecturerID := c.Locals("lecturer_id").(string)
		studentIDs, _ := s.StudentRepo.GetStudentIDsByAdvisor(lecturerID)

		allowed := false
		for _, id := range studentIDs {
			if id == targetStudentID {
				allowed = true
				break
			}
		}

		if !allowed {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	refs, err := s.PostgresRepo.GetByStudentID(targetStudentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	stats := map[string]int{
		"draft":     0,
		"submitted": 0,
		"verified":  0,
		"rejected":  0,
		"deleted":   0,
	}

	for _, r := range refs {
		stats[r.Status]++
	}

	return c.JSON(fiber.Map{
		"student_id":        targetStudentID,
		"total_achievements": len(refs),
		"by_status":         stats,
		"achievements":      refs,
	})
}
