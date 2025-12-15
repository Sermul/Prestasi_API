package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"prestasi_api/app/model"
	// "prestasi_api/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

//
// ======================================================
// MOCK USER REPOSITORIES (DIPISAH PER SKENARIO)
// ======================================================
//

// USER TIDAK DITEMUKAN
type MockUserRepoNotFound struct{}

func (m *MockUserRepoNotFound) GetByUsername(username string) (*model.User, error) {
	return nil, errors.New("user not found")
}
func (m *MockUserRepoNotFound) Create(user *model.User) error               { return nil }
func (m *MockUserRepoNotFound) GetByEmail(email string) (*model.User, error) { return nil, nil }
func (m *MockUserRepoNotFound) GetByID(id string) (*model.User, error)       { return nil, nil }
func (m *MockUserRepoNotFound) GetAll() ([]model.User, error)                { return nil, nil }
func (m *MockUserRepoNotFound) Update(id string, user *model.User) error     { return nil }
func (m *MockUserRepoNotFound) Delete(id string) error                       { return nil }
func (m *MockUserRepoNotFound) UpdateRole(id string, roleID string) error    { return nil }

// PASSWORD SALAH
type MockUserRepoWrongPassword struct{}

func (m *MockUserRepoWrongPassword) GetByUsername(username string) (*model.User, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password-benar"), bcrypt.DefaultCost)
	return &model.User{
		ID:           "user-1",
		Username:     "user",
		PasswordHash: string(hashed),
		RoleID:       "role-1",
	}, nil
}
func (m *MockUserRepoWrongPassword) Create(user *model.User) error               { return nil }
func (m *MockUserRepoWrongPassword) GetByEmail(email string) (*model.User, error) { return nil, nil }
func (m *MockUserRepoWrongPassword) GetByID(id string) (*model.User, error)       { return nil, nil }
func (m *MockUserRepoWrongPassword) GetAll() ([]model.User, error)                { return nil, nil }
func (m *MockUserRepoWrongPassword) Update(id string, user *model.User) error     { return nil }
func (m *MockUserRepoWrongPassword) Delete(id string) error                       { return nil }
func (m *MockUserRepoWrongPassword) UpdateRole(id string, roleID string) error    { return nil }

// LOGIN SUKSES
type MockUserRepoSuccess struct{}

func (m *MockUserRepoSuccess) GetByUsername(username string) (*model.User, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	return &model.User{
		ID:           "user-1",
		Username:     "admin",
		FullName:     "Admin Test",
		PasswordHash: string(hashed),
		RoleID:       "role-1",
	}, nil
}
func (m *MockUserRepoSuccess) Create(user *model.User) error               { return nil }
func (m *MockUserRepoSuccess) GetByEmail(email string) (*model.User, error) { return nil, nil }
func (m *MockUserRepoSuccess) GetByID(id string) (*model.User, error)       { return nil, nil }
func (m *MockUserRepoSuccess) GetAll() ([]model.User, error)                { return nil, nil }
func (m *MockUserRepoSuccess) Update(id string, user *model.User) error     { return nil }
func (m *MockUserRepoSuccess) Delete(id string) error                       { return nil }
func (m *MockUserRepoSuccess) UpdateRole(id string, roleID string) error    { return nil }

//
// ======================================================
// MOCK ROLE / STUDENT / LECTURER
// ======================================================
//

type MockRoleRepo struct{}

func (m *MockRoleRepo) GetByID(id string) (*model.Role, error) {
	return &model.Role{ID: "role-1", Name: "Admin"}, nil
}
func (m *MockRoleRepo) GetByName(name string) (*model.Role, error) { return nil, nil }
func (m *MockRoleRepo) GetAll() ([]model.Role, error)              { return []model.Role{}, nil }
func (m *MockRoleRepo) Create(role *model.Role) error              { return nil }

type MockStudentRepo struct{}

func (m *MockStudentRepo) GetByUserID(userID string) (*model.Student, error) {
	return nil, errors.New("student not found")
}
func (m *MockStudentRepo) GetStudentIDsByAdvisor(advisorID string) ([]string, error) {
	return []string{}, nil
}
func (m *MockStudentRepo) GetByID(studentID string) (*model.Student, error) { return nil, nil }
func (m *MockStudentRepo) Create(student *model.Student) error               { return nil }
func (m *MockStudentRepo) UpdateAdvisor(studentID, lecturerID string) error  { return nil }
func (m *MockStudentRepo) GetAll() ([]*model.Student, error)                  { return []*model.Student{}, nil }
func (m *MockStudentRepo) GetStudentAchievements(studentID string) ([]*model.AchievementReference, error) {
	return []*model.AchievementReference{}, nil
}
func (m *MockStudentRepo) CountAll() (int, error) { return 0, nil }
func (m *MockStudentRepo) IsStudentOfAdvisor(studentID, lecturerUserID string) (bool, error) {
	return false, nil
}

type MockLecturerRepo struct{}

func (m *MockLecturerRepo) Create(l *model.Lecturer) error { return nil }
func (m *MockLecturerRepo) GetByUserID(userID string) (*model.Lecturer, error) {
	return nil, errors.New("lecturer not found")
}
func (m *MockLecturerRepo) GetByLecturerID(lecturerID string) (*model.Lecturer, error) {
	return nil, nil
}
func (m *MockLecturerRepo) List() ([]model.Lecturer, error)  { return []model.Lecturer{}, nil }
func (m *MockLecturerRepo) Detail(id string) (*model.Lecturer, error) {
	return nil, nil
}

//
// ======================================================
// TEST CASES
// ======================================================
//

func TestAuthLogin_UserNotFound(t *testing.T) {
	app := fiber.New()
	auth := &AuthService{
		UserRepo:     &MockUserRepoNotFound{},
		RoleRepo:     &MockRoleRepo{},
		StudentRepo:  &MockStudentRepo{},
		LecturerRepo: &MockLecturerRepo{},
	}

	app.Post("/login", auth.Login)

	body, _ := json.Marshal(map[string]string{
		"username": "salah",
		"password": "123",
	})

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestAuthLogin_WrongPassword(t *testing.T) {
	app := fiber.New()
	auth := &AuthService{
		UserRepo:     &MockUserRepoWrongPassword{},
		RoleRepo:     &MockRoleRepo{},
		StudentRepo:  &MockStudentRepo{},
		LecturerRepo: &MockLecturerRepo{},
	}

	app.Post("/login", auth.Login)

	body, _ := json.Marshal(map[string]string{
		"username": "user",
		"password": "salah",
	})

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestAuthLogin_Success(t *testing.T) {
	app := fiber.New()
	auth := &AuthService{
		UserRepo:     &MockUserRepoSuccess{},
		RoleRepo:     &MockRoleRepo{},
		StudentRepo:  &MockStudentRepo{},
		LecturerRepo: &MockLecturerRepo{},
	}

	app.Post("/login", auth.Login)

	body, _ := json.Marshal(map[string]string{
		"username": "admin",
		"password": "123",
	})

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthRefresh_TokenMissing(t *testing.T) {
	app := fiber.New()
	auth := &AuthService{}

	app.Post("/refresh", auth.Refresh)

	req := httptest.NewRequest("POST", "/refresh", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthRefresh_InvalidToken(t *testing.T) {
	app := fiber.New()
	auth := &AuthService{}

	app.Post("/refresh", auth.Refresh)

	req := httptest.NewRequest("POST", "/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "token-salah"})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthLogout_Success(t *testing.T) {
	app := fiber.New()
	auth := &AuthService{}

	app.Post("/logout", auth.Logout)

	req := httptest.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer token-palsu")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthProfile_Unauthorized(t *testing.T) {
	app := fiber.New()
	auth := &AuthService{
		UserRepo: &MockUserRepoNotFound{},
	}

	app.Get("/profile", auth.Profile)

	req := httptest.NewRequest("GET", "/profile", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
func TestAuthRefresh_ValidToken(t *testing.T) {
	app := fiber.New()

	authService := &AuthService{
		UserRepo:     &MockUserRepoSuccess{},
		RoleRepo:     &MockRoleRepo{},
		StudentRepo:  &MockStudentRepo{},
		LecturerRepo: &MockLecturerRepo{},
	}

	app.Post("/refresh", authService.Refresh)

	// token palsu tapi format JWT valid (helper.ParseToken yg handle)
	validToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user-1",
		"role_id": "role-1",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte("SECRET_KEY"))

	req := httptest.NewRequest("POST", "/refresh", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: validToken,
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
func TestAuthProfile_Authorized(t *testing.T) {
	app := fiber.New()

	authService := &AuthService{
		UserRepo: &MockUserRepoSuccess{},
	}

	// middleware palsu isi user_id
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "user-1")
		return c.Next()
	})

	app.Get("/profile", authService.Profile)

	req := httptest.NewRequest("GET", "/profile", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
type MockRoleRepoMahasiswa struct{}

func (m *MockRoleRepoMahasiswa) GetByID(id string) (*model.Role, error) {
	return &model.Role{
		ID:   "role-mhs",
		Name: "Mahasiswa",
	}, nil
}

func (m *MockRoleRepoMahasiswa) GetByName(name string) (*model.Role, error) {
	return &model.Role{
		ID:   "role-mhs",
		Name: "Mahasiswa",
	}, nil
}

func (m *MockRoleRepoMahasiswa) GetAll() ([]model.Role, error) {
	return []model.Role{
		{
			ID:   "role-mhs",
			Name: "Mahasiswa",
		},
	}, nil
}

// func TestAuthLogin_Mahasiswa(t *testing.T) {
// 	app := fiber.New()

// 	authService := &AuthService{
// 		UserRepo:     &MockUserRepoSuccess{},
// 		RoleRepo:     &MockRoleRepoMahasiswa{},
// 		StudentRepo:  &MockStudentRepo{},
// 		LecturerRepo: &MockLecturerRepo{},
// 	}

// 	app.Post("/login", authService.Login)

// 	body := map[string]string{
// 		"username": "admin",
// 		"password": "123",
// 	}

// 	jsonBody, _ := json.Marshal(body)
// 	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := app.Test(req)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 200, resp.StatusCode)
// }
type MockStudentRepoSuccess struct{}

func (m *MockStudentRepoSuccess) GetByUserID(userID string) (*model.Student, error) {
	return &model.Student{
		ID:     "student-1",
		UserID: userID,
	}, nil
}

func (m *MockStudentRepoSuccess) GetStudentIDsByAdvisor(advisorID string) ([]string, error) {
	return []string{}, nil
}
func (m *MockStudentRepoSuccess) GetByID(studentID string) (*model.Student, error) {
	return nil, nil
}
func (m *MockStudentRepoSuccess) Create(student *model.Student) error { return nil }
func (m *MockStudentRepoSuccess) UpdateAdvisor(studentID, lecturerID string) error {
	return nil
}
func (m *MockStudentRepoSuccess) GetAll() ([]*model.Student, error) {
	return []*model.Student{}, nil
}
func (m *MockStudentRepoSuccess) GetStudentAchievements(studentID string) ([]*model.AchievementReference, error) {
	return []*model.AchievementReference{}, nil
}
func (m *MockStudentRepoSuccess) CountAll() (int, error) { return 0, nil }
func (m *MockStudentRepoSuccess) IsStudentOfAdvisor(studentID, lecturerUserID string) (bool, error) {
	return false, nil
}
func TestAuthLogin_Mahasiswa(t *testing.T) {
	app := fiber.New()

	authService := &AuthService{
		UserRepo:     &MockUserRepoSuccess{},
		RoleRepo:     &MockRoleRepoMahasiswa{},
		StudentRepo:  &MockStudentRepoSuccess{}, // 🔥 GANTI INI
		LecturerRepo: &MockLecturerRepo{},
	}

	app.Post("/login", authService.Login)

	body := map[string]string{
		"username": "admin",
		"password": "123",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
