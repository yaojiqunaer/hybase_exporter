package main

import (
	"fmt"
	"hybase_exporter/collector"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init()  {
     //注册自身采集器
    prometheus.MustRegister(collector.NewNodeCollector())
}
func main() {
    http.Handle("/metrics", promhttp.Handler())
    if err := http.ListenAndServe(":9555", nil); err != nil {
        fmt.Printf("Error occur when start server %v", err)
    }
}