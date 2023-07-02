package dto

import (
	"github.com/ensiouel/basket/internal/model"
	"mime/multipart"
)

type GetFileResponse struct {
	model.File
}

type UpdateFileRequest struct {
	Title       string `json:"title"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateFileResponse struct {
	model.File
}

type UploadFileRequest struct {
	FileHeader *multipart.FileHeader `form:"file"`
}

type UploadFileResponse struct {
	model.File
}

type DeleteFileResponse int
