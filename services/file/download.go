package file

import (
	"fmt"
	"os"

	"github.com/alexshv/file-storage/types"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (s *fileService) validateDownloadRequest(key string) (*types.File, error) {
	repository := s.fileRepository

	if s.validateUuid(key) == false {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid key")
	}

	file, err := repository.GetFileByKey(key)

	if err != nil {
		return nil, fmt.Errorf("get file by key database error: %w", err)
	}

	if file == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	return file, nil
}

func (s *fileService) Download(requestId interface{}, key string) (string, error) {
	log := s.log
	repository := s.fileRepository

	file, err := repository.GetFileByKey(key)

	if err != nil {
		return "", err
	}

	if file == nil {
		return "", fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	filepath := s.getFilepath(file)

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
		return "", err
	}

	return filepath, nil
}
