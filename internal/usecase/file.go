package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
)

type fileUsecase struct {
	repo           domain.FileReposirory
	contextTimeout time.Duration
}

func NewFileUsecase(repo domain.FileReposirory, contextTimeout time.Duration) domain.FileUsecase {
	return &fileUsecase{repo: repo, contextTimeout: contextTimeout}
}

func (f *fileUsecase) Upload(c context.Context, file *multipart.FileHeader) (string, error) {
	file.Filename = f.uniqueFileName(file)
	// dst := filepath.Join("file", file.Filename)
	return "", nil
}

func (f *fileUsecase) FileList(c context.Context, fileType, fileSource string, page, size int) ([]domain.File, error) {
	return nil, nil
}

func (f *fileUsecase) uniqueFileName(file *multipart.FileHeader) string {
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
}

func (f *fileUsecase) checkFileMd5(file *multipart.FileHeader) (bool, error) {
	return false, nil
}
