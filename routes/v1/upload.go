package v1

import (
	"backend-brevet/controllers"
	"backend-brevet/dto"
	"backend-brevet/middlewares"
	"backend-brevet/services"

	"github.com/gofiber/fiber/v2"
)

// RegisterUploadRoutes registers upload routes
func RegisterUploadRoutes(r fiber.Router) {

	fileService := services.NewFileService()
	uploadController := controllers.NewFileController(fileService)

	// Routes untuk upload
	imageGroup := r.Group("/images")
	imageGroup.Post("/", middlewares.ValidateBody[dto.UploadRequest](), uploadController.UploadImage)

	docGroup := r.Group("/documents")
	docGroup.Post("/", middlewares.ValidateBody[dto.UploadRequest](), uploadController.UploadDocument)

	// Hapus file (umum untuk semua jenis file)
	r.Delete("/", middlewares.ValidateBody[dto.DeleteRequest](), uploadController.DeleteFile)
}
