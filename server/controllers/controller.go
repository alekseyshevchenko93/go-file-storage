package controllers

import (
	"fmt"

	container "github.com/alexshv/file-storage/container"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type FileController struct{}

func NewFileController() *FileController {
	return &FileController{}
}

func (ctr *FileController) DownloadFile(c *fiber.Ctx) error {
	key := c.Params("key")
	requestId := c.Locals("requestid")
	container := c.Locals("container").(*container.Container)
	log := container.GetLogger()
	fileService := container.GetFileService()

	log.WithFields(logrus.Fields{
		"requestId": requestId,
		"key":       key,
	}).Info("handlers.downloadHandler.request")

	filepath, err := fileService.Download(requestId, key)

	if err != nil {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err.Error(),
		}).Info("handlers.downloadHandler.error")

		return err
	}

	log.WithFields(logrus.Fields{
		"requestId": requestId,
	}).Info("handlers.downloadHandler.success")

	return c.SendFile(filepath)
}

func (ctr *FileController) UploadFile(c *fiber.Ctx) error {
	uuid := uuid.New()
	clientChecksum := c.Query("checksum")
	contentType := c.Get("Content-Type")

	fmt.Println("in upload", c.Locals("container"))

	requestId := c.Locals("requestid")
	container := c.Locals("container").(*container.Container)

	bodyStream := c.Context().RequestBodyStream()

	log := container.GetLogger()
	fileService := container.GetFileService()

	log.WithFields(logrus.Fields{
		"requestId":      requestId,
		"uuid":           uuid,
		"clientChecksum": clientChecksum,
	}).Info("handlers.uploadHandler.request")

	if err := fileService.Upload(requestId, uuid, clientChecksum, contentType, bodyStream); err != nil {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err.Error(),
		}).Error("handlers.uploadHandler.error")

		return err
	}

	log.WithFields(logrus.Fields{
		"requestId": requestId,
	}).Info("handlers.uploadHandler.success")

	return c.JSON(fiber.Map{
		"key": uuid,
	})
}

func (ctr *FileController) DeleteFile(c *fiber.Ctx) error {
	key := c.Params("key")
	requestId := c.Locals("requestid")
	container := c.Locals("container").(*container.Container)
	log := container.GetLogger()
	fileService := container.GetFileService()

	log.WithFields(logrus.Fields{
		"requestId": requestId,
		"key":       key,
	}).Info("handlers.deleteFileHandler.request")

	err := fileService.Delete(requestId, key)

	if err != nil {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err.Error(),
		}).Info("handlers.deleteFileHandler.error")

		return err
	}

	log.WithFields(logrus.Fields{
		"requestId": requestId,
	}).Info("handlers.deleteFileHandler.success")

	c.Status(fiber.StatusAccepted)

	return nil
}
