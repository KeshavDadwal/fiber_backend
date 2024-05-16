package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keshav/fiber/controllers"
	"github.com/keshav/fiber/middleware"
)

func SetupAdminRoutes(app *fiber.App) {

	app.Post("/login", controllers.Adminlogin)
	app.Post("/role",controllers.VariousRole)
	app.Post("/signup", controllers.AdminSignUp)
	app.Get("/users",middleware.RequireAuth, controllers.GetAllUser)
	app.Post("/uploadfile",controllers.UploadImage)
	app.Post("/uploadVideo",controllers.UploadHandler)
}
