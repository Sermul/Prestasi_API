package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"prestasi_api/app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)


// MOCK REPOSITORIES
type MockReportMongoRepo struct {
	MockAchievementMongoRepo
}


func (m *MockReportMongoRepo) GetAllForReport() ([]model.AchievementMongo, error) {
	return []model.AchievementMongo{
		{
			StudentID:       "student-1",
			AchievementType: "competition",
			Points:          10,
			CreatedAt:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Details: model.AchievementDetails{
				CompetitionLevel: "national",
			},
		},
		{
			StudentID:       "student-2",
			AchievementType: "seminar",
			Points:          5,
			CreatedAt:       time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}, nil
}


type MockReportStudentRepo struct{}

func (m *MockReportStudentRepo) GetStudentIDsByAdvisor(advisorID string) ([]string, error) {
	if advisorID == "lect-1" {
		return []string{"student-1"}, nil
	}
	return []string{}, nil
}

/* ===== dummy  ===== */
func (m *MockReportStudentRepo) GetByID(string) (*model.Student, error) { return nil, nil }
func (m *MockReportStudentRepo) GetByUserID(string) (*model.Student, error) {
	return nil, nil
}
func (m *MockReportStudentRepo) Create(*model.Student) error { return nil }
func (m *MockReportStudentRepo) UpdateAdvisor(string, string) error {
	return nil
}
func (m *MockReportStudentRepo) GetAll() ([]*model.Student, error) {
	return nil, nil
}
func (m *MockReportStudentRepo) GetStudentAchievements(string) ([]*model.AchievementReference, error) {
	return nil, nil
}
func (m *MockReportStudentRepo) CountAll() (int, error) { return 0, nil }
func (m *MockReportStudentRepo) IsStudentOfAdvisor(string, string) (bool, error) {
	return false, nil
}


// SETUP
func setupReportService() (*ReportService, *fiber.App) {
	app := fiber.New()

	svc := &ReportService{
		MongoRepo:   &MockReportMongoRepo{},
		StudentRepo: &MockReportStudentRepo{},
	}

	return svc, app
}

// TEST STATISTICS

func TestReport_Statistics_Admin(t *testing.T) {
	svc, app := setupReportService()

	app.Get("/report/statistics", func(c *fiber.Ctx) error {
		c.Locals("role", "Admin")
		return svc.Statistics(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/report/statistics", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReport_Statistics_DosenWali(t *testing.T) {
	svc, app := setupReportService()

	app.Get("/report/statistics", func(c *fiber.Ctx) error {
		c.Locals("role", "Dosen Wali")
		c.Locals("lecturer_id", "lect-1")
		return svc.Statistics(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/report/statistics", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReport_Statistics_Mahasiswa(t *testing.T) {
	svc, app := setupReportService()

	app.Get("/report/statistics", func(c *fiber.Ctx) error {
		c.Locals("role", "Mahasiswa")
		c.Locals("student_id", "student-1")
		return svc.Statistics(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/report/statistics", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReport_Statistics_Forbidden(t *testing.T) {
	svc, app := setupReportService()

	app.Get("/report/statistics", func(c *fiber.Ctx) error {
		c.Locals("role", "Guest")
		return svc.Statistics(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/report/statistics", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}

// =====================
// TEST STUDENT REPORT
// =====================

func TestReport_Student_Admin(t *testing.T) {
	svc, app := setupReportService()

	app.Get("/report/student/:id", func(c *fiber.Ctx) error {
		c.Locals("role", "Admin")
		return svc.StudentReport(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/report/student/student-1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReport_Student_Mahasiswa_Self(t *testing.T) {
	svc, app := setupReportService()

	app.Get("/report/student/:id", func(c *fiber.Ctx) error {
		c.Locals("role", "Mahasiswa")
		c.Locals("student_id", "student-1")
		return svc.StudentReport(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/report/student/student-1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestReport_Student_Mahasiswa_Forbidden(t *testing.T) {
	svc, app := setupReportService()

	app.Get("/report/student/:id", func(c *fiber.Ctx) error {
		c.Locals("role", "Mahasiswa")
		c.Locals("student_id", "student-1")
		return svc.StudentReport(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/report/student/student-2", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}

func TestReport_Student_DosenWali_Success(t *testing.T) {
	svc, app := setupReportService()

	app.Get("/report/student/:id", func(c *fiber.Ctx) error {
		c.Locals("role", "Dosen Wali")
		c.Locals("lecturer_id", "lect-1")
		return svc.StudentReport(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/report/student/student-1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestReport_Student_DosenWali_Forbidden(t *testing.T) {
	svc, app := setupReportService()

	app.Get("/report/student/:id", func(c *fiber.Ctx) error {
		c.Locals("role", "Dosen Wali")
		c.Locals("lecturer_id", "lect-1")
		return svc.StudentReport(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/report/student/student-2", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}
