package service

import (
    "prestasi_api/app/repository"
    "github.com/gofiber/fiber/v2"
)

type StudentService struct {
    StudentRepo  repository.StudentPostgresRepository
    LecturerRepo repository.LecturerPostgresRepository
}

func NewStudentService(studentRepo repository.StudentPostgresRepository, lecturerRepo repository.LecturerPostgresRepository) *StudentService {
    return &StudentService{
        StudentRepo:  studentRepo,
        LecturerRepo: lecturerRepo,
    }
}

func (s *StudentService) AssignAdvisor(c *fiber.Ctx) error {
    studentID := c.Params("id")

    var payload struct {
        LecturerID string `json:"lecturer_id"`
    }

    if err := c.BodyParser(&payload); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
    }

    // Ambil lecturer berdasarkan lecturer_id (contoh: DOSEN001)
    lec, err := s.LecturerRepo.GetByLecturerID(payload.LecturerID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Lecturer not found"})
    }

    // Update advisor memakai LECTURER UUID (bukan lecturer_id)
    err = s.StudentRepo.UpdateAdvisor(studentID, lec.ID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "message":     "Advisor assigned successfully",
        "student_id":  studentID,
        "lecturer_id": payload.LecturerID,
    })
}
