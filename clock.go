package flextime

import "time"

// Clock is an interface that implements the functions of the standard time package
type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration

	After(d time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) Timer
	NewTimer(d time.Duration) Timer
	NewTicker(d time.Duration) Ticker
	Tick(d time.Duration) <-chan time.Time
}

// Ticker has an interface similar to the standard time.Ticker
type Ticker interface {
	C() <-chan time.Time
	Stop()
}

// Timer has an interface similar to the standard time.Timer
type Timer interface {
	C() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}
