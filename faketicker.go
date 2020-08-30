package flextime

import (
	"errors"
	"sync"
	"time"
)

type fakeTicker struct {
	Timer timerIface

	ch   chan time.Time
	done chan struct{}

	dur   time.Duration
	durMu sync.RWMutex
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
		dur:   d,
	}
	go func() {
		c := ftick.Timer.C()
		for {
			select {
			case ti := <-c:
				ftick.ch <- ti
				ftick.Timer.Reset(ftick.getDur())
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

func (ftick *fakeTicker) Reset(d time.Duration) {
	ftick.setDur(d)
	ftick.Timer.Reset(ftick.getDur())
}

func (ftick *fakeTicker) setDur(d time.Duration) {
	ftick.durMu.Lock()
	defer ftick.durMu.Unlock()
	ftick.dur = d
}

func (ftick *fakeTicker) getDur() time.Duration {
	ftick.durMu.RLock()
	defer ftick.durMu.RUnlock()
	return ftick.dur
}
