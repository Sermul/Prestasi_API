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


  // ASSIGN ADVISOR

func (s *StudentService) AssignAdvisor(c *fiber.Ctx) error {
    studentID := c.Params("id")

    var body struct {
        LecturerID string `json:"lecturer_id"`
    }

    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
    }

    lec, err := s.LecturerRepo.GetByLecturerID(body.LecturerID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Lecturer not found"})
    }

    if err := s.StudentRepo.UpdateAdvisor(studentID, lec.ID); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to update advisor"})
    }

    return c.JSON(fiber.Map{
        "message": "Advisor assigned successfully",
    })
}


   //LIST STUDENTS
func (s *StudentService) List(c *fiber.Ctx) error {
    students, err := s.StudentRepo.GetAll()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch students"})
    }
    return c.JSON(students)
}


   //DETAIL STUDENT
func (s *StudentService) Detail(c *fiber.Ctx) error {
    id := c.Params("id")

    student, err := s.StudentRepo.GetByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Student not found"})
    }

    return c.JSON(student)
}


  // STUDENT ACHIEVEMENTS
func (s *StudentService) Achievements(c *fiber.Ctx) error {
    id := c.Params("id")

    // cek dulu apakah student ada
    _, err := s.StudentRepo.GetByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Student not found"})
    }

    refs, err := s.StudentRepo.GetStudentAchievements(id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch achievements"})
    }

    // NOTE:
    // Sesuai modul, StudentService TIDAK ambil detail Mongo di sini.
    // Hanya mengembalikan references dari PostgreSQL.

    return c.JSON(refs)
}
