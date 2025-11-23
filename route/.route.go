package route

import (
	"app/service"
	"github.com/gofiber/fiber/v2"
)

func AchievementRoute(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Post("/achievements", service.CreateAchievementService)
	api.Post("/achievements/:refId/submit", service.SubmitAchievementService)
	api.Delete("/achievements/:refId", service.DeleteAchievementService)
}
