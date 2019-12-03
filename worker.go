package memstats

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/go-x-pkg/bufpool"
	"github.com/go-x-pkg/log"
)

type Worker struct {
	stop chan struct{}
	done chan struct{}
	cfg  Config
}

func (it *Worker) Stop() { it.stop <- struct{}{} }
func (it *Worker) Done() { <-it.done }
func (it *Worker) DoneWithContext(ctx context.Context) {
	select {
	case <-it.done:
	case <-ctx.Done():
	}
}

func (it *Worker) Start() {
	ctx := context.Background()
	it.start(ctx)
}

func (it *Worker) StartWithCtx(ctx context.Context) {
	it.start(ctx)
}

func (it *Worker) start(ctx context.Context) {
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
