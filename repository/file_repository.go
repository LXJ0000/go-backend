package repository

import (
	"context"
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/orm"
	"mime/multipart"
)

type fileRepository struct {
	dao orm.Database
}

func NewFileRepository(dao orm.Database) domain.FileRepository {
	return &fileRepository{dao: dao}
}

func (f fileRepository) UploadOne(c context.Context, file *multipart.FileHeader) (domain.File, error) {
	//TODO implement me
	panic("implement me")
}

func (f fileRepository) UploadMany(c context.Context, files []*multipart.FileHeader) ([]domain.File, error) {
	//TODO implement me
	panic("implement me")
}
