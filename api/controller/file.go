package controller

import (
	"net/http"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type FileController struct {
	FileUsecase domain.FileUsecase
}

func (col *FileController) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	resp, err := col.FileUsecase.Upload(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Upload File Fail", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(resp))
}

func (col *FileController) Uploads(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	files := form.File["file"]

	resp, err := col.FileUsecase.Uploads(c, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Upload File Fail", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(resp))
}

func (col *FileController) FileList(c *gin.Context) {
	var req domain.FileListRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResp("Bad Params", err))
		return
	}
	resp, cnt, err := col.FileUsecase.FileList(c, req.Type, req.Source, req.Page, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResp("Upload File Fail", err))
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResp(map[string]interface{}{
		"count":     cnt,
		"file_list": resp,
	}))
}
