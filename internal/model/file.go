package model

import (
	"github.com/google/uuid"
	"time"
)

type File struct {
	ID            uuid.UUID `db:"id" json:"id"`
	SourceID      string    `db:"source_id" json:"-"`
	Title         string    `db:"title" json:"title"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	Size          int64     `db:"size" json:"size"`
	DownloadCount int       `db:"download_count" json:"download_count"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
