package flextime_test

import (
	"testing"
	"time"

	"github.com/Songmu/flextime"
)

type virtualNS struct {
}

var _ flextime.NowSleeper = (*virtualNS)(nil)

func (vns *virtualNS) Now() time.Time {
	return time.Now()
}

func (nvs *virtualNS) Sleep(d time.Duration) {
	time.Sleep(d)
}

func TestClock_NewTimer(t *testing.T) {
	restore := flextime.Switch(flextime.NewFakeClock(&virtualNS{}))
	defer restore()

	ti := flextime.NewTimer(10 * time.Millisecond)
	if !ti.Stop() {
		t.Errorf("ti.Stop() should be true (active)")
	}
	if ti.Stop() {
		t.Errorf("ti.Stop() should be false after Stop() called once")
	}

	var blocked bool
	select {
	case <-ti.C:
	default:
		blocked = true
	}
	if !blocked {
		t.Errorf("channel of stopped Timer should be blocked")
	}

	if ti.Reset(50 * time.Millisecond) {
		t.Errorf("ti.Reset() should be false with stopped timer")
	}
	if !ti.Reset(50 * time.Millisecond) {
		t.Errorf("ti.Reset() should be true with active timer")
	}

	select {
	case <-time.After(60 * time.Millisecond):
		t.Errorf("ti.C should not be blocked but blocked")
	case <-ti.C:
	}
}
