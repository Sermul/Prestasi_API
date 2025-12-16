package service

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"prestasi_api/app/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

/* =========================
   MOCK STUDENT REPO
========================= */

type MockStudentRepoLecturer struct{}

func (m *MockStudentRepoLecturer) GetStudentIDsByAdvisor(advisorID string) ([]string, error) {
	if advisorID == "empty" {
		return []string{}, nil
	}
	if advisorID == "error" {
		return nil, assert.AnError
	}
	return []string{"stu-1", "stu-2"}, nil
}

func (m *MockStudentRepoLecturer) GetByID(id string) (*model.Student, error) {
	return &model.Student{
		ID:        id,
		StudentID: "20201234",
	}, nil
}

/* unused but required */
func (m *MockStudentRepoLecturer) GetByUserID(string) (*model.Student, error) {
	return nil, nil
}
func (m *MockStudentRepoLecturer) Create(*model.Student) error {
	return nil
}
func (m *MockStudentRepoLecturer) UpdateAdvisor(string, string) error {
	return nil
}
func (m *MockStudentRepoLecturer) GetAll() ([]*model.Student, error) {
	return nil, nil
}
func (m *MockStudentRepoLecturer) GetStudentAchievements(string) ([]*model.AchievementReference, error) {
	return nil, nil
}
func (m *MockStudentRepoLecturer) CountAll() (int, error) {
	return 0, nil
}
func (m *MockStudentRepoLecturer) IsStudentOfAdvisor(string, string) (bool, error) {
	return false, nil
}

/* =========================
   MOCK LECTURER REPO
========================= */

type MockLecturerRepoLecturer struct{}

func (m *MockLecturerRepoLecturer) List() ([]model.Lecturer, error) {
	return []model.Lecturer{
		{ID: "lect-1", LecturerID: "L001"},
	}, nil
}

func (m *MockLecturerRepoLecturer) Detail(id string) (*model.Lecturer, error) {
	if id == "not-found" {
		return nil, assert.AnError
	}
	return &model.Lecturer{
		ID:         id,
		LecturerID: "L001",
	}, nil
}

/* unused but required */
func (m *MockLecturerRepoLecturer) Create(*model.Lecturer) error {
	return nil
}
func (m *MockLecturerRepoLecturer) GetByUserID(string) (*model.Lecturer, error) {
	return nil, nil
}
func (m *MockLecturerRepoLecturer) GetByLecturerID(string) (*model.Lecturer, error) {
	return nil, nil
}

/* =========================
   SETUP
========================= */

func setupLecturerService() (*fiber.App, *LecturerService) {
	app := fiber.New()

	svc := NewLecturerService(
		&MockStudentRepoLecturer{},
		&MockLecturerRepoLecturer{},
	)

	return app, svc
}

/* =========================
   TEST ListAdvisees
========================= */

func TestLecturer_ListAdvisees_Success(t *testing.T) {
	app, svc := setupLecturerService()

	app.Get("/lecturers/:id/students", svc.ListAdvisees)

	req := httptest.NewRequest("GET", "/lecturers/lect-1/students", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestLecturer_ListAdvisees_Empty(t *testing.T) {
	app, svc := setupLecturerService()

	app.Get("/lecturers/:id/students", svc.ListAdvisees)

	req := httptest.NewRequest("GET", "/lecturers/empty/students", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)

	var result []interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.Len(t, result, 0)
}

func TestLecturer_ListAdvisees_Error(t *testing.T) {
	app, svc := setupLecturerService()

	app.Get("/lecturers/:id/students", svc.ListAdvisees)

	req := httptest.NewRequest("GET", "/lecturers/error/students", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

/* =========================
   TEST List
========================= */

func TestLecturer_List_Success(t *testing.T) {
	app, svc := setupLecturerService()

	app.Get("/lecturers", svc.List)

	req := httptest.NewRequest("GET", "/lecturers", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

/* =========================
   TEST Detail
========================= */

func TestLecturer_Detail_Success(t *testing.T) {
	app, svc := setupLecturerService()

	app.Get("/lecturers/:id", svc.Detail)

	req := httptest.NewRequest("GET", "/lecturers/lect-1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestLecturer_Detail_NotFound(t *testing.T) {
	app, svc := setupLecturerService()

	app.Get("/lecturers/:id", svc.Detail)

	req := httptest.NewRequest("GET", "/lecturers/not-found", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}
