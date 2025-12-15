package service

import (
	// "bytes"
	// "encoding/json"
	// "errors"
	"net/http/httptest"
	"testing"
	"strings"
	"prestasi_api/app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

//
// ======================================================
// MOCK STUDENT REPOSITORY
// ======================================================
//
type MockStudentRepoStudent struct{}


func (m *MockStudentRepoStudent) UpdateAdvisor(studentID, lecturerID string) error {
	return nil
}

func (m *MockStudentRepoStudent) GetByID(id string) (*model.Student, error) {
	return &model.Student{
		ID:        id,
		AdvisorID: "lect-1",
	}, nil
}

func (m *MockStudentRepoStudent) GetAll() ([]*model.Student, error) {
	return []*model.Student{
		{ID: "s1"},
		{ID: "s2"},
	}, nil
}

func (m *MockStudentRepoStudent) GetStudentIDsByAdvisor(advisorID string) ([]string, error) {
	return []string{"s1"}, nil
}

func (m *MockStudentRepoStudent) GetStudentAchievements(studentID string) ([]*model.AchievementReference, error) {
	return []*model.AchievementReference{}, nil
}

func (m *MockStudentRepoStudent) Create(*model.Student) error { return nil }
func (m *MockStudentRepoStudent) CountAll() (int, error)      { return 0, nil }
func (m *MockStudentRepoStudent) IsStudentOfAdvisor(string, string) (bool, error) {
	return true, nil
}
func (m *MockStudentRepoStudent) GetByUserID(userID string) (*model.Student, error) {
	return nil, nil
}
//
// ======================================================
// MOCK LECTURER REPOSITORY
// ======================================================
//
type MockLecturerRepoStudent struct{}

func (m *MockLecturerRepoStudent) GetByLecturerID(id string) (*model.Lecturer, error) {
	return &model.Lecturer{
		ID: id,
	}, nil
}

func (m *MockLecturerRepoStudent) Detail(id string) (*model.Lecturer, error) {
	return &model.Lecturer{
		ID: id,
	}, nil
}

func (m *MockLecturerRepoStudent) Create(*model.Lecturer) error { return nil }
func (m *MockLecturerRepoStudent) GetByUserID(string) (*model.Lecturer, error) {
	return nil, nil
}
func (m *MockLecturerRepoStudent) List() ([]model.Lecturer, error) {
	return []model.Lecturer{}, nil
}


//
// ======================================================
// TEST: AssignAdvisor
// ======================================================
//
func TestStudent_AssignAdvisor_Success(t *testing.T) {
	app := fiber.New()

	service := &StudentService{
		StudentRepo:  &MockStudentRepoStudent{},
		LecturerRepo: &MockLecturerRepoStudent{},
	}

	app.Put("/students/:id/advisor", service.AssignAdvisor)

	body := `{"lecturer_id":"lect-1"}`
	req := httptest.NewRequest("PUT", "/students/stu-1/advisor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

//
// ======================================================
// TEST: List (ADMIN)
// ======================================================
//
func TestStudent_List_Admin(t *testing.T) {
	app := fiber.New()

	service := &StudentService{
		StudentRepo:  &MockStudentRepo{},
		LecturerRepo: &MockLecturerRepo{},
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "Admin")
		return c.Next()
	})

	app.Get("/students", service.List)

	req := httptest.NewRequest("GET", "/students", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

//
// ======================================================
// TEST: Detail (ADMIN)
// ======================================================
//
func TestStudent_Detail_Admin(t *testing.T) {
	app := fiber.New()

	service := &StudentService{
		StudentRepo:  &MockStudentRepo{},
		LecturerRepo: &MockLecturerRepo{},
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "Admin")
		return c.Next()
	})

	app.Get("/students/:id", service.Detail)

	req := httptest.NewRequest("GET", "/students/student-1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

//
// ======================================================
// TEST: Achievements (ADMIN)
// ======================================================
//
func TestStudent_Achievements_Admin(t *testing.T) {
	app := fiber.New()

	service := &StudentService{
		StudentRepo:  &MockStudentRepo{},
		LecturerRepo: &MockLecturerRepo{},
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "Admin")
		return c.Next()
	})

	app.Get("/students/:id/achievements", service.Achievements)

	req := httptest.NewRequest("GET", "/students/student-1/achievements", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
