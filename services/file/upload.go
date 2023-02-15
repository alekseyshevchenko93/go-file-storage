package file

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alexshv/file-storage/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (s *fileService) isValidSha1(value string) (bool, error) {
	match, err := regexp.MatchString("^([0-9A-Fa-f]{2}[:]){19}([0-9A-Fa-f]{2})$", value)

	if err != nil {
		return false, fmt.Errorf("isValidSha1 error: %w", err)
	}

	return match, nil
}

func (s *fileService) validateUploadParams(
	requestId interface{},
	contentType string,
	clientChecksum string,
) (string, error) {
	if contentType == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "Content type is missing")
	}

	if clientChecksum != "" {
		valid, err := s.isValidSha1(clientChecksum)

		if err != nil {
			return "", err
		}

		if valid == false {
			return "", fiber.NewError(fiber.StatusBadRequest, "Checksum is not a sha1")
		}
	}

	mediaType, params, err := mime.ParseMediaType(contentType)

	if err != nil || mediaType != fiber.MIMEMultipartForm {
		return "", fiber.NewError(fiber.StatusBadRequest, "Invalid media type")
	}

	return params["boundary"], nil
}

func (s *fileService) Upload(
	requestId interface{},
	uuid uuid.UUID,
	clientChecksum string,
	contentType string,
	bodyStream io.Reader,
) error {
	log := s.log
	repository := s.fileRepository

	boundary, err := s.validateUploadParams(requestId, contentType, clientChecksum)

	if err != nil {
		return err
	}

	multipartReader := multipart.NewReader(bodyStream, boundary)
	part, err := multipartReader.NextPart()

	if err != nil && err != io.EOF {
		return fmt.Errorf("multipartReader error, failed to read first part: %w", err)
	}

	extension := strings.Trim(filepath.Ext(part.FileName()), ".")
	path := fmt.Sprintf("%s/%s.%s", os.Getenv("STORAGE_PATH"), uuid, extension)
	hasher := sha1.New()

	fd, err := os.Create(path)
	defer fd.Close()

	if err != nil {
		return fmt.Errorf("failed to open file error: %w", err)
	}

	if err = s.processPart(requestId, part, hasher, fd); err != nil {
		return fmt.Errorf("failed to process first part error: %w", err)
	}

	for {
		part, err := multipartReader.NextPart()

		if err != nil && err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read part error: %w", err)
		}

		if err = s.processPart(requestId, part, hasher, fd); err != nil {
			return fmt.Errorf("failed to process part error: %w", err)
		}
	}

	serverChecksum := hex.EncodeToString(hasher.Sum(nil))

	if clientChecksum != "" && clientChecksum != serverChecksum {
		err := os.Remove(path)

		if err != nil {
			return fmt.Errorf("failed to remove file error: %w", err)
		}

		log.WithFields(logrus.Fields{
			"requestId":      requestId,
			"clientChecksum": clientChecksum,
			"serverChecksum": serverChecksum,
		}).Warn("fileService.upload.checksumsDontMatch")

		return fiber.NewError(fiber.StatusUnprocessableEntity, "Please upload again, client checksum is not requal to server checksum")
	}

	databaseFile := types.File{
		Key:       uuid.String(),
		Extension: extension,
	}

	err = repository.CreateFile(&databaseFile)

	if err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"requestId":      requestId,
		"serverChecksum": serverChecksum,
	}).Info("fileService.upload.success")

	return nil
}

func (s *fileService) processPart(requestId interface{}, part *multipart.Part, hasher hash.Hash, file *os.File) error {
	buffer := make([]byte, 1024*1024)

	for {
		read, err := part.Read(buffer)

		if err != nil && err != io.EOF {
			return fmt.Errorf("processPart, failed to read part: %w", err)
		}

		_, hasherErr := hasher.Write(buffer[:read])

		if hasherErr != nil {
			return fmt.Errorf("processPart, failed to write to hash: %w", err)
		}

		_, writeErr := file.WriteString(string(buffer[:read]))

		if writeErr != nil {
			return fmt.Errorf("processPart, failed to write file: %w", err)
		}

		if err != nil && err == io.EOF {
			break
		}
	}

	return nil
}
