package fileService

import (
	"fmt"
	"os"

	log "github.com/alexshv/file-storage/logger"
	"github.com/alexshv/file-storage/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Download(requestId interface{}, key string) (string, error) {
	file, err := repository.GetFileByKey(key)

	if err != nil {
		log.GetLogger().WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err,
		}).Error("fileService.download.databaseError")
	}

	if file == nil {
		log.GetLogger().WithFields(logrus.Fields{
			"requestId": requestId,
		}).Warn("fileService.download.fileNotFoundInDatabase")
		return "", fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	filepath := fmt.Sprintf("%s/%s.%s", os.Getenv("STORAGE_PATH"), file.Key, file.Extension)

	_, err = os.Stat(filepath)

	if err != nil && err == os.ErrNotExist {
		log.GetLogger().WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err,
		}).Warn("fileService.download.fileNotFoundInStorage")

		return "", fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	if err != nil {
		log.GetLogger().WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err,
		}).Error("fileService.download.failedToFindFile")

		return "", fiber.NewError(fiber.StatusInternalServerError)
	}

	return filepath, nil
}
