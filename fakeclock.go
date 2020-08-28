package flextime

import "time"

// NowSleeper is, as the name implies, an interface with Now and Sleep methods.
// By simply implementing these two methods, we can create an object with a Clock
// interface by combining it with the NewFakeClock function.
type NowSleeper interface {
	Now() time.Time
	Sleep(d time.Duration)
}

type nowSleeperImpl struct {
	now   func() time.Time
	sleep func(d time.Duration)
}

var _ NowSleeper = (*nowSleeperImpl)(nil)

func (ns *nowSleeperImpl) Now() time.Time {
	return ns.now()
}

func (ns *nowSleeperImpl) Sleep(d time.Duration) {
	ns.sleep(d)
}

func newNowSleeper(now func() time.Time, sleep func(d time.Duration)) NowSleeper {
	return &nowSleeperImpl{now: now, sleep: sleep}
}

type fakeClock struct {
	ns NowSleeper
}

// NewFakeClock accepts a NowSleeper interface and returns an object with a Clock interface.
func NewFakeClock(ns NowSleeper) Clock {
	return &fakeClock{ns: ns}
}

var _ Clock = (*fakeClock)(nil)

func (fc *fakeClock) Now() time.Time {
	return fc.ns.Now()
}

func (fc *fakeClock) Sleep(d time.Duration) {
	fc.ns.Sleep(d)
}

func (fc *fakeClock) Since(t time.Time) time.Duration {
	return fc.Now().Sub(t)
}

func (fc *fakeClock) Until(t time.Time) time.Duration {
	return t.Sub(fc.Now())
}

func (fc *fakeClock) After(d time.Duration) <-chan time.Time {
	return fc.NewTimer(d).C
}

func (fc *fakeClock) AfterFunc(d time.Duration, f func()) *Timer {
	return createTimer(newFakeTimer(fc.ns, d, f))
}

func (fc *fakeClock) NewTicker(d time.Duration) *Ticker {
	ti := newFakeTimer(fc.ns, d, nil)
	ti.IsTicker = true
	return newFakeTicker(ti, d)
}

func (fc *fakeClock) NewTimer(d time.Duration) *Timer {
	return createTimer(newFakeTimer(fc.ns, d, nil))
}

func (fc *fakeClock) Tick(d time.Duration) <-chan time.Time {
	if d <= 0 {
		return nil
	}
	return fc.NewTicker(d).C
}
