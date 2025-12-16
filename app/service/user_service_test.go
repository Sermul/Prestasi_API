package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"prestasi_api/app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)


// MOCK USER REPOSITORY (KHUSUS USER SERVICE)
//
type MockUserRepoUser struct{}

func (m *MockUserRepoUser) GetAll() ([]model.User, error) {
	return []model.User{
		{ID: "u1", Username: "user1"},
	}, nil
}

func (m *MockUserRepoUser) GetByID(id string) (*model.User, error) {
	if id == "not-found" {
		return nil, assert.AnError
	}
	return &model.User{
		ID:       id,
		Username: "user",
		RoleID:   "role-1",
	}, nil
}

func (m *MockUserRepoUser) Create(user *model.User) error { return nil }
func (m *MockUserRepoUser) Update(id string, user *model.User) error {
	return nil
}
func (m *MockUserRepoUser) Delete(id string) error { return nil }
func (m *MockUserRepoUser) UpdateRole(id, roleID string) error {
	return nil
}
func (m *MockUserRepoUser) GetByUsername(string) (*model.User, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	return &model.User{
		ID:           "u1",
		Username:     "user",
		PasswordHash: string(hashed),
		RoleID:       "role-1",
	}, nil
}
func (m *MockUserRepoUser) GetByEmail(string) (*model.User, error) {
	return nil, nil
}


// MOCK ROLE REPOSITORY (KHUSUS USER SERVICE)
type MockRoleRepoUser struct{}

func (m *MockRoleRepoUser) GetByID(id string) (*model.Role, error) {
	if id == "invalid" {
		return nil, assert.AnError
	}
	return &model.Role{
		ID:   id,
		Name: "Dosen Wali",
	}, nil
}

func (m *MockRoleRepoUser) GetByName(string) (*model.Role, error) {
	return nil, nil
}
func (m *MockRoleRepoUser) GetAll() ([]model.Role, error) {
	return []model.Role{}, nil
}
func (m *MockRoleRepoUser) Create(*model.Role) error { return nil }


// MOCK STUDENT REPOSITORY (KHUSUS USER SERVICE)

type MockStudentRepoUser struct{}

func (m *MockStudentRepoUser) GetByUserID(userID string) (*model.Student, error) {
	return &model.Student{
		ID:     "student-1",
		UserID: userID,
	}, nil
}
func (m *MockStudentRepoUser) GetStudentIDsByAdvisor(string) ([]string, error) {
	return []string{}, nil
}
func (m *MockStudentRepoUser) GetByID(string) (*model.Student, error) {
	return nil, nil
}
func (m *MockStudentRepoUser) Create(*model.Student) error { return nil }
func (m *MockStudentRepoUser) UpdateAdvisor(string, string) error {
	return nil
}
func (m *MockStudentRepoUser) GetAll() ([]*model.Student, error) {
	return []*model.Student{}, nil
}
func (m *MockStudentRepoUser) GetStudentAchievements(string) ([]*model.AchievementReference, error) {
	return []*model.AchievementReference{}, nil
}
func (m *MockStudentRepoUser) CountAll() (int, error) { return 0, nil }
func (m *MockStudentRepoUser) IsStudentOfAdvisor(string, string) (bool, error) {
	return false, nil
}


// MOCK LECTURER REPOSITORY (KHUSUS USER SERVICE)

type MockLecturerRepoUser struct{}

func (m *MockLecturerRepoUser) Create(*model.Lecturer) error { return nil }
func (m *MockLecturerRepoUser) GetByUserID(string) (*model.Lecturer, error) {
	return nil, nil
}
func (m *MockLecturerRepoUser) GetByLecturerID(string) (*model.Lecturer, error) {
	return nil, nil
}
func (m *MockLecturerRepoUser) List() ([]model.Lecturer, error) {
	return []model.Lecturer{}, nil
}
func (m *MockLecturerRepoUser) Detail(string) (*model.Lecturer, error) {
	return nil, nil
}


// SETUP

func setupUserService() (*UserService, *fiber.App) {
	app := fiber.New()

	svc := &UserService{
		UserRepo:     &MockUserRepoUser{},
		RoleRepo:     &MockRoleRepoUser{},
		StudentRepo:  &MockStudentRepoUser{},
		LecturerRepo: &MockLecturerRepoUser{},
	}

	return svc, app
}


// TEST CASES

// ---------- LIST ----------
func TestUser_List_Success(t *testing.T) {
	svc, app := setupUserService()

	app.Get("/users", svc.List)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// ---------- DETAIL ----------
func TestUser_Detail_Success(t *testing.T) {
	svc, app := setupUserService()

	app.Get("/users/:id", svc.Detail)

	req := httptest.NewRequest(http.MethodGet, "/users/u1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestUser_Detail_NotFound(t *testing.T) {
	svc, app := setupUserService()

	app.Get("/users/:id", svc.Detail)

	req := httptest.NewRequest(http.MethodGet, "/users/not-found", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

// ---------- CREATE ----------
func TestUser_Create_Success(t *testing.T) {
	svc, app := setupUserService()

	app.Post("/users", svc.Create)

	body, _ := json.Marshal(map[string]string{
		"username":  "newuser",
		"email":     "test@mail.com",
		"password":  "123",
		"full_name": "Test User",
		"role_id":   "role-1",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

// ---------- UPDATE ----------
func TestUser_Update_Success(t *testing.T) {
	svc, app := setupUserService()

	app.Put("/users/:id", svc.Update)

	body, _ := json.Marshal(map[string]interface{}{
		"full_name": "Updated Name",
		"is_active": true,
	})

	req := httptest.NewRequest(http.MethodPut, "/users/u1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

// ---------- DELETE ----------
func TestUser_Delete_Success(t *testing.T) {
	svc, app := setupUserService()

	app.Delete("/users/:id", svc.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/u1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// ---------- CHANGE ROLE ----------
func TestUser_ChangeRole_DosenWali(t *testing.T) {
	svc, app := setupUserService()

	app.Put("/users/:id/role", svc.ChangeRole)

	body, _ := json.Marshal(map[string]string{
		"role_id": "role-dosen",
	})

	req := httptest.NewRequest(http.MethodPut, "/users/u1/role", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}
