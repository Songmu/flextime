package flextime

import (
	"sync"
	"time"
)

func newOffsetClock(t time.Time) Clock {
	return NewFakeClock(&offsetNS{
		offset: t.Sub(time.Now()),
		loc:    t.Location(),
	})
}

type offsetNS struct {
	offset time.Duration
	loc    *time.Location

	mu sync.Mutex
}

var _ NowSleeper = (*offsetNS)(nil)

func (oc *offsetNS) Now() time.Time {
	return time.Now().Add(oc.offset).In(oc.loc)
}

func (oc *offsetNS) Sleep(d time.Duration) {
	if d <= 0 {
		return
	}
	oc.mu.Lock()
	defer oc.mu.Unlock()
	oc.offset = oc.offset + d
}
