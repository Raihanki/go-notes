package helpers

import "github.com/gofiber/fiber/v2"

type ApiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	FailedField string      `json:"key"`
	Tag         string      `json:"tag"`
	Message     interface{} `json:"message"`
}

func Response(ctx *fiber.Ctx, statusCode int, message string, data interface{}) error {
	response := ApiResponse{
		Message: message,
		Data:    data,
	}

	return ctx.Status(statusCode).JSON(response)
}
