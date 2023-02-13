package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	handlers "github.com/alexshv/file-storage/handlers"
	log "github.com/alexshv/file-storage/logger"
	"github.com/alexshv/file-storage/postgres"
	"github.com/alexshv/file-storage/redis"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.GetLogger().WithField("message", err).Info("init.dotenv.error")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Init()

	redis.Init(ctx)
	defer redis.Shutdown()

	if err := os.MkdirAll(os.Getenv("STORAGE_PATH"), os.ModeDir); err != nil {
		log.GetLogger().WithField("message", err).Info("init.storagePath.error")
		os.Exit(1)
	}

	postgresClient := postgres.New()
	defer postgresClient.Close()

	app := fiber.New(fiber.Config{
		DisableStartupMessage:        true,
		DisablePreParseMultipartForm: true,
		ErrorHandler:                 handlers.ErrorHandler,
	})
	app.Server().StreamRequestBody = true

	app.Use(requestid.New())

	v1 := app.Group("/api/v1")
	v1.Get("/download/:key", handlers.DownloadHandler)
	v1.Post("/upload", handlers.UploadHandler)

	port := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))

	go func() {
		log.GetLogger().WithField("port", port).Info("server.listening")

		defer app.Shutdown()

		if err := app.Listen(port); err != nil {
			log.GetLogger().WithField("message", err.Error()).Error("server.start.error")
			os.Exit(1)
		}
	}()

	select {
	case <-ctx.Done():
		log.GetLogger().Info("server.shutdown")
	}
}
