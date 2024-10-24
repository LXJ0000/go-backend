package bootstrap

import (
	"log"
	"time"

	domain "github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/pkg/orm"
	_prometheus "github.com/prometheus/client_golang/prometheus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
)

func NewOrmDatabase(env *Env) orm.Database {
	//db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", env.DBName)), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Info),
	//})
	// In WSL how to connect sqlite ?
	// move go-backend.db to /mnt/c/Users/JANNAN/Desktop/go-backend.db then
	// ln -s /mnt/c/Users/JANNAN/Desktop/go-backend.db ./go-backend.db
	//dsn := "root:root@tcp(127.0.0.1:3306)/go-backend?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(env.MySQLAddress), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	initPrometheus(db)

	if err = db.AutoMigrate(
		&domain.Post{},
		&domain.User{},
		&domain.Task{},
		&domain.Interaction{},
		&domain.UserLike{},
		&domain.UserCollect{},
		&domain.Comment{},
		&domain.Relation{},
		&domain.File{},
		&domain.TagBiz{},
		&domain.Tag{},
		&domain.FeedPull{},
		&domain.FeedPush{},
	); err != nil {
		log.Fatal(err)
	}
	database := orm.NewDatabase(db)

	return database
}

func initPrometheus(db *gorm.DB) {
	if err := db.Use(prometheus.New(prometheus.Config{
		DBName:          "go_backend",
		RefreshInterval: 600, // TODO: smaller interval
		StartServer:     false,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		},
	})); err != nil {
		log.Fatal(err)
	}
	// 监控查询时间
	before := func(db *gorm.DB) {
		begin := time.Now()
		db.Set("begin", begin)
	}
	if err := db.Callback().Create().Before("*").Register("prometheus_create_before", before); err != nil {
		log.Fatal(err)
	}
	if err := db.Callback().Query().Before("*").Register("prometheus_query_before", before); err != nil {
		log.Fatal(err)
	}
	if err := db.Callback().Delete().Before("*").Register("prometheus_delete_before", before); err != nil {
		log.Fatal(err)
	}
	if err := db.Callback().Update().Before("*").Register("prometheus_update_before", before); err != nil {
		log.Fatal(err)
	}
	if err := db.Callback().Row().Before("*").Register("prometheus_raw_before", before); err != nil {
		log.Fatal(err)
	}
	if err := db.Callback().Raw().Before("*").Register("prometheus_row_before", before); err != nil {
		log.Fatal(err)
	}

	vector := _prometheus.NewSummaryVec(_prometheus.SummaryOpts{
		Namespace: "lxj0000",
		Subsystem: "go_backend",
		Name:      "gorm_response_time",
		Help:      "Statistics gorm interface",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"type", "table"})
	after := func(typeName string) func(db *gorm.DB) {
		return func(db *gorm.DB) {
			raw, _ := db.Get("begin")
			if begin, ok := raw.(time.Time); ok {
				table := db.Statement.Table
				if table == "" {
					table = "UnKnow"
				}
				vector.WithLabelValues(typeName, table).Observe(float64(time.Since(begin).Milliseconds()))
			}
		}
	}
	if err := db.Callback().Create().After("*").Register("prometheus_create_after", after("Create")); err != nil {
		log.Fatal(err)
	}
	if err := db.Callback().Query().After("*").Register("prometheus_query_after", after("Query")); err != nil {
		log.Fatal(err)
	}
	if err := db.Callback().Delete().After("*").Register("prometheus_delete_after", after("Delete")); err != nil {
		log.Fatal(err)
	}
	if err := db.Callback().Update().After("*").Register("prometheus_update_after", after("Update")); err != nil {
		log.Fatal(err)
	}
	if err := db.Callback().Raw().After("*").Register("prometheus_raw_after", after("Raw")); err != nil {
		log.Fatal(err)
	}
	if err := db.Callback().Row().After("*").Register("prometheus_row_after", after("Row")); err != nil {
		log.Fatal(err)
	}
	_prometheus.MustRegister(vector)
}
