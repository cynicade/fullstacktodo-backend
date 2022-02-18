package routes

import (
	"github.com/cynicade/todo/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Get("/login", controllers.Login)
	app.Get("/logout", controllers.Logout)
	app.Get("/todos", controllers.GetAll)
	app.Post("/new", controllers.NewTodo)
	app.Post("/todo/:id", controllers.Update)
	app.Delete("/todo/:id", controllers.Delete)
	app.Get("/ping", controllers.Ping)
}
