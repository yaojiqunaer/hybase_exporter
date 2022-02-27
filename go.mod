module hybase_exporter

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/prometheus/client_golang v1.12.1
	//借助gopsutil 采集主机指标
	github.com/shirou/gopsutil v2.21.11+incompatible
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/ugorji/go v1.2.7 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sys v0.0.0-20220224120231-95c6836cb0e7 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

go 1.15
