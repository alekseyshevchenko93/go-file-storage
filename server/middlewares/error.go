package middlewares

import (
	"errors"
	"fmt"

	"github.com/alexshv/file-storage/container"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	fmt.Println("in error handler")
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var e *fiber.Error

	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	log := c.Locals("container").(*container.Container).GetLogger()
	requestId := c.Locals("requestid")

	log.WithFields(logrus.Fields{
		"requestId": requestId,
		"status":    code,
		"message":   message,
	}).Info("http.request.error")

	return c.Status(code).JSON(fiber.Map{
		"message": message,
	})
}
