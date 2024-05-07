package main

import (
	"log"

	"github.com/Raihanki/go-notes/config"
	"github.com/Raihanki/go-notes/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.LoadConfig()
	config.ConnectDB()

	app := fiber.New(fiber.Config{})
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.ENV.APP_FRONTEND_URL,
	}))
	routes.Router(app)

	errListen := app.Listen("localhost:" + config.ENV.APP_PORT)
	if errListen != nil {
		log.Fatal(errListen)
	}
}
