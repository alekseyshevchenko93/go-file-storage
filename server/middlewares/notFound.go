package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func NotFoundHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "Not found",
	})
}
