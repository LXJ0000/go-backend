package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type Config struct {
	Namespace  string
	Subsystem  string
	Name       string
	Help       string
	InstanceID string
}

func DefaultConfig() Config {
	return Config{
		Namespace:  "lxj0000",
		Subsystem:  "go_backend",
		Name:       "gin_http",
		Help:       "Statistics gin http interface",
		InstanceID: "1",
	}
}

func New(config Config) gin.HandlerFunc {
	label := []string{"method", "patter", "status"}
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: config.Namespace,
		Subsystem: config.Subsystem,
		Name:      config.Name + "_response_time",
		Help:      config.Help,
		ConstLabels: map[string]string{
			"instance_id": config.InstanceID,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, label)
	prometheus.MustRegister(summary)
	return func(context *gin.Context) {
		begin := time.Now()
		defer func() {
			duration := time.Since(begin)
			pattern := context.Request.Method
			if pattern == "" {
				pattern = "UnKnow"
			}
			summary.WithLabelValues(
				context.Request.Method, pattern, strconv.Itoa(context.Writer.Status()),
			).Observe(float64(duration))
		}()
		context.Next()
	}
}

func PrometheusMiddleware() gin.HandlerFunc {
	return New(DefaultConfig())
}
