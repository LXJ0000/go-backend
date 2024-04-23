package usecase

import (
	"context"
	"github.com/LXJ0000/go-backend/domain"
	"mime/multipart"
)

type fileUsecase struct {
	repo domain.FileRepository
}

func NewFileUsecase(repo domain.FileRepository) domain.FileUsecase {
	return &fileUsecase{repo: repo}
}

func (f fileUsecase) UploadOne(c context.Context, file *multipart.FileHeader) (domain.File, error) {
	//TODO implement me
	panic("implement me")
}

func (f fileUsecase) UploadMany(c context.Context, files []*multipart.FileHeader) ([]domain.File, error) {
	//TODO implement me
	panic("implement me")
}
