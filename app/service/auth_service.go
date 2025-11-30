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
) *AuthService {
	return &AuthService{
		UserRepo:           userRepo,
		RoleRepo:           roleRepo,
		PermissionRepo:     permissionRepo,
		RolePermissionRepo: rolePermissionRepo,
	}
}

// =========================
// REGISTER
// =========================
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

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 12)

	// Ambil role default "Mahasiswa"
	role, err := s.RoleRepo.GetByName("Mahasiswa")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Role not found"})
	}

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

	if err := s.UserRepo.Create(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "User registered",
	})
}

// =========================
// LOGIN
// =========================
func (s *AuthService) Login(c *fiber.Ctx) error {
    var body struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
    }

    user, err := s.UserRepo.GetByEmail(body.Email)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "User not found"})
    }

    if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)) != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Wrong password"})
    }

    role, err := s.RoleRepo.GetByID(user.RoleID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to load role"})
    }

    // ==== ACCESS TOKEN (baru & benar) ====
accessClaims := jwt.MapClaims{
    "user_id": user.ID,
    "role":    role.Name, // ⬅ WAJIB! Untuk RoleGuard
    "exp":     time.Now().Add(24 * time.Hour).Unix(),
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

    // Simpan refresh token ke cookie
    c.Cookie(&fiber.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        HTTPOnly: true,
        Secure:   false,
        Path:     "/",
        MaxAge:   7 * 24 * 3600,
    })

    return c.JSON(fiber.Map{
        "access_token": accessToken,
        "user":         user,
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

    newAccessToken, err := helper.GenerateToken(userID, roleID)
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
