package route

import (
	"prestasi_api/service"

	"github.com/gofiber/fiber/v2"
)

var AchievementRouter = func(app *fiber.App, svc *service.AchievementService) {

	api := app.Group("/api/v1/achievements")

	api.Post("/", svc.Create)               // FR-003
	api.Post("/:refId/submit", svc.Submit)  // FR-004
	api.Delete("/:refId", svc.Delete)       // FR-005
}
