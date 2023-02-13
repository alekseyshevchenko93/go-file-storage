package types

import "time"

type File struct {
	Id               int        `db:"id"`
	Key              string     `db:"key"`
	Extension        string     `db:"extension"`
	LastDownloadedAt *time.Time `db:"last_downloaded_at"`
	CreatedAt        *time.Time `db:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at"`
}
