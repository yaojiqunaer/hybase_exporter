// Package collector
package collector

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"hybase_exporter/client"
	"hybase_exporter/common/log"
	"strconv"
	"sync"
)

// HyBaseCollector hybase struct implements github.com/prometheus/client_golang@v1.12.1/prometheus/collector.go:27
type HyBaseCollector struct {
	trsHybaseServerInfo   *prometheus.Desc // Gauge
	trsHybaseServerStatus *prometheus.Desc // Gauge
	trsHybaseCpuUsage     *prometheus.Desc // Gauge
	trsHybaseDiskUsage    *prometheus.Desc // Gauge
	trsHybaseDiskSize     *prometheus.Desc // Gauge
	trsHybaseDiskFree     *prometheus.Desc // Gauge
	trsHybaseDiskIoOut    *prometheus.Desc // Gauge
	trsHybaseDiskIoIn     *prometheus.Desc // Gauge
	trsHybaseNicByteOut   *prometheus.Desc // Gauge
	trsHybaseNicByteIn    *prometheus.Desc // Gauge
	trsHybaseSearchRate   *prometheus.Desc // Gauge
	mutex                 sync.Mutex       // sync
}

var hbCli client.HyBaseClient

// Describe returns all descriptions of the collector.
//实现采集器Describe接口
func (n *HyBaseCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.trsHybaseServerInfo
	ch <- n.trsHybaseServerStatus
}

// Collect returns the current state of all metrics of the collector.
// 实现采集器Collect接口,真正采集动作
func (n *HyBaseCollector) Collect(ch chan<- prometheus.Metric) {
	n.mutex.Lock()
	publicStatus := client.GetPublicStatus(hbCli)
	bytes, _ := json.Marshal(publicStatus)
	if publicStatus.Code != 0 {
		log.Error.Log("msg", "code not match, the response is "+string(bytes))
		return
	}
	data := publicStatus.Data
	if len(data) == 0 {
		log.Error.Log("msg", "hybase data size is 0")
		return
	}
	var node = data[0]
	var ip string
	var id string
	var sysName string
	if len(data) > 1 {
		log.Warn.Log("msg", "current version only get one node info to metrics")
	}
	ip = node.Ip
	id = node.Id
	sysName = node.Sysname
	ch <- prometheus.MustNewConstMetric(n.trsHybaseServerInfo, prometheus.GaugeValue, 1, ip, sysName, id)
	var items = node.Itemlist
	for _, item := range items {
		key := item.Key
		value := item.Value
		if key == "server_status" {
			ch <- prometheus.MustNewConstMetric(n.trsHybaseServerStatus, prometheus.GaugeValue, 1, ip, item.Value)
		}
		if key == "cpu_usage" {
			floatV, _ := strconv.ParseFloat(value, 64)
			ch <- prometheus.MustNewConstMetric(n.trsHybaseCpuUsage, prometheus.GaugeValue, floatV, ip)
		}
		if key == "disk_usage" {
			floatV, _ := strconv.ParseFloat(value, 64)
			ch <- prometheus.MustNewConstMetric(n.trsHybaseDiskUsage, prometheus.GaugeValue, floatV, ip)
		}
		if key == "disk_size" {
			floatV, _ := strconv.ParseFloat(value, 64)
			ch <- prometheus.MustNewConstMetric(n.trsHybaseDiskSize, prometheus.GaugeValue, floatV, ip)
		}
		if key == "disk_free" {
			floatV, _ := strconv.ParseFloat(value, 64)
			ch <- prometheus.MustNewConstMetric(n.trsHybaseDiskFree, prometheus.GaugeValue, floatV, ip)
		}
		if key == "disk_io_out" {
			floatV, _ := strconv.ParseFloat(value, 32)
			ch <- prometheus.MustNewConstMetric(n.trsHybaseDiskIoOut, prometheus.GaugeValue, floatV, ip)
		}
		if key == "disk_io_in" {
			floatV, _ := strconv.ParseFloat(value, 32)
			ch <- prometheus.MustNewConstMetric(n.trsHybaseDiskIoIn, prometheus.GaugeValue, floatV, ip)
		}
		if key == "nic_byte_in" {
			floatV, _ := strconv.ParseFloat(value, 32)
			ch <- prometheus.MustNewConstMetric(n.trsHybaseNicByteIn, prometheus.GaugeValue, floatV, ip)
		}
		if key == "nic_byte_out" {
			floatV, _ := strconv.ParseFloat(value, 32)
			ch <- prometheus.MustNewConstMetric(n.trsHybaseNicByteOut, prometheus.GaugeValue, floatV, ip)
		}

	}
	n.mutex.Unlock()
}

// NewHyBaseStatusCollector Collector
func NewHyBaseStatusCollector(client *client.HyBaseClient) prometheus.Collector {
	hbCli = *client
	collector := HyBaseCollector{
		trsHybaseServerInfo:   prometheus.NewDesc("trs_hybase_server_info", "海贝节点信息", []string{"ip", "sysname", "id"}, nil),
		trsHybaseServerStatus: prometheus.NewDesc("trs_hybase_server_status", "海贝节点状态，1正常，0异常", []string{"ip", "status"}, nil),
		trsHybaseCpuUsage:     prometheus.NewDesc("trs_hybase_cpu_usage", "CPU使用率", []string{"ip"}, nil),
		trsHybaseDiskUsage:    prometheus.NewDesc("trs_hybase_disk_usage", "磁盘使用率", []string{"ip"}, nil),
		trsHybaseDiskSize:     prometheus.NewDesc("trs_hybase_disk_size", "磁盘总空间 Byte", []string{"ip"}, nil),
		trsHybaseDiskFree:     prometheus.NewDesc("trs_hybase_disk_free", "磁盘剩余空间 Byte", []string{"ip"}, nil),
		trsHybaseDiskIoOut:    prometheus.NewDesc("trs_hybase_disk_io_out", "磁盘读速率 Byte/s", []string{"ip"}, nil),
		trsHybaseDiskIoIn:     prometheus.NewDesc("trs_hybase_disk_io_in", "磁盘写速率 Byte/s", []string{"ip"}, nil),
		trsHybaseNicByteOut:   prometheus.NewDesc("trs_hybase_nic_byte_out", "网卡上行速率 Byte/s", []string{"ip"}, nil),
		trsHybaseNicByteIn:    prometheus.NewDesc("trs_hybase_nic_byte_in", "网卡下行速率 Byte/s", []string{"ip"}, nil),
		trsHybaseSearchRate:   prometheus.NewDesc("trs_hybase_search_rate", "检索效率 次/s", []string{"ip"}, nil),
	}
	return &collector
}
