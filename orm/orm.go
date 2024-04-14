package orm

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database interface {
	FindOne(c context.Context, model interface{}, filter interface{}) (interface{}, error)
	FindMany(c context.Context, model interface{}, filter interface{}, page, size int, items interface{}) error
	InsertOne(c context.Context, model interface{}, item interface{}) (interface{}, error)
	DeleteOne(c context.Context, model interface{}, item interface{}) (interface{}, error)
	UpdateOne(c context.Context, model interface{}, filter interface{}, update interface{}) (interface{}, error)
	Upsert(c context.Context, model interface{}, update map[string]interface{}, create interface{}) error
	Transaction(c context.Context, fn func(tx *gorm.DB) error) error
}

type database struct {
	db *gorm.DB
}

func NewDatabase(db *gorm.DB) Database {
	return &database{db: db}
}

func (dao *database) FindOne(c context.Context, model interface{}, filter interface{}) (interface{}, error) {
	err := dao.db.WithContext(c).Model(model).Where(filter).First(filter).Error
	return filter, err
}

func (dao *database) FindMany(c context.Context, model interface{}, filter interface{}, page, size int, items interface{}) error {
	err := dao.db.WithContext(c).Model(model).Where(filter).Offset((page - 1) * size).Limit(size).Find(items).Error
	return err
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

func (dao *database) Upsert(c context.Context, model interface{}, update map[string]interface{}, create interface{}) error {
	err := dao.db.WithContext(c).Model(model).Clauses(clause.OnConflict{
		//Columns:  // mysql 可以不写
		DoUpdates: clause.Assignments(update),
	}).Create(create).Error
	return err
}

func (dao *database) Transaction(c context.Context, fn func(tx *gorm.DB) error) error {
	err := dao.db.WithContext(c).Transaction(fn)
	return err
}
