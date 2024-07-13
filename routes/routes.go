package routes

import(
	"github.com/gofiber/fiber/v2"
	"github.com/devmehta22/JWT-Auth/controllers"
)

func SetUpRoutes(app *fiber.App)  {
	app.Post("/api/register",controllers.Register)
	app.Post("/api/login",controllers.Login)
	app.Get("/api/protected",controllers.Protected)
	app.Get("/api/logout",controllers.Logout)
}