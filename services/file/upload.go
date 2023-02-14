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
	"strings"

	"github.com/alexshv/file-storage/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (s *fileService) Upload(
	requestId interface{},
	uuid uuid.UUID,
	clientChecksum string,
	contentType string,
	bodyStream io.Reader,
) error {
	log := s.log
	repository := s.fileRepository

	if contentType == "" {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
		}).Warn("fileService.upload.emptyContentType")

		return fiber.NewError(fiber.StatusBadRequest, "Content type is missing")
	}

	mediaType, params, err := mime.ParseMediaType(contentType)

	if err != nil || mediaType != fiber.MIMEMultipartForm {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err,
		}).Warn("fileService.upload.badMediaType")

		return fiber.NewError(fiber.StatusBadRequest)
	}

	boundary := params["boundary"]
	multipartReader := multipart.NewReader(bodyStream, boundary)
	part, err := multipartReader.NextPart()

	if err != nil && err != io.EOF {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err,
		}).Error("fileService.upload.failedToReadFirstPart")

		return fiber.NewError(fiber.StatusInternalServerError)
	}

	extension := strings.Trim(filepath.Ext(part.FileName()), ".")
	path := fmt.Sprintf("%s/%s.%s", os.Getenv("STORAGE_PATH"), uuid, extension)
	hasher := sha1.New()

	fd, err := os.Create(path)
	defer fd.Close()

	if err != nil {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err,
		}).Error("fileService.upload.failedToOpenFile")

		return fiber.NewError(fiber.StatusInternalServerError)
	}

	if err = s.processPart(requestId, part, hasher, fd); err != nil {
		log.WithFields(logrus.Fields{
			"requestId": requestId,
			"message":   err,
		}).Error("fileService.upload.failedToProcessFromBeginningOfPart")

		return fiber.NewError(fiber.StatusInternalServerError)
	}

	for {
		part, err := multipartReader.NextPart()

		if err != nil && err == io.EOF {
			break
		}

		if err != nil {
			log.WithFields(logrus.Fields{
				"requestId": requestId,
				"message":   err,
			}).Error("fileService.upload.failedToReadFromPart")

			return fiber.NewError(fiber.StatusInternalServerError)
		}

		if err = s.processPart(requestId, part, hasher, fd); err != nil {
			log.WithFields(logrus.Fields{
				"requestId": requestId,
				"message":   err,
			}).Error("fileService.upload.failedToProcessPart")

			return fiber.NewError(fiber.StatusInternalServerError)
		}
	}

	serverChecksum := hex.EncodeToString(hasher.Sum(nil))

	if clientChecksum != "" && clientChecksum != serverChecksum {
		err := os.Remove(path)

		if err != nil {
			log.WithFields(logrus.Fields{
				"requestId": requestId,
				"message":   err,
			}).Error("fileService.upload.failedToRemoveFile")

			return fiber.NewError(fiber.StatusInternalServerError)
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
	log := s.log
	buffer := make([]byte, 1024*1024)

	for {
		read, err := part.Read(buffer)

		if err != nil && err != io.EOF {
			log.WithFields(logrus.Fields{
				"requestId": requestId,
				"message":   err,
			}).Error("handlers.processPart.failedToReadPart")

			return fiber.NewError(fiber.StatusInternalServerError)
		}

		_, hasherErr := hasher.Write(buffer[:read])

		if hasherErr != nil {
			log.WithFields(logrus.Fields{
				"requestId": requestId,
				"message":   err,
			}).Error("handlers.processPart.failedToWriteToHash")

			return fiber.NewError(fiber.StatusInternalServerError)
		}

		_, writeErr := file.WriteString(string(buffer[:read]))

		if writeErr != nil {
			log.WithFields(logrus.Fields{
				"requestId": requestId,
				"message":   err,
			}).Error("handlers.processPart.failedToWriteToFile")

			return fiber.NewError(fiber.StatusInternalServerError)
		}

		if err != nil && err == io.EOF {
			break
		}
	}

	return nil
}
