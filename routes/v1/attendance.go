package v1

import (
	"backend-brevet/controllers"
	"backend-brevet/middlewares"
	"backend-brevet/repository"
	"backend-brevet/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterAttendanceRoutes register attendance routes
func RegisterAttendanceRoutes(r fiber.Router, db *gorm.DB) {

	meetingRepository := repository.NewMeetingRepository(db)
	purchaseRepository := repository.NewPurchaseRepository(db)

	attendanceRepository := repository.NewAttendanceRepository(db)
	attendanceService := services.NewAttendanceService(attendanceRepository, meetingRepository, purchaseRepository, db)
	attendanceController := controllers.NewAttendanceController(attendanceService, db)

	r.Get("/", middlewares.RequireAuth(),
		middlewares.RequireRole([]string{"admin"}), attendanceController.GetAllAttendances)
	r.Get("/:id", middlewares.RequireAuth(),
		middlewares.RequireRole([]string{"admin"}), attendanceController.GetAttendanceByID)
}
