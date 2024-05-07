package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func Guest(ctx *fiber.Ctx) error {
	return ctx.Next()
}
