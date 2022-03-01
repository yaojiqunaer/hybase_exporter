package main

import (
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"hybase_exporter/collector"
	"net/http"
	"os"
)

func init() {
	//注册自身采集器
	prometheus.Register(collector.NewHyBaseStatusCollector())
	prometheus.MustRegister()
}

func main() {
	var (
		metricsPath   = kingpin.Flag("metrics_path", "Path under which to expose metrics.").Default("/metrics").String()
		hybaseAddress = kingpin.Flag("HYBASE_ADDRESS", "Hybase address.").Default("127.0.0.1:5555").String()
	)
	promLogConfig := &promlog.Config{}
	logger := promlog.New(promLogConfig)
	flag.AddFlags(kingpin.CommandLine, promLogConfig)
	kingpin.Version(version.Print("hybase_exporter"))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	info := level.Info(logger)
	info.Log("msg", "Starting hybase_exporter", "version", version.Info())
	info.Log("msg", "metrics path is "+*metricsPath)
	info.Log("msg", "hybase address is "+*hybaseAddress)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		info.Log("client", r.RemoteAddr, "uri", r.RequestURI)
		w.Write([]byte(`<html>
			<head><title>HyBase Exporter</title></head>
			<body>
			<h1>HyBase Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	http.Handle(*metricsPath, promhttp.Handler())
	err := web.ListenAndServe(&http.Server{Addr: ":9555"}, "", logger)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
