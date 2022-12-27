package memstats

import (
	"time"

	"github.com/go-x-pkg/log"
)

const (
	StringLogger LoggerType = iota
	ZapLogger
)

var (
	defaultFnPeriod = func() time.Duration { return 60 * time.Second }
	defaultFnLog    = log.LogStd
)
