package config

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var Logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.Kitchen,
	Prefix:          "Referal system 🍪 ",
})

func LogFatal(s string, data ...any) {
	Logger.Fatal(s, data)
}
func LogInfo(s string, data ...any) {
	Logger.Info(s, data)
}
func LogDebug(s string, data ...any) {
	Logger.Debug(s, data)
}
