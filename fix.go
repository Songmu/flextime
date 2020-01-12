package flextime

import (
	"sync"
	"time"
)

func newFixedClock(t time.Time) Clock {
	return NewFakeClock(&fixedNS{t: t})
}

type fixedNS struct {
	t  time.Time
	mu sync.RWMutex
}

var _ NowSleeper = (*fixedNS)(nil)

func (fi *fixedNS) Now() time.Time {
	fi.mu.RLock()
	defer fi.mu.RUnlock()
	return fi.t
}

func (fi *fixedNS) Sleep(d time.Duration) {
	if d <= 0 {
		return
	}
	fi.mu.Lock()
	defer fi.mu.Unlock()
	fi.t = fi.t.Add(d)
}
