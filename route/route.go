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

// ACHIEVEMENT ROUTER (Mahasiswa)
func AchievementRouter(app *fiber.App, svc *service.AchievementService) {
    api := app.Group("/api/v1/achievements",
        middleware.JWTMiddleware(),
    )

    // Mahasiswa only
    api.Post("/", middleware.RoleGuard("Mahasiswa"), svc.Create)
    api.Put("/:id", middleware.RoleGuard("Mahasiswa"), svc.Update)
    api.Delete("/:id", middleware.RoleGuard("Mahasiswa"), svc.Delete)
    api.Post("/:id/submit", middleware.RoleGuard("Mahasiswa"), svc.Submit)
    api.Post("/:id/attachments", middleware.RoleGuard("Mahasiswa"), svc.UploadAttachment)

    // DOSEN WALI
    api.Post("/:id/verify", middleware.RoleGuard("Dosen Wali"), svc.Verify)
    api.Post("/:id/reject", middleware.RoleGuard("Dosen Wali"), svc.Reject)

    // Semua role bisa lihat
    api.Get("/", svc.ListOwn) // atau ListAll di modul
    api.Get("/:id", svc.Detail)
    api.Get("/:id/history", svc.History)
}


// Dosen Wali
func LecturerRouter(app *fiber.App, svc *service.LecturerService) {
    api := app.Group("/api/v1/lecturers",
        middleware.JWTMiddleware(),
        middleware.RoleGuard("Admin", "Dosen Wali"),
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

func StudentRouter(app *fiber.App, svc *service.StudentService) {
    api := app.Group("/api/v1/students",
        middleware.JWTMiddleware(),
        middleware.RoleGuard("Admin"), // Admin yang assign doswal
    )

    api.Put("/:id/advisor", svc.AssignAdvisor)
}



// REPORT ROUTER
func ReportRouter(app *fiber.App, svc *service.ReportService) {
	api := app.Group("/api/v1/reports",
		middleware.JWTMiddleware(),
	)

	api.Get("/statistics", svc.Statistics)
	api.Get("/student/:id", svc.StudentReport)
}
