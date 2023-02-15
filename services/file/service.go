package file

import (
	"io"

	repository "github.com/alexshv/file-storage/repository"
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
	Delete(requestId interface{}, key string) error
}

type fileService struct {
	log            *logrus.Logger
	fileRepository repository.FileRepository
}

func NewFileService(log *logrus.Logger, repository repository.FileRepository) *fileService {
	return &fileService{
		log,
		repository,
	}
}
