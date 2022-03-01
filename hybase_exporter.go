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
	"strconv"
	"strings"
)

func main() {
	var (
		metricsPath   = kingpin.Flag("metrics.path", "Path under which to expose metrics.").Default(getEnv("METRICS_PATH", "/metrics")).String()
		hybaseAddress = kingpin.Flag("hybase.address", "Hybase address.").Default(getEnv("HYBASE_ADDRESS", "127.0.0.1:5555")).String()
		hybaseUser    = kingpin.Flag("hybase.user", "Hybase user.").Default(getEnv("HYBASE_USER", "admin")).String()
		_             = kingpin.Flag("hybase.pwd", "Hybase password.").Default(getEnv("HYBASE_PWD", "admin")).String()
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
	info.Log("msg", "hybase user is "+*hybaseUser)
	info.Log("msg", "register hybase status collector")
	hyBaseClient := parseAddress(*hybaseAddress)
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

func getEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if envVal, ok := os.LookupEnv(key); ok {
		envBool, err := strconv.ParseBool(envVal)
		if err == nil {
			return envBool
		}
	}
	return defaultVal
}

func getEnvInt64(key string, defaultVal int64) int64 {
	if envVal, ok := os.LookupEnv(key); ok {
		envInt64, err := strconv.ParseInt(envVal, 10, 64)
		if err == nil {
			return envInt64
		}
	}
	return defaultVal
}

func parseAddress(address string) client.HyBaseClient {
	split := strings.Split(address, ":")
	port, _ := strconv.ParseInt(split[1], 10, 32)
	return client.HyBaseClient{Protocol: "http", Ip: split[0], Port: int(port)}
}
