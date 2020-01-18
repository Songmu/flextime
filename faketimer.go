package flextime

import (
	"sync"
	"sync/atomic"
	"time"
)

type fakeTimer struct {
	T        NowSleeper
	IsTicker bool
	fun      func()

	resetMu sync.Mutex

	once       sync.Once
	ch, inch   chan time.Time
	active     int32
	stop, done chan struct{}
	triggerAt  time.Time
}

var _ timerIface = (*fakeTimer)(nil)

func newFakeTimer(c NowSleeper, d time.Duration, f func()) *fakeTimer {
	fti := &fakeTimer{
		T:    c,
		ch:   make(chan time.Time),
		inch: make(chan time.Time),
		stop: make(chan struct{}, 1),
		fun:  f,
	}
	fti.Reset(d)
	return fti
}

func (fti *fakeTimer) isActive() bool {
	return atomic.LoadInt32(&fti.active) > 0
}

func (fti *fakeTimer) C() <-chan time.Time {
	return fti.ch
}

func (fti *fakeTimer) send() {
	fti.once.Do(func() {
		go func() {
			for t := range fti.inch {
				if fti.fun != nil {
					go fti.fun()
				} else {
					fti.ch <- t
				}
			}
		}()
	})
	atomic.StoreInt32(&fti.active, 1)
	go func() {
		select {
		case fti.inch <- func() time.Time {
			fti.T.Sleep(fti.triggerAt.Sub(fti.T.Now()))
			return fti.triggerAt
		}():
		case <-fti.stop:
		}
		atomic.StoreInt32(&fti.active, 0)
		close(fti.done)
	}()
}

func (fti *fakeTimer) Reset(d time.Duration) bool {
	fti.resetMu.Lock()
	defer fti.resetMu.Unlock()
	if d < 0 {
		d = 0
	}
	active := fti.Stop()
	fti.done = make(chan struct{})
	if fti.IsTicker && !fti.triggerAt.IsZero() {
		// to keep base time
		now := fti.T.Now()
		nextDur := d - (now.Sub(fti.triggerAt) % d)
		fti.triggerAt = now.Add(nextDur)
	} else {
		fti.triggerAt = fti.T.Now().Add(d)
	}
	fti.send()
	return active
}

func (fti *fakeTimer) Stop() bool {
	active := fti.isActive()
	if active {
		fti.stop <- struct{}{}
		<-fti.done
		select {
		case <-fti.stop:
		default:
		}
	}
	return active
}
