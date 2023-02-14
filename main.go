package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"

	"github.com/alexshv/file-storage/container"
	"github.com/alexshv/file-storage/controllers"
	"github.com/alexshv/file-storage/logger"
	"github.com/alexshv/file-storage/middlewares"
	"github.com/alexshv/file-storage/postgres"
	repository "github.com/alexshv/file-storage/repository"
	fileServicePackage "github.com/alexshv/file-storage/services/file"
)

func main() {
	log := logger.New()

	if err := godotenv.Load(); err != nil {
		log.WithField("message", err).Info("init.dotenv.error")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := os.MkdirAll(os.Getenv("STORAGE_PATH"), os.ModeDir); err != nil {
		log.WithField("message", err).Info("init.storagePath.error")
		os.Exit(1)
	}

	postgresClient, err := postgres.New(log)

	if err != nil {
		log.WithField("message", err.Error()).Info("init.initDatabase.error")
	}

	defer postgresClient.Shutdown()

	fileController := controllers.NewFileController()
	fileRepository := repository.NewFileRepository(postgresClient)
	fileService := fileServicePackage.NewFileService(log, fileRepository)
	cnt := container.New(log, fileService)

	app := fiber.New(fiber.Config{
		DisableStartupMessage:        true,
		DisablePreParseMultipartForm: true,
		ErrorHandler:                 middlewares.ErrorHandler,
	})

	defer app.Shutdown()

	app.Server().StreamRequestBody = true
	app.Use(requestid.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("container", cnt)
		return c.Next()
	})

	v1 := app.Group("/api/v1")
	v1.Get("/download/:key", fileController.Download)
	v1.Post("/upload", fileController.Upload)
	app.Use(middlewares.NotFoundHandler)

	go func() {
		port := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))
		log.WithField("port", port).Info("server.listening")

		if err := app.Listen(port); err != nil {
			log.WithField("message", err.Error()).Error("server.start.error")
			os.Exit(1)
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("server.shutdown")
		os.Exit(0)
	}
}
