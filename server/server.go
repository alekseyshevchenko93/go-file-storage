package server

import (
	"github.com/alexshv/file-storage/container"
	"github.com/alexshv/file-storage/server/controllers"
	"github.com/alexshv/file-storage/server/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type server struct {
	port           string
	fiber          *fiber.App
	container      *container.Container
	fileController *controllers.FileController
}

func (s *server) Stop() error {
	return s.fiber.Shutdown()
}

func (s *server) Start() error {
	return s.fiber.Listen(s.port)
}

func NewServer(port string, container *container.Container, fileController *controllers.FileController) *server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage:        true,
		DisablePreParseMultipartForm: true,
		ErrorHandler:                 middlewares.ErrorHandler,
	})

	app.Server().StreamRequestBody = true
	app.Use(requestid.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("container", container)
		return c.Next()
	})

	v1 := app.Group("/api/v1")
	v1.Get("/download/:key", fileController.DownloadFile)
	v1.Delete("/download/:key", fileController.DeleteFile)
	v1.Post("/upload", fileController.UploadFile)

	app.Use(middlewares.NotFoundHandler)

	return &server{
		port,
		app,
		container,
		fileController,
	}
}
