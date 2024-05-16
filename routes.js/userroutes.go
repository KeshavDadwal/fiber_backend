package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keshav/fiber/controllers"
	"github.com/keshav/fiber/middleware"
)

func SetupUserRoutes(app *fiber.App) {
	
	app.Post("/signup", controllers.UserSignUp)
	app.Post("/login", controllers.UserLogin)
	app.Put("/users/:id", middleware.RequireAuth, controllers.UpdateOneUser)
	app.Delete("/users/:id", middleware.RequireAuth, controllers.Deletedata)
	app.Get("/verifymail/:id/:secret_code", controllers.VerifyApi)
}
