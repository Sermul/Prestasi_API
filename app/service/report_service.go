package service

import (
	"prestasi_api/app/repository"

	"github.com/gofiber/fiber/v2"
)

type ReportService struct {
	MongoRepo   repository.AchievementMongoRepository
	StudentRepo repository.StudentPostgresRepository
}

func NewReportService(
	mongo repository.AchievementMongoRepository,
	student repository.StudentPostgresRepository,
) *ReportService {
	return &ReportService{
		MongoRepo:   mongo,
		StudentRepo: student,
	}
}

// ======================================================
// FR-011 Achievement Statistics (FIXED & FINAL)
// ======================================================
func (s *ReportService) Statistics(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	var scope string
	var allowedStudentIDs map[string]bool = nil

	// =========================
	// ROLE & FILTER STUDENT ID
	// =========================
	switch role {
	case "Admin":
		scope = "all"
		// admin lihat semua → no filter

	case "Dosen Wali":
		scope = "advisees"

		lecturerID := c.Locals("lecturer_id").(string)
		studentIDs, err := s.StudentRepo.GetStudentIDsByAdvisor(lecturerID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		allowedStudentIDs = map[string]bool{}
		for _, id := range studentIDs {
			allowedStudentIDs[id] = true
		}

	case "Mahasiswa":
		scope = "self"

		studentID := c.Locals("student_id").(string)
		allowedStudentIDs = map[string]bool{
			studentID: true,
		}

	default:
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	// =========================
	// AMBIL DATA MONGO
	// =========================
	achievements, err := s.MongoRepo.GetAllForReport()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// =========================
	// AGGREGATION
	// =========================
	totalByType := map[string]int{}
	totalByPeriod := map[string]int{}
	competitionLevel := map[string]int{}
	studentStats := map[string]struct {
		Points int
		Count  int
	}{}

	for _, a := range achievements {

		// =========================
		// FILTER BERDASARKAN ROLE
		// =========================
		if allowedStudentIDs != nil {
			if !allowedStudentIDs[a.StudentID] {
				continue
			}
		}

		// by type
		if a.AchievementType != "" {
			totalByType[a.AchievementType]++
		}

		// by year
		year := a.CreatedAt.Format("2006")
		totalByPeriod[year]++

		// competition level
		if a.AchievementType == "competition" {
			level := a.Details.CompetitionLevel
			if level != "" {
				competitionLevel[level]++
			}
		}

		// top students
		stat := studentStats[a.StudentID]
		stat.Points += a.Points
		stat.Count++
		studentStats[a.StudentID] = stat
	}

	// =========================
	// TOP STUDENTS RESPONSE
	// =========================
	var topStudents []fiber.Map
	for sid, stat := range studentStats {
		topStudents = append(topStudents, fiber.Map{
			"student_id":         sid,
			"total_achievements": stat.Count,
			"total_points":       stat.Points,
		})
	}

	// =========================
	// RESPONSE FINAL
	// =========================
	return c.JSON(fiber.Map{
		"scope":                          scope,
		"total_by_type":                 totalByType,
		"total_by_period":               totalByPeriod,
		"competition_level_distribution": competitionLevel,
		"top_students":                  topStudents,
	})
}

// ======================================================
// Report Detail per Student (TIDAK DIUBAH)
// ======================================================
func (s *ReportService) StudentReport(c *fiber.Ctx) error {
	targetStudentID := c.Params("id")
	role := c.Locals("role").(string)

	// =====================
	// RBAC (AUTORISASI)
	// =====================
	if role == "Mahasiswa" {
		if c.Locals("student_id") != targetStudentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	if role == "Dosen Wali" {
		lecturerID := c.Locals("lecturer_id").(string)
		studentIDs, err := s.StudentRepo.GetStudentIDsByAdvisor(lecturerID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

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

	// =====================
	// AMBIL DATA MONGO
	// =====================
	achievements, err := s.MongoRepo.GetAllForReport()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// =====================
	// AGGREGATION
	// =====================
	totalAchievements := 0
	totalPoints := 0
	byType := map[string]int{}
	byYear := map[string]int{}

	for _, a := range achievements {

		// filter mahasiswa target
		if a.StudentID != targetStudentID {
			continue
		}

		// total
		totalAchievements++
		totalPoints += a.Points

		// by type
		if a.AchievementType != "" {
			byType[a.AchievementType]++
		}

		// by year (pakai createdAt)
		year := a.CreatedAt.Format("2006")
		byYear[year]++
	}

	// =====================
	// RESPONSE FINAL
	// =====================
	return c.JSON(fiber.Map{
		"student_id": targetStudentID,
		"summary": fiber.Map{
			"total_achievements": totalAchievements,
			"total_points":      totalPoints,
		},
		"by_type": byType,
		"by_year": byYear,
	})
}
