package flextime

import "time"

type realClock struct{}

var _ Clock = (*realClock)(nil)

func (clock *realClock) Now() time.Time {
	return time.Now()
}

func (clock *realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (clock *realClock) Until(t time.Time) time.Duration {
	return time.Until(t)
}

func (clock *realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (clock *realClock) After(d time.Duration) <-chan time.Time {
	return clock.NewTimer(d).C()
}

func (clock *realClock) AfterFunc(d time.Duration, f func()) Timer {
	t := time.AfterFunc(d, f)
	return &realTimer{
		t: t,
	}
}

func (clock *realClock) NewTimer(d time.Duration) Timer {
	return &realTimer{
		t: time.NewTimer(d),
	}
}

func (clock *realClock) NewTicker(d time.Duration) Ticker {
	return &realTicker{
		t: time.NewTicker(d),
	}
}

func (clock *realClock) Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

type realTimer struct {
	t *time.Timer
}

func (t *realTimer) C() <-chan time.Time {
	return t.t.C
}

func (t *realTimer) Reset(d time.Duration) bool {
	return t.t.Reset(d)
}

func (t *realTimer) Stop() bool {
	return t.t.Stop()
}

type realTicker struct {
	t *time.Ticker
}

func (t *realTicker) C() <-chan time.Time {
	return t.t.C
}

func (t *realTicker) Stop() {
	t.t.Stop()
}
