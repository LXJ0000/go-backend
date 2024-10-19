package bootstrap

import (
	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/internal/event"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"github.com/LXJ0000/go-backend/utils/logutil"
	"github.com/LXJ0000/go-backend/utils/prometheusutil"
	"github.com/LXJ0000/go-backend/utils/snowflakeutil"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/robfig/cron/v3"
)

type Application struct {
	Env *Env
	//Mongo mongo.Client
	Orm        orm.Database
	Cache      cache.RedisCache
	LocalCache cache.LocalCache

	Producer event.Producer

	SaramaClient sarama.Client

	Cron *cron.Cron

	SMSAliyunClient *sms.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	//app.Mongo = NewMongoDatabase(app.Env)
	app.Orm = NewOrmDatabase(app.Env)
	app.Cache = NewRedisCache(app.Env)
	// app.LocalCache = NewLocalCache(app.Env)
	logutil.Init(app.Env.AppEnv)
	snowflakeutil.Init(app.Env.SnowflakeStartTime, app.Env.SnowflakeMachineID)
	prometheusutil.Init(app.Env.PrometheusAddress)

	app.Producer = NewProducer(app.Env)
	app.SaramaClient = NewSaramaClient(app.Env)

	app.Cron = NewCron(app.LocalCache, app.Cache, app.Orm)

	app.SMSAliyunClient = NewAliyunClient(app.Env)

	return *app
}

//func (app *Application) CloseDBConnection() {
//	CloseMongoDBConnection(app.Mongo)
//}
