package repository

import (
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"golang.org/x/net/context"
)

type fileReposirory struct {
	dao orm.Database
}

func NewFileReposirory(dao orm.Database) domain.FileReposirory {
	return &fileReposirory{dao: dao}
}

func (f *fileReposirory) Upload(c context.Context, file domain.File) error {
	return f.dao.InsertOne(c, &domain.File{}, &file)
}

func (f *fileReposirory) FileList(c context.Context, fileType, fileSource string, page, size int) ([]domain.File, error) {
	filter := map[string]interface{}{
		"type":   fileType,
		"source": fileSource,
	}
	var items []domain.File
	err := f.dao.FindPage(c, &domain.File{}, filter, page, size, &items)
	return items, err
}
