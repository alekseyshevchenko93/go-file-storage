package file

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func (s *fileService) Delete(requestId interface{}, key string) error {
	repository := s.fileRepository

	file, err := repository.GetFileByKey(key)

	if err != nil {
		return fmt.Errorf("get file by key database error: %w", err)
	}

	if file == nil {
		return fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	repository.DeleteFile(file)

	filepath := fmt.Sprintf("%s/%s.%s", os.Getenv("STORAGE_PATH"), file.Key, file.Extension)

	_, err = os.Stat(filepath)

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
