package memstats

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/go-x-pkg/bufpool"
	"github.com/go-x-pkg/log"
	"go.uber.org/zap"
)

type Worker struct {
	stop chan struct{}
	done chan struct{}
	cfg  Config
}

func (it *Worker) Stop() { it.stop <- struct{}{} }
func (it *Worker) Done() { <-it.done }
func (it *Worker) DoneContext(ctx context.Context) {
	select {
	case <-it.done:
	case <-ctx.Done():
	}
}

func (it *Worker) Start() {
	ctx := context.Background()
	it.startMonitor(ctx)
}

func (it *Worker) StartContext(ctx context.Context) {
	it.startMonitor(ctx)
}

func (it *Worker) startMonitor(ctx context.Context) {
	if ctx == nil { // parent-context
		ctx = context.TODO()
	}

	it.cfg.fnLog(log.Info, "started")

	defer func() {
		it.cfg.fnLog(log.Info, "finished")

		select {
		case it.done <- struct{}{}:
		default:
		}
	}()

	period := it.cfg.fnPeriod()

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-it.stop:
			it.cfg.fnLog(log.Info, "stoped")
			return
		case <-ctx.Done():
			it.cfg.fnLog(log.Info, "stoped")
			return
		case <-ticker.C:
			it.Perform()
		}
	}
}

func (it *Worker) Perform() {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	if it.cfg.loggerType == ZapLogger {
		it.cfg.fnLog(log.Info, "memstat statistics",
			zap.Int("gorutines", runtime.NumGoroutine()),
			zap.Uint32("numGC", memStats.NumGC),
			zap.Uint64("alloc", memStats.Alloc),
			zap.Uint64("mallocs", memStats.Mallocs),
			zap.Uint64("frees", memStats.Frees),
			zap.Uint64("heapAlloc", memStats.HeapAlloc),
			zap.Uint64("stackInuse", memStats.StackInuse))
		return
	}

	buf := bufpool.NewBuf()
	defer buf.Release()

	buf.WriteString(fmt.Sprintf("(:gorutines %d :num-gc %d :alloc %d :mallocs %d :frees %d :heap-alloc %d :stack-inuse %d)",
		runtime.NumGoroutine(),
		memStats.NumGC,
		memStats.Alloc,
		memStats.Mallocs,
		memStats.Frees,
		memStats.HeapAlloc,
		memStats.StackInuse,
	))

	it.cfg.fnLog(log.Info, buf.String())
}

func (it *Worker) Initialize(fnArgs ...Arg) {
	it.stop = make(chan struct{})
	it.done = make(chan struct{}, 1)

	it.cfg.defaultize()

	for _, fn := range fnArgs {
		fn(&it.cfg)
	}
}
