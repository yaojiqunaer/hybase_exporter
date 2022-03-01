// Package collector
package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

//HyBaseCollector hybase struct implements github.com/prometheus/client_golang@v1.12.1/prometheus/collector.go:27
type HyBaseCollector struct {
	trsHybaseServerInfo *prometheus.Desc //Gauge
	mutex               sync.Mutex
}

func NewHyBaseStatusCollector() prometheus.Collector {
	collector := HyBaseCollector{
		trsHybaseServerInfo: prometheus.NewDesc(
			"trs_hybase_server_info",
			"trs_hybase_server_info hybase信息",
			[]string{"ip", "sysname", "id"},
			nil),
	}
	return &collector
}

// Describe returns all descriptions of the collector.
//实现采集器Describe接口
func (n *HyBaseCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.trsHybaseServerInfo
}

// Collect returns the current state of all metrics of the collector.
//实现采集器Collect接口,真正采集动作
func (n *HyBaseCollector) Collect(ch chan<- prometheus.Metric) {
	n.mutex.Lock()
	ch <- prometheus.MustNewConstMetric(n.trsHybaseServerInfo, prometheus.GaugeValue, 0, "127.0.0.1", "hybase1", "hybase1")
	n.mutex.Unlock()
}
