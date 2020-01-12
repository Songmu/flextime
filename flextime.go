package flextime

import (
	"sync"
	"time"
)

var (
	backend   Clock = &realClock{}
	backendMu sync.RWMutex
)

func getBackend() Clock {
	backendMu.RLock()
	defer backendMu.RUnlock()
	return backend
}

// Switch switches backend Clock. Is it useful for testing.
func Switch(c Clock) (restore func()) {
	backendMu.Lock()
	defer backendMu.Unlock()
	orig := backend
	backend = c
	return func() { Switch(orig) }
}

// Restore the default real Clock
func Restore() {
	Switch(&realClock{})
}

// Now returns the current time from backend Clock.
func Now() time.Time {
	return getBackend().Now()
}

// Sleep pauses the current process for at least the duration d using backend Clock.
func Sleep(d time.Duration) {
	getBackend().Sleep(d)
}

// Since returns the time elapsed since t using backend Clock.
func Since(t time.Time) time.Duration {
	return getBackend().Since(t)
}

// Until returns the duration until t using backend Clock.
func Until(t time.Time) time.Duration {
	return getBackend().Until(t)
}

// After waits for the duration to elapse and then sends the current time on the returned
// channel using backend Clock.
func After(d time.Duration) <-chan time.Time {
	return getBackend().After(d)
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine using
// backend Clock. It returns a Timer that can be used to cancel the call using its Stop method.
func AfterFunc(d time.Duration, f func()) Timer {
	return getBackend().AfterFunc(d, f)
}

// NewTimer creates a new Timer that will send the current time on its channel after at
// least duration d using backend Clock.
func NewTimer(d time.Duration) Timer {
	return getBackend().NewTimer(d)
}

// NewTicker returns a new Ticker containing a channel that will send the time with a period
// specified by the duration argument using backend Clock.
func NewTicker(d time.Duration) Ticker {
	return getBackend().NewTicker(d)
}

// Tick is a convenience wrapper for NewTicker providing access to the ticking channel only
// using backend Clock.
func Tick(d time.Duration) <-chan time.Time {
	return getBackend().Tick(d)
}
