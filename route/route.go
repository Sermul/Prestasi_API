package route

import (
	"prestasi_api/app/service"
	"github.com/gofiber/fiber/v2"
)

func AchievementRouter(app *fiber.App, svc *service.AchievementService) {

	api := app.Group("/api/v1/achievements")

	api.Post("/", svc.Create)
	api.Post("/:refId/submit", svc.Submit)
	api.Delete("/:refId", svc.Delete)

	advisor := app.Group("/api/v1/advisor/achievements")
	advisor.Get("/", svc.GetAdviseeAchievementsHandler)
	advisor.Post("/:refId/verify", svc.VerifyHandler)
	advisor.Post("/:refId/reject", svc.RejectHandler)
}
