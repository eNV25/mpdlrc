package timerpool

import (
	"sync"
	"time"
)

// Get returns a timer for the given duration d from the pool.
func Get(d time.Duration) *time.Timer {
	if v := timerPool.Get(); v != nil {
		timer := v.(*time.Timer)
		timer.Reset(d)
		return timer
	}
	return time.NewTimer(d)
}

// Put returns t to the pool. Set consumed to true if you have received from C.
//
// timer cannot be accessed after returning to the pool.
func Put(timer *time.Timer, consumed bool) {
	if !consumed && !timer.Stop() {
		<-timer.C
	}
	timerPool.Put(timer)
}

var timerPool sync.Pool
