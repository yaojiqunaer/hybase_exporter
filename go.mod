module hybase_exporter

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/prometheus/client_golang v1.12.1
	//借助gopsutil 采集主机指标
	github.com/shirou/gopsutil v2.21.11+incompatible
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
)

go 1.15
