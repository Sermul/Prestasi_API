package service

import (
    "prestasi_api/app/repository"

    "github.com/gofiber/fiber/v2"
)

type LecturerService struct {
    StudentRepo repository.StudentPostgresRepository
}

func NewLecturerService(studentRepo repository.StudentPostgresRepository) *LecturerService {
    return &LecturerService{
        StudentRepo: studentRepo,
    }
}

func (s *LecturerService) ListAdvisees(c *fiber.Ctx) error {
    lecturerID := c.Params("id")

    students, err := s.StudentRepo.GetStudentIDsByAdvisor(lecturerID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(students)
}
