package service

import (
	"time"

	    "strings"

	"prestasi_api/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	
	"prestasi_api/helper"
	"golang.org/x/crypto/bcrypt"
)

// AuthService sekarang memuat semua repository RBAC yang diperlukan
type AuthService struct {
	UserRepo           repository.UserPostgresRepository
	RoleRepo           repository.RolePostgresRepository
    StudentRepo        repository.StudentPostgresRepository   
	LecturerRepo       repository.LecturerPostgresRepository

    // JWT JWTService  // ⬅️ Tambahkan ini
}



// CONSTRUCTOR (RBAC penuh)
func NewAuthService(
    userRepo repository.UserPostgresRepository,
    roleRepo repository.RolePostgresRepository,
    studentRepo repository.StudentPostgresRepository,
    lecturerRepo repository.LecturerPostgresRepository,
) *AuthService {
    return &AuthService{
        UserRepo:           userRepo,
        RoleRepo:           roleRepo,
        StudentRepo:        studentRepo,
        LecturerRepo:       lecturerRepo,  
    }
}






// LOGIN
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




    // student berdasarkan user_id 
  var studentID string

//  Student ID (jika Mahasiswa) 
if role.Name == "Mahasiswa" {
    student, err := s.StudentRepo.GetByUserID(user.ID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "student not found"})
    }
    studentID = student.ID
}

//  Lecturer ID (jika Dosen Wali) 
var lecturerID string
if role.Name == "Dosen Wali" {
    lecturer, err := s.LecturerRepo.GetByUserID(user.ID)
    if err == nil {
        lecturerID = lecturer.ID
    }
}

// ACCESS TOKEN 
accessClaims := jwt.MapClaims{
    "user_id":     user.ID,
    "role_id":     role.ID,      
    "role":        role.Name,
    "student_id":  studentID,
    "lecturer_id": lecturerID,
    "exp":         time.Now().Add(24 * time.Hour).Unix(),
}

accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
accessToken, _ := accessTokenObj.SignedString([]byte("SECRET_KEY"))

// REFRESH TOKEN 
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

    // CEK ROLE DULU 
    role, _ := s.RoleRepo.GetByID(roleID)
    studentID := ""
    lecturerID := ""

    if role != nil && role.Name == "Mahasiswa" {
        student, _ := s.StudentRepo.GetByUserID(userID)
        if student != nil {
            studentID = student.ID
        }
    }

    if role != nil && role.Name == "Dosen Wali" {
        lecturer, _ := s.LecturerRepo.GetByUserID(userID)
        if lecturer != nil {
            lecturerID = lecturer.ID
        }
    }
  

    // Generate token baru
    newAccessToken, err := helper.GenerateFullToken(userID, roleID, studentID, lecturerID, role.Name)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to generate access token"})
    }

    return c.JSON(fiber.Map{
        "access_token": newAccessToken,
    })
}



func (s *AuthService) Logout(c *fiber.Ctx) error {

    refresh := c.Cookies("refresh_token")
    access := c.Get("Authorization")

    if refresh != "" {
        helper.AddToBlacklist(refresh)
    }

    if access != "" && strings.HasPrefix(access, "Bearer ") {
        access = strings.TrimPrefix(access, "Bearer ")
        helper.AddToBlacklist(access)
    }

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
