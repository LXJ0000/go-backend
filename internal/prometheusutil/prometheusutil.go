package prometheusutil

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func Init(addr string) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(addr, nil)
	}()

}
