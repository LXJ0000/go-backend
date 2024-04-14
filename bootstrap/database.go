package bootstrap

import (
	"fmt"
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/orm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func NewOrmDatabase(env *Env) orm.Database {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", env.DBName)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	// In WSL how to connect sqlite ?
	// move go-backend.db to /mnt/c/Users/JANNAN/Desktop/go-backend.db then
	// ln -s /mnt/c/Users/JANNAN/Desktop/go-backend.db ./go-backend.db
	//dsn := "root:root@tcp(127.0.0.1:3306)/go-backend?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Info),
	//})
	if err != nil {
		log.Fatal(err)
	}
	if err = db.AutoMigrate(
		&domain.Post{},
		&domain.User{},
		&domain.Task{},
		&domain.Interaction{},
		&domain.UserLike{},
	); err != nil {
		log.Fatal(err)
	}
	database := orm.NewDatabase(db)

	return database
}
