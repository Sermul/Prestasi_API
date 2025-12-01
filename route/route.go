package route

import (
	"prestasi_api/app/service"
	"prestasi_api/middleware"

	"github.com/gofiber/fiber/v2"
)
// AUTH ROUTER
func AuthRouter(app *fiber.App, svc *service.AuthService) {
	api := app.Group("/api/v1/auth")
	api.Post("/register", svc.Register)
	api.Post("/login", svc.Login)
	// Tambahan sesuai SRS
	api.Post("/refresh", svc.Refresh)
	api.Post("/logout", middleware.JWTMiddleware(), svc.Logout)
	api.Get("/profile", middleware.JWTMiddleware(), svc.Profile)
}
// ACHIEVEMENT ROUTER (Mahasiswa)
func AchievementRouter(app *fiber.App, svc *service.AchievementService) {
	api := app.Group("/api/v1/achievements",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("mahasiswa"),)
	api.Post("/", svc.Create)               
	api.Post("/:refId/submit", svc.Submit)  
	api.Delete("/:refId", svc.Delete)       
	// Mahasiswa
	api.Get("/", svc.ListOwn)                     // list prestasi mahasiswa
	api.Get("/:refId", svc.Detail)                // detail prestasi
	api.Put("/:refId", svc.Update)                // update draft
	api.Post("/:refId/attachments", svc.UploadAttachment) // upload file
	api.Get("/:refId/history", svc.History)       // riwayat status
	// ACHIEVEMENT ROUTER (Dosen Wali)
	advisor := app.Group("/api/v1/advisor/achievements",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("dosen_wali"),)
	advisor.Get("/", svc.GetAdviseeAchievements)
	advisor.Post("/:refId/verify", svc.Verify)
	advisor.Post("/:refId/reject", svc.Reject)
}
// LECTURER ROUTER (Dosen Wali)
func LecturerRouter(app *fiber.App, svc *service.LecturerService) {api := app.Group("/api/v1/lecturers",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("dosen_wali"),
	)
	// Tambahan sesuai SRS
	api.Get("/:id/advisees", svc.ListAdvisees)
}
// USER ROUTER (Admin)
func UserRouter(app *fiber.App, svc *service.UserService) {
	api := app.Group("/api/v1/users",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("admin"),
	)
	// CRUD Users 
	api.Get("/", svc.List)
	api.Get("/:id", svc.Detail)
	api.Post("/", svc.Create)
	api.Put("/:id", svc.Update)
	api.Delete("/:id", svc.Delete)
	api.Put("/:id/role", svc.ChangeRole)
}
// ADMIN 
func AdminAchievementRouter(app *fiber.App, svc *service.AchievementService) {
	api := app.Group("/api/v1/admin/achievements",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("admin"),
	)
	// List semua prestasi 
	api.Get("/", svc.AdminList)}
// REPORT 
func ReportRouter(app *fiber.App, svc *service.ReportService) {
	api := app.Group("/api/v1/reports",
		middleware.JWTMiddleware(),
	)

	api.Get("/statistics", svc.Statistics)
	api.Get("/student/:id", svc.StudentReport)
}
