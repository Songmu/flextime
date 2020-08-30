package flextime

import "time"

// Clock is an interface that implements the functions of the standard time package
type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration

	After(d time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) *Timer
	NewTimer(d time.Duration) *Timer
	NewTicker(d time.Duration) *Ticker
	Tick(d time.Duration) <-chan time.Time
}

// The Timer type represents a single event. It has same API with time.Timer
type Timer struct {
	C     <-chan time.Time
	timer timerIface
}

func createTimer(ti timerIface) *Timer {
	return &Timer{
		C:     ti.C(),
		timer: ti,
	}
}

// Stop prevents the Timer from firing.
func (ti *Timer) Stop() bool {
	return ti.timer.Stop()
}

// Reset changes the timer to expire after duration d.
func (ti *Timer) Reset(d time.Duration) bool {
	return ti.timer.Reset(d)
}

// timerIface has an interface similar to the standard time.Timer
type timerIface interface {
	C() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

// A Ticker holds a channel that delivers `ticks' of a clock at intervals.
// It has same API with time.Ticker
type Ticker struct {
	C      <-chan time.Time
	ticker tickerIface
}

func createTicker(ti tickerIface) *Ticker {
	return &Ticker{
		C:      ti.C(),
		ticker: ti,
	}
}

// Reset stops a ticker and resets its period to the specified duration
// The next tick will arrive after the new period elapses.
func (ti *Ticker) Reset(d time.Duration) {
	ti.ticker.Reset(d)
}

// Stop turns off a ticker.
func (ti *Ticker) Stop() {
	ti.ticker.Stop()
}

// tickerIface has an interface similar to the standard time.Ticker
type tickerIface interface {
	C() <-chan time.Time
	Reset(d time.Duration)
	Stop()
}
