package routes

import (
	"github.com/Raihanki/go-notes/controllers"
	"github.com/Raihanki/go-notes/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	api := app.Group("/api")

	// topic routes
	topicController := controllers.NewTopicController()
	topicRoute := api.Group("/topics")
	topicRoute.Get("/", topicController.Index)

	// user routes
	userController := controllers.NewUserController()
	userRoute := api.Group("/users")
	userRoute.Get("/", middleware.Auth, userController.Show)
	userRoute.Post("/register", middleware.Guest, userController.Register)
	userRoute.Post("/login", middleware.Guest, userController.Login)

	// note routes
	noteController := controllers.NewNoteController()
	noteRoute := api.Group("/notes")
	noteRoute.Get("/", noteController.Index)
	noteRoute.Post("/", middleware.Auth, noteController.Store)
	noteRoute.Get("/:id", noteController.Show)
	noteRoute.Put("/:id", middleware.Auth, noteController.Update)
	noteRoute.Delete("/:id", middleware.Auth, noteController.Destroy)
}
