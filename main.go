package main

import (
	"log"

	"github.com/Raihanki/go-notes/config"
	"github.com/Raihanki/go-notes/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadConfig()
	config.ConnectDB()

	app := fiber.New(fiber.Config{})
	routes.Router(app)

	errListen := app.Listen("localhost:" + config.ENV.APP_PORT)
	if errListen != nil {
		log.Fatal(errListen)
	}
}
