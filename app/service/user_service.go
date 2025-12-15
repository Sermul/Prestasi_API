package service

import (
    "prestasi_api/app/model"
    "prestasi_api/app/repository"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
      "golang.org/x/crypto/bcrypt" 
)

type UserService struct {
    UserRepo repository.UserPostgresRepository
    RoleRepo repository.RolePostgresRepository
    StudentRepo  repository.StudentPostgresRepository
    LecturerRepo repository.LecturerPostgresRepository
}

func NewUserService(
    userRepo repository.UserPostgresRepository,
    roleRepo repository.RolePostgresRepository,
    studentRepo repository.StudentPostgresRepository,
    lecturerRepo repository.LecturerPostgresRepository,
) *UserService {
    return &UserService{
        UserRepo:     userRepo,
        RoleRepo:     roleRepo,
        StudentRepo:  studentRepo,
        LecturerRepo: lecturerRepo,
    }
}

// LIST USERS
func (s *UserService) List(c *fiber.Ctx) error {
    users, err := s.UserRepo.GetAll()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(users)
}

// DETAIL USER
func (s *UserService) Detail(c *fiber.Ctx) error {
    id := c.Params("id")
    user, err := s.UserRepo.GetByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "user not found"})
    }
    return c.JSON(user)
}

// CREATE USER
func (s *UserService) Create(c *fiber.Ctx) error {
    var body struct {
        Username  string `json:"username"`
        Email     string `json:"email"`
        Password  string `json:"password"`
        FullName  string `json:"full_name"`
        RoleID    string `json:"role_id"`
    }

    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    // Hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
    }

    user := model.User{
        ID:           uuid.New().String(),
        Username:     body.Username,
        Email:        body.Email,
        PasswordHash: string(hash),
        FullName:     body.FullName,
        RoleID:       body.RoleID,
        IsActive:     true,                // default sesuai modul
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    if err := s.UserRepo.Create(&user); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(user)
}


// UPDATE USER
func (s *UserService) Update(c *fiber.Ctx) error {
    id := c.Params("id")

    var body struct {
        Username *string `json:"username"`
        Email    *string `json:"email"`
        FullName *string `json:"full_name"`
        IsActive *bool   `json:"is_active"`
    }

    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    // Ambil user lama dari DB
    user, err := s.UserRepo.GetByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "user not found"})
    }

    // Update hanya field yg dikirim
    if body.Username != nil {
        user.Username = *body.Username
    }
    if body.Email != nil {
        user.Email = *body.Email
    }
    if body.FullName != nil {
        user.FullName = *body.FullName
    }
    if body.IsActive != nil {
        user.IsActive = *body.IsActive
    }

    // Password TIDAK DIUBAH !!!
    // user.PasswordHash TETAP seperti sebelumnya

    user.UpdatedAt = time.Now()

    if err := s.UserRepo.Update(id, user); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(user)
}



// DELETE USER
func (s *UserService) Delete(c *fiber.Ctx) error {
    id := c.Params("id")

    if err := s.UserRepo.Delete(id); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "user deleted"})
}

// CHANGE ROLE
func (s *UserService) ChangeRole(c *fiber.Ctx) error {
    id := c.Params("id")

    var body struct {
        RoleID string `json:"role_id"`
    }
    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    // 1️⃣ ambil role
    role, err := s.RoleRepo.GetByID(body.RoleID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid role_id"})
    }

    // 2️⃣ update role user
    if err := s.UserRepo.UpdateRole(id, role.ID); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    // 3️⃣ JIKA ROLE = DOSEN WALI → TAMBAH KE TABEL LECTURERS
    if role.Name == "Dosen Wali" {

        user, err := s.UserRepo.GetByID(id)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "user not found"})
        }

lecturer := model.Lecturer{
    ID:         uuid.New().String(),
    UserID:     user.ID,
    LecturerID: user.Username,      // atau kode dosen jika ada
    Department: "Teknik Informatika",
    CreatedAt:  time.Now(),
}


        if err := s.LecturerRepo.Create(&lecturer); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
    }

    return c.JSON(fiber.Map{
        "message": "role updated",
    })
}

