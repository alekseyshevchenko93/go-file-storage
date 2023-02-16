package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/alexshv/file-storage/container"
	"github.com/alexshv/file-storage/logger"
	"github.com/alexshv/file-storage/postgres"
	repository "github.com/alexshv/file-storage/repository"
	"github.com/alexshv/file-storage/server"
	"github.com/alexshv/file-storage/server/controllers"
	fileServicePackage "github.com/alexshv/file-storage/services/file"
	"github.com/alexshv/file-storage/workers"
)

func main() {
	log := logger.New()

	if err := godotenv.Load(); err != nil {
		log.WithField("message", err).Info("init.dotenv.error")
		os.Exit(1)
	}

	if err := os.MkdirAll(os.Getenv("STORAGE_PATH"), os.ModeDir); err != nil {
		log.WithField("message", err).Info("init.storagePath.error")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	postgresClient, err := postgres.New(log)

	if err != nil {
		log.WithField("message", err.Error()).Info("init.initDatabase.error")
	}

	defer postgresClient.Shutdown()

	fileController := controllers.NewFileController()
	fileRepository := repository.NewFileRepository(postgresClient)
	fileService := fileServicePackage.NewFileService(log, fileRepository)
	container := container.New(log, fileService)

	worker := workers.NewWorker(ctx, log, fileRepository)
	worker.Start()
	defer worker.Stop()

	port := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))
	server := server.NewServer(
		port,
		container,
		fileController,
	)

	defer server.Stop()

	go func() {
		log.WithField("port", port).Info("server.started")

		if err := server.Start(); err != nil {
			log.WithField("message", err.Error()).Error("server.start.error")
			return
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("server.shutdown")
	}
}
