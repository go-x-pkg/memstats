package memstats

import (
	"time"

	"github.com/go-x-pkg/log"
)

type Config struct {
	fnPeriod func() time.Duration
	fnLog    func(log.Level, string)
}

func (it *Config) defaultize() {
	it.fnPeriod = defaultFnPeriod
	it.fnLog = defaultFnLog
}

type Arg func(*Config)

func FnLog(v func(log.Level, string)) Arg {
	return func(cfg *Config) { cfg.fnLog = v }
}

func FnPeriod(v func() time.Duration) Arg {
	return func(cfg *Config) { cfg.fnPeriod = v }
}
