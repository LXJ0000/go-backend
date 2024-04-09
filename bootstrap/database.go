package bootstrap

import (
	"context"
	"fmt"
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/orm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"time"

	"github.com/LXJ0000/go-backend/mongo"
)

func NewMongoDatabase(env *Env) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass

	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)

	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	client, err := mongo.NewClient(mongodbURI)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CloseMongoDBConnection(client mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Connection to MongoDB closed.")
}

func NewOrmDatabase(env *Env) orm.Database {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", env.DBName)), &gorm.Config{})
	//dsn := "root:root@tcp(127.0.0.1:3306)/go-backend?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err = db.AutoMigrate(&domain.Post{}); err != nil {
		log.Fatal(err)
	}
	db.Debug()
	database := orm.NewDatabase(db)

	return database
}
