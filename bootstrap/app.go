package bootstrap

import (
	snowflake "github.com/LXJ0000/go-backend/internal/snowflakeutil"
	"github.com/LXJ0000/go-backend/orm"
)

type Application struct {
	Env *Env
	//Mongo mongo.Client
	Orm orm.Database
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	//app.Mongo = NewMongoDatabase(app.Env)
	app.Orm = NewOrmDatabase(app.Env)
	snowflake.Init("2023-01-01", 1)
	return *app
}

//func (app *Application) CloseDBConnection() {
//	CloseMongoDBConnection(app.Mongo)
//}
