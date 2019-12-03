package memstats

import (
	"time"

	"github.com/go-x-pkg/log"
)

var (
	defaultFnPeriod = func() time.Duration { return 60 * time.Second }
	defaultFnLog    = log.LogStd
)
