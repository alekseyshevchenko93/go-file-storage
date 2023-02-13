package handlers

import (
	"errors"

	log "github.com/alexshv/file-storage/logger"
	fileService "github.com/alexshv/file-storage/services/fileService"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var e *fiber.Error

	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	log.GetLogger().WithFields(logrus.Fields{
		"status":  code,
		"message": message,
	}).Info("request.error")

	return ctx.Status(code).JSON(fiber.Map{
		"message": message,
	})
}

func DownloadHandler(c *fiber.Ctx) error {
	key := c.Params("key")
	requestId := c.Locals("requestid")

	log.GetLogger().WithFields(logrus.Fields{
		"requestId": requestId,
		"key":       key,
	}).Info("handlers.downloadHandler.request")

	filepath, err := fileService.Download(requestId, key)

	if err != nil {
		log.GetLogger().WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err,
		}).Info("handlers.downloadHandler.error")

		return err
	}

	log.GetLogger().WithFields(logrus.Fields{
		"requestId": requestId,
	}).Info("handlers.downloadHandler.success")

	return c.SendFile(filepath)
}

func UploadHandler(c *fiber.Ctx) error {
	uuid := uuid.New()
	requestId := c.Locals("requestid")
	clientChecksum := c.Query("checksum")
	contentType := c.Get("Content-Type")
	bodyStream := c.Context().RequestBodyStream()

	log.GetLogger().WithFields(logrus.Fields{
		"requestId":      requestId,
		"uuid":           uuid,
		"clientChecksum": clientChecksum,
	}).Info("handlers.uploadHandler.request")

	if err := fileService.Upload(requestId, uuid, clientChecksum, contentType, bodyStream); err != nil {
		log.GetLogger().WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err,
		}).Error("handlers.uploadHandler.error")

		return err
	}

	log.GetLogger().WithFields(logrus.Fields{
		"requestId": requestId,
	}).Info("handlers.uploadHandler.success")

	return c.JSON(fiber.Map{
		"key": uuid,
	})
}
