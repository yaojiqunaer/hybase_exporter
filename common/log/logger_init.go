package log

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var Debug log.Logger
var Info log.Logger
var Warn log.Logger
var Error log.Logger

func init() {
	promLogConfig := &promlog.Config{}
	logger := promlog.New(promLogConfig)
	flag.AddFlags(kingpin.CommandLine, promLogConfig)
	// Stdout support
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	Debug = level.Debug(logger)
	Info = level.Info(logger)
	Warn = level.Warn(logger)
	Error = level.Error(logger)
	fmt.Println("init logger debug|info|warn|error ...")
}
