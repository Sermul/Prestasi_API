package service

import (
    "prestasi_api/app/repository"

    "github.com/gofiber/fiber/v2"
)

type ReportService struct {
    PostgresRepo repository.AchievementPostgresRepository
    StudentRepo  repository.StudentPostgresRepository
}

func NewReportService(pg repository.AchievementPostgresRepository, student repository.StudentPostgresRepository) *ReportService {
    return &ReportService{
        PostgresRepo: pg,
        StudentRepo:  student,
    }
}

// Statistik umum
func (s *ReportService) Statistics(c *fiber.Ctx) error {
    data := map[string]interface{}{
        "total_students": 0,
        "total_achievements": 0,
    }

    return c.JSON(data)
}

// Report per student
func (s *ReportService) StudentReport(c *fiber.Ctx) error {
    id := c.Params("id")

    refs, err := s.PostgresRepo.GetByStudentID(id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(refs)
}
