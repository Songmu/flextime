package flextime

import (
	"sync"
	"time"
)

type fakeTimer struct {
	T        NowSleeper
	IsTicker bool
	fun      func()

	resetMu sync.Mutex

	ch, inch   chan time.Time
	stop, done chan struct{}
	triggerAt  time.Time
}

var _ timerIface = (*fakeTimer)(nil)

func newFakeTimer(c NowSleeper, d time.Duration, f func()) *fakeTimer {
	fti := &fakeTimer{
		T:    c,
		ch:   make(chan time.Time, 1),
		inch: make(chan time.Time),
		stop: make(chan struct{}, 1),
		fun:  f,
	}
	fti.Reset(d)
	return fti
}

func (fti *fakeTimer) isActive() bool {
	if fti.done == nil {
		return false
	}
	select {
	case <-fti.done:
		return false
	default:
		return true
	}
}

func (fti *fakeTimer) C() <-chan time.Time {
	return fti.ch
}

// The `send` is called only inside `Reset` and exclusive control is performed on the `Reset` side,
// so the `send` itself need not do exclusive control.
func (fti *fakeTimer) send() {
	fti.done = make(chan struct{})

	go func() {
		select {
		case t := <-fti.inch:
			if fti.fun != nil {
				go fti.fun()
			} else {
				fti.ch <- t
			}
		case <-fti.done:
		}
	}()

	go func() {
		select {
		case fti.inch <- func() time.Time {
			fti.T.Sleep(fti.triggerAt.Sub(fti.T.Now()))
			return fti.triggerAt
		}():
		case <-fti.stop:
		}
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
	// If multiple `Reset` are called concurrently, this termination process would run at the same time
	// and it returns `true` for each call, but it is no problem because time.Timer of the standard package
	// behaves like that.
	if active {
		fti.stop <- struct{}{}
		<-fti.done
		// The Timer may be fired at the same timing as the Stop. Also, the multiple `Reset` may be called
		// concurrently. In that case, `struct{}{}` could be accumulated in the channel, so drain it here.
		select {
		case <-fti.stop:
		default:
		}
	}
	return active
}
