package main

import (
	"github.com/LXJ0000/go-backend/api/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"

	route "github.com/LXJ0000/go-backend/api/route"
	"github.com/LXJ0000/go-backend/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {

	app := bootstrap.App()

	env := app.Env

	db := app.Orm
	cache := app.Cache

	producer := app.Producer
	saramaClient := app.SaramaClient

	timeout := time.Duration(env.ContextTimeout) * time.Second // TODO

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())
	server.Use(middleware.RateLimitMiddleware())
	server.Use(ResponseTimeMiddleware())
	route.Setup(env, timeout, db, cache, server, producer, saramaClient)
	initPrometheus()
	_ = server.Run(env.ServerAddress)
}

func initPrometheus() {
	go func() {
		// 专门给 prometheus 用的端口
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

func ResponseTimeMiddleware() gin.HandlerFunc {
	labels := []string{"method", "pattern", "status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "jannan",
		Subsystem: "go_backend",
		Help:      "response time",
		Name:      "gin_http_resp_time",
		ConstLabels: map[string]string{
			"instance_id": "",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, labels)
	prometheus.MustRegister(vector)
	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Milliseconds()
			method := ctx.Request.Method
			pattern := ctx.FullPath()
			status := ctx.Writer.Status()
			vector.WithLabelValues(method, pattern, strconv.Itoa(status)).
				Observe(float64(duration))
		}()
		ctx.Next()
	}
}
