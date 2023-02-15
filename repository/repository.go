package repository

import (
	"fmt"
	"time"

	"github.com/alexshv/file-storage/postgres"
	"github.com/alexshv/file-storage/types"
)

type FileRepository interface {
	GetFileByKey(key string) (*types.File, error)
	CreateFile(file *types.File) error
	DeleteFile(file *types.File) error
	UpdateFileLastDownloadedAt(file *types.File) error
	GetLeastUsedFilesIds(date *time.Time) ([]int, error)
}

type fileRepository struct {
	db postgres.PostgresClient
}

func NewFileRepository(db postgres.PostgresClient) *fileRepository {
	return &fileRepository{
		db,
	}
}

func (r *fileRepository) GetFileByKey(key string) (*types.File, error) {
	client := r.db.GetClient()

	var f []types.File

	stmt, err := client.PrepareNamed("SELECT * FROM files WHERE key = :key")

	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	params := map[string]interface{}{
		"key": key,
	}

	if err := stmt.Select(&f, params); err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	if len(f) == 0 {
		return nil, nil
	}

	return &f[0], nil
}

func (r *fileRepository) CreateFile(file *types.File) error {
	client := r.db.GetClient()

	params := map[string]interface{}{
		"key":       file.Key,
		"extension": file.Extension,
	}

	if _, err := client.NamedExec("INSERT INTO files(key, extension) VALUES(:key, :extension)", params); err != nil {
		return fmt.Errorf("insert error: %w", err)
	}

	return nil
}

func (r *fileRepository) DeleteFile(file *types.File) error {
	client := r.db.GetClient()

	params := map[string]interface{}{
		"id": file.Id,
	}

	if _, err := client.NamedExec("DELETE FROM files WHERE id = :id", params); err != nil {
		return fmt.Errorf("delete error: %w", err)
	}

	return nil
}

func (r *fileRepository) UpdateFileLastDownloadedAt(file *types.File) error {
	client := r.db.GetClient()

	params := map[string]interface{}{
		"id":           file.Id,
		"downloadedAt": time.Now().Format(time.RFC3339),
	}

	_, err := client.NamedExec(`UPDATE files SET last_downloaded_at = :downloadedAt WHERE id = :id`, params)

	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	return nil
}

func (r *fileRepository) GetLeastUsedFilesIds(date *time.Time) ([]int, error) {
	client := r.db.GetClient()

	params := map[string]interface{}{
		"date": date,
	}

	var ids []int

	stmt, err := client.PrepareNamed(`
	SELECT 
		id 
	FROM 
		files 
	WHERE 
		last_downloaded_at IS NULL OR last_downloaded_at < :date
`)

	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	if err := stmt.Select(&ids, params); err != nil {
		return nil, fmt.Errorf("failed to select error: %w", err)
	}

	return ids, nil
}

func (r *fileRepository) DeleteFilesByIds(ids []int) {

}
