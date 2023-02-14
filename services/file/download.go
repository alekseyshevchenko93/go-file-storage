package file

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (s *fileService) Download(requestId interface{}, key string) (string, error) {
	log := s.log
	repository := s.fileRepository

	file, err := repository.GetFileByKey(key)

	if err != nil {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err.Error(),
		}).Error("fileService.download.databaseError")

		return "", fiber.NewError(fiber.StatusInternalServerError)
	}

	if file == nil {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
		}).Warn("fileService.download.fileNotFoundInDatabase")

		return "", fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	filepath := fmt.Sprintf("%s/%s.%s", os.Getenv("STORAGE_PATH"), file.Key, file.Extension)

	_, err = os.Stat(filepath)

	if err != nil && err == os.ErrNotExist {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err.Error(),
		}).Warn("fileService.download.fileNotFoundInStorage")

		return "", fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	if err != nil {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err.Error(),
		}).Error("fileService.download.failedToFindFile")

		return "", fiber.NewError(fiber.StatusInternalServerError)
	}

	return filepath, nil
}
