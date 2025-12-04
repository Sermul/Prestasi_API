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
	api.Post("/refresh", svc.Refresh)
	api.Post("/logout", middleware.JWTMiddleware(), svc.Logout)
	api.Get("/profile", middleware.JWTMiddleware(), svc.Profile)
}

// ACHIEVEMENT ROUTER (Mahasiswa)
func AchievementRouter(app *fiber.App, svc *service.AchievementService) {
	api := app.Group("/api/v1/achievements",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("Mahasiswa"),
	)

	api.Post("/", svc.Create)
	api.Post("/:refId/submit", svc.Submit)
	api.Delete("/:refId", svc.Delete)

	// Mahasiswa
	api.Get("/", svc.ListOwn)
	api.Get("/:refId", svc.Detail)
	api.Put("/:refId", svc.Update)
	api.Post("/:refId/attachments", svc.UploadAttachment)
	api.Get("/:refId/history", svc.History)

	// ACHIEVEMENT ROUTER (Dosen Wali)
	advisor := app.Group("/api/v1/advisor/achievements",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("Dosen Wali"),
	)
	advisor.Get("/", svc.GetAdviseeAchievements)
	advisor.Post("/:refId/verify", svc.Verify)
	advisor.Post("/:refId/reject", svc.Reject)
}

// LECTURER ROUTER (Dosen Wali)
func LecturerRouter(app *fiber.App, svc *service.LecturerService) {
	api := app.Group("/api/v1/lecturers",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("Dosen Wali"),
	)

	api.Get("/:id/advisees", svc.ListAdvisees)
}

// USER ROUTER (Admin)
func UserRouter(app *fiber.App, svc *service.UserService) {
	api := app.Group("/api/v1/users",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("Admin"),
	)

	api.Get("/", svc.List)
	api.Get("/:id", svc.Detail)
	api.Post("/", svc.Create)
	api.Put("/:id", svc.Update)
	api.Delete("/:id", svc.Delete)
	api.Put("/:id/role", svc.ChangeRole)
}

// ADMIN ACHIEVEMENT ROUTER
func AdminAchievementRouter(app *fiber.App, svc *service.AchievementService) {
	api := app.Group("/api/v1/admin/achievements",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("Admin"),
	)

	api.Get("/", svc.AdminList)
}

// REPORT ROUTER
func ReportRouter(app *fiber.App, svc *service.ReportService) {
	api := app.Group("/api/v1/reports",
		middleware.JWTMiddleware(),
	)

	api.Get("/statistics", svc.Statistics)
	api.Get("/student/:id", svc.StudentReport)
}
