package service

import (
	"time"

	"prestasi_api/app/model"
	"prestasi_api/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"prestasi_api/helper"
	"golang.org/x/crypto/bcrypt"
)

// AuthService sekarang memuat semua repository RBAC yang diperlukan
type AuthService struct {
	UserRepo           repository.UserPostgresRepository
	RoleRepo           repository.RolePostgresRepository
	PermissionRepo     repository.PermissionPostgresRepository
	RolePermissionRepo repository.RolePermissionPostgresRepository
    StudentRepo        repository.StudentPostgresRepository   // ⬅ WAJIB
	 LecturerRepo       repository.LecturerPostgresRepository
    // JWT JWTService  // ⬅️ Tambahkan ini
}


// =========================
// CONSTRUCTOR (RBAC penuh)
// =========================
func NewAuthService(
    userRepo repository.UserPostgresRepository,
    roleRepo repository.RolePostgresRepository,
    permissionRepo repository.PermissionPostgresRepository,
    rolePermissionRepo repository.RolePermissionPostgresRepository,
    studentRepo repository.StudentPostgresRepository,
    lecturerRepo repository.LecturerPostgresRepository,
) *AuthService {
    return &AuthService{
        UserRepo:           userRepo,
        RoleRepo:           roleRepo,
        PermissionRepo:     permissionRepo,
        RolePermissionRepo: rolePermissionRepo,
        StudentRepo:        studentRepo,
        LecturerRepo:       lecturerRepo,  // <--- WAJIB
    }
}


func (s *AuthService) Register(c *fiber.Ctx) error {
    var body struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
        FullName string `json:"full_name"`
    }

    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
    }

    // Hash password
    hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 12)

    // Ambil role "Mahasiswa"
    role, err := s.RoleRepo.GetByName("Mahasiswa")
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Role not found"})
    }

    // Buat user baru
    user := model.User{
        ID:           uuid.New().String(),
        Username:     body.Username,
        Email:        body.Email,
        PasswordHash: string(hash),
        FullName:     body.FullName,
        RoleID:       role.ID,
        IsActive:     true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    // Simpan user
    if err := s.UserRepo.Create(&user); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    // ==========================================
    // AUTO CREATE STUDENT (SOLUSINYA DI SINI)
    // ==========================================
    if role.Name == "Mahasiswa" {
        now := time.Now()
        student := model.Student{
            ID:           uuid.New().String(),
            UserID:       user.ID,
            StudentID:    body.Username,   // atau generate sendiri jika mau
            ProgramStudy: "Teknik Informatika", // default
            AcademicYear: "2023",               // default
            AdvisorID:    nil,                  // bisa diisi nanti
            CreatedAt:    &now,
            UpdatedAt:    &now,
        }

        if err := s.StudentRepo.Create(&student); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to create student"})
        }
    }

    return c.JSON(fiber.Map{
        "message": "User registered successfully",
    })
}


// =========================
// LOGIN
// =========================
func (s *AuthService) Login(c *fiber.Ctx) error {
  var body struct {
    Username string `json:"username"`
    Password string `json:"password"`
}


    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
    }

   user, err := s.UserRepo.GetByUsername(body.Username)

    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "User not found"})
    }

    if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)) != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Wrong password"})
    }

   // Ambil role
role, err := s.RoleRepo.GetByID(user.RoleID)
if err != nil {
    return c.Status(500).JSON(fiber.Map{"error": "failed to load role"})
}

// Ambil permissions
permissions, err := s.PermissionRepo.GetByRoleID(role.ID)
if err != nil {
    return c.Status(500).JSON(fiber.Map{"error": "failed to load permissions"})
}


    // === Ambil student berdasarkan user_id ===
  var studentID string

// === Student ID (jika Mahasiswa) ===
if role.Name == "Mahasiswa" {
    student, err := s.StudentRepo.GetByUserID(user.ID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "student not found"})
    }
    studentID = student.ID
}

// === Lecturer ID (jika Dosen Wali) ===
var lecturerID string
if role.Name == "Dosen Wali" {
    lecturer, err := s.LecturerRepo.GetByUserID(user.ID)
    if err == nil {
        lecturerID = lecturer.ID
    }
}

// ==== ACCESS TOKEN ====
accessClaims := jwt.MapClaims{
    "user_id":     user.ID,
    "role_id":     role.ID,      
    "role":        role.Name,
    "permissions": permissions,
    "student_id":  studentID,
    "lecturer_id": lecturerID,
    "exp":         time.Now().Add(24 * time.Hour).Unix(),
}

accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
accessToken, _ := accessTokenObj.SignedString([]byte("SECRET_KEY"))

// ==== REFRESH TOKEN ====
refreshClaims := jwt.MapClaims{
    "user_id": user.ID,
    "role_id": role.ID,
    "exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
}
refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
refreshToken, _ := refreshTokenObj.SignedString([]byte("SECRET_KEY"))

c.Cookie(&fiber.Cookie{
    Name:     "refresh_token",
    Value:    refreshToken,
    HTTPOnly: true,
    Secure:   false,
    Path:     "/",
    MaxAge:   7 * 24 * 3600,
})

return c.JSON(fiber.Map{
    "status": "success",
    "data": fiber.Map{
        "token":        accessToken,
        "refreshToken": refreshToken,
        "user": fiber.Map{
            "id":          user.ID,
            "username":    user.Username,
            "fullName":    user.FullName,
            "role":        role.Name,
            "permissions": permissions,
            "student_id":  studentID,
            "lecturer_id": lecturerID,
        },
    },
})
}



func (s *AuthService) Refresh(c *fiber.Ctx) error {
    refresh := c.Cookies("refresh_token")
    if refresh == "" {
        return c.Status(401).JSON(fiber.Map{"error": "refresh token missing"})
    }

    token, err := helper.ParseToken(refresh)
    if err != nil || !token.Valid {
        return c.Status(401).JSON(fiber.Map{"error": "invalid refresh token"})
    }

    claims := token.Claims.(jwt.MapClaims)
    userID := claims["user_id"].(string)
    roleID := claims["role_id"].(string)

    // ============== FIX DI SINI ==============
    studentID := ""
    if roleID != "" {
        student, _ := s.StudentRepo.GetByUserID(userID)
        if student != nil {
            studentID = student.ID
        }
    }
    // ==========================================

    newAccessToken, err := helper.GenerateToken(userID, roleID, studentID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to generate access token"})
    }

    return c.JSON(fiber.Map{
        "access_token": newAccessToken,
    })
}


func (s *AuthService) Logout(c *fiber.Ctx) error {
    c.ClearCookie("access_token")
    c.ClearCookie("refresh_token")

    return c.JSON(fiber.Map{
        "message": "Logout successful",
    })
}
func (s *AuthService) Profile(c *fiber.Ctx) error {
    userID := c.Locals("user_id")
    if userID == nil {
        return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
    }

    user, err := s.UserRepo.GetByID(userID.(string))
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "user not found"})
    }

    return c.JSON(user)
}
