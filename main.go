package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"prestasi_api/app/repository"
	"prestasi_api/app/service"
	"prestasi_api/database"
	"prestasi_api/route"

	
)

func main() {
	  // Load .env
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  .env file tidak ditemukan atau gagal dibaca")
    } else {
        log.Println("📄 .env berhasil dimuat")
    }
	// Connect database
	// Connect database
if err := database.ConnectPostgres(); err != nil {
    log.Fatal(err)
}

if err := database.ConnectMongo(); err != nil {
    log.Fatal(err)
}

	app := fiber.New()

	// ===== REPOSITORY =====
	achievementMongoRepo := repository.NewAchievementMongoRepository()
	achievementPostgresRepo := repository.NewAchievementPostgresRepository()
	studentRepo := repository.NewStudentPostgresRepository()
	userRepo := repository.NewUserPostgresRepository()
	roleRepo := repository.NewRolePostgresRepository()
	permissionRepo := repository.NewPermissionPostgresRepository()
	rolePermissionRepo := repository.NewRolePermissionPostgresRepository()
	lecturerRepo := repository.NewLecturerPostgresRepository()
	// ===== SERVICE =====
	authSvc := service.NewAuthService(
    userRepo,
    roleRepo,
    permissionRepo,
    rolePermissionRepo,
    studentRepo,
    lecturerRepo,  // <--- TAMBAH INI
)

	achievementSvc := &service.AchievementService{
		MongoRepo:    achievementMongoRepo,
		PostgresRepo: achievementPostgresRepo,
		StudentRepo:  studentRepo,
		
	}
// === Tambahkan service lainnya sesuai modul ===
userSvc := service.NewUserService(userRepo, roleRepo)
lecturerSvc := service.NewLecturerService(studentRepo)
reportSvc := service.NewReportService(achievementPostgresRepo, studentRepo)

	// ===== ROUTES =====
	route.AuthRouter(app, authSvc)
route.AchievementRouter(app, achievementSvc)
route.UserRouter(app, userSvc)
route.LecturerRouter(app, lecturerSvc)
route.AdminAchievementRouter(app, achievementSvc)
route.ReportRouter(app, reportSvc)


	log.Println("🚀 Prestasi API berjalan di port 3000...")
	app.Listen(":3000")
}
