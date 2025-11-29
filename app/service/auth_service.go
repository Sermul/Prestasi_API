package service

import (
	"time"

	"prestasi_api/app/model"
	"prestasi_api/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService sekarang memuat semua repository RBAC yang diperlukan
type AuthService struct {
	UserRepo           repository.UserPostgresRepository
	RoleRepo           repository.RolePostgresRepository
	PermissionRepo     repository.PermissionPostgresRepository
	RolePermissionRepo repository.RolePermissionPostgresRepository
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

	// Ambil role_name dari tabel roles
	role, err := s.RoleRepo.GetByID(user.RoleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to load role"})
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    role.Name,  // ⬅️ penting! sekarang role name
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("SECRET_KEY"))

	return c.JSON(fiber.Map{
		"token": signed,
		"user":  user,
	})
}

