package memstats

import (
	"time"

	"github.com/go-x-pkg/log"
)

type LoggerType int

type Config struct {
	loggerType LoggerType
	fnPeriod   func() time.Duration
	fnLog      func(log.Level, string, ...interface{})
}

func (it *Config) defaultize() {
	it.loggerType = StringLogger
	it.fnPeriod = defaultFnPeriod
	it.fnLog = defaultFnLog
}

type Arg func(*Config)

func FnLog(v func(log.Level, string, ...interface{})) Arg {
	return func(cfg *Config) { cfg.fnLog = v }
}

func FnPeriod(v func() time.Duration) Arg {
	return func(cfg *Config) { cfg.fnPeriod = v }
}

func FnLoggerType(v LoggerType) Arg {
	return func(cfg *Config) { cfg.loggerType = v }
}
