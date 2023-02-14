package middlewares

import (
	"github.com/alexshv/file-storage/container"
	"github.com/gofiber/fiber/v2"
)

func SetContainer(container *container.Container) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals("container", container)
		return c.Next()
	}
}
