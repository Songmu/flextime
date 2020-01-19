package flextime

import (
	"errors"
	"time"
)

type fakeTicker struct {
	Timer timerIface

	ch   chan time.Time
	done chan struct{}
}

var _ tickerIface = (*fakeTicker)(nil)

func newFakeTicker(t timerIface, d time.Duration) *Ticker {
	if d <= 0 {
		// I don't want to panic, but the standard package is too.
		panic(errors.New("non-positive interval for NewTicker"))
	}
	ftick := &fakeTicker{
		Timer: t,
		ch:    make(chan time.Time, 1),
		done:  make(chan struct{}),
	}
	go func() {
		c := t.C()
		for {
			select {
			case ti := <-c:
				ftick.ch <- ti
				ftick.Timer.Reset(d)
			case <-ftick.done:
				return
			}
		}
	}()
	return createTicker(ftick)
}

func (ftick *fakeTicker) C() <-chan time.Time {
	return ftick.ch
}

func (ftick *fakeTicker) Stop() {
	ftick.Timer.Stop()
	close(ftick.done)
}
