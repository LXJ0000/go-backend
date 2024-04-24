package bootstrap

import (
	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/cache"
	"github.com/LXJ0000/go-backend/event"
	"github.com/LXJ0000/go-backend/internal/logutil"
	"github.com/LXJ0000/go-backend/internal/prometheusutil"
	"github.com/LXJ0000/go-backend/internal/snowflakeutil"
	"github.com/LXJ0000/go-backend/orm"
	"github.com/robfig/cron/v3"
)

type Application struct {
	Env *Env
	//Mongo mongo.Client
	Orm   orm.Database
	Cache cache.Cache

	Producer event.Producer

	SaramaClient sarama.Client

	Cron *cron.Cron
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	//app.Mongo = NewMongoDatabase(app.Env)
	app.Orm = NewOrmDatabase(app.Env)
	app.Cache = NewRedisCache(app.Env)

	logutil.Init(app.Env.AppEnv)
	snowflakeutil.Init(app.Env.SnowflakeStartTime, app.Env.SnowflakeMachineID)
	prometheusutil.Init(app.Env.PrometheusAddress)

	app.Producer = NewProducer(app.Env)
	app.SaramaClient = NewSaramaClient(app.Env)

	return *app
}

//func (app *Application) CloseDBConnection() {
//	CloseMongoDBConnection(app.Mongo)
//}
