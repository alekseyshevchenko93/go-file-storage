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
		return "", fmt.Errorf("get file by key database error: %w", err)
	}

	if file == nil {
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
		return "", fmt.Errorf("find file in file system error: %w", err)
	}

	if err := repository.UpdateFileLastDownloadedAt(file); err != nil {
		return "", fmt.Errorf("repository.UpdateFileLastDownloadedAt err: %w", err)
	}

	return filepath, nil
}
