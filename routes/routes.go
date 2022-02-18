package routes

import (
	"github.com/cynicade/todo/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/todo/api/register", controllers.Register)
	app.Post("/todo/api/login", controllers.Login)
	app.Get("/todo/api/login", controllers.Login)
	app.Get("/todo/api/logout", controllers.Logout)
	app.Get("/todo/api/todos", controllers.GetAll)
	app.Post("/todo/api/new", controllers.NewTodo)
	app.Post("/todo/api/todo/:id", controllers.Update)
	app.Delete("/todo/api/todo/:id", controllers.Delete)
	app.Get("/todo/api/ping", controllers.Ping)
}
