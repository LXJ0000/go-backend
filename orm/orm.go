package orm

import (
	"context"
	"gorm.io/gorm"
)

type Database interface {
	FindOne(c context.Context, model interface{}, item interface{}) (interface{}, error)
	InsertOne(c context.Context, model interface{}, item interface{}) (interface{}, error)
	DeleteOne(c context.Context, model interface{}, item interface{}) (interface{}, error)
	UpdateOne(c context.Context, model interface{}, filter interface{}, update interface{}) (interface{}, error)
}

type database struct {
	db *gorm.DB
}

func NewDatabase(db *gorm.DB) Database {
	return &database{db: db}
}

func (dao *database) FindOne(c context.Context, model interface{}, item interface{}) (interface{}, error) {
	err := dao.db.WithContext(c).Model(model).First(item).Error
	return item, err
}

func (dao *database) InsertOne(c context.Context, model interface{}, item interface{}) (interface{}, error) {
	err := dao.db.WithContext(c).Model(model).Create(item).Error
	return nil, err
}

func (dao *database) DeleteOne(c context.Context, model interface{}, item interface{}) (interface{}, error) {
	err := dao.db.WithContext(c).Model(model).Where(item).Delete(item).Error
	return nil, err
}

func (dao *database) UpdateOne(c context.Context, model interface{}, filter interface{}, update interface{}) (interface{}, error) {
	err := dao.db.WithContext(c).Model(model).Where(filter).Updates(update).Error
	return nil, err
}
