package flextime

import (
	"errors"
	"sync"
	"time"
)

type fakeTicker struct {
	Timer Timer
	Dur   time.Duration

	ch   chan time.Time
	once sync.Once
	done chan struct{}
}

var _ Ticker = (*fakeTicker)(nil)

func newFakeTicker(t Timer, d time.Duration) Ticker {
	if d <= 0 {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	return &fakeTicker{
		Timer: t,
		Dur:   d,
		ch:    make(chan time.Time),
		done:  make(chan struct{}),
	}
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
