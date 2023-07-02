package handler

import (
	"github.com/ensiouel/basket/internal/dto"
	"github.com/ensiouel/basket/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"net/http"
)

type FileHandler struct {
	fileService service.FileService
}

func NewFileHandler(fileService service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (handler *FileHandler) Register(router gin.IRoutes) {
	router.GET("/:file_id", handler.get)
	router.PATCH("/:file_id", handler.update)
	router.DELETE("/:file_id", handler.delete)
	router.POST("/upload", handler.upload)
	router.GET("/:file_id/download", handler.download)
}

func (handler *FileHandler) get(c *gin.Context) {
	fileID, err := uuid.Parse(c.Param("file_id"))
	if err != nil {
		c.Error(err)
		return
	}

	fileInfo, err := handler.fileService.Get(c, fileID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": dto.GetFileResponse{File: fileInfo}})
}

func (handler *FileHandler) update(c *gin.Context) {
	fileID, err := uuid.Parse(c.Param("file_id"))
	if err != nil {
		c.Error(err)
		return
	}

	var request dto.UpdateFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	file, err := handler.fileService.Update(c, fileID, request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": dto.UpdateFileResponse{File: file}})
}

func (handler *FileHandler) delete(c *gin.Context) {
	fileID, err := uuid.Parse(c.Param("file_id"))
	if err != nil {
		c.Error(err)
		return
	}

	err = handler.fileService.Delete(c, fileID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": dto.DeleteFileResponse(1)})
}

func (handler *FileHandler) upload(c *gin.Context) {
	var request dto.UploadFileRequest
	if err := c.ShouldBindWith(&request, binding.FormMultipart); err != nil {
		c.Error(err)
		return
	}

	file, err := handler.fileService.Upload(c, request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": dto.UploadFileResponse{File: file}})
}

func (handler *FileHandler) download(c *gin.Context) {
	fileID, err := uuid.Parse(c.Param("file_id"))
	if err != nil {
		c.Error(err)
		return
	}

	err = handler.fileService.Download(c, fileID, c.Writer)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
