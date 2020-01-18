package flextime

import (
	"errors"
	"sync"
	"time"
)

type fakeTicker struct {
	Timer timerIface
	Dur   time.Duration

	ch   chan time.Time
	once sync.Once
	done chan struct{}
}

var _ tickerIface = (*fakeTicker)(nil)

func newFakeTicker(t timerIface, d time.Duration) *Ticker {
	if d <= 0 {
		// I don't want to panic, but the standard package is too.
		panic(errors.New("non-positive interval for NewTicker"))
	}
	return createTicker(&fakeTicker{
		Timer: t,
		Dur:   d,
		ch:    make(chan time.Time),
		done:  make(chan struct{}),
	})
}

func (ftick *fakeTicker) C() <-chan time.Time {
	ftick.once.Do(func() {
		c := ftick.Timer.C()
		go func() {
			for {
				select {
				case t := <-c:
					ftick.ch <- t
					ftick.Timer.Reset(ftick.Dur)
				case <-ftick.done:
					return
				}
			}
		}()
	})
	return ftick.ch
}

func (ftick *fakeTicker) Stop() {
	ftick.Timer.Stop()
	close(ftick.done)
}
