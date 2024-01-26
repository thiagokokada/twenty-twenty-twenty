package sound

import (
	"sync"
	"sync/atomic"
)

// WaitGroup with count, to make it easier to visualise the wg state in logs
// WARN: please don't rely in count for any business logic, since the state of
// the count may have changed when the wg itself is called
// See: https://stackoverflow.com/a/68995552
type wgCount struct {
	sync.WaitGroup
	c atomic.Int64
}

func (wg *wgCount) add(delta int) {
	wg.c.Add(int64(delta))
	wg.WaitGroup.Add(delta)
}

func (wg *wgCount) done() {
	wg.c.Add(-1)
	wg.WaitGroup.Done()
}

func (wg *wgCount) count() int {
	return int(wg.c.Load())
}
