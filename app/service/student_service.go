package service

import (
    "prestasi_api/app/repository"
    "github.com/gofiber/fiber/v2"
    "prestasi_api/app/model"
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
        AdvisorID  string `json:"advisor_id"`
    }

    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
    }

    // Accept both lecturer_id and advisor_id
    lecID := body.LecturerID
    if lecID == "" {
        lecID = body.AdvisorID
    }

    if lecID == "" {
        return c.Status(400).JSON(fiber.Map{"error": "lecturer_id or advisor_id required"})
    }

    // Try to find lecturer by lecturer_id field first, then by ID
    var lecturer *model.Lecturer
    var err error
    
    lecturer, err = s.LecturerRepo.GetByLecturerID(lecID)
    if err != nil {
        // If not found by lecturer_id, try by ID
        lecturer, err = s.LecturerRepo.Detail(lecID)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "Lecturer not found"})
        }
    }

    if err := s.StudentRepo.UpdateAdvisor(studentID, lecturer.ID); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to update advisor"})
    }

    return c.JSON(fiber.Map{
        "message": "Advisor assigned successfully",
    })
}


   //LIST STUDENTS
func (s *StudentService) List(c *fiber.Ctx) error {
    role := c.Locals("role").(string)

    // =====================
    // ADMIN → lihat semua
    // =====================
    if role == "Admin" {
        students, err := s.StudentRepo.GetAll()
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch students"})
        }
        return c.JSON(students)
    }

    // =====================
    // DOSEN → mahasiswa bimbingan saja
    // =====================
    if role == "Dosen Wali" {
        lecturerID := c.Locals("lecturer_id").(string)

        studentIDs, err := s.StudentRepo.GetStudentIDsByAdvisor(lecturerID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch advisees"})
        }

        var students []*model.Student
        for _, id := range studentIDs {
            sdt, err := s.StudentRepo.GetByID(id)
            if err == nil {
                students = append(students, sdt)
            }
        }

        if students == nil {
            students = []*model.Student{}
        }

        return c.JSON(students)
    }

    return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
}


   //DETAIL STUDENT
//DETAIL STUDENT
func (s *StudentService) Detail(c *fiber.Ctx) error {
    id := c.Params("id")

    student, err := s.StudentRepo.GetByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Student not found"})
    }

    // =====================
    // AUTHORIZATION CHECK
    // =====================
    role := c.Locals("role").(string)
    
    // Admin bisa lihat semua
    if role == "Admin" {
        return c.JSON(student)
    }

    // Dosen hanya bisa lihat mahasiswa bimbingannnya
    if role == "Dosen Wali" {
        lecturerID := c.Locals("lecturer_id").(string)
        if student.AdvisorID != lecturerID {
            return c.Status(403).JSON(fiber.Map{"error": "Anda bukan pembimbing mahasiswa ini"})
        }
        return c.JSON(student)
    }

    return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
}


  // STUDENT ACHIEVEMENTS
func (s *StudentService) Achievements(c *fiber.Ctx) error {
    id := c.Params("id")

    // cek dulu apakah student ada
    student, err := s.StudentRepo.GetByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Student not found"})
    }

    // Authorization: only Admin or the student's advisor (Dosen Wali) can view
    role := c.Locals("role").(string)
    if role != "Admin" {
        if role == "Dosen Wali" {
            lecturerID := c.Locals("lecturer_id").(string)
            if student.AdvisorID != lecturerID {
                return c.Status(403).JSON(fiber.Map{"error": "Anda bukan pembimbing mahasiswa ini"})
            }
        } else {
            return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
        }
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
