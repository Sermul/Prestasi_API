package route

import (
	"prestasi_api/service"

	"github.com/gofiber/fiber/v2"
)

var AchievementRouter = func(app *fiber.App, svc *service.AchievementService) {

	api := app.Group("/api/v1/achievements")

	// FR-003 — Create Achievement
	api.Post("/", svc.Create)

	// FR-004 — Submit Achievement
	api.Post("/:refId/submit", svc.Submit)

	// FR-005 — Delete Achievement
	api.Delete("/:refId", svc.Delete)

	// FR-006 — Get Achievements for Advisee Students
	// (dosen wali akses daftar prestasi mahasiswa bimbingannya)
	api.Get("/advisor/list", svc.GetAdviseeAchievements)

	// FR-007 — Verify Achievement
	api.Post("/:refId/verify", svc.Verify)

	// FR-008 — Reject Achievement
	api.Post("/:refId/reject", svc.Reject)
}
