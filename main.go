package main

import (
	"prestasi_api/app/repository"
	"prestasi_api/app/service"
	"prestasi_api/database"
	"prestasi_api/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Connect database
	database.ConnectPostgres()
	database.ConnectMongo()

	app := fiber.New()

	// ===== REPOSITORY =====
	achievementMongoRepo := repository.NewAchievementMongoRepository()
	achievementPostgresRepo := repository.NewAchievementPostgresRepository()
	studentRepo := repository.NewStudentPostgresRepository()
	userRepo := repository.NewUserPostgresRepository()
	roleRepo := repository.NewRolePostgresRepository()
	permissionRepo := repository.NewPermissionPostgresRepository()
	rolePermissionRepo := repository.NewRolePermissionPostgresRepository()

	// ===== SERVICE =====
	authSvc := service.NewAuthService(userRepo, roleRepo, permissionRepo, rolePermissionRepo)
	achievementSvc := &service.AchievementService{
		MongoRepo:    achievementMongoRepo,
		PostgresRepo: achievementPostgresRepo,
		StudentRepo:  studentRepo,
	}

	// ===== ROUTES =====
	route.AuthRouter(app, authSvc)
	route.AchievementRouter(app, achievementSvc)

	println("Prestasi API berjalan di port 3000...")
	app.Listen(":3000")
}
