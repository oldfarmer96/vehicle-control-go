// Package response - respuestas formao universal
package response

import "github.com/gofiber/fiber/v3"

type Response struct {
	Success bool   `json:"success"`
	Res     any    `json:"res,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Success(c fiber.Ctx, data any) error {
	return c.JSON(Response{
		Success: true,
		Res:     data,
	})
}

func Error(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Error:   message,
	})
}
