package file

import (
	"io"

	postgresRepository "github.com/alexshv/file-storage/postgres/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FileService interface {
	Download(requestId interface{}, key string) (string, error)
	Upload(requestId interface{},
		uuid uuid.UUID,
		clientChecksum string,
		contentType string,
		bodyStream io.Reader,
	) error
}

type fileService struct {
	log            *logrus.Logger
	fileRepository postgresRepository.FileRepository
}

func NewFileService(log *logrus.Logger, repository postgresRepository.FileRepository) *fileService {
	return &fileService{
		log,
		repository,
	}
}
