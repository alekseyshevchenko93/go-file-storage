package container

import (
	"github.com/alexshv/file-storage/services/file"
	"github.com/sirupsen/logrus"
)

type Container struct {
	log         *logrus.Logger
	fileService file.FileService
}

func (c *Container) GetLogger() *logrus.Logger {
	return c.log
}

func (c *Container) GetFileService() file.FileService {
	return c.fileService
}

func New(log *logrus.Logger, fs file.FileService) *Container {
	return &Container{
		log:         log,
		fileService: fs,
	}
}
