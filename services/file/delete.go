package file

import (
	"fmt"
	"os"

	"github.com/alexshv/file-storage/types"
	"github.com/gofiber/fiber/v2"
)

func (s *fileService) validateDeleteRequest(key string) (*types.File, error) {
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

func (s *fileService) deleteFromDisk(file *types.File) error {
	filepath := s.getFilepath(file)

	_, err := os.Stat(filepath)

	if err != nil && err == os.ErrNotExist {
		return fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	if err != nil {
		return fmt.Errorf("find file in file storage error: %w", err)
	}

	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("os.Remove error: %w", err)
	}

	return nil
}

func (s *fileService) Delete(requestId interface{}, key string) error {
	repository := s.fileRepository

	file, err := s.validateDeleteRequest(key)

	if err != nil {
		return err
	}

	repository.DeleteFile(file)

	if err := s.deleteFromDisk(file); err != nil {
		return err
	}

	return nil
}
