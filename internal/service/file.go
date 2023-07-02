package service

import (
	"context"
	"github.com/ensiouel/apperror"
	pb_static "github.com/ensiouel/basket-contract/gen/go/static/v1"
	"github.com/ensiouel/basket/internal/dto"
	"github.com/ensiouel/basket/internal/model"
	"github.com/ensiouel/basket/internal/storage"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"time"
)

type FileService interface {
	Get(ctx context.Context, fileID uuid.UUID) (model.File, error)
	Update(ctx context.Context, fileID uuid.UUID, request dto.UpdateFileRequest) (model.File, error)
	Upload(ctx context.Context, request dto.UploadFileRequest) (model.File, error)
	Download(ctx context.Context, fileID uuid.UUID, writer http.ResponseWriter) error
	Delete(ctx context.Context, fileID uuid.UUID) error
}

type FileServiceImpl struct {
	staticClient pb_static.StaticClient
	storage      storage.FileStorage
	maxFileSize  int64
}

func NewFileService(staticClient pb_static.StaticClient, storage storage.FileStorage, maxFileSize int64) *FileServiceImpl {
	return &FileServiceImpl{
		staticClient: staticClient,
		storage:      storage,
		maxFileSize:  maxFileSize,
	}
}

func (service *FileServiceImpl) Get(ctx context.Context, fileID uuid.UUID) (model.File, error) {
	file, err := service.storage.Get(ctx, fileID)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.NotFound); ok {
			return model.File{}, apperr.WithMessage("file not found")
		}

		return model.File{}, err
	}

	return file, nil
}

func (service *FileServiceImpl) Update(ctx context.Context, fileID uuid.UUID, request dto.UpdateFileRequest) (model.File, error) {
	file, err := service.Get(ctx, fileID)
	if err != nil {
		return model.File{}, err
	}

	if request.Title != "" {
		file.Title = request.Title
	}

	if request.Name != "" {
		file.Name = request.Name
	}

	if request.Description != "" {
		file.Description = request.Description
	}

	file.UpdatedAt = time.Now()

	err = service.storage.Update(ctx, file)
	if err != nil {
		return model.File{}, err
	}

	return file, nil
}

func (service *FileServiceImpl) Upload(ctx context.Context, request dto.UploadFileRequest) (model.File, error) {
	if request.FileHeader == nil {
		return model.File{}, apperror.BadRequest.WithMessage("file is required")
	}

	if request.FileHeader.Filename == "" {
		return model.File{}, apperror.BadRequest.WithMessage("file name is required")
	}

	if request.FileHeader.Size > service.maxFileSize {
		return model.File{}, apperror.BadRequest.WithMessage("file size is too large")
	}

	source, err := request.FileHeader.Open()
	if err != nil {
		return model.File{}, apperror.Internal.WithError(err)
	}
	defer source.Close()

	var client pb_static.Static_UploadClient
	client, err = service.staticClient.Upload(ctx)
	if err != nil {
		return model.File{}, err
	}

	var (
		buffer = make([]byte, 32*1024)
		n      int
	)
	for {
		n, err = source.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}

			return model.File{}, err
		}

		err = client.Send(&pb_static.UploadRequest{
			Data: buffer[:n],
		})
		if err != nil {
			return model.File{}, err
		}
	}

	var recv *pb_static.UploadResponse
	recv, err = client.CloseAndRecv()
	if err != nil {
		return model.File{}, err
	}

	id := uuid.New()
	sourceID := recv.GetSourceId()
	title, _, _ := strings.Cut(request.FileHeader.Filename, ".")
	name := request.FileHeader.Filename
	size := request.FileHeader.Size

	now := time.Now()

	file := model.File{
		ID:            id,
		SourceID:      sourceID,
		Title:         title,
		Name:          name,
		Description:   "",
		Size:          size,
		DownloadCount: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = service.storage.Create(ctx, file)
	if err != nil {
		return model.File{}, err
	}

	return file, nil
}

func (service *FileServiceImpl) Download(ctx context.Context, fileID uuid.UUID, writer http.ResponseWriter) error {
	file, err := service.Get(ctx, fileID)
	if err != nil {
		return err
	}

	client, err := service.staticClient.Download(ctx, &pb_static.DownloadRequest{
		SourceId: file.SourceID,
	})
	if err != nil {
		return err
	}

	writer.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
	writer.Header().Set("Content-Type", "application/octet-stream")

	var recv *pb_static.DownloadResponse
	for {
		recv, err = client.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		_, err = writer.Write(recv.GetData())
		if err != nil {
			return err
		}
	}

	err = client.CloseSend()
	if err != nil {
		return err
	}

	file.DownloadCount++

	err = service.storage.Update(ctx, file)
	if err != nil {
		return err
	}

	return nil
}

func (service *FileServiceImpl) Delete(ctx context.Context, fileID uuid.UUID) error {
	file, err := service.Get(ctx, fileID)
	if err != nil {
		return err
	}

	err = service.storage.Delete(ctx, file.ID)
	if err != nil {
		return err
	}

	exists, err := service.storage.ExistsBySourceID(ctx, file.SourceID)
	if err != nil {
		return err
	}

	// если source_id больше никем не используется - удаляем source
	if !exists {
		_, err = service.staticClient.Delete(ctx, &pb_static.DeleteRequest{
			SourceId: file.SourceID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
