package service

import (
	"log"
	"os"

	"github.com/hashicorp/go-hclog"
)

type GranularityLevel string

const (
	Trace GranularityLevel = "[TRACE] "
	Debug GranularityLevel = "[DEBUG] "
	Info  GranularityLevel = "[INFO] "
	Warn  GranularityLevel = "[WARN] "
	Error GranularityLevel = "[ERROR] "
)

// Since go-plugin handle logs as is described on https://github.com/hashicorp/go-plugin#features under (Built-in Logging)
// is necessary to send the logs as JSON format or with a prefix as is show here
// https://github.com/hashicorp/go-plugin/blob/master/client.go#L1007

// NewLoggerHelper Create a Logger with JSON format allowing user to have a more granular support for example
// NewLoggerHelper().Info("octant-sample-plugin is starting")
func NewLoggerHelper() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Output:     os.Stderr,
		JSONFormat: true,
	})
}

// SetupPluginLogger Set prefix so go-plugin recognize the level of granularity and applied to the host process
// Unfortunately logs go through hclog and the through Zap, for this reason it looks like this
// "2021-08-23T11:35:02.714-0500    INFO    octant-sample-plugin    plugin/logger.go:43     [INFO] 2021/08/23 11:35:02 ..."
func SetupPluginLogger(g GranularityLevel) {
	log.SetPrefix(string(g))
}
