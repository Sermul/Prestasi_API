package route

import (
	"prestasi_api/app/service"
	"prestasi_api/middleware"
	"github.com/gofiber/fiber/v2"
)
// AUTH ROUTER
func AuthRouter(app *fiber.App, svc *service.AuthService) {
	api := app.Group("/api/v1/auth")
	api.Post("/login", svc.Login)
	api.Post("/refresh", svc.Refresh)
	api.Post("/logout", middleware.JWTMiddleware(), svc.Logout)
	api.Get("/profile", middleware.JWTMiddleware(), svc.Profile)
}
// ACHIEVEMENT (Mahasiswa)
func AchievementRouter(app *fiber.App, svc *service.AchievementService) {
    api := app.Group("/api/v1/achievements",
        middleware.JWTMiddleware(),
    )
    // Create => Mahasiswa & Admin
    api.Post("/", middleware.RoleGuard("Mahasiswa", "Admin"), svc.Create)
    // Update => Mahasiswa only
    // Update => Mahasiswa & Admin
api.Put("/:refId", middleware.RoleGuard("Mahasiswa", "Admin"), svc.Update)
    // Delete => Mahasiswa & Admin
    api.Delete("/:refId", middleware.RoleGuard("Mahasiswa", "Admin"), svc.Delete)
    // Submit => Mahasiswa & Admin
    api.Post("/:refId/submit", middleware.RoleGuard("Mahasiswa", "Admin"), svc.Submit)
    // Upload attachment => Mahasiswa & Admin
    api.Post("/:refId/attachments", middleware.RoleGuard("Mahasiswa", "Admin"), svc.UploadAttachment)
    // Verify / Reject => Dosen Wali & Admin
    api.Post("/:refId/verify", middleware.RoleGuard("Dosen Wali", "Admin"), svc.Verify)
    api.Post("/:refId/reject", middleware.RoleGuard("Dosen Wali", "Admin"), svc.Reject)
    // Everyone with token can read
    api.Get("/", svc.List)
    api.Get("/:refId", svc.Detail)
    api.Get("/:refId/history", svc.History)
}
// Dosen Wali
func LecturerRouter(app *fiber.App, svc *service.LecturerService) {
    api := app.Group("/api/v1/lecturers",
        middleware.JWTMiddleware(),
        middleware.RoleGuard("Admin"),
    )
    api.Get("/", svc.List)
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
//Students & Lecturers 
func StudentRouter(app *fiber.App, svc *service.StudentService) {
    api := app.Group("/api/v1/students",
        middleware.JWTMiddleware(),
        middleware.RoleGuard("Admin", "Dosen Wali"),
    )
    api.Get("/", svc.List)                  
    api.Get("/:id", svc.Detail)
    api.Get("/:id/achievements", svc.Achievements)
    api.Put("/:id/advisor", middleware.RoleGuard("Admin"), svc.AssignAdvisor)
}

