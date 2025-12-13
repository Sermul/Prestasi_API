package service

import (
    "prestasi_api/app/repository"

    "github.com/gofiber/fiber/v2"
)

type LecturerService struct {
    StudentRepo repository.StudentPostgresRepository
     LecturerRepo repository.LecturerPostgresRepository
}

func NewLecturerService(
    studentRepo repository.StudentPostgresRepository,
    lecturerRepo repository.LecturerPostgresRepository,
) *LecturerService {
    return &LecturerService{
        StudentRepo:  studentRepo,
        LecturerRepo: lecturerRepo,
    }
}

func (s *LecturerService) ListAdvisees(c *fiber.Ctx) error {
    lecturerID := c.Params("id")

    ids, err := s.StudentRepo.GetStudentIDsByAdvisor(lecturerID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    var students []interface{}
    for _, id := range ids {
        st, err := s.StudentRepo.GetByID(id)
        if err != nil {
            // skip not found entries
            continue
        }
        students = append(students, st)
    }

    if students == nil {
        students = []interface{}{}
    }

    return c.JSON(students)
}
func (s *LecturerService) List(c *fiber.Ctx) error {
	list, err := s.LecturerRepo.List()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

func (s *LecturerService) Detail(c *fiber.Ctx) error {
	id := c.Params("id")
	l, err := s.LecturerRepo.Detail(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Lecturer not found"})
	}
	return c.JSON(l)
}
