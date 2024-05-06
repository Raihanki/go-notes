package routes

import (
	"github.com/Raihanki/go-notes/controllers"
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
	userRoute.Post("/register", userController.Register)
	userRoute.Post("/login", userController.Login)

	// note routes
	noteController := controllers.NewNoteController()
	noteRoute := api.Group("/notes")
	noteRoute.Get("/", noteController.Index)
	noteRoute.Post("/", noteController.Store)
	noteRoute.Get("/:id", noteController.Show)
	noteRoute.Put("/:id", noteController.Update)
	noteRoute.Delete("/:id", noteController.Destroy)
}
