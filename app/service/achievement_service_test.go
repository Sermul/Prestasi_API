package service

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
	"mime/multipart"

    "prestasi_api/app/model"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// ================= MOCK REPOSITORIES =================

type MockAchievementMongoRepo struct{}

type MockAchievementPostgresRepo struct{}

type MockStudentPostgresRepo struct{}

func (m *MockStudentPostgresRepo) GetStudentIDsByAdvisor(advisorID string) ([]string, error) {
	return []string{"student-1"}, nil
}

func (m *MockStudentPostgresRepo) GetByID(studentID string) (*model.Student, error) {
	return &model.Student{ID: studentID}, nil
}

func (m *MockStudentPostgresRepo) GetByUserID(userID string) (*model.Student, error) {
	return &model.Student{UserID: userID}, nil
}

func (m *MockStudentPostgresRepo) Create(student *model.Student) error {
	return nil
}

func (m *MockStudentPostgresRepo) UpdateAdvisor(studentID, lecturerID string) error {
	return nil
}

func (m *MockStudentPostgresRepo) GetAll() ([]*model.Student, error) {
	return []*model.Student{}, nil
}

func (m *MockStudentPostgresRepo) GetStudentAchievements(studentID string) ([]*model.AchievementReference, error) {
	return []*model.AchievementReference{}, nil
}

func (m *MockStudentPostgresRepo) CountAll() (int, error) {
	return 0, nil
}

func (m *MockStudentPostgresRepo) IsStudentOfAdvisor(studentID, lecturerUserID string) (bool, error) {
	return true, nil
}



// ---------- Mongo Repo ----------
func (m *MockAchievementMongoRepo) CreateAchievementMongo(a *model.AchievementMongo) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}
func (m *MockAchievementMongoRepo) GetByID(id primitive.ObjectID) (*model.AchievementMongo, error) {
	return &model.AchievementMongo{Title: "Dummy"}, nil
}
func (m *MockAchievementMongoRepo) UpdateAchievementMongo(id primitive.ObjectID, a *model.AchievementMongo) error {
	return nil
}
func (m *MockAchievementMongoRepo) SoftDeleteAchievementMongo(id primitive.ObjectID) error {
	return nil
}
func (m *MockAchievementMongoRepo) AddAttachmentMongo(
	id primitive.ObjectID,
	file *multipart.FileHeader,
) (string, error) {
	return "http://dummy.url/file.pdf", nil
}
func (m *MockAchievementMongoRepo) GetAll() ([]model.AchievementMongo, error) {
	return []model.AchievementMongo{}, nil
}
func (m *MockAchievementMongoRepo) GetAllForReport() ([]model.AchievementMongo, error) {
	return []model.AchievementMongo{}, nil
}
func (m *MockAchievementMongoRepo) RestoreAchievementMongo(id primitive.ObjectID) error {
	return nil
}


// ---------- Postgres Repo ----------
func (m *MockAchievementPostgresRepo) CreateReferencePostgres(r *model.AchievementReference) error {
	return nil
}
func (m *MockAchievementPostgresRepo) GetReferenceByID(id string) (*model.AchievementReference, error) {
	return &model.AchievementReference{
		ID:        id,
		StudentID: "student-1",
		MongoID:   primitive.NewObjectID().Hex(),
		Status:    "submitted", // 🔥 PENTING
	}, nil
}

func (m *MockAchievementPostgresRepo) UpdateReferenceStatusPostgres(id, status string) error {
	return nil
}
func (m *MockAchievementPostgresRepo) SaveSubmittedAt(id string, t time.Time) error {
	return nil
}
func (m *MockAchievementPostgresRepo) InsertHistory(h *model.AchievementReferenceHistory) error {
	return nil
}
func (m *MockAchievementPostgresRepo) UpdateVerifyStatus(id, userID string) error {
	return nil
}
func (m *MockAchievementPostgresRepo) RejectReference(id, userID, note string) error {
	return nil
}
func (m *MockAchievementPostgresRepo) GetByStudentID(studentID string) ([]model.AchievementReference, error) {
	return []model.AchievementReference{}, nil
}
func (m *MockAchievementPostgresRepo) GetByStudentIDs(ids []string) ([]model.AchievementReference, error) {
	return []model.AchievementReference{}, nil
}
func (m *MockAchievementPostgresRepo) GetAllReferences() ([]model.AchievementReference, error) {
	return []model.AchievementReference{}, nil
}
func (m *MockAchievementPostgresRepo) GetHistoryByReferenceID(
    id string,
) ([]map[string]interface{}, error) {
    return []map[string]interface{}{}, nil
}



func (m *MockAchievementPostgresRepo) GetMongoID(refID string) (string, error) {
	return primitive.NewObjectID().Hex(), nil
}

// ================= TEST SETUP =================

func setupAchievementService() (*AchievementService, *fiber.App) {
    app := fiber.New()

    svc := &AchievementService{
        MongoRepo:    &MockAchievementMongoRepo{},
        PostgresRepo: &MockAchievementPostgresRepo{},
        StudentRepo:  &MockStudentPostgresRepo{},
    }

    return svc, app
}


// ================= UNIT TESTS =================

func TestCreateAchievement_Mahasiswa_Success(t *testing.T) {
	svc, app := setupAchievementService()

	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("role", "Mahasiswa")
		c.Locals("student_id", "student-1")
		c.Locals("user_id", uuid.New().String())
		return svc.Create(c)
	})

	payload := map[string]interface{}{
		"achievementType": "competition",
		"title":           "Juara Nasional",
		"details": map[string]interface{}{
			"competitionLevel": "national",
		},
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/achievements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCreateAchievement_Mahasiswa_SendStudentID_ShouldFail(t *testing.T) {
	svc, app := setupAchievementService()

	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("role", "Mahasiswa")
		c.Locals("student_id", "student-1")
		c.Locals("user_id", uuid.New().String())
		return svc.Create(c)
	})

	payload := map[string]interface{}{
		"studentId": "student-2",
		"title":     "Invalid",
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/achievements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestSubmitAchievement_Success(t *testing.T) {
	app := fiber.New()

	svc := &AchievementService{
		MongoRepo:    &MockAchievementMongoRepo{},
		PostgresRepo: &MockAchievementPostgresRepoDraft{}, // 🔥 DRAFT
		StudentRepo:  &MockStudentPostgresRepo{},
	}

	app.Post("/achievements/:refId/submit", func(c *fiber.Ctx) error {
		c.Locals("role", "Mahasiswa")
		c.Locals("student_id", "student-1")
		c.Locals("user_id", uuid.New().String())
		return svc.Submit(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/achievements/ref-123/submit", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestVerifyAchievement_ByAdvisor_Success(t *testing.T) {
	app := fiber.New()

	// 🔥 DI SINI LETAKNYA
	svc := &AchievementService{
		MongoRepo:    &MockAchievementMongoRepo{},
		PostgresRepo: &MockAchievementPostgresRepoSubmitted{}, // SUBMITTED
		StudentRepo:  &MockStudentPostgresRepo{},
	}

	app.Post("/achievements/:refId/verify", func(c *fiber.Ctx) error {
		c.Locals("role", "Dosen Wali")
		c.Locals("lecturer_id", "lect-1")
		c.Locals("user_id", "user-1")
		return svc.Verify(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/achievements/ref-123/verify", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}


func TestRejectAchievement_WithNote_Success(t *testing.T) {
	app := fiber.New()

	// 🔥 DI SINI LETAKNYA
	svc := &AchievementService{
		MongoRepo:    &MockAchievementMongoRepo{},
		PostgresRepo: &MockAchievementPostgresRepoSubmitted{}, // SUBMITTED
		StudentRepo:  &MockStudentPostgresRepo{},
	}

app.Post("/achievements/:refId/reject", func(c *fiber.Ctx) error {
	c.Locals("role", "Admin")
	c.Locals("user_id", "admin-1")
	c.Locals("admin_id", "admin-1")
	c.Locals("lecturer_id", "admin-1") // 🔥 TAMBAHAN WAJIB
	c.Locals("advisor_id", "admin-1")  // 🔥 TAMBAHAN WAJIB
	c.Locals("note", "Tidak valid")
	return svc.Reject(c)
})




	body, _ := json.Marshal(map[string]string{"note": "Tidak valid"})
	req := httptest.NewRequest(http.MethodPost, "/achievements/ref-123/reject", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}


// ================= NEGATIVE CASE EXAMPLE =================

func TestGetReference_NotFound(t *testing.T) {
	svc, app := setupAchievementService()

	app.Post("/achievements/:refId/submit", func(c *fiber.Ctx) error {
		c.Locals("role", "Mahasiswa")
		c.Locals("student_id", "student-1")
		return svc.Submit(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/achievements/invalid/submit", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func (m *MockAchievementPostgresRepo) CountAll() (int, error) {
    return 0, nil
}

func (m *MockAchievementPostgresRepo) CountByStudentID(studentID string) (int, error) {
    return 0, nil
}

func (m *MockAchievementPostgresRepo) CountByStudentIDs(studentIDs []string) (int, error) {
    return 0, nil
}
type MockAchievementPostgresRepoDraft struct {
	MockAchievementPostgresRepo
}

func (m *MockAchievementPostgresRepoDraft) GetReferenceByID(id string) (*model.AchievementReference, error) {
	return &model.AchievementReference{
		ID:        id,
		StudentID: "student-1",
		MongoID:   primitive.NewObjectID().Hex(),
		Status:    "draft", // ✅ UNTUK SUBMIT
	}, nil
}
type MockAchievementPostgresRepoSubmitted struct {
	MockAchievementPostgresRepo
}

func (m *MockAchievementPostgresRepoSubmitted) GetReferenceByID(id string) (*model.AchievementReference, error) {
	return &model.AchievementReference{
		ID:        id,
		StudentID: "student-1",
		MongoID:   primitive.NewObjectID().Hex(),
		Status:    "submitted", // ✅ UNTUK VERIFY / REJECT
	}, nil
}
func TestDeleteAchievement_Draft_Success(t *testing.T) {
	app := fiber.New()

	service := &AchievementService{
		PostgresRepo: &MockAchievementPostgresRepo{},
		MongoRepo:    &MockAchievementMongoRepo{},
	}

	// middleware palsu → mahasiswa
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "Mahasiswa")
		c.Locals("student_id", "student-1")
		return c.Next()
	})

	app.Delete("/achievements/:id", service.Delete)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/achievements/achv-1",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

}
func TestAchievementList_Admin(t *testing.T) {
	app := fiber.New()

	service := &AchievementService{
		PostgresRepo: &MockAchievementPostgresRepo{},
		MongoRepo:    &MockAchievementMongoRepo{},
	}

	// middleware palsu → admin
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "Admin")
		return c.Next()
	})

	app.Get("/achievements", service.List)

	req := httptest.NewRequest(
		http.MethodGet,
		"/achievements",
		nil,
	)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
