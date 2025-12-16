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
)



// -------- Student Repo --------
type MockStudentRepoSS struct{}

func (m *MockStudentRepoSS) GetStudentIDsByAdvisor(advisorID string) ([]string, error) {
	return []string{"student-1"}, nil
}

func (m *MockStudentRepoSS) GetByID(id string) (*model.Student, error) {
	if id == "not-found" {
		return nil, assert.AnError
	}
	return &model.Student{
		ID:        id,
		AdvisorID: "lect-1",
	}, nil
}

func (m *MockStudentRepoSS) GetByUserID(userID string) (*model.Student, error) {
	return &model.Student{ID: "student-1", UserID: userID}, nil
}

func (m *MockStudentRepoSS) Create(*model.Student) error {
	return nil
}

func (m *MockStudentRepoSS) UpdateAdvisor(studentID, lecturerID string) error {
	return nil
}

func (m *MockStudentRepoSS) GetAll() ([]*model.Student, error) {
	return []*model.Student{
		{ID: "student-1", AdvisorID: "lect-1"},
	}, nil
}

func (m *MockStudentRepoSS) GetStudentAchievements(studentID string) ([]*model.AchievementReference, error) {
	return []*model.AchievementReference{}, nil
}

func (m *MockStudentRepoSS) CountAll() (int, error) {
	return 1, nil
}

func (m *MockStudentRepoSS) IsStudentOfAdvisor(studentID, lecturerUserID string) (bool, error) {
	return true, nil
}

// -------- Lecturer Repo --------
type MockLecturerRepoSS struct{}

func (m *MockLecturerRepoSS) Create(*model.Lecturer) error {
	return nil
}

func (m *MockLecturerRepoSS) Detail(id string) (*model.Lecturer, error) {
	return &model.Lecturer{ID: "lect-1"}, nil
}

func (m *MockLecturerRepoSS) GetByLecturerID(id string) (*model.Lecturer, error) {
	if id == "not-found" {
		return nil, assert.AnError
	}
	return &model.Lecturer{ID: "lect-1"}, nil
}

func (m *MockLecturerRepoSS) GetAll() ([]*model.Lecturer, error) {
	return []*model.Lecturer{}, nil
}
func (m *MockLecturerRepoSS) GetByUserID(userID string) (*model.Lecturer, error) {
	if userID == "not-found" {
		return nil, assert.AnError
	}
	return &model.Lecturer{
		ID:     "lect-1",
		UserID: userID,
	}, nil
}
func (m *MockLecturerRepoSS) List() ([]model.Lecturer, error) {
	return []model.Lecturer{}, nil
}



// ======================= SETUP =======================

func setupStudentService() (*StudentService, *fiber.App) {
	app := fiber.New()

	svc := &StudentService{
		StudentRepo:  &MockStudentRepoSS{},
		LecturerRepo: &MockLecturerRepoSS{},
	}

	return svc, app
}



// ======================= TESTS =======================

// ---------- ASSIGN ADVISOR ----------
func TestAssignAdvisor_Success(t *testing.T) {
	svc, app := setupStudentService()

	app.Post("/students/:id/advisor", func(c *fiber.Ctx) error {
		return svc.AssignAdvisor(c)
	})

	body, _ := json.Marshal(map[string]string{
		"lecturer_id": "lect-1",
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/students/student-1/advisor",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAssignAdvisor_LecturerNotFound(t *testing.T) {
	svc, app := setupStudentService()

	app.Post("/students/:id/advisor", func(c *fiber.Ctx) error {
		return svc.AssignAdvisor(c)
	})

	body, _ := json.Marshal(map[string]string{
		"lecturer_id": "not-found",
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/students/student-1/advisor",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 404, resp.StatusCode)
}

// ---------- LIST ----------
func TestStudentList_Admin(t *testing.T) {
	svc, app := setupStudentService()

	app.Get("/students", func(c *fiber.Ctx) error {
		c.Locals("role", "Admin")
		return svc.List(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/students", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestStudentList_DosenWali(t *testing.T) {
	svc, app := setupStudentService()

	app.Get("/students", func(c *fiber.Ctx) error {
		c.Locals("role", "Dosen Wali")
		c.Locals("lecturer_id", "lect-1")
		return svc.List(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/students", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// ---------- DETAIL ----------
func TestStudentDetail_Admin(t *testing.T) {
	svc, app := setupStudentService()

	app.Get("/students/:id", func(c *fiber.Ctx) error {
		c.Locals("role", "Admin")
		return svc.Detail(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/students/student-1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestStudentDetail_DosenForbidden(t *testing.T) {
	svc, app := setupStudentService()

	app.Get("/students/:id", func(c *fiber.Ctx) error {
		c.Locals("role", "Dosen Wali")
		c.Locals("lecturer_id", "lect-x")
		return svc.Detail(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/students/student-1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}

// ---------- ACHIEVEMENTS ----------
func TestStudentAchievements_Admin(t *testing.T) {
	svc, app := setupStudentService()

	app.Get("/students/:id/achievements", func(c *fiber.Ctx) error {
		c.Locals("role", "Admin")
		return svc.Achievements(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/students/student-1/achievements", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestStudentAchievements_StudentNotFound(t *testing.T) {
	svc, app := setupStudentService()

	app.Get("/students/:id/achievements", func(c *fiber.Ctx) error {
		c.Locals("role", "Admin")
		return svc.Achievements(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/students/not-found/achievements", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}
