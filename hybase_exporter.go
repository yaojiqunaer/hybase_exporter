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
	"hybase_exporter/client"
	"hybase_exporter/collector"
	"net/http"
	"os"
)

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
	info.Log("msg", "register hybase status collector")
	hyBaseClient := client.HyBaseClient{
		Protocol: "http",
		Ip:       "127.0.0.1",
		Port:     5555,
	}
	http.Handle(*metricsPath, hybaseHandler(&hyBaseClient))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		info.Log("client", r.RemoteAddr, "uri", r.RequestURI)
		w.Write(index(*metricsPath))
	})
	err := web.ListenAndServe(&http.Server{Addr: ":9555"}, "", logger)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}

//hybaseHandler hybase collector handler, this can remove the metrics those start with go_*, such as go_gc_*
func hybaseHandler(client *client.HyBaseClient) http.Handler {
	// disable go language build env metrics, so use NewRegistry instead of prometheus.Register()
	r := prometheus.NewRegistry()
	// registry self define collector
	r.MustRegister(collector.NewHyBaseStatusCollector(client))
	// [...] update metrics within a goroutine
	return promhttp.HandlerFor(r, promhttp.HandlerOpts{})
}

//index 未知地址 首页内容
func index(metricPath string) []byte {
	return []byte(`<html>
			<head><title>HyBase Exporter</title></head>
			<body>
			<h1>HyBase Exporter</h1>
			<p><a href="` + metricPath + `">Metrics</a></p>
			</body>
			</html>`)
}
