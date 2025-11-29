package route

import (
	"prestasi_api/app/service"
	"prestasi_api/middleware"

	"github.com/gofiber/fiber/v2"
)

// =============================
// AUTH ROUTER
// =============================
func AuthRouter(app *fiber.App, svc *service.AuthService) {
	api := app.Group("/api/v1/auth")

	api.Post("/register", svc.Register)
	api.Post("/login", svc.Login)
}

// =============================
// ACHIEVEMENT ROUTER
// =============================
func AchievementRouter(app *fiber.App, svc *service.AchievementService) {

	api := app.Group("/api/v1/achievements",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("mahasiswa"),
	)

	api.Post("/", svc.Create)               // mahasiswa create
	api.Post("/:refId/submit", svc.Submit)  // mahasiswa submit
	api.Delete("/:refId", svc.Delete)       // mahasiswa delete draft

	advisor := app.Group("/api/v1/advisor/achievements",
		middleware.JWTMiddleware(),
		middleware.RoleGuard("dosen_wali"),
	)

	advisor.Get("/", svc.GetAdviseeAchievements) // dosen wali melihat list
	advisor.Post("/:refId/verify", svc.Verify)   // dosen wali verifikasi
	advisor.Post("/:refId/reject", svc.Reject)   // dosen wali reject
}
