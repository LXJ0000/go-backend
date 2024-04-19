package bootstrap

import (
	"github.com/LXJ0000/go-backend/internal/logutil"
	"github.com/LXJ0000/go-backend/internal/snowflakeutil"
	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/event"
	"github.com/LXJ0000/go-backend/orm"
	"github.com/LXJ0000/go-backend/redis"
)

type Application struct {
	Env *Env
	//Mongo mongo.Client
	Orm   orm.Database
	Cache redis.Cache

	Producer event.Producer

	SaramaClient sarama.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	//app.Mongo = NewMongoDatabase(app.Env)
	app.Orm = NewOrmDatabase(app.Env)
	app.Cache = NewRedisCache(app.Env)

	snowflakeutil.Init(app.Env.SnowflakeStartTime, app.Env.SnowflakeMachineID)
	logutil.Init(app.Env.AppEnv)

	app.Producer = NewProducer(app.Env)
	app.SaramaClient = NewSaramaClient(app.Env)

	return *app
}

//func (app *Application) CloseDBConnection() {
//	CloseMongoDBConnection(app.Mongo)
//}
