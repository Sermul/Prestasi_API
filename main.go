package main

import (
	"prestasi_api/app/repository"
	"prestasi_api/app/service"
	"prestasi_api/database"
	"prestasi_api/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Koneksi database
	database.ConnectPostgres()
	database.ConnectMongo()

	// Inisialisasi Fiber
	app := fiber.New()

	// Inisialisasi repository
	mongoRepo := repository.NewAchievementMongoRepository()
	postgresRepo := repository.NewAchievementPostgresRepository()
	studentRepo := repository.NewStudentPostgresRepository()

	// Inisialisasi service
	svc := &service.AchievementService{
		MongoRepo:    mongoRepo,
		PostgresRepo: postgresRepo,
		StudentRepo:  studentRepo,
	}

	// Register route
	route.AchievementRouter(app, svc)

	// Jalankan server
	println("Prestasi API siap jalan di port 3000...")
	app.Listen(":3000")
}
