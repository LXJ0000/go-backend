package repository

import (
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type fileRepository struct {
	dao orm.Database
}

func NewFileRepository(dao orm.Database) domain.FileRepository {
	return &fileRepository{dao: dao}
}

func (f *fileRepository) Upload(c context.Context, file domain.File) error {
	return f.dao.InsertOne(c, &domain.File{}, &file)
}

func (f *fileRepository) FileList(c context.Context, fileType, fileSource string, page, size int) ([]domain.File, int, error) {
	filter := map[string]interface{}{}
	if fileType != "" {
		filter["type"] = fileType
	}
	if fileSource != "" {
		filter["source"] = fileSource
	}
	var (
		items []domain.File
		cnt   int64
	)
	g := errgroup.Group{}
	g.Go(func() error {
		return f.dao.FindPage(c, &domain.File{}, filter, page, size, &items)
	})
	g.Go(func() error {
		var err error
		cnt, err = f.dao.Count(c, &domain.File{}, filter)
		return err
	})
	if err := g.Wait(); err != nil {
		return nil, 0, err
	}
	return items, int(cnt), nil
}

func (f *fileRepository) FindByHash(c context.Context, hash string) (domain.File, error) {
	filter := map[string]interface{}{
		"hash": hash,
	}
	var item domain.File
	// TODO cache
	err := f.dao.FindOne(c, &domain.File{}, filter, &item)
	return item, err
}
