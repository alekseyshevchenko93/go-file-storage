package repository

import (
	log "github.com/alexshv/file-storage/logger"
	"github.com/alexshv/file-storage/postgres"
	"github.com/alexshv/file-storage/types"
	"github.com/sirupsen/logrus"
)

func GetFileByKey(key string) (*types.File, error) {
	client := postgres.GetClient()

	var f []types.File

	stmt, err := client.PrepareNamed("SELECT * FROM files WHERE key = :key")

	if err != nil {
		log.GetLogger().WithFields(logrus.Fields{
			"message": err.Error(),
		}).Error("postgresRepository.getFileByKey.prepareStatementError")

		return nil, err
	}

	params := map[string]interface{}{
		"key": key,
	}

	if err := stmt.Select(&f, params); err != nil {
		log.GetLogger().WithFields(logrus.Fields{
			"message": err.Error(),
		}).Error("postgresRepository.getFileByKey.queryError")

		return nil, err
	}

	if len(f) == 0 {
		return nil, nil
	}

	return &f[0], nil
}

func CreateFile(file types.File) error {
	client := postgres.GetClient()

	params := map[string]interface{}{
		"key":       file.Key,
		"extension": file.Extension,
	}

	if _, err := client.NamedExec("INSERT INTO files(key, extension) VALUES(:key, :extension)", params); err != nil {
		log.GetLogger().WithFields(logrus.Fields{
			"message": err.Error(),
		}).Error("postgresRepository.createfile.error")

		return err
	}

	return nil
}

func DeleteFile() {

}

func UpdateFileLastDownloadedAt() {

}

func GetLeastUsedFiles() {

}

func DeleteLeastUsedFiles() {

}
