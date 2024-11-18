package bootstrap

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/LXJ0000/go-backend/internal/event"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/chat"
	"github.com/LXJ0000/go-backend/pkg/file"
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

	MinioClient file.FileStorage

	DoubaoChat chat.Chat
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
	prometheusutil.Init(app.Env.PrometheusAddr)

	app.Producer = NewProducer(app.Env)
	app.SaramaClient = NewSaramaClient(app.Env)

	app.Cron = NewCron(app.LocalCache, app.Cache, app.Orm)

	app.SMSAliyunClient = NewAliyunClient(app.Env)

	app.MinioClient = NewMinio()

	app.DoubaoChat = NewDoubaoChat()

	return *app
}

//func (app *Application) CloseDBConnection() {
//	CloseMongoDBConnection(app.Mongo)
//}

func init() {
	fmt.Println("=====================================")
	fmt.Println("若要使用 sms 服务，请配置好环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	fmt.Println("具体配置请联系项目管理员")
	fmt.Println("=====================================")
}
