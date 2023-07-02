package storage

import (
	"context"
	"errors"
	"github.com/ensiouel/apperror"
	"github.com/ensiouel/basket/internal/model"
	"github.com/ensiouel/basket/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FileStorage interface {
	Create(ctx context.Context, file model.File) error
	Get(ctx context.Context, fileID uuid.UUID) (model.File, error)
	Update(ctx context.Context, file model.File) error
	Delete(ctx context.Context, fileID uuid.UUID) error
	ExistsBySourceID(ctx context.Context, sourceID string) (bool, error)
}

type FileStorageImpl struct {
	client postgres.Client
}

func NewFileStorage(client postgres.Client) *FileStorageImpl {
	return &FileStorageImpl{client: client}
}

func (storage *FileStorageImpl) Create(ctx context.Context, file model.File) error {
	q := `
INSERT INTO file (id, source_id, title, name, description, size, download_count, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

	_, err := storage.client.Exec(ctx, q,
		file.ID,
		file.SourceID,
		file.Title,
		file.Name,
		file.Description,
		file.Size,
		file.DownloadCount,
		file.CreatedAt,
		file.UpdatedAt,
	)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return nil
}

func (storage *FileStorageImpl) Get(ctx context.Context, fileID uuid.UUID) (model.File, error) {
	q := `
SELECT * 
FROM file
WHERE id = $1
`

	var file model.File
	err := storage.client.Get(ctx, &file, q, fileID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.File{}, apperror.NotFound.WithError(err)
		}

		return model.File{}, apperror.Internal.WithError(err)
	}

	return file, nil
}

func (storage *FileStorageImpl) Update(ctx context.Context, file model.File) error {
	q := `
UPDATE file
SET title          = $1,
    name           = $2,
    description    = $3,
    download_count = $4,
    updated_at     = $5
WHERE id = $6
`

	_, err := storage.client.Exec(ctx, q,
		file.Title,
		file.Name,
		file.Description,
		file.DownloadCount,
		file.UpdatedAt,
		file.ID,
	)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return nil
}

func (storage *FileStorageImpl) Delete(ctx context.Context, fileID uuid.UUID) error {
	q := `
DELETE
FROM file
WHERE id = $1
`

	_, err := storage.client.Exec(ctx, q, fileID)
	if err != nil {
		return apperror.Internal.WithError(err)
	}

	return nil
}

func (storage *FileStorageImpl) ExistsBySourceID(ctx context.Context, sourceID string) (bool, error) {
	q := `
SELECT EXISTS (SELECT 1
               FROM file
               WHERE source_id = $1)
`

	var exists bool
	err := storage.client.Get(ctx, &exists, q, sourceID)
	if err != nil {
		return false, apperror.Internal.WithError(err)
	}

	return exists, err
}
